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

package main_test

import (
	"bufio"
	"os/exec"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"testing"
)

func TestHelloapp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HelloApp Suite")
}

var (
	service *gexec.Session
	port    int
)

var _ = BeforeSuite(func() {
	service, port = StartService()
})

var _ = AfterSuite(func() {
	if service != nil {
		service.Kill()
	}
})

// StartService starts the service on a random port.
// Returns the session and the port.
func StartService() (session *gexec.Session, port int) {
	exe, err := gexec.Build(
		"github.com/smartedgemec/controller-ce/cmd/helloapp",
	)
	Expect(err).ToNot(HaveOccurred(), "Problem building service")

	cmd := exec.Command(exe, "-port", "0")

	session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred(), "Problem starting service")

	Eventually(session.Err, 3).Should(gbytes.Say("listening on port:"),
		"Service did not start in time")

	// Scan the next word for the port
	scanner := bufio.NewScanner(session.Err)
	scanner.Split(bufio.ScanWords)
	scanner.Scan()
	Expect(scanner.Err()).ToNot(HaveOccurred(), "Couldn't scan for port")
	port, err = strconv.Atoi(scanner.Text())
	Expect(err).ToNot(HaveOccurred(), "Couldn't parse port")

	return session, port
}
