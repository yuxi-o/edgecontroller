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

var _ = Describe("/dns_configs", func() {
	Describe("POST /dns_configs", func() {
		DescribeTable("201 Created",
			func(req string) {
				By("Sending a POST /dns_configs request")
				resp, err := http.Post(
					"http://127.0.0.1:8080/dns_configs",
					"application/json",
					strings.NewReader(req))
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
				"POST /dns_configs",
				`
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
				}`),
		)

		DescribeTable("400 Bad Request",
			func(req, expectedResp string) {
				By("Sending a POST /dns_configs request")
				resp, err := http.Post(
					"http://127.0.0.1:8080/dns_configs",
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
				"POST /dns_configs with id",
				`
				{
					"id": "123"
				}`,
				"Validation failed: id cannot be specified in POST request"),
			Entry(
				"POST /dns_configs without name",
				`
				{
				}`,
				"Validation failed: name cannot be empty"),
			Entry(
				"POST /dns_configs without a_records|forwarders",
				`
				{
					"name": "dns config 123"
				}`,
				"Validation failed: a_records|forwarders cannot both be empty"),
			Entry(
				"POST /dns_configs without a_records[0].name",
				`
				{
					"name": "dns config 123",
					"a_records": [{
					}]
				}`,
				"Validation failed: a_records[0].name cannot be empty"),
			Entry(
				"POST /dns_configs without a_records[0].description",
				`
				{
					"name": "dns config 123",
					"a_records": [{
						"name": "a record 1"
					}]
				}`,
				"Validation failed: a_records[0].description cannot be empty"),
			Entry(
				"POST /dns_configs with invalid a_records[0].ips[0]",
				`
				{
					"name": "dns config 123",
					"a_records": [{
						"name": "a record 1",
						"description": "description 1",
						"ips": [
							"1724.16.55.43",
							"172.16.55.44"
						]
					}]
				}`,
				"Validation failed: a_records[0].ips[0] could not be parsed"),
			Entry(
				"POST /dns_configs without forwarders[0].name",
				`
				{
					"name": "dns config 123",
					"forwarders": [{
					}]
				}`,
				"Validation failed: forwarders[0].name cannot be empty"),
			Entry(
				"POST /dns_configs without forwarders[0].description",
				`
				{
					"name": "dns config 123",
					"forwarders": [{
						"name": "forwarder 1"
					}]
				}`,
				"Validation failed: forwarders[0].description cannot be empty"),
			Entry(
				"POST /dns_configs with invalid forwarders[0].ip",
				`
				{
					"name": "dns config 123",
					"forwarders": [{
						"name": "forwarder 1",
						"description": "description 1",
						"ip": "888.8.8.8"
					}]
				}`,
				"Validation failed: forwarders[0].ip could not be parsed"),
		)
	})

	Describe("GET /dns_configs", func() {
		var (
			dnsConfigID  string
			dnsConfig2ID string
		)

		BeforeEach(func() {
			dnsConfigID = postDNSConfigs()
			dnsConfig2ID = postDNSConfigs()
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a GET /dns_configs request")
				resp, err := http.Get("http://127.0.0.1:8080/dns_configs")

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				var dnsConfigs []cce.DNSConfig

				By("Unmarshalling the response")
				Expect(json.Unmarshal(body, &dnsConfigs)).To(Succeed())

				By("Verifying the 2 created DNS configs were returned")
				Expect(dnsConfigs).To(ContainElement(
					cce.DNSConfig{
						ID:   dnsConfigID,
						Name: "dns config 123",
						ARecords: []*cce.DNSARecord{
							{
								Name:        "a record 1",
								Description: "description 1",
								IPs: []string{
									"172.16.55.43",
									"172.16.55.44",
								},
							},
						},
						Forwarders: []*cce.DNSForwarder{
							{
								Name:        "forwarder 1",
								Description: "description 1",
								IP:          "8.8.8.8",
							},
							{
								Name:        "forwarder 2",
								Description: "description 2",
								IP:          "1.1.1.1",
							},
						},
					}))
				Expect(dnsConfigs).To(ContainElement(
					cce.DNSConfig{

						ID:   dnsConfig2ID,
						Name: "dns config 123",
						ARecords: []*cce.DNSARecord{
							{
								Name:        "a record 1",
								Description: "description 1",
								IPs: []string{
									"172.16.55.43",
									"172.16.55.44",
								},
							},
						},
						Forwarders: []*cce.DNSForwarder{
							{
								Name:        "forwarder 1",
								Description: "description 1",
								IP:          "8.8.8.8",
							},
							{
								Name:        "forwarder 2",
								Description: "description 2",
								IP:          "1.1.1.1",
							},
						},
					}))
			},
			Entry("GET /dns_configs"),
		)
	})

	Describe("GET /dns_configs/{id}", func() {
		var (
			dnsConfigID string
		)

		BeforeEach(func() {
			dnsConfigID = postDNSConfigs()
		})

		DescribeTable("200 OK",
			func() {
				dnsConfig := getDNSConfig(dnsConfigID)

				By("Verifying the created DNS config was returned")
				Expect(dnsConfig).To(Equal(
					&cce.DNSConfig{
						ID:   dnsConfigID,
						Name: "dns config 123",
						ARecords: []*cce.DNSARecord{
							{
								Name:        "a record 1",
								Description: "description 1",
								IPs: []string{
									"172.16.55.43",
									"172.16.55.44",
								},
							},
						},
						Forwarders: []*cce.DNSForwarder{
							{
								Name:        "forwarder 1",
								Description: "description 1",
								IP:          "8.8.8.8",
							},
							{
								Name:        "forwarder 2",
								Description: "description 2",
								IP:          "1.1.1.1",
							},
						},
					},
				))
			},
			Entry("GET /dns_configs/{id}"),
		)

		DescribeTable("404 Not Found",
			func() {
				By("Sending a GET /dns_configs/{id} request")
				resp, err := http.Get(
					fmt.Sprintf("http://127.0.0.1:8080/dns_configs/%s",
						uuid.New()))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("GET /dns_configs/{id} with nonexistent ID"),
		)
	})

	Describe("PATCH /dns_configs", func() {
		var (
			dnsConfigID string
		)

		BeforeEach(func() {
			dnsConfigID = postDNSConfigs()
		})

		DescribeTable("204 No Content",
			func(reqStr string, expectedConfig *cce.DNSConfig) {
				By("Sending a PATCH /dns_configs request")
				req, err := http.NewRequest(
					http.MethodPatch,
					"http://127.0.0.1:8080/dns_configs",
					strings.NewReader(fmt.Sprintf(reqStr, dnsConfigID)))
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 204 No Content response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				By("Getting the updated DNS config")
				updatedConfig := getDNSConfig(dnsConfigID)

				By("Verifying the DNS config was updated")
				expectedConfig.SetID(dnsConfigID)
				Expect(updatedConfig).To(Equal(expectedConfig))
			},
			Entry(
				"PATCH /dns_configs/{id}",
				`
				[
					{
						"id": "%s",
						"name": "dns config 123456",
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
					}
				]`,
				&cce.DNSConfig{
					Name: "dns config 123456",
					ARecords: []*cce.DNSARecord{
						{
							Name:        "a record 1",
							Description: "description 1",
							IPs: []string{
								"172.16.55.43",
								"172.16.55.44",
							},
						},
					},
					Forwarders: []*cce.DNSForwarder{
						{
							Name:        "forwarder 1",
							Description: "description 1",
							IP:          "8.8.8.8",
						},
						{
							Name:        "forwarder 2",
							Description: "description 2",
							IP:          "1.1.1.1",
						},
					},
				}),
		)

		DescribeTable("400 Bad Request",
			func(reqStr string, expectedResp string) {
				By("Sending a PATCH /dns_configs request")
				if strings.Contains(reqStr, "%s") {
					reqStr = fmt.Sprintf(reqStr, dnsConfigID)
				}
				req, err := http.NewRequest(
					http.MethodPatch,
					"http://127.0.0.1:8080/dns_configs",
					strings.NewReader(reqStr))
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)

				By("Verifying a 400 Bad Request")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the response body")
				Expect(string(body)).To(Equal(expectedResp))
			},
			// Don't repeat all the validation testing we did in POST, just
			// one for ID and another one as a sanity check.
			Entry(
				"PATCH /dns_configs without id",
				`
				[{}]`,
				"Validation failed: id not a valid uuid"),
			Entry(
				"PATCH /dns_configs without name",
				`
				[
					{
						"id": "%s"
					}
				]`,
				"Validation failed: name cannot be empty"),
		)
	})

	Describe("DELETE /dns_configs/{id}", func() {
		var (
			dnsConfigID string
		)

		BeforeEach(func() {
			dnsConfigID = postDNSConfigs()
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a DELETE /dns_configs/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("http://127.0.0.1:8080/dns_configs/%s",
						dnsConfigID),
					nil)
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Verifying the DNS config was deleted")

				By("Sending a GET /dns_configs/{id} request")
				resp, err = http.Get(
					fmt.Sprintf("http://127.0.0.1:8080/dns_configs/%s",
						dnsConfigID))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("DELETE /dns_configs/{id}"),
		)

		DescribeTable("404 Not Found",
			func(id string) {
				By("Sending a DELETE /dns_configs/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("http://127.0.0.1:8080/dns_configs/%s", id),
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
				"DELETE /dns_configs/{id} with nonexistent ID",
				uuid.New()),
		)

		DescribeTable("422 Unprocessable Entity",
			func(resource, expectedResp string) {
				switch resource {
				case "dns_configs_dns_app_aliases":
					postDNSConfigsDNSAppAliases(dnsConfigID,
						postDNSAppAliases(postApps("container")))
				case "dns_configs_dns_vnf_aliases":
					postDNSConfigsDNSVNFAliases(dnsConfigID,
						postDNSVNFAliases(postVNFs("container")))
				}

				By("Sending a DELETE /dns_configs/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("http://127.0.0.1:8080/dns_configs/%s",
						dnsConfigID),
					nil)
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 422 response")
				Expect(resp.StatusCode).To(Equal(
					http.StatusUnprocessableEntity))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the response body")
				Expect(string(body)).To(Equal(
					fmt.Sprintf(expectedResp, dnsConfigID)))
			},
			Entry(
				"DELETE /dns_configs/{id} with dns_configs_dns_app_aliases record", //nolint:lll
				"dns_configs_dns_app_aliases",
				"cannot delete dns_config_id %s: record in use in dns_configs_dns_app_aliases", //nolint:lll
			),
			Entry(
				"DELETE /dns_configs/{id} with dns_configs_dns_vnf_aliases record", //nolint:lll
				"dns_configs_dns_vnf_aliases",
				"cannot delete dns_config_id %s: record in use in dns_configs_dns_vnf_aliases", //nolint:lll
			),
		)
	})
})
