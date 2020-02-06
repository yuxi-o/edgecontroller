// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package nfd

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	logger "github.com/otcshare/common/log"
	"github.com/otcshare/edgecontroller/mysql"
	"github.com/otcshare/edgecontroller/pki"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"path/filepath"
	pb "sigs.k8s.io/node-feature-discovery/pkg/labeler"
	"time"
)

var log = logger.DefaultLogger.WithField("nfd-master", nil)

type ServerNFD struct {
	persistenceService *mysql.PersistenceService
	Endpoint           int
	CaCertPath         string
	CaKeyPath          string
	Sni                string
	Dsn                string
}

type labeler struct {
}

func (s ServerNFD) connectDB() error {
	db, err := sql.Open("mysql", s.Dsn)
	if err != nil {
		log.Errf("Error opening db: %v", err)
		return err
	}

	for {
		if err = db.Ping(); err != nil {
			log.Errf("DB ping failed: %v", err)
			time.Sleep(10 * time.Second)
		} else {
			log.Info("DB connection established")
			break
		}
	}

	s.persistenceService = &mysql.PersistenceService{DB: db}
	return nil
}

func (s ServerNFD) createCredentialsTLS() (credentials.TransportCredentials, error) {

	var caCert *x509.Certificate
	var caKey crypto.PrivateKey
	var err error

	caCert, err = pki.LoadCertificate(filepath.Clean(s.CaCertPath))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load CA certificate")
	}

	caKey, err = pki.LoadKey(filepath.Clean(s.CaKeyPath))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load CA key")
	}

	nfdKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate server key")
	}

	rootCA := pki.RootCA{
		Cert: caCert,
		Key:  caKey,
	}

	nfdCert, err := rootCA.NewTLSServerCert(nfdKey, s.Sni)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate server certificate")
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(caCert)

	return credentials.NewTLS(&tls.Config{
		ClientAuth: tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{nfdCert.Raw},
			PrivateKey:  nfdKey,
		}},
		ClientCAs: certPool,
	}), nil
}

func (s ServerNFD) ServeGRPC(ctx context.Context) error {

	creds, err := s.createCredentialsTLS()
	if err != nil {
		log.Errf("Failed to create TLS credentials: %v", err)
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Endpoint))
	if err != nil {
		log.Errf("net.Listen error: %+v", err)
		return err
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	labelerServer := labeler{}
	pb.RegisterLabelerServer(grpcServer, &labelerServer)

	go func() {
		<-ctx.Done()
		log.Info("Executing graceful stop")
		grpcServer.GracefulStop()
	}()

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Errf("Failed grpcServe %v", err)
		return err
	}

	log.Info("NFD Master stopped serving")
	return nil
}

func (l labeler) SetLabels(c context.Context, r *pb.SetLabelsRequest) (*pb.SetLabelsReply, error) {
	log.Infof("REQUEST Node: %s NFD-version: %s Labels: %s", r.NodeName, r.NfdVersion, r.Labels)
	return &pb.SetLabelsReply{}, nil
}
