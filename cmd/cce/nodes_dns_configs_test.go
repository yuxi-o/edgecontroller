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
	"strings"

	cce "github.com/smartedgemec/controller-ce"
	"github.com/smartedgemec/controller-ce/uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("/nodes_dns_configs", func() {
	var (
		nodeID      string
		dnsConfigID string
	)

	BeforeEach(func() {
		nodeID = postNodes()
		dnsConfigID = postDNSConfigs()
	})

	Describe("POST /nodes_dns_configs", func() {
		DescribeTable("201 Created",
			func() {
				By("Sending a POST /nodes_dns_configs request")
				resp, err := http.Post(
					"http://127.0.0.1:8080/nodes_dns_configs",
					"application/json",
					strings.NewReader(fmt.Sprintf(
						`
						{
							"node_id": "%s",
							"dns_config_id": "%s"
						}`, nodeID, dnsConfigID)))
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 201 response")
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				var respBody struct {
					ID string
				}

				By("Unmarshalling the response")
				Expect(json.Unmarshal(body, &respBody)).To(Succeed())

				By("Verifying a UUID was returned")
				Expect(uuid.IsValid(respBody.ID)).To(BeTrue())
			},
			Entry(
				"POST /nodes_dns_configs"),
		)

		DescribeTable("400 Bad Request",
			func(req, expectedResp string) {
				By("Sending a POST /nodes_dns_configs request")
				resp, err := http.Post(
					"http://127.0.0.1:8080/nodes_dns_configs",
					"application/json",
					strings.NewReader(req))
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 400 Bad Request response")
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the response body")
				Expect(string(body)).To(Equal(expectedResp))
			},
			Entry(
				"POST /nodes_dns_configs with id",
				`
				{
					"id": "123"
				}`,
				"Validation failed: id cannot be specified in POST request"),
			Entry(
				"POST /nodes_dns_configs without node_id",
				`
				{
				}`,
				"Validation failed: node_id not a valid uuid"),
			Entry(
				"POST /nodes_dns_configs without dns_config_id",
				fmt.Sprintf(`
				{
					"node_id": "%s"
				}`, uuid.New()),
				"Validation failed: dns_config_id not a valid uuid"),
		)

		DescribeTable("422 Unprocessable Entity",
			func() {
				var (
					resp *http.Response
					err  error
				)

				By("Sending a POST /nodes_dns_configs request")
				postNodesDNSConfigs(nodeID, dnsConfigID)

				By("Repeating the first POST /nodes_dns_configs request")
				resp, err = http.Post(
					"http://127.0.0.1:8080/nodes_dns_configs",
					"application/json",
					strings.NewReader(fmt.Sprintf(
						`
						{
							"node_id": "%s",
							"dns_config_id": "%s"
						}`, nodeID, dnsConfigID)))
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 422 response")
				Expect(resp.StatusCode).To(Equal(
					http.StatusUnprocessableEntity))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the response body")
				Expect(string(body)).To(Equal(fmt.Sprintf(
					"duplicate record detected for node_id %s and "+
						"dns_config_id %s",
					nodeID,
					dnsConfigID)))
			},
			Entry("POST /nodes_dns_configs with duplicate node_id"),
		)
	})

	Describe("GET /nodes_dns_configs", func() {
		var (
			nodeDNSConfigID  string
			node2ID          string
			dnsConfig2ID     string
			nodeDNSConfig2ID string
		)

		BeforeEach(func() {
			nodeDNSConfigID = postNodesDNSConfigs(nodeID, dnsConfigID)
			node2ID = postNodes()
			dnsConfig2ID = postDNSConfigs()
			nodeDNSConfig2ID = postNodesDNSConfigs(node2ID, dnsConfig2ID)
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a GET /nodes_dns_configs request")
				resp, err := http.Get(
					"http://127.0.0.1:8080/nodes_dns_configs")

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				var nodeDNSConfigs []cce.NodeDNSConfig

				By("Unmarshalling the response")
				Expect(json.Unmarshal(body, &nodeDNSConfigs)).
					To(Succeed())

				By("Verifying the 2 created node <-> DNS configs were returned")
				Expect(nodeDNSConfigs).To(ContainElement(
					cce.NodeDNSConfig{
						ID:          nodeDNSConfigID,
						NodeID:      nodeID,
						DNSConfigID: dnsConfigID,
					}))
				Expect(nodeDNSConfigs).To(ContainElement(
					cce.NodeDNSConfig{
						ID:          nodeDNSConfig2ID,
						NodeID:      node2ID,
						DNSConfigID: dnsConfig2ID,
					}))
			},
			Entry("GET /nodes_dns_configs"),
		)
	})

	Describe("GET /nodes_dns_configs/{id}", func() {
		var (
			nodeDNSConfigID string
		)

		BeforeEach(func() {
			nodeDNSConfigID = postNodesDNSConfigs(nodeID, dnsConfigID)
		})

		DescribeTable("200 OK",
			func() {
				nodeDNSConfig := getNodeDNSConfig(nodeDNSConfigID)

				By("Verifying the created node <-> DNS config was returned")
				Expect(nodeDNSConfig).To(Equal(
					&cce.NodeDNSConfig{
						ID:          nodeDNSConfigID,
						NodeID:      nodeID,
						DNSConfigID: dnsConfigID,
					},
				))
			},
			Entry("GET /nodes_dns_configs/{id}"),
		)

		DescribeTable("404 Not Found",
			func() {
				By("Sending a GET /nodes_dns_configs/{id} request")
				resp, err := http.Get(
					fmt.Sprintf(
						"http://127.0.0.1:8080/nodes_dns_configs/%s",
						uuid.New()))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("GET /nodes_dns_configs/{id} with nonexistent ID"),
		)
	})

	Describe("DELETE /nodes_dns_configs/{id}", func() {
		var (
			nodeDNSConfigID string
		)

		BeforeEach(func() {
			nodeDNSConfigID = postNodesDNSConfigs(nodeID, dnsConfigID)
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a DELETE /nodes_dns_configs/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf(
						"http://127.0.0.1:8080/nodes_dns_configs/%s",
						nodeDNSConfigID),
					nil)
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Verifying the node <-> DNS config was deleted")

				By("Sending a GET /nodes_dns_configs/{id} request")
				resp, err = http.Get(
					fmt.Sprintf(
						"http://127.0.0.1:8080/nodes_dns_configs/%s",
						nodeDNSConfigID))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("DELETE /nodes_dns_configs/{id}"),
		)

		DescribeTable("404 Not Found",
			func(id string) {
				By("Sending a DELETE /nodes_dns_configs/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf(
						"http://127.0.0.1:8080/nodes_dns_configs/%s",
						id),
					nil)
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry(
				"DELETE /nodes_dns_configs/{id} with nonexistent ID",
				uuid.New()),
		)
	})
})
