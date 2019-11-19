// Copyright 2019 Intel Corporation. All rights reserved
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
	"flag"
	"fmt"
	"os"

	"github.com/otcshare/edgecontroller/edgednscli"
)

func main() {
	addr := flag.String("address", ":4204", "EdgeDNS API address")
	set := flag.String("set", "",
		"Path to JSON file containing HostRecordSet for set operation")
	del := flag.String("del", "",
		"Path to JSON file containing RecordSet for del operation")

	pkiCrtPath := flag.String("cert", "certs/cert.pem", "PKI Cert Path")
	pkiKeyPath := flag.String("key", "certs/key.pem", "PKI Key Path")
	pkiCAPath := flag.String("ca", "certs/root.pem", "PKI CA Path")
	serverNameOverride := flag.String("name", "",
		"PKI Server Name to override while grpc connection")

	flag.Parse()

	pki := cli.PKIPaths{
		CrtPath:            *pkiCrtPath,
		KeyPath:            *pkiKeyPath,
		CAPath:             *pkiCAPath,
		ServerNameOverride: *serverNameOverride}

	cfg := cli.AppFlags{
		Address: *addr,
		Set:     *set,
		Del:     *del,
		PKI:     &pki}

	if cfg.Set == "" && cfg.Del == "" {
		fmt.Println("No 'set' or 'del' command specified. Please use -h or -help")
		os.Exit(-1)
	}

	if err := cli.ExecuteCommands(&cfg); err != nil {
		fmt.Printf("Execution failed: %v\n", err)
		os.Exit(-1)
	}
}
