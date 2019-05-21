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

var _ = Describe("/dns_configs_app_aliases", func() {
	var (
		dnsConfigID string
		appID       string
	)

	BeforeEach(func() {
		dnsConfigID = postDNSConfigs()
		appID = postApps("container")
	})

	Describe("POST /dns_configs_app_aliases", func() {
		DescribeTable("201 Created",
			func() {
				By("Sending a POST /dns_configs_app_aliases request")
				resp, err := apiCli.Post(
					"http://127.0.0.1:8080/dns_configs_app_aliases",
					"application/json",
					strings.NewReader(fmt.Sprintf(
						`
						{
							"dns_config_id": "%s",
							"name": "dns config app alias",
							"description": "my dns config app alias",
							"app_id": "%s"
						}`, dnsConfigID, appID)))
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
				"POST /dns_configs_app_aliases"),
		)

		DescribeTable("400 Bad Request",
			func(req, expectedResp string) {
				By("Sending a POST /dns_configs_app_aliases request")
				resp, err := apiCli.Post(
					"http://127.0.0.1:8080/dns_configs_app_aliases",
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
				"POST /dns_configs_app_aliases with id",
				`
				{
					"id": "123"
				}`,
				"Validation failed: id cannot be specified in POST request"),
			Entry(
				"POST /dns_configs_app_aliases without dns_config_id",
				`
				{
				}`,
				"Validation failed: dns_config_id not a valid uuid"),
			Entry(
				"POST /dns_configs_app_aliases without name",
				fmt.Sprintf(`
				{
					"dns_config_id": "%s"
				}`, uuid.New()),
				"Validation failed: name cannot be empty"),
			Entry(
				"POST /dns_configs_app_aliases without description",
				fmt.Sprintf(`
				{
					"dns_config_id": "%s",
					"name": "dns config app alias"
				}`, uuid.New()),
				"Validation failed: description cannot be empty"),
			Entry(
				"POST /dns_configs_app_aliases without app_id",
				fmt.Sprintf(`
				{
					"dns_config_id": "%s",
					"name": "dns config app alias",
					"description": "my dns config app alias"
				}`, uuid.New()),
				"Validation failed: app_id not a valid uuid"),
		)

		DescribeTable("422 Unprocessable Entity",
			func() {
				var (
					resp *http.Response
					err  error
				)

				By("Sending a POST /dns_configs_app_aliases request")
				postDNSConfigsAppAliases(dnsConfigID, appID)

				By("Repeating the first POST /dns_configs_app_aliases request")
				resp, err = apiCli.Post(
					"http://127.0.0.1:8080/dns_configs_app_aliases",
					"application/json",
					strings.NewReader(fmt.Sprintf(
						`
						{
							"dns_config_id": "%s",
							"name": "dns config app alias",
							"description": "my dns config app alias",
							"app_id": "%s"
						}`, dnsConfigID, appID)))
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
						"app_id %s",
					dnsConfigID,
					appID)))
			},
			Entry("POST /dns_configs_app_aliases with duplicate dns_config_id/app_id"),
		)
	})

	Describe("GET /dns_configs_app_aliases", func() {
		var (
			dnsConfigAppAliasID  string
			app2ID               string
			dnsConfigAppAlias2ID string
		)

		BeforeEach(func() {
			dnsConfigAppAliasID = postDNSConfigsAppAliases(dnsConfigID, appID)
			app2ID = postApps("container")
			dnsConfigAppAlias2ID = postDNSConfigsAppAliases(dnsConfigID, app2ID)
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a GET /dns_configs_app_aliases request")
				resp, err := apiCli.Get(
					"http://127.0.0.1:8080/dns_configs_app_aliases")
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 200 OK response")
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				var dnsConfigAppAliases []*cce.DNSConfigAppAlias

				By("Unmarshalling the response")
				Expect(json.Unmarshal(body, &dnsConfigAppAliases)).
					To(Succeed())

				By("Verifying the 2 created app aliases were returned")
				Expect(dnsConfigAppAliases).To(ContainElement(
					&cce.DNSConfigAppAlias{
						ID:          dnsConfigAppAliasID,
						DNSConfigID: dnsConfigID,
						Name:        "dns config app alias",
						Description: "my dns config app alias",
						AppID:       appID,
					}))
				Expect(dnsConfigAppAliases).To(ContainElement(
					&cce.DNSConfigAppAlias{
						ID:          dnsConfigAppAlias2ID,
						DNSConfigID: dnsConfigID,
						Name:        "dns config app alias",
						Description: "my dns config app alias",
						AppID:       app2ID,
					}))
			},
			Entry("GET /dns_configs_app_aliases"),
		)
	})

	Describe("GET /dns_configs_app_aliases/{id}", func() {
		var (
			dnsConfigAppAliasID string
		)

		BeforeEach(func() {
			dnsConfigAppAliasID =
				postDNSConfigsAppAliases(dnsConfigID, appID)
		})

		DescribeTable("200 OK",
			func() {
				dnsConfigAppAlias := getDNSConfigsAppAlias(dnsConfigAppAliasID)

				By("Verifying the created DNS config app alias was returned")
				Expect(dnsConfigAppAlias).To(Equal(
					&cce.DNSConfigAppAlias{
						ID:          dnsConfigAppAliasID,
						DNSConfigID: dnsConfigID,
						Name:        "dns config app alias",
						Description: "my dns config app alias",
						AppID:       appID,
					},
				))
			},
			Entry("GET /dns_configs_app_aliases/{id}"),
		)

		DescribeTable("404 Not Found",
			func() {
				By("Sending a GET /dns_configs_app_aliases/{id} request")
				resp, err := apiCli.Get(
					fmt.Sprintf(
						"http://127.0.0.1:8080/dns_configs_app_aliases/%s",
						uuid.New()))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 404 Not Found response")
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("GET /dns_configs_app_aliases/{id} with nonexistent ID"),
		)
	})

	Describe("DELETE /dns_configs_app_aliases/{id}", func() {
		var (
			dnsConfigAppAliasID string
		)

		BeforeEach(func() {
			dnsConfigAppAliasID =
				postDNSConfigsAppAliases(dnsConfigID, appID)
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a DELETE /dns_configs_app_aliases/{id} request")
				resp, err := apiCli.Delete(
					fmt.Sprintf(
						"http://127.0.0.1:8080/dns_configs_app_aliases/%s",
						dnsConfigAppAliasID))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 200 OK response")
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Verifying the DNS config app alias was deleted")

				By("Sending a GET /dns_configs_app_aliases/{id} request")
				resp, err = apiCli.Get(
					fmt.Sprintf(
						"http://127.0.0.1:8080/dns_configs_app_aliases/%s",
						dnsConfigAppAliasID))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 404 Not Found response")
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("DELETE /dns_configs_app_aliases/{id}"),
		)

		DescribeTable("404 Not Found",
			func(id string) {
				By("Sending a DELETE /dns_configs_app_aliases/{id} request")
				resp, err := apiCli.Delete(
					fmt.Sprintf(
						"http://127.0.0.1:8080/dns_configs_app_aliases/%s",
						id))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 404 Not Found response")
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry(
				"DELETE /dns_configs_app_aliases/{id} with nonexistent ID",
				uuid.New()),
		)
	})
})
