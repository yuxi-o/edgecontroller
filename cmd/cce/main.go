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
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	cce "github.com/smartedgemec/controller-ce"
	"github.com/smartedgemec/controller-ce/gorilla"
	"github.com/smartedgemec/controller-ce/grpc"
	"github.com/smartedgemec/controller-ce/http"
	"github.com/smartedgemec/controller-ce/mysql"
	"github.com/smartedgemec/controller-ce/pki"
)

const certsDir = "./certificates"

func main() {
	var (
		err error

		// flags
		dsn      string
		httpPort int
		grpcPort int

		rootCA *pki.RootCA

		db         *sql.DB
		controller *cce.Controller
		nodeMap    map[string]*cce.Node = make(map[string]*cce.Node)

		httpListener net.Listener
		g            *gorilla.Gorilla
		httpServer   *http.Server

		grpcListener net.Listener
		grpcServer   *grpc.Server
	)

	log.Print("Controller CE starting")

	// CLI flags
	flag.StringVar(&dsn, "dsn", "", "Data source name")
	flag.IntVar(&httpPort, "httpPort", 8080, "Controller HTTP port")
	flag.IntVar(&grpcPort, "grpcPort", 8081, "Controller gRPC port")
	flag.Parse()

	// Connect to the db
	if db, err = sql.Open("mysql", dsn); err != nil {
		log.Fatal("Error opening db: ", err)
	}

	// Verify connection
	if err = db.Ping(); err != nil {
		log.Fatal("DB ping failed: ", err)
	}

	log.Print("DB connection established")

	if rootCA, err = pki.InitRootCA(
		filepath.Join(certsDir, "ca"),
	); err != nil {
		log.Fatal("Error initializing Controller CA: ", err)
	}

	log.Print("Initialized Controller CA")

	controller = &cce.Controller{
		PersistenceService: &mysql.PersistenceService{DB: db},
		AuthorityService:   rootCA,
	}

	// Setup http server tcp listener
	if httpListener, err = net.Listen(
		"tcp",
		fmt.Sprintf(":%d", httpPort),
	); err != nil {
		log.Fatal("Could not listen on : ", err)
	}
	defer httpListener.Close()

	// Setup grpc server tcp listener
	if grpcListener, err = net.Listen(
		"tcp",
		fmt.Sprintf(":%d", grpcPort),
	); err != nil {
		log.Fatal("Could not listen on : ", err)
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
	g = gorilla.NewGorilla(controller, nodeMap)

	log.Println("HTTP handler ready")

	// Configure http server
	httpServer = http.NewServer(g)

	// Start the http server
	log.Printf("Starting HTTP server on port %d\n", httpPort)
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
			log.Println("HTTP graceful shutdown exceeded timeout, using force")
			httpServer.Close()
		}
	}()

	// Configure grpc server
	grpcServer = grpc.NewServer(controller)

	// Start the grpc server
	log.Printf("Starting gRPC server on port %d\n", httpPort)
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
			log.Println("gRPC server shutdown exceeded timeout, using force")
			grpcServer.Stop()
		case <-stopped:
			return
		}
	}()

	log.Println("Controller CE ready")
	if err = eg.Wait(); err != nil {
		log.Fatal(err)
	}
}
