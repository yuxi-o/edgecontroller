// Copyright 2019 Intel Corporation. All rights reserved.
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

package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"os"
	"testing"
	"time"
)

var (
	Iserv InterfaceServiceServer //fake server

	HelpOut = `
    Get or attach/detach network interfaces to OVS on remote edge node

    -endpoint      Endpoint to be requested
    -servicename   Name to be used as server name for TLS handshake
    -cmd           Supported commands: get, attach, detach
    -val           PCI address for attach and detach commands. Multiple addresses can be passed
                   and must be separated by commas: -val=0000:00:00.0,0000:00:00.1
    -certsdir      Directory where cert.pem and key.pem for client and root.pem for CA resides   
    -timeout       Timeout value [s] for grpc requests

	`

	WarningOut = "Unrecognized action: " + "test123\n" + HelpOut
)

func TestCli(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "interfaceservicecli suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(GinkgoWriter)

	CertsDir := "./certs"
	err := os.MkdirAll(CertsDir, os.ModePerm)
	Expect(err).ShouldNot(HaveOccurred())

	Expect(prepareTestCredentials(CertsDir)).ToNot(HaveOccurred())
	Iserv = InterfaceServiceServer{
		Endpoint: "localhost:2020",
	}
	Expect(Iserv.StartServer()).ToNot(HaveOccurred())
	time.Sleep(1 * time.Second)
})

var _ = AfterSuite(func() {
	err := os.RemoveAll("./certs")
	Expect(err).ShouldNot(HaveOccurred())
	Expect(Iserv.GracefulStop()).ToNot(HaveOccurred())
})
