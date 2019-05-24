// Copyright 2019 Smart-Edge.com, Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"

	"golang.org/x/sync/errgroup"

	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	cce "github.com/smartedgemec/controller-ce"
	"github.com/smartedgemec/controller-ce/gorilla"
	"github.com/smartedgemec/controller-ce/grpc"
	"github.com/smartedgemec/controller-ce/http"
	"github.com/smartedgemec/controller-ce/jose"
	"github.com/smartedgemec/controller-ce/k8s"
	"github.com/smartedgemec/controller-ce/mysql"
	"github.com/smartedgemec/controller-ce/pki"
	logger "github.com/smartedgemec/log"
)

const certsDir = "./certificates"

var log = logger.DefaultLogger.WithField("pkg", "main")

// CLI flags
var (
	dsn       string
	adminPass string
	logLevel  string
	httpPort  int
	grpcPort  int
	orchMode  string
	k8sClient k8s.Client
)

func init() {
	flag.StringVar(&dsn, "dsn", "", "Data source name")
	flag.StringVar(&adminPass, "adminPass", "", "Admin user password")
	flag.StringVar(&logLevel, "log-level", "info", "Syslog level")
	flag.IntVar(&httpPort, "httpPort", 8080, "Controller HTTP port")
	flag.IntVar(&grpcPort, "grpcPort", 8081, "Controller gRPC port")

	// application orchestration mode
	flag.StringVar(&orchMode, "orchestration-mode", "native", "Orchestration mode. options [native, kubernetes] ")

	// k8s
	flag.StringVar(&k8sClient.CAFile, "k8s-client-ca-path", "", "Kubernetes root certificate path")
	flag.StringVar(&k8sClient.CertFile, "k8s-client-cert-path", "", "Kubernetes client certificate path")
	flag.StringVar(&k8sClient.KeyFile, "k8s-client-key-path", "", "Kubernetes client private key path")
	flag.StringVar(&k8sClient.Host, "k8s-master-host", "", "Kubernetes master host")
	flag.StringVar(&k8sClient.APIPath, "k8s-api-path", "", "Kubernetes api path")
	flag.StringVar(&k8sClient.Username, "k8s-master-user", "", "Kubernetes default user")
}

func main() { // nolint: gocyclo
	flag.Parse()
	log.Info("Controller CE starting")

	// Set log level
	lvl, err := logger.ParseLevel(logLevel)
	if err != nil {
		log.Alert("Bad log level %q: %v", logLevel, err)
		os.Exit(1)
	}
	logger.SetLevel(lvl)

	// Setup orchestrator
	var orchestrationMode cce.OrchestrationMode
	switch orchMode {
	case "native":
		orchestrationMode = cce.OrchestrationModeNative
	case "kubernetes":
		orchestrationMode = cce.OrchestrationModeKubernetes
		if err = k8sClient.Ping(); err != nil {
			log.Alertf("Error configuring kubernetes client: %v", err)
			os.Exit(1)
		}
	default:
		log.Alertf("Invalid orchestration mode %q", orchMode)
		os.Exit(1)
	}

	// Connect to the db and verify
	if adminPass == "" {
		log.Alert("User admin password cannot be empty")
		os.Exit(1)
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Alertf("Error opening db: %v", err)
		os.Exit(1)
	}
	if err = db.Ping(); err != nil {
		log.Alertf("DB ping failed: %v", err)
		os.Exit(1)
	}
	log.Info("DB connection established")

	// Initialize self-signed root CA
	rootCA, err := pki.InitRootCA(filepath.Join(certsDir, "ca"))
	if err != nil {
		log.Alertf("Error initializing Controller CA: %v", err)
		os.Exit(1)
	}

	log.Info("Initialized Controller CA")

	// Print self-signed Controller CA. This is used to manually configure the
	// Appliance by adding the Controller to its trust anchor pool for TLS
	// connections.
	//
	// TODO: Replace printing to STDERR with writing to a file or making the
	// certificate available via an HTTP endpoint.
	log.Info("Root CA:\n%s", string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: rootCA.Cert.Raw,
		},
	)))

	// Generate a key for signing authentication tokens. The key is only stored
	// in memory and will be re-generated upon Controller restart.
	//
	// TODO: Persist the key to avoid having API/UI users to have to login and
	// get a new token every time the Controller is restarted.
	joseKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		log.Alertf("error generating token signing key: %v", err)
		os.Exit(1)
	}

	controller := &cce.Controller{
		PersistenceService: &mysql.PersistenceService{DB: db},
		AuthorityService:   rootCA,
		TokenService: &jose.JWSTokenIssuer{
			Key:          joseKey,
			KeyAlgorithm: "ES384",
		},
		AdminCreds: &cce.AuthCreds{
			Username: "admin",
			Password: adminPass,
		},
		OrchestrationMode: orchestrationMode,
		KubernetesClient:  &k8sClient,
	}

	// Setup http server tcp listener
	httpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", httpPort))
	if err != nil {
		log.Alertf("Could not listen on port %d: %v", httpPort, err)
		os.Exit(1)
	}
	defer httpListener.Close()

	// Setup grpc server tcp listener
	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Alertf("Could not listen on port %d: %v", grpcPort, err)
		os.Exit(1)
	}
	defer grpcListener.Close()

	// Create an error group to manage server goroutines
	eg, ctx := errgroup.WithContext(context.Background())

	// Catch exit signals
	eg.Go(func() error {
		ch := make(chan os.Signal, 2)

		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case signal := <-ch:
			return errors.New(signal.String())
		}
	})

	// Create the gorilla and feed it a controller and its nodes
	koko := gorilla.NewGorilla(controller)

	log.Info("HTTP handler ready")

	// Configure http server
	httpServer := http.NewServer(koko)

	// Start the http server
	log.Infof("Starting HTTP server on port %d", httpPort)
	eg.Go(func() error {
		return httpServer.Serve(httpListener)
	})

	// Shutdown http server on exit signal
	go func() {
		<-ctx.Done()

		ctxShutdown, cancel := context.WithTimeout(context.TODO(), time.Minute)
		defer cancel()

		err = httpServer.Shutdown(ctxShutdown)
		if err != nil {
			log.Info("HTTP graceful shutdown exceeded timeout, using force")
			httpServer.Close()
		}
	}()

	// Configure grpc server
	serverConf, err := newTLSConf(rootCA, grpc.SNI)
	if err != nil {
		log.Alertf("Error creating TLS config for gRPC server: %v", err)
		os.Exit(1)
	}
	serverConf.NextProtos = []string{"h2"}
	serverConf.ClientAuth = tls.RequireAndVerifyClientCert
	enrollmentConf, err := newTLSConf(rootCA, grpc.EnrollmentSNI)
	if err != nil {
		log.Alertf("Error creating TLS config for gRPC server: %v", err)
		os.Exit(1)
	}
	enrollmentConf.NextProtos = []string{"h2"}
	enrollmentConf.ClientAuth = tls.NoClientCert
	grpcServer := grpc.NewServer(controller, &tls.Config{
		GetConfigForClient: func(
			hello *tls.ClientHelloInfo,
		) (*tls.Config, error) {
			switch hello.ServerName {
			case grpc.SNI:
				return serverConf, nil
			case grpc.EnrollmentSNI:
				return enrollmentConf, nil
			default:
				return nil, fmt.Errorf("unexpected server name: %s", hello.ServerName)
			}
		},
	})

	// Start the grpc server
	log.Infof("Starting gRPC server on port %d", grpcPort)
	eg.Go(func() error {
		return grpcServer.Serve(grpcListener)
	})

	// Shutdown grpc server on exit signal
	go func() {
		<-ctx.Done()

		// Try to gracefully shutdown
		stopped := make(chan struct{})
		go func() {
			grpcServer.GracefulStop()
			close(stopped)
		}()

		select {
		case <-time.After(time.Minute):
			log.Info("gRPC server shutdown exceeded timeout, using force")
			grpcServer.Stop()
		case <-stopped:
			return
		}
	}()

	log.Info("Controller CE ready")
	if err := eg.Wait(); err != nil {
		log.Alert(err)
		os.Exit(1)
	}
}

func newTLSConf(rootCA *pki.RootCA, sni string) (*tls.Config, error) {
	tlsKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("error generating TLS key: %v", err)
	}
	tlsCert, err := rootCA.NewTLSServerCert(tlsKey, sni)
	if err != nil {
		return nil, fmt.Errorf("error generating TLS cert: %v", err)
	}
	tlsCAChain, err := rootCA.CAChain()
	if err != nil {
		return nil, fmt.Errorf("error getting TLS CA chain: %v", err)
	}
	tlsChain := [][]byte{tlsCert.Raw}
	for _, caCert := range tlsCAChain {
		tlsChain = append(tlsChain, caCert.Raw)
	}
	tlsRoots := x509.NewCertPool()
	tlsRoots.AddCert(tlsCAChain[len(tlsCAChain)-1])
	return &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: tlsChain,
			PrivateKey:  tlsKey,
			Leaf:        tlsCert,
		}},
		ClientCAs: tlsRoots,
	}, nil
}
