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

package clients_test

import (
	"bufio"
	"context"
	"os/exec"
	"strconv"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/smartedgemec/controller-ce/grpc"
	gclients "github.com/smartedgemec/controller-ce/grpc/clients"
)

func TestApplicationClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Application Client Suite")
}

var (
	// server
	service *gexec.Session
	port    int

	// client
	ctx                   = context.Background()
	conn                  *grpc.ClientConn
	appDeploySvcCli       *gclients.ApplicationDeploymentServiceClient
	appLifeSvcCli         *gclients.ApplicationLifecycleServiceClient
	appPolicySvcCli       *gclients.ApplicationPolicyServiceClient
	vnfDeploySvcCli       *gclients.VNFDeploymentServiceClient
	vnfLifeSvcCli         *gclients.VNFLifecycleServiceClient
	interfaceSvcCli       *gclients.InterfaceServiceClient
	zoneSvcCli            *gclients.ZoneServiceClient
	interfacePolicySvcCli *gclients.InterfacePolicyServiceClient
)

var _ = BeforeSuite(func() {
	By("Building the mock grpc server")
	exe, err := gexec.Build(
		"github.com/smartedgemec/controller-ce/test/node/grpc",
	)
	Expect(err).ToNot(HaveOccurred(), "Problem building service")

	By("Starting the server on a random port")
	service, err = gexec.Start(
		exec.Command(exe, "-port", "0"), GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred(), "Problem starting service")
	Eventually(service.Err, 3).Should(gbytes.Say("listening on port:"),
		"Service did not start in time")

	By("Scanning for the server's port")
	scanner := bufio.NewScanner(service.Err)
	scanner.Split(bufio.ScanWords)
	scanner.Scan()
	Expect(scanner.Err()).ToNot(HaveOccurred(), "Couldn't scan for port")

	By("Parsing the server's port")
	port, err = strconv.Atoi(scanner.Text())
	Expect(err).ToNot(HaveOccurred(), "Couldn't parse port")

	By("Dialing the server")
	conn, err = grpc.Dial(ctx, "127.0.0.1", port)
	Expect(err).ToNot(HaveOccurred(), "Dial failed: %v", err)

	By("Creating the clients")
	appDeploySvcCli = gclients.NewApplicationDeploymentServiceClient(conn)
	appLifeSvcCli = gclients.NewApplicationLifecycleServiceClient(conn)
	appPolicySvcCli = gclients.NewApplicationPolicyServiceClient(conn)
	vnfDeploySvcCli = gclients.NewVNFDeploymentServiceClient(conn)
	vnfLifeSvcCli = gclients.NewVNFLifecycleServiceClient(conn)
	interfaceSvcCli = gclients.NewInterfaceServiceClient(conn)
	zoneSvcCli = gclients.NewZoneServiceClient(conn)
	interfacePolicySvcCli = gclients.NewInterfacePolicyServiceClient(conn)
})

var _ = AfterSuite(func() {
	if service != nil {
		By("Stopping the service")
		service.Kill()
	}
	if conn != nil {
		By("Closing the client connection")
		conn.Close()
	}
})
