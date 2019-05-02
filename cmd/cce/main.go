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
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"path/filepath"

	cce "github.com/smartedgemec/controller-ce"
	"github.com/smartedgemec/controller-ce/gorilla"
	"github.com/smartedgemec/controller-ce/http"
	"github.com/smartedgemec/controller-ce/mysql"
	"github.com/smartedgemec/controller-ce/pki"
)

const certsDir = "./certificates"

func main() {
	var (
		err error

		// flags
		dsn  string
		port int

		rootCA *pki.RootCA

		db         *sql.DB
		controller *cce.Controller
		nodeMap    map[string]*cce.Node = make(map[string]*cce.Node)

		listener net.Listener
		g        *gorilla.Gorilla
		srv      *http.Server
	)

	log.Print("Controller CE starting")

	// CLI flags
	flag.StringVar(&dsn, "dsn", "", "Data source name")
	flag.IntVar(&port, "port", 8080, "Host port")
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

	// Listen on a local network address
	if listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal("Could not listen on : ", err)
	}
	defer listener.Close()

	log.Printf("Listener ready on tcp port %d", port)

	// Create the gorilla and feed it a controller and its nodes
	g = gorilla.NewGorilla(controller, nodeMap)

	log.Print("Handler ready, starting server")

	// Start the server
	srv = http.NewServer(g)
	log.Fatal(srv.Serve(listener))
}
