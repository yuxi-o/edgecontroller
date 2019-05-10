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

var _ = Describe("/dns_configs_dns_app_aliases", func() {
	var (
		dnsConfigID   string
		appID         string
		dnsAppAliasID string
	)

	BeforeEach(func() {
		dnsConfigID = postDNSConfigs()
		appID = postApps("container")
		dnsAppAliasID = postDNSAppAliases(appID)
	})

	Describe("POST /dns_configs_dns_app_aliases", func() {
		DescribeTable("201 Created",
			func() {
				By("Sending a POST /dns_configs_dns_app_aliases request")
				resp, err := http.Post(
					"http://127.0.0.1:8080/dns_configs_dns_app_aliases",
					"application/json",
					strings.NewReader(fmt.Sprintf(
						`
						{
							"dns_config_id": "%s",
							"dns_app_alias_id": "%s"
						}`, dnsConfigID, dnsAppAliasID)))
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
				"POST /dns_configs_dns_app_aliases"),
		)

		DescribeTable("400 Bad Request",
			func(req, expectedResp string) {
				By("Sending a POST /dns_configs_dns_app_aliases request")
				resp, err := http.Post(
					"http://127.0.0.1:8080/dns_configs_dns_app_aliases",
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
				"POST /dns_configs_dns_app_aliases with id",
				`
				{
					"id": "123"
				}`,
				"Validation failed: id cannot be specified in POST request"),
			Entry(
				"POST /dns_configs_dns_app_aliases without dns_config_id",
				`
				{
				}`,
				"Validation failed: dns_config_id not a valid uuid"),
			Entry(
				"POST /dns_configs_dns_app_aliases without dns_app_alias_id",
				fmt.Sprintf(`
				{
					"dns_config_id": "%s"
				}`, uuid.New()),
				"Validation failed: dns_app_alias_id not a valid uuid"),
		)

		DescribeTable("422 Unprocessable Entity",
			func() {
				var (
					resp *http.Response
					err  error
				)

				By("Sending a POST /dns_configs_dns_app_aliases request")
				postDNSConfigsDNSAppAliases(dnsConfigID, dnsAppAliasID)

				By("Repeating the first POST /dns_configs_dns_app_aliases request") //nolint:lll
				resp, err = http.Post(
					"http://127.0.0.1:8080/dns_configs_dns_app_aliases",
					"application/json",
					strings.NewReader(fmt.Sprintf(
						`
						{
							"dns_config_id": "%s",
							"dns_app_alias_id": "%s"
						}`, dnsConfigID, dnsAppAliasID)))
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 422 response")
				Expect(resp.StatusCode).To(Equal(
					http.StatusUnprocessableEntity))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the response body")
				Expect(string(body)).To(Equal(fmt.Sprintf(
					"duplicate record detected for dns_config_id %s and "+
						"dns_app_alias_id %s",
					dnsConfigID,
					dnsAppAliasID)))
			},
			Entry("POST /dns_configs_dns_app_aliases with duplicate dns_config_id/dns_app_alias_id"), //nolint:lll
		)
	})

	Describe("GET /dns_configs_dns_app_aliases", func() {
		var (
			dnsConfigDNSAppAliasID  string
			dnsAppAlias2ID          string
			dnsConfigDNSAppAlias2ID string
		)

		BeforeEach(func() {
			dnsConfigDNSAppAliasID =
				postDNSConfigsDNSAppAliases(dnsConfigID, dnsAppAliasID)
			dnsAppAlias2ID = postDNSAppAliases(appID)
			dnsConfigDNSAppAlias2ID =
				postDNSConfigsDNSAppAliases(dnsConfigID, dnsAppAlias2ID)
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a GET /dns_configs_dns_app_aliases request")
				resp, err := http.Get(
					"http://127.0.0.1:8080/dns_configs_dns_app_aliases")

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				var dnsConfigDNSAppAliases []cce.DNSConfigDNSAppAlias

				By("Unmarshalling the response")
				Expect(json.Unmarshal(body, &dnsConfigDNSAppAliases)).
					To(Succeed())

				By("Verifying the 2 created DNS app aliases were returned")
				Expect(dnsConfigDNSAppAliases).To(ContainElement(
					cce.DNSConfigDNSAppAlias{
						ID:            dnsConfigDNSAppAliasID,
						DNSConfigID:   dnsConfigID,
						DNSAppAliasID: dnsAppAliasID,
					}))
				Expect(dnsConfigDNSAppAliases).To(ContainElement(
					cce.DNSConfigDNSAppAlias{
						ID:            dnsConfigDNSAppAlias2ID,
						DNSConfigID:   dnsConfigID,
						DNSAppAliasID: dnsAppAlias2ID,
					}))
			},
			Entry("GET /dns_configs_dns_app_aliases"),
		)
	})

	Describe("GET /dns_configs_dns_app_aliases/{id}", func() {
		var (
			dnsConfigDNSAppAliasID string
		)

		BeforeEach(func() {
			dnsConfigDNSAppAliasID =
				postDNSConfigsDNSAppAliases(dnsConfigID, dnsAppAliasID)
		})

		DescribeTable("200 OK",
			func() {
				dnsConfigDNSAppAlias :=
					getDNSConfigsDNSAppAlias(
						dnsConfigDNSAppAliasID)

				By("Verifying the created DNS config <-> DNS app alias was returned") //nolint:lll
				Expect(dnsConfigDNSAppAlias).To(Equal(
					&cce.DNSConfigDNSAppAlias{
						ID:            dnsConfigDNSAppAliasID,
						DNSConfigID:   dnsConfigID,
						DNSAppAliasID: dnsAppAliasID,
					},
				))
			},
			Entry("GET /dns_configs_dns_app_aliases/{id}"),
		)

		DescribeTable("404 Not Found",
			func() {
				By("Sending a GET /dns_configs_dns_app_aliases/{id} request")
				resp, err := http.Get(
					fmt.Sprintf(
						"http://127.0.0.1:8080/dns_configs_dns_app_aliases/%s",
						uuid.New()))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("GET /dns_configs_dns_app_aliases/{id} with nonexistent ID"),
		)
	})

	Describe("DELETE /dns_configs_dns_app_aliases/{id}", func() {
		var (
			dnsConfigDNSAppAliasID string
		)

		BeforeEach(func() {
			dnsConfigDNSAppAliasID =
				postDNSConfigsDNSAppAliases(dnsConfigID, dnsAppAliasID)
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a DELETE /dns_configs_dns_app_aliases/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf(
						"http://127.0.0.1:8080/dns_configs_dns_app_aliases/%s",
						dnsConfigDNSAppAliasID),
					nil)
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Verifying the DNS config <-> DNS app alias was deleted")

				By("Sending a GET /dns_configs_dns_app_aliases/{id} request")
				resp, err = http.Get(
					fmt.Sprintf(
						"http://127.0.0.1:8080/dns_configs_dns_app_aliases/%s",
						dnsConfigDNSAppAliasID))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("DELETE /dns_configs_dns_app_aliases/{id}"),
		)

		DescribeTable("404 Not Found",
			func(id string) {
				By("Sending a DELETE /dns_configs_dns_app_aliases/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf(
						"http://127.0.0.1:8080/dns_configs_dns_app_aliases/%s",
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
				"DELETE /dns_configs_dns_app_aliases/{id} with nonexistent ID",
				uuid.New()),
		)
	})
})
