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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"testing"

	cce "github.com/smartedgemec/controller-ce"

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
	service, port = startService()
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

func startService() (session *gexec.Session, port int) {
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

func postContainerApps() (id string) {
	By("Sending a POST /container_apps request")
	resp, err := http.Post(
		"http://127.0.0.1:8080/container_apps",
		"application/json",
		strings.NewReader(`
            {
                "name": "container app",
                "vendor": "smart edge",
                "description": "my container app",
                "image": "http://www.test.com/my_container_app.tar.gz",
                "cores": 4,
                "memory": 1024
            }`))

	By("Verifying a 201 Created response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var respBody struct {
		ID string
	}

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &respBody)).To(Succeed())

	return respBody.ID
}

func getContainerApp(id string) *cce.ContainerApp {
	By("Sending a GET /container_apps/{id} request")
	resp, err := http.Get(
		fmt.Sprintf("http://127.0.0.1:8080/container_apps/%s", id))

	By("Verifying a 200 OK response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var containerApp cce.ContainerApp

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &containerApp)).To(Succeed())

	return &containerApp
}

func postContainerVNFs() (id string) {
	By("Sending a POST /container_vnfs request")
	resp, err := http.Post(
		"http://127.0.0.1:8080/container_vnfs",
		"application/json",
		strings.NewReader(`
            {
                "name": "container vnf",
                "vendor": "smart edge",
                "description": "my container vnf",
                "image": "http://www.test.com/my_container_vnf.tar.gz",
                "cores": 4,
                "memory": 1024
            }`))

	By("Verifying a 201 Created response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var respBody struct {
		ID string
	}

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &respBody)).To(Succeed())

	return respBody.ID
}

func getContainerVNF(id string) *cce.ContainerVNF {
	By("Sending a GET /container_vnfs/{id} request")
	resp, err := http.Get(
		fmt.Sprintf("http://127.0.0.1:8080/container_vnfs/%s", id))

	By("Verifying a 200 OK response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var containerVNF cce.ContainerVNF

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &containerVNF)).To(Succeed())

	return &containerVNF
}

func postDNSConfigs() (id string) {
	By("Sending a POST /dns_configs request")
	resp, err := http.Post(
		"http://127.0.0.1:8080/dns_configs",
		"application/json",
		strings.NewReader(`
            {
                "name": "dns config 123",
                "a_records": [{
                    "name": "a record 1",
                    "description": "description 1",
                    "ips": [
                        "172.16.55.43",
                        "172.16.55.44"
                    ]
                }],
                "forwarders": [{
                    "name": "forwarder 1",
                    "description": "description 1",
                    "ip": "8.8.8.8"
                }, {
                    "name": "forwarder 2",
                    "description": "description 2",
                    "ip": "1.1.1.1"
                }]
            }`))

	By("Verifying a 201 Created response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var respBody struct {
		ID string
	}

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &respBody)).To(Succeed())

	return respBody.ID
}

func getDNSConfig(id string) *cce.DNSConfig {
	By("Sending a GET /dns_configs/{id} request")
	resp, err := http.Get(
		fmt.Sprintf("http://127.0.0.1:8080/dns_configs/%s", id))

	By("Verifying a 200 OK response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var dnsConfig cce.DNSConfig

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &dnsConfig)).To(Succeed())

	return &dnsConfig
}

func postDNSContainerAppAliases(containerAppID string) (id string) {
	By("Sending a POST /dns_container_app_aliases request")
	resp, err := http.Post(
		"http://127.0.0.1:8080/dns_container_app_aliases",
		"application/json",
		strings.NewReader(fmt.Sprintf(`
            {
                "name": "dns container app alias 123",
                "description": "description 1",
                "container_app_id": "%s"
            }`, containerAppID)))

	By("Verifying a 201 Created response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var respBody struct {
		ID string
	}

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &respBody)).To(Succeed())

	return respBody.ID
}

func getDNSContainerAppAlias(id string) *cce.DNSContainerAppAlias {
	By("Sending a GET /dns_container_app_aliases/{id} request")
	resp, err := http.Get(
		fmt.Sprintf("http://127.0.0.1:8080/dns_container_app_aliases/%s", id))

	By("Verifying a 200 OK response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var dnsContainerAppAlias cce.DNSContainerAppAlias

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &dnsContainerAppAlias)).To(Succeed())

	return &dnsContainerAppAlias
}

func postDNSContainerVNFAliases(containerVNFID string) (id string) {
	By("Sending a POST /dns_container_vnf_aliases request")
	resp, err := http.Post(
		"http://127.0.0.1:8080/dns_container_vnf_aliases",
		"application/json",
		strings.NewReader(fmt.Sprintf(`
            {
                "name": "dns container vnf alias 123",
                "description": "description 1",
                "container_vnf_id": "%s"
            }`, containerVNFID)))

	By("Verifying a 201 Created response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var respBody struct {
		ID string
	}

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &respBody)).To(Succeed())

	return respBody.ID
}

func getDNSContainerVNFAlias(id string) *cce.DNSContainerVNFAlias {
	By("Sending a GET /dns_container_vnf_aliases/{id} request")
	resp, err := http.Get(
		fmt.Sprintf("http://127.0.0.1:8080/dns_container_vnf_aliases/%s", id))

	By("Verifying a 200 OK response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var dnsContainerVNFAlias cce.DNSContainerVNFAlias

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &dnsContainerVNFAlias)).To(Succeed())

	return &dnsContainerVNFAlias
}

func postDNSVMAppAliases(vmAppID string) (id string) {
	By("Sending a POST /dns_vm_app_aliases request")
	resp, err := http.Post(
		"http://127.0.0.1:8080/dns_vm_app_aliases",
		"application/json",
		strings.NewReader(fmt.Sprintf(`
            {
                "name": "dns vm app alias 123",
                "description": "description 1",
                "vm_app_id": "%s"
            }`, vmAppID)))

	By("Verifying a 201 Created response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var respBody struct {
		ID string
	}

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &respBody)).To(Succeed())

	return respBody.ID
}

func getDNSVMAppAlias(id string) *cce.DNSVMAppAlias {
	By("Sending a GET /dns_vm_app_aliases/{id} request")
	resp, err := http.Get(
		fmt.Sprintf("http://127.0.0.1:8080/dns_vm_app_aliases/%s", id))

	By("Verifying a 200 OK response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var dnsVMAppAlias cce.DNSVMAppAlias

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &dnsVMAppAlias)).To(Succeed())

	return &dnsVMAppAlias
}

func postDNSVMVNFAliases(vmVNFID string) (id string) {
	By("Sending a POST /dns_vm_vnf_aliases request")
	resp, err := http.Post(
		"http://127.0.0.1:8080/dns_vm_vnf_aliases",
		"application/json",
		strings.NewReader(fmt.Sprintf(`
            {
                "name": "dns vm vnf alias 123",
                "description": "description 1",
                "vm_vnf_id": "%s"
            }`, vmVNFID)))

	By("Verifying a 201 Created response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var respBody struct {
		ID string
	}

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &respBody)).To(Succeed())

	return respBody.ID
}

func getDNSVMVNFAlias(id string) *cce.DNSVMVNFAlias {
	By("Sending a GET /dns_vm_vnf_aliases/{id} request")
	resp, err := http.Get(
		fmt.Sprintf("http://127.0.0.1:8080/dns_vm_vnf_aliases/%s", id))

	By("Verifying a 200 OK response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var dnsVMVNFAlias cce.DNSVMVNFAlias

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &dnsVMVNFAlias)).To(Succeed())

	return &dnsVMVNFAlias
}

func postNodes() (id string) {
	By("Sending a POST /nodes request")
	resp, err := http.Post(
		"http://127.0.0.1:8080/nodes",
		"application/json",
		strings.NewReader(`
            {
                "name": "node123",
                "location": "smart edge lab",
                "serial": "abc123"
            }`))

	By("Verifying a 201 Created response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var respBody struct {
		ID string
	}

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &respBody)).To(Succeed())

	return respBody.ID
}

func getNode(id string) *cce.Node {
	By("Sending a GET /nodes/{id} request")
	resp, err := http.Get(
		fmt.Sprintf("http://127.0.0.1:8080/nodes/%s", id))

	By("Verifying a 200 OK response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var node cce.Node

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &node)).To(Succeed())

	return &node
}

func postTrafficPolicies() (id string) {
	By("Sending a POST /traffic_policies request")
	resp, err := http.Post(
		"http://127.0.0.1:8080/traffic_policies",
		"application/json",
		strings.NewReader(`
        {
            "rules": [{
                "description": "test-rule-1",
                "priority": 1,
                "source": {
                    "description": "test-source-1",
                    "macs": {
                        "mac_addresses": [
                            "F0-59-8E-7B-36-8A",
                            "23-20-8E-15-89-D1",
                            "35-A4-38-73-35-45"
                        ]
                    },
                    "ip": {
                        "address": "223.1.1.0",
                        "mask": 16,
                        "begin_port": 2000,
                        "end_port": 2012,
                        "protocol": "tcp"
                    },
                    "gtp": {
                        "address": "10.6.7.2",
                        "mask": 12,
                        "imsis": [
                            "310150123456789",
                            "310150123456790",
                            "310150123456791"
                        ]
                    }
                },
                "destination": {
                    "description": "test-destination-1",
                    "macs": {
                        "mac_addresses": [
                            "7D-C2-3A-1C-63-D9",
                            "E9-6B-D1-D2-1A-6B",
                            "C8-32-A9-43-85-55"
                        ]
                    },
                    "ip": {
                        "address": "64.1.1.0",
                        "mask": 16,
                        "begin_port": 1000,
                        "end_port": 1012,
                        "protocol": "tcp"
                    },
                    "gtp": {
                        "address": "108.6.7.2",
                        "mask": 4,
                        "imsis": [
                            "310150123456792",
                            "310150123456793",
                            "310150123456794"
                        ]
                    }
                },
                "target": {
                    "description": "test-target-1",
                    "action": "accept",
                    "mac": {
                        "mac_address": "C7-5A-E7-98-1B-A3"
                    },
                    "ip": {
                        "address": "123.2.3.4",
                        "port": 1600
                    }
                }
            }]
        }`))

	By("Verifying a 201 Created response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var respBody struct {
		ID string
	}

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &respBody)).To(Succeed())

	return respBody.ID
}

func getTrafficPolicy(id string) *cce.TrafficPolicy {
	By("Sending a GET /traffic_policies/{id} request")
	resp, err := http.Get(
		fmt.Sprintf("http://127.0.0.1:8080/traffic_policies/%s", id))

	By("Verifying a 200 OK response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var trafficPolicy cce.TrafficPolicy

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &trafficPolicy)).To(Succeed())

	return &trafficPolicy
}

func postVMApps() (id string) {
	By("Sending a POST /vm_apps request")
	resp, err := http.Post(
		"http://127.0.0.1:8080/vm_apps",
		"application/json",
		strings.NewReader(`
            {
                "name": "vm app",
                "vendor": "smart edge",
                "description": "my vm app",
                "image": "http://www.test.com/my_vm_app.tar.gz",
                "cores": 4,
                "memory": 1024
            }`))

	By("Verifying a 201 Created response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var respBody struct {
		ID string
	}

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &respBody)).To(Succeed())

	return respBody.ID
}

func getVMApp(id string) *cce.VMApp {
	By("Sending a GET /vm_apps/{id} request")
	resp, err := http.Get(
		fmt.Sprintf("http://127.0.0.1:8080/vm_apps/%s", id))

	By("Verifying a 200 OK response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var vmApp cce.VMApp

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &vmApp)).To(Succeed())

	return &vmApp
}

func postVMVNFs() (id string) {
	By("Sending a POST /vm_vnfs request")
	resp, err := http.Post(
		"http://127.0.0.1:8080/vm_vnfs",
		"application/json",
		strings.NewReader(`
            {
                "name": "vm vnf",
                "vendor": "smart edge",
                "description": "my vm vnf",
                "image": "http://www.test.com/my_vm_vnf.tar.gz",
                "cores": 4,
                "memory": 1024
            }`))

	By("Verifying a 201 Created response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var respBody struct {
		ID string
	}

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &respBody)).To(Succeed())

	return respBody.ID
}

func getVMVNF(id string) *cce.VMVNF {
	By("Sending a GET /vm_vnfs/{id} request")
	resp, err := http.Get(
		fmt.Sprintf("http://127.0.0.1:8080/vm_vnfs/%s", id))

	By("Verifying a 200 OK response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var vmVNF cce.VMVNF

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &vmVNF)).To(Succeed())

	return &vmVNF
}
