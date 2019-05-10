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

var _ = Describe("/dns_vnf_aliases", func() {
	var (
		vnfID string
	)

	BeforeEach(func() {
		vnfID = postVNFs("container")
	})

	Describe("POST /dns_vnf_aliases", func() {
		DescribeTable("201 Created",
			func(req string) {
				By("Sending a POST /dns_vnf_aliases request")
				resp, err := http.Post(
					"http://127.0.0.1:8080/dns_vnf_aliases",
					"application/json",
					strings.NewReader(fmt.Sprintf(req, vnfID)))
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
				"POST /dns_vnf_aliases",
				`
				{
					"name": "dns vnf alias 123",
					"description": "description 1",
					"vnf_id": "%s"
				}`),
		)

		DescribeTable("400 Bad Request",
			func(req, expectedResp string) {
				By("Sending a POST /dns_vnf_aliases request")
				resp, err := http.Post(
					"http://127.0.0.1:8080/dns_vnf_aliases",
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
				"POST /dns_vnf_aliases with id",
				`
				{
					"id": "123"
				}`,
				"Validation failed: id cannot be specified in POST request"),
			Entry(
				"POST /dns_vnf_aliases without name",
				`
				{
					"description": "description 1",
					"vnf_id": "123"
				}`,
				"Validation failed: name cannot be empty"),
			Entry(
				"POST /dns_vnf_aliases without description",
				`
				{
					"name": "dns vnf alias 123",
					"vnf_id": "123"
				}`,
				"Validation failed: description cannot be empty"),
			Entry(
				"POST /dns_vnf_aliases without vnf_id",
				`
				{
					"name": "dns vnf alias 123",
					"description": "description 1"
				}`,
				"Validation failed: vnf_id not a valid uuid"),
		)
	})

	Describe("GET /dns_vnf_aliases", func() {
		var (
			dnsVNFAliasID  string
			dnsVNFAlias2ID string
		)

		BeforeEach(func() {
			dnsVNFAliasID = postDNSVNFAliases(vnfID)
			dnsVNFAlias2ID = postDNSVNFAliases(vnfID)
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a GET /dns_vnf_aliases request")
				resp, err := http.Get(
					"http://127.0.0.1:8080/dns_vnf_aliases")

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				var dnsVNFAliases []cce.DNSVNFAlias

				By("Unmarshalling the response")
				Expect(json.Unmarshal(body, &dnsVNFAliases)).
					To(Succeed())

				By("Verifying the 2 created DNS VNF aliases were returned")
				Expect(dnsVNFAliases).To(ContainElement(
					cce.DNSVNFAlias{
						ID:          dnsVNFAliasID,
						Name:        "dns vnf alias 123",
						Description: "description 1",
						VNFID:       vnfID,
					}))
				Expect(dnsVNFAliases).To(ContainElement(
					cce.DNSVNFAlias{
						ID:          dnsVNFAlias2ID,
						Name:        "dns vnf alias 123",
						Description: "description 1",
						VNFID:       vnfID,
					}))
			},
			Entry("GET /dns_vnf_aliases"),
		)
	})

	Describe("GET /dns_vnf_aliases/{id}", func() {
		var (
			dnsVNFAliasID string
		)

		BeforeEach(func() {
			dnsVNFAliasID = postDNSVNFAliases(vnfID)
		})

		DescribeTable("200 OK",
			func() {
				dnsVNFAlias := getDNSVNFAlias(
					dnsVNFAliasID)

				By("Verifying the created DNS VNF alias was returned")
				Expect(dnsVNFAlias).To(Equal(
					&cce.DNSVNFAlias{
						ID:          dnsVNFAliasID,
						Name:        "dns vnf alias 123",
						Description: "description 1",
						VNFID:       vnfID,
					},
				))
			},
			Entry("GET /dns_vnf_aliases/{id}"),
		)

		DescribeTable("404 Not Found",
			func() {
				By("Sending a GET /dns_vnf_aliases/{id} request")
				resp, err := http.Get(
					fmt.Sprintf(
						"http://127.0.0.1:8080/dns_vnf_aliases/%s",
						uuid.New()))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("GET /dns_vnf_aliases/{id} with nonexistent ID"),
		)
	})

	Describe("PATCH /dns_vnf_aliases", func() {
		var (
			dnsVNFAliasID string
		)

		BeforeEach(func() {
			dnsVNFAliasID = postDNSVNFAliases(vnfID)
		})

		DescribeTable("204 No Content",
			func(reqStr string, expectedAlias *cce.DNSVNFAlias) {
				By("Sending a PATCH /dns_vnf_aliases request")
				req, err := http.NewRequest(
					http.MethodPatch,
					"http://127.0.0.1:8080/dns_vnf_aliases",
					strings.NewReader(fmt.Sprintf(reqStr,
						dnsVNFAliasID, vnfID)))
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 204 No Content response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				By("Getting the updated VNF alias")
				updatedAlias := getDNSVNFAlias(dnsVNFAliasID)

				By("Verifying the DNS VNF alias was updated")
				expectedAlias.SetID(dnsVNFAliasID)
				expectedAlias.VNFID = vnfID
				Expect(updatedAlias).To(Equal(expectedAlias))
			},
			Entry(
				"PATCH /dns_vnf_aliases",
				`
				[
					{
						"id": "%s",
						"name": "dns vnf alias 123456",
						"description": "description 1",
						"vnf_id": "%s"
					}
				]`,
				&cce.DNSVNFAlias{
					Name:        "dns vnf alias 123456",
					Description: "description 1",
				}),
		)

		DescribeTable("400 Bad Request",
			func(reqStr string, expectedResp string) {
				By("Sending a PATCH /dns_vnf_aliases request")
				if strings.Contains(reqStr, "%s") {
					reqStr = fmt.Sprintf(reqStr, dnsVNFAliasID)
				}
				req, err := http.NewRequest(
					http.MethodPatch,
					"http://127.0.0.1:8080/dns_vnf_aliases",
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
			Entry(
				"PATCH /dns_vnf_aliases without id",
				`
				[{}]`,
				"Validation failed: id not a valid uuid"),
			Entry(
				"PATCH /dns_vnf_aliases without name",
				`
				[
					{
						"id": "%s",
						"description": "description 1",
						"vnf_id": "123"
					}
				]`,
				"Validation failed: name cannot be empty"),
			Entry(
				"PATCH /dns_vnf_aliases without description",
				`
				[
					{
						"id": "%s",
						"name": "dns vnf alias 123",
						"vnf_id": "123"
					}
				]`,
				"Validation failed: description cannot be empty"),
			Entry(
				"PATCH /dns_vnf_aliases without vnf_id",
				`
				[
					{
						"id": "%s",
						"name": "dns vnf alias 123",
						"description": "description 1"
					}
				]`,
				"Validation failed: vnf_id not a valid uuid"),
		)
	})

	Describe("DELETE /dns_vnf_aliases/{id}", func() {
		var (
			dnsVNFAliasID string
		)

		BeforeEach(func() {
			dnsVNFAliasID = postDNSVNFAliases(vnfID)
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a DELETE /dns_vnf_aliases/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf(
						"http://127.0.0.1:8080/dns_vnf_aliases/%s",
						dnsVNFAliasID),
					nil)
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Verifying the DNS VNF alias was deleted")

				By("Sending a GET /dns_vnf_aliases/{id} request")
				resp, err = http.Get(
					fmt.Sprintf(
						"http://127.0.0.1:8080/dns_vnf_aliases/%s",
						dnsVNFAliasID))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("DELETE /dns_vnf_aliases/{id}"),
		)

		DescribeTable("404 Not Found",
			func(id string) {
				By("Sending a DELETE /dns_vnf_aliases/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf(
						"http://127.0.0.1:8080/dns_vnf_aliases/%s",
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
				"DELETE /dns_vnf_aliases/{id} with nonexistent ID",
				uuid.New()),
		)

		DescribeTable("422 Unprocessable Entity",
			func() {
				postDNSConfigsDNSVNFAliases(
					postDNSConfigs(),
					dnsVNFAliasID)

				By("Sending a DELETE /dns_vnf_aliases/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf(
						"http://127.0.0.1:8080/dns_vnf_aliases/%s",
						dnsVNFAliasID),
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
				Expect(string(body)).To(Equal(fmt.Sprintf(
					"cannot delete dns_vnf_alias_id %s: record in "+
						"use in dns_configs_dns_vnf_aliases",
					dnsVNFAliasID)))
			},
			Entry("DELETE /dns_vnf_aliases/{id} with dns_configs_dns_vnf_aliases record"), //nolint:lll
		)
	})
})
