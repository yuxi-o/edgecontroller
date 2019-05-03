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

var _ = Describe("/dns_vm_vnf_aliases", func() {
	var (
		vmVNFID string
	)

	BeforeEach(func() {
		vmVNFID = postVMVNFs()
	})

	Describe("POST /dns_vm_vnf_aliases", func() {
		DescribeTable("201 Created",
			func(req string) {
				By("Sending a POST /dns_vm_vnf_aliases request")
				resp, err := http.Post(
					"http://127.0.0.1:8080/dns_vm_vnf_aliases",
					"application/json",
					strings.NewReader(fmt.Sprintf(req, vmVNFID)))
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
				"POST /dns_vm_vnf_aliases",
				`
                {
                    "name": "dns vm vnf alias 123",
                    "description": "description 1",
                    "vm_vnf_id": "%s"
                }`),
		)

		DescribeTable("400 Bad Request",
			func(req, expectedResp string) {
				By("Sending a POST /dns_vm_vnf_aliases request")
				resp, err := http.Post(
					"http://127.0.0.1:8080/dns_vm_vnf_aliases",
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
				"POST /dns_vm_vnf_aliases with id",
				`
                {
                    "id": "123"
                }`,
				"Validation failed: id cannot be specified in POST request"),
			Entry(
				"POST /dns_vm_vnf_aliases without name",
				`
                {
                    "description": "description 1",
                    "vm_vnf_id": "123"
                }`,
				"Validation failed: name cannot be empty"),
			Entry(
				"POST /dns_vm_vnf_aliases without description",
				`
                {
                    "name": "dns vm vnf alias 123",
                    "vm_vnf_id": "123"
                }`,
				"Validation failed: description cannot be empty"),
			Entry(
				"POST /dns_vm_vnf_aliases without vm_vnf_id",
				`
                {
                    "name": "dns vm vnf alias 123",
                    "description": "description 1"
                }`,
				"Validation failed: vm_vnf_id not a valid uuid"),
		)
	})

	Describe("GET /dns_vm_vnf_aliases", func() {
		var (
			dnsVMVNFAliasID  string
			dnsVMVNFAlias2ID string
		)

		BeforeEach(func() {
			dnsVMVNFAliasID = postDNSVMVNFAliases(vmVNFID)
			dnsVMVNFAlias2ID = postDNSVMVNFAliases(vmVNFID)
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a GET /dns_vm_vnf_aliases request")
				resp, err := http.Get(
					"http://127.0.0.1:8080/dns_vm_vnf_aliases")

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				var dnsVMVNFAliases []cce.DNSVMVNFAlias

				By("Unmarshalling the response")
				Expect(json.Unmarshal(body, &dnsVMVNFAliases)).
					To(Succeed())

				By("Verifying the 2 created DNS vm vnf aliases were returned") //nolint:lll
				Expect(dnsVMVNFAliases).To(ContainElement(
					cce.DNSVMVNFAlias{
						ID:          dnsVMVNFAliasID,
						Name:        "dns vm vnf alias 123",
						Description: "description 1",
						VMVNFID:     vmVNFID,
					}))
				Expect(dnsVMVNFAliases).To(ContainElement(
					cce.DNSVMVNFAlias{

						ID:          dnsVMVNFAlias2ID,
						Name:        "dns vm vnf alias 123",
						Description: "description 1",
						VMVNFID:     vmVNFID,
					}))
			},
			Entry("GET /dns_vm_vnf_aliases"),
		)
	})

	Describe("GET /dns_vm_vnf_aliases/{id}", func() {
		var (
			dnsVMVNFAliasID string
		)

		BeforeEach(func() {
			dnsVMVNFAliasID = postDNSVMVNFAliases(vmVNFID)
		})

		DescribeTable("200 OK",
			func() {
				dnsVMVNFAlias := getDNSVMVNFAlias(
					dnsVMVNFAliasID)

				By("Verifying the created DNS vm vnf alias was returned")
				Expect(dnsVMVNFAlias).To(Equal(
					&cce.DNSVMVNFAlias{
						ID:          dnsVMVNFAliasID,
						Name:        "dns vm vnf alias 123",
						Description: "description 1",
						VMVNFID:     vmVNFID,
					},
				))
			},
			Entry("GET /dns_vm_vnf_aliases/{id}"),
		)

		DescribeTable("404 Not Found",
			func() {
				By("Sending a GET /dns_vm_vnf_aliases/{id} request")
				resp, err := http.Get(
					fmt.Sprintf(
						"http://127.0.0.1:8080/dns_vm_vnf_aliases/%s",
						uuid.New()))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("GET /dns_vm_vnf_aliases/{id} with nonexistent ID"),
		)
	})

	Describe("PATCH /dns_vm_vnf_aliases", func() {
		var (
			dnsVMVNFAliasID string
		)

		BeforeEach(func() {
			dnsVMVNFAliasID = postDNSVMVNFAliases(vmVNFID)
		})

		DescribeTable("204 No Content",
			func(reqStr string, expectedAlias *cce.DNSVMVNFAlias) {
				By("Sending a PATCH /dns_vm_vnf_aliases request")
				req, err := http.NewRequest(
					http.MethodPatch,
					"http://127.0.0.1:8080/dns_vm_vnf_aliases",
					strings.NewReader(fmt.Sprintf(reqStr,
						dnsVMVNFAliasID, vmVNFID)))
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 204 No Content response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				By("Getting the updated VNF alias")
				updatedAlias := getDNSVMVNFAlias(dnsVMVNFAliasID)

				By("Verifying the DNS vm vnf alias was updated")
				expectedAlias.SetID(dnsVMVNFAliasID)
				expectedAlias.VMVNFID = vmVNFID
				Expect(updatedAlias).To(Equal(expectedAlias))
			},
			Entry(
				"PATCH /dns_vm_vnf_aliases",
				`
                [
                    {
                        "id": "%s",
                        "name": "dns vm vnf alias 123456",
                        "description": "description 1",
                        "vm_vnf_id": "%s"
                    }
                ]`,
				&cce.DNSVMVNFAlias{
					Name:        "dns vm vnf alias 123456",
					Description: "description 1",
				}),
		)

		DescribeTable("400 Bad Request",
			func(reqStr string, expectedResp string) {
				By("Sending a PATCH /dns_vm_vnf_aliases request")
				if strings.Contains(reqStr, "%s") {
					reqStr = fmt.Sprintf(reqStr, dnsVMVNFAliasID)
				}
				req, err := http.NewRequest(
					http.MethodPatch,
					"http://127.0.0.1:8080/dns_vm_vnf_aliases",
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
				"PATCH /dns_vm_vnf_aliases without id",
				`
                [{}]`,
				"Validation failed: id not a valid uuid"),
			Entry(
				"PATCH /dns_vm_vnf_aliases without name",
				`
                [
                    {
                        "id": "%s",
                        "description": "description 1",
                        "vm_vnf_id": "123"
                    }
                ]`,
				"Validation failed: name cannot be empty"),
			Entry(
				"PATCH /dns_vm_vnf_aliases without description",
				`
                [
                    {
                        "id": "%s",
                        "name": "dns vm vnf alias 123",
                        "vm_vnf_id": "123"
                    }
                ]`,
				"Validation failed: description cannot be empty"),
			Entry(
				"PATCH /dns_vm_vnf_aliases without vm_vnf_id",
				`
                [
                    {
                        "id": "%s",
                        "name": "dns vm vnf alias 123",
                        "description": "description 1"
                    }
                ]`,
				"Validation failed: vm_vnf_id not a valid uuid"),
		)
	})

	Describe("DELETE /dns_vm_vnf_aliases/{id}", func() {
		var (
			dnsVMVNFAliasID string
		)

		BeforeEach(func() {
			dnsVMVNFAliasID = postDNSVMVNFAliases(vmVNFID)
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a DELETE /dns_vm_vnf_aliases/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf(
						"http://127.0.0.1:8080/dns_vm_vnf_aliases/%s",
						dnsVMVNFAliasID),
					nil)
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Verifying the DNS vm vnf alias was deleted")

				By("Sending a GET /dns_vm_vnf_aliases/{id} request")
				resp, err = http.Get(
					fmt.Sprintf(
						"http://127.0.0.1:8080/dns_vm_vnf_aliases/%s",
						dnsVMVNFAliasID))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("DELETE /dns_vm_vnf_aliases/{id}"),
		)

		DescribeTable("404 Not Found",
			func(id string) {
				By("Sending a DELETE /dns_vm_vnf_aliases/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf(
						"http://127.0.0.1:8080/dns_vm_vnf_aliases/%s",
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
				"DELETE /dns_vm_vnf_aliases/{id} with nonexistent ID",
				uuid.New()),
		)
	})
})
