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
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

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

func TestApplicationClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller CE API Suite")
}

// StartService starts the service on a random port.
// Returns the session and the port.
func StartService() (session *gexec.Session, port int) {
	exe, err := gexec.Build("github.com/smartedgemec/controller-ce/cmd/cce")
	Expect(err).ToNot(HaveOccurred(), "Problem building service")

	cmd := exec.Command(exe,
		"-dsn", "root:beer@tcp(:8083)/controller_ce",
		"-port", "8080")

	session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred(), "Problem starting service")

	Eventually(session.Err, 3).Should(gbytes.Say(
		"Handler ready, starting server"),
		"Service did not start in time")

	return session, port
}
