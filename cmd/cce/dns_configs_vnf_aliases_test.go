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

var _ = Describe("/dns_configs_vnf_aliases", func() {
	var (
		dnsConfigID string
		vnfID       string
	)

	BeforeEach(func() {
		dnsConfigID = postDNSConfigs()
		vnfID = postVNFs("container")
	})

	Describe("POST /dns_configs_vnf_aliases", func() {
		DescribeTable("201 Created",
			func() {
				By("Sending a POST /dns_configs_vnf_aliases request")
				resp, err := apiCli.Post(
					"http://127.0.0.1:8080/dns_configs_vnf_aliases",
					"application/json",
					strings.NewReader(fmt.Sprintf(
						`
						{
							"dns_config_id": "%s",
							"name": "dns config vnf alias",
							"description": "my dns config vnf alias",
							"vnf_id": "%s"
						}`, dnsConfigID, vnfID)))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

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
				"POST /dns_configs_vnf_aliases"),
		)

		DescribeTable("400 Bad Request",
			func(req, expectedResp string) {
				By("Sending a POST /dns_configs_vnf_aliases request")
				resp, err := apiCli.Post(
					"http://127.0.0.1:8080/dns_configs_vnf_aliases",
					"application/json",
					strings.NewReader(req))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 400 Bad Request response")
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the response body")
				Expect(string(body)).To(Equal(expectedResp))
			},
			Entry(
				"POST /dns_configs_vnf_aliases with id",
				`
				{
					"id": "123"
				}`,
				"Validation failed: id cannot be specified in POST request"),
			Entry(
				"POST /dns_configs_vnf_aliases without dns_config_id",
				`
				{
				}`,
				"Validation failed: dns_config_id not a valid uuid"),
			Entry(
				"POST /dns_configs_vnf_aliases without name",
				fmt.Sprintf(`
				{
					"dns_config_id": "%s"
				}`, uuid.New()),
				"Validation failed: name cannot be empty"),
			Entry(
				"POST /dns_configs_vnf_aliases without description",
				fmt.Sprintf(`
				{
					"dns_config_id": "%s",
					"name": "dns config vnf alias"
				}`, uuid.New()),
				"Validation failed: description cannot be empty"),
			Entry(
				"POST /dns_configs_vnf_aliases without vnf_id",
				fmt.Sprintf(`
				{
					"dns_config_id": "%s",
					"name": "dns config vnf alias",
					"description": "my dns config vnf alias"
				}`, uuid.New()),
				"Validation failed: vnf_id not a valid uuid"),
		)

		DescribeTable("422 Unprocessable Entity",
			func() {
				var (
					resp *http.Response
					err  error
				)

				By("Sending a POST /dns_configs_vnf_aliases request")
				postDNSConfigsVNFAliases(dnsConfigID, vnfID)

				By("Repeating the first POST /dns_configs_vnf_aliases request")
				resp, err = apiCli.Post(
					"http://127.0.0.1:8080/dns_configs_vnf_aliases",
					"application/json",
					strings.NewReader(fmt.Sprintf(
						`
						{
							"dns_config_id": "%s",
							"name": "dns config vnf alias",
							"description": "my dns config vnf alias",
							"vnf_id": "%s"
						}`, dnsConfigID, vnfID)))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 422 response")
				Expect(resp.StatusCode).To(Equal(
					http.StatusUnprocessableEntity))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the response body")
				Expect(string(body)).To(Equal(fmt.Sprintf(
					"duplicate record detected for dns_config_id %s and "+
						"vnf_id %s",
					dnsConfigID,
					vnfID)))
			},
			Entry("POST /dns_configs_vnf_aliases with duplicate dns_config_id/vnf_id"),
		)
	})

	Describe("GET /dns_configs_vnf_aliases", func() {
		var (
			dnsConfigVNFAliasID  string
			vnf2ID               string
			dnsConfigVNFAlias2ID string
		)

		BeforeEach(func() {
			dnsConfigVNFAliasID = postDNSConfigsVNFAliases(dnsConfigID, vnfID)
			vnf2ID = postVNFs("container")
			dnsConfigVNFAlias2ID = postDNSConfigsVNFAliases(dnsConfigID, vnf2ID)
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a GET /dns_configs_vnf_aliases request")
				resp, err := apiCli.Get(
					"http://127.0.0.1:8080/dns_configs_vnf_aliases")
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 200 OK response")
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				var dnsConfigVNFAliases []*cce.DNSConfigVNFAlias

				By("Unmarshalling the response")
				Expect(json.Unmarshal(body, &dnsConfigVNFAliases)).
					To(Succeed())

				By("Verifying the 2 created VNF aliases were returned")
				Expect(dnsConfigVNFAliases).To(ContainElement(
					&cce.DNSConfigVNFAlias{
						ID:          dnsConfigVNFAliasID,
						DNSConfigID: dnsConfigID,
						Name:        "dns config vnf alias",
						Description: "my dns config vnf alias",
						VNFID:       vnfID,
					}))
				Expect(dnsConfigVNFAliases).To(ContainElement(
					&cce.DNSConfigVNFAlias{
						ID:          dnsConfigVNFAlias2ID,
						DNSConfigID: dnsConfigID,
						Name:        "dns config vnf alias",
						Description: "my dns config vnf alias",
						VNFID:       vnf2ID,
					}))
			},
			Entry("GET /dns_configs_vnf_aliases"),
		)
	})

	Describe("GET /dns_configs_vnf_aliases/{id}", func() {
		var (
			dnsConfigVNFAliasID string
		)

		BeforeEach(func() {
			dnsConfigVNFAliasID =
				postDNSConfigsVNFAliases(dnsConfigID, vnfID)
		})

		DescribeTable("200 OK",
			func() {
				dnsConfigVNFAlias := getDNSConfigsVNFAlias(dnsConfigVNFAliasID)

				By("Verifying the created DNS config VNF alias was returned")
				Expect(dnsConfigVNFAlias).To(Equal(
					&cce.DNSConfigVNFAlias{
						ID:          dnsConfigVNFAliasID,
						DNSConfigID: dnsConfigID,
						Name:        "dns config vnf alias",
						Description: "my dns config vnf alias",
						VNFID:       vnfID,
					},
				))
			},
			Entry("GET /dns_configs_vnf_aliases/{id}"),
		)

		DescribeTable("404 Not Found",
			func() {
				By("Sending a GET /dns_configs_vnf_aliases/{id} request")
				resp, err := apiCli.Get(
					fmt.Sprintf(
						"http://127.0.0.1:8080/dns_configs_vnf_aliases/%s",
						uuid.New()))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 404 Not Found response")
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("GET /dns_configs_vnf_aliases/{id} with nonexistent ID"),
		)
	})

	Describe("DELETE /dns_configs_vnf_aliases/{id}", func() {
		var (
			dnsConfigVNFAliasID string
		)

		BeforeEach(func() {
			dnsConfigVNFAliasID =
				postDNSConfigsVNFAliases(dnsConfigID, vnfID)
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a DELETE /dns_configs_vnf_aliases/{id} request")
				resp, err := apiCli.Delete(
					fmt.Sprintf(
						"http://127.0.0.1:8080/dns_configs_vnf_aliases/%s",
						dnsConfigVNFAliasID))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 200 OK response")
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Verifying the DNS config VNF alias was deleted")

				By("Sending a GET /dns_configs_vnf_aliases/{id} request")
				resp, err = apiCli.Get(
					fmt.Sprintf(
						"http://127.0.0.1:8080/dns_configs_vnf_aliases/%s",
						dnsConfigVNFAliasID))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 404 Not Found response")
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("DELETE /dns_configs_vnf_aliases/{id}"),
		)

		DescribeTable("404 Not Found",
			func(id string) {
				By("Sending a DELETE /dns_configs_vnf_aliases/{id} request")
				resp, err := apiCli.Delete(
					fmt.Sprintf(
						"http://127.0.0.1:8080/dns_configs_vnf_aliases/%s",
						id))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 404 Not Found response")
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry(
				"DELETE /dns_configs_vnf_aliases/{id} with nonexistent ID",
				uuid.New()),
		)
	})
})
