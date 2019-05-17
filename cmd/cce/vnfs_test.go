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

var _ = Describe("/vnfs", func() {
	Describe("POST /vnfs", func() {
		DescribeTable("201 Created",
			func(req string) {
				By("Sending a POST /vnfs request")
				resp, err := apiCli.Post(
					"http://127.0.0.1:8080/vnfs",
					"application/json",
					strings.NewReader(req))
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
				"POST /vnfs",
				`
				{
					"type": "container",
					"name": "container vnf",
					"version": "latest",
					"vendor": "smart edge",
					"description": "my container vnf",	
					"cores": 4,
					"memory": 1024,
					"source": "http://www.test.com/my_container_vnf.tar.gz"
				}`),
			Entry(
				"POST /vnfs without description",
				`
				{
					"type": "container",
					"name": "container vnf",
					"version": "latest",
					"vendor": "smart edge",
					"cores": 4,
					"memory": 1024,
					"source": "http://www.test.com/my_container_vnf.tar.gz"
				}`),
		)

		DescribeTable("400 Bad Request",
			func(req, expectedResp string) {
				By("Sending a POST /vnfs request")
				resp, err := apiCli.Post(
					"http://127.0.0.1:8080/vnfs",
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
				"POST /vnfs with id",
				`
				{
					"id": "123"
				}`,
				"Validation failed: id cannot be specified in POST request"),
			Entry(
				"POST /vnfs without type",
				`
				{
					"name": "container vnf",
					"version": "latest",
					"vendor": "smart edge",
					"description": "my container vnf",
					"cores": 4,
					"memory": 1024,
					"source": "http://www.test.com/my_container_vnf.tar.gz"
				}`,
				`Validation failed: type must be either "container" or "vm"`),
			Entry(
				"POST /vnfs without name",
				`
				{
					"type": "container",
					"version": "latest",
					"vendor": "smart edge",
					"description": "my container vnf",
					"cores": 4,
					"memory": 1024,
					"source": "http://www.test.com/my_container_vnf.tar.gz"
				}`,
				"Validation failed: name cannot be empty"),
			Entry(
				"POST /vnfs without version",
				`
						{
							"type": "container",
							"name": "container vnf",
							"vendor": "smart edge",
							"description": "my container vnf",
							"cores": 4,
							"memory": 1024,
							"source": "http://www.test.com/my_container_vnf.tar.gz"
						}`,
				"Validation failed: version cannot be empty"),
			Entry(
				"POST /vnfs without vendor",
				`
				{
					"type": "container",
					"name": "container vnf",
					"version": "latest",
					"description": "my container vnf",
					"cores": 4,
					"memory": 1024,
					"source": "http://www.test.com/my_container_vnf.tar.gz"
				}`,
				"Validation failed: vendor cannot be empty"),
			Entry("POST /vnfs with cores not in [1..8]",
				`
				{
					"type": "container",
					"name": "container vnf",
					"version": "latest",
					"vendor": "smart edge",
					"description": "my container vnf",
					"cores": 9,
					"memory": 1024,
					"source": "http://www.test.com/my_container_vnf.tar.gz"
				}`,
				"Validation failed: cores must be in [1..8]"),
			Entry("POST /vnfs with memory not in [1..16384]",
				`
				{
					"type": "container",
					"name": "container vnf",
					"version": "latest",
					"vendor": "smart edge",
					"description": "my container vnf",
					"cores": 8,
					"memory": 16385,
					"source": "http://www.test.com/my_container_vnf.tar.gz"
				}`,
				"Validation failed: memory must be in [1..16384]"),
			Entry(
				"POST /vnfs without source",
				`
					{
						"type": "container",
						"name": "container vnf",
						"version": "latest",
						"vendor": "smart edge",
						"description": "my container vnf",
						"cores": 4,
						"memory": 1024
					}`,
				"Validation failed: source cannot be empty"),
			Entry(
				"POST /vnfs without source",
				`
						{
							"type": "container",
							"name": "container vnf",
							"version": "latest",
							"vendor": "smart edge",
							"description": "my container vnf",
							"cores": 4,
							"memory": 1024,
							"source": "invalid.url"
						}`,
				"Validation failed: source cannot be parsed as a URI"),
		)
	})

	Describe("GET /vnfs", func() {
		var (
			containerVNFID string
			vmVNFID        string
		)

		BeforeEach(func() {
			containerVNFID = postVNFs("container")
			vmVNFID = postVNFs("vm")
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a GET /vnfs request")
				resp, err := apiCli.Get("http://127.0.0.1:8080/vnfs")
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 200 OK response")
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				var vnfs []cce.VNF

				By("Unmarshalling the response")
				Expect(json.Unmarshal(body, &vnfs)).To(Succeed())

				By("Verifying the 2 created VNFs were returned")
				Expect(vnfs).To(ContainElement(
					cce.VNF{
						ID:          containerVNFID,
						Type:        "container",
						Name:        "container vnf",
						Vendor:      "smart edge",
						Description: "my container vnf",
						Version:     "latest",
						Cores:       4,
						Memory:      1024,
						Source:      "http://www.test.com/my_container_vnf.tar.gz",
					}))
				Expect(vnfs).To(ContainElement(
					cce.VNF{
						ID:          vmVNFID,
						Type:        "vm",
						Name:        "vm vnf",
						Vendor:      "smart edge",
						Description: "my vm vnf",
						Version:     "latest",
						Cores:       4,
						Memory:      1024,
						Source:      "http://www.test.com/my_vm_vnf.tar.gz",
					}))
			},
			Entry("GET /vnfs"),
		)
	})

	Describe("GET /vnfs/{id}", func() {
		var (
			containerVNFID string
		)

		BeforeEach(func() {
			containerVNFID = postVNFs("container")
		})

		DescribeTable("200 OK",
			func() {
				vnf := getVNF(containerVNFID)

				By("Verifying the created VNF was returned")
				Expect(vnf).To(Equal(
					&cce.VNF{
						ID:          containerVNFID,
						Type:        "container",
						Name:        "container vnf",
						Vendor:      "smart edge",
						Description: "my container vnf",
						Version:     "latest",
						Cores:       4,
						Memory:      1024,
						Source:      "http://www.test.com/my_container_vnf.tar.gz",
					},
				))
			},
			Entry("GET /vnfs/{id}"),
		)

		DescribeTable("404 Not Found",
			func() {
				By("Sending a GET /vnfs/{id} request")
				resp, err := apiCli.Get(
					fmt.Sprintf("http://127.0.0.1:8080/vnfs/%s",
						uuid.New()))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 404 Not Found response")
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("GET /vnfs/{id} with nonexistent ID"),
		)
	})

	Describe("PATCH /vnfs", func() {
		var (
			containerVNFID string
		)

		BeforeEach(func() {
			containerVNFID = postVNFs("container")
		})

		DescribeTable("204 No Content",
			func(reqStr string, expectedVNF *cce.VNF) {
				By("Sending a PATCH /vnfs request")
				resp, err := apiCli.Patch(
					"http://127.0.0.1:8080/vnfs",
					"application/json",
					strings.NewReader(fmt.Sprintf(reqStr, containerVNFID)))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 204 No Content response")
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				By("Getting the updated VNF")
				updatedVNF := getVNF(containerVNFID)

				By("Verifying the VNF was updated")
				expectedVNF.SetID(containerVNFID)
				Expect(updatedVNF).To(Equal(expectedVNF))
			},
			Entry(
				"PATCH /vnfs",
				`
				[
					{
						"id": "%s",
						"type": "container",
						"name": "container vnf2",
						"version": "latest",
						"vendor": "smart edge",
						"description": "my container vnf",
						"cores": 4,
						"memory": 1024,
						"source": "http://www.test.com/my_container_vnf.tar.gz"
					}
				]`,
				&cce.VNF{
					Type:        "container",
					Name:        "container vnf2",
					Version:     "latest",
					Vendor:      "smart edge",
					Description: "my container vnf",
					Cores:       4,
					Memory:      1024,
					Source:      "http://www.test.com/my_container_vnf.tar.gz",
				}),
			Entry("PATCH /vnfs with no description",
				`
				[
					{
						"id": "%s",
						"type": "container",
						"name": "container vnf2",
						"version": "latest",
						"vendor": "smart edge",
						"cores": 4,
						"memory": 1024,
						"source": "http://www.test.com/my_container_vnf.tar.gz"
					}
				]`,
				&cce.VNF{
					Type:        "container",
					Name:        "container vnf2",
					Version:     "latest",
					Vendor:      "smart edge",
					Description: "",
					Cores:       4,
					Memory:      1024,
					Source:      "http://www.test.com/my_container_vnf.tar.gz",
				}),
		)

		DescribeTable("400 Bad Request",
			func(reqStr string, expectedResp string) {
				By("Sending a PATCH /vnfs request")
				if strings.Contains(reqStr, "%s") {
					reqStr = fmt.Sprintf(reqStr, containerVNFID)
				}
				resp, err := apiCli.Patch(
					"http://127.0.0.1:8080/vnfs",
					"application/json",
					strings.NewReader(reqStr))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 400 Bad Request")
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the response body")
				Expect(string(body)).To(Equal(expectedResp))
			},
			Entry(
				"PATCH /vnfs without id",
				`
				[
					{
						"type": "container",
						"name": "container vnf2",
						"version": "latest",
						"vendor": "smart edge",
						"description": "my container vnf",
						"cores": 4,
						"memory": 1024,
						"source": "http://www.test.com/my_container_vnf.tar.gz"
					}
				]`,
				"Validation failed: id not a valid uuid"),
			Entry(
				"PATCH /vnfs without type",
				`
				[
					{
						"id": "%s",
						"name": "container vnf2",
						"version": "latest",
						"vendor": "smart edge",
						"description": "my container vnf",
						"cores": 4,
						"memory": 1024,
						"source": "http://www.test.com/my_container_vnf.tar.gz"
					}
				]`,
				`Validation failed: type must be either "container" or "vm"`),
			Entry(
				"PATCH /vnfs without name",
				`
				[
					{
						"id": "%s",
						"type": "container",
						"version": "latest",
						"vendor": "smart edge",
						"description": "my container vnf",
						"cores": 4,
						"memory": 1024,
						"source": "http://www.test.com/my_container_vnf.tar.gz"
					}
				]`,
				"Validation failed: name cannot be empty"),
			Entry("PATCH /vnfs without version",
				`
				[
					{
						"id": "%s",
						"type": "container",
						"name": "container vnf2",
						"vendor": "smart edge",
						"description": "my container vnf",
						"cores": 4,
						"memory": 1024,
						"source": "http://www.test.com/my_container_vnf.tar.gz"
					}
				]`,
				"Validation failed: version cannot be empty"),
			Entry("PATCH /vnfs without vendor",
				`
				[
					{
						"id": "%s",
						"type": "container",
						"name": "container vnf2",
						"version": "latest",
						"description": "my container vnf",
						"cores": 4,
						"memory": 1024,
						"source": "http://www.test.com/my_container_vnf.tar.gz"
					}
				]`,
				"Validation failed: vendor cannot be empty"),
			Entry("PATCH /vnfs with cores not in [1..8]",
				`
				[
					{
						"id": "%s",
						"type": "container",
						"name": "container vnf2",
						"version": "latest",
						"vendor": "smart edge",
						"description": "my container vnf",
						"cores": 9,
						"memory": 1024,
						"source": "http://www.test.com/my_container_vnf.tar.gz"
					}
				]`,
				"Validation failed: cores must be in [1..8]"),
			Entry("PATCH /vnfs with memory not in [1..16384]",
				`
				[
					{
						"id": "%s",
						"type": "container",
						"name": "container vnf2",
						"version": "latest",
						"vendor": "smart edge",
						"description": "my container vnf",
						"cores": 4,
						"memory": 16385,
						"source": "http://www.test.com/my_container_vnf.tar.gz"
					}
				]`,
				"Validation failed: memory must be in [1..16384]"),
			Entry("PATCH /vnfs without source",
				`
				[
					{
						"id": "%s",
						"type": "container",
						"name": "container vnf2",
						"version": "latest",
						"vendor": "smart edge",
						"description": "my container vnf",
						"cores": 4,
						"memory": 1024
					}
				]`,
				"Validation failed: source cannot be empty"),
			Entry("PATCH /vnfs without source",
				`
				[
					{
						"id": "%s",
						"type": "container",
						"name": "container vnf2",
						"version": "latest",
						"vendor": "smart edge",
						"description": "my container vnf",
						"cores": 4,
						"memory": 1024,
						"source" : "invalid.url"
					}
				]`,
				"Validation failed: source cannot be parsed as a URI"),
		)
	})

	Describe("DELETE /vnfs/{id}", func() {
		var (
			containerVNFID string
		)

		BeforeEach(func() {
			containerVNFID = postVNFs("container")
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a DELETE /vnfs/{id} request")
				resp, err := apiCli.Delete(
					fmt.Sprintf("http://127.0.0.1:8080/vnfs/%s",
						containerVNFID))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 200 OK response")
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Verifying the VNF was deleted")

				By("Sending a GET /vnfs/{id} request")
				resp, err = apiCli.Get(
					fmt.Sprintf("http://127.0.0.1:8080/vnfs/%s",
						containerVNFID))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 404 Not Found response")
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("DELETE /vnfs/{id}"),
		)

		DescribeTable("404 Not Found",
			func(id string) {
				By("Sending a DELETE /vnfs/{id} request")
				resp, err := apiCli.Delete(
					fmt.Sprintf("http://127.0.0.1:8080/vnfs/%s", id))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 404 Not Found response")
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry(
				"DELETE /vnfs/{id} with nonexistent ID",
				uuid.New()),
		)

		DescribeTable("422 Unprocessable Entity",
			func(resource, expectedResp string) {
				switch resource {
				case "dns_configs_vnf_aliases":
					postDNSConfigsVNFAliases(
						postDNSConfigs(),
						containerVNFID)
				case "nodes_vnfs":
					postNodesVNFs(
						postNodes(),
						containerVNFID)
				}

				By("Sending a DELETE /vnfs/{id} request")
				resp, err := apiCli.Delete(
					fmt.Sprintf(
						"http://127.0.0.1:8080/vnfs/%s",
						containerVNFID))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 422 response")
				Expect(resp.StatusCode).To(Equal(
					http.StatusUnprocessableEntity))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the response body")
				Expect(string(body)).To(Equal(
					fmt.Sprintf(expectedResp, containerVNFID)))
			},
			Entry(
				"DELETE /vnfs/{id} with dns_configs_vnf_aliases record",
				"dns_configs_vnf_aliases",
				"cannot delete vnf_id %s: record in use in dns_configs_vnf_aliases",
			),
			Entry(
				"DELETE /vnfs/{id} with nodes_vnfs record",
				"nodes_vnfs",
				"cannot delete vnf_id %s: record in use in nodes_vnfs",
			),
		)
	})
})
