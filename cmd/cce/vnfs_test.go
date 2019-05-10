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
				resp, err := http.Post(
					"http://127.0.0.1:8080/vnfs",
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
				"POST /vnfs",
				`
				{
					"type": "container",
					"name": "container vnf",
					"vendor": "smart edge",
					"description": "my container vnf",
					"image": "http://www.test.com/my_container_vnf.tar.gz",
					"cores": 4,
					"memory": 1024
				}`),
			Entry(
				"POST /vnfs without description",
				`
				{
					"type": "container",
					"name": "container vnf",
					"vendor": "smart edge",
					"image": "http://www.test.com/my_container_vnf.tar.gz",
					"cores": 4,
					"memory": 1024
				}`),
		)

		DescribeTable("400 Bad Request",
			func(req, expectedResp string) {
				By("Sending a POST /vnfs request")
				resp, err := http.Post(
					"http://127.0.0.1:8080/vnfs",
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
					"vendor": "smart edge",
					"description": "my container vnf",
					"image": "http://www.test.com/my_container_vnf.tar.gz",
					"cores": 4,
					"memory": 1024
				}`,
				`Validation failed: type must be either "container" or "vm"`),
			Entry(
				"POST /vnfs without name",
				`
				{
					"type": "container",
					"vendor": "smart edge",
					"description": "my container vnf",
					"image": "http://www.test.com/my_container_vnf.tar.gz",
					"cores": 4,
					"memory": 1024
				}`,
				"Validation failed: name cannot be empty"),
			Entry(
				"POST /vnfs without vendor",
				`
				{
					"type": "container",
					"name": "container vnf",
					"description": "my container vnf",
					"image": "http://www.test.com/my_container_vnf.tar.gz",
					"cores": 4,
					"memory": 1024
				}`,
				"Validation failed: vendor cannot be empty"),
			Entry(
				"POST /vnfs without image",
				`
				{
					"type": "container",
					"name": "container vnf",
					"vendor": "smart edge",
					"description": "my container vnf",
					"cores": 4,
					"memory": 1024
				}`,
				"Validation failed: image cannot be empty"),
			Entry("POST /vnfs with cores not in [1..8]",
				`
				{
					"type": "container",
					"name": "container vnf",
					"vendor": "smart edge",
					"description": "my container vnf",
					"image": "http://www.test.com/my_container_vnf.tar.gz",
					"cores": 9,
					"memory": 1024
				}`,
				"Validation failed: cores must be in [1..8]"),
			Entry("POST /vnfs with memory not in [1..16384]",
				`
				{
					"type": "container",
					"name": "container vnf",
					"vendor": "smart edge",
					"description": "my container vnf",
					"image": "http://www.test.com/my_container_vnf.tar.gz",
					"cores": 8,
					"memory": 16385
				}`,
				"Validation failed: memory must be in [1..16384]"),
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
				resp, err := http.Get("http://127.0.0.1:8080/vnfs")

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
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
						Image:       "http://www.test.com/my_container_vnf.tar.gz", //nolint:lll
						Cores:       4,
						Memory:      1024,
					}))
				Expect(vnfs).To(ContainElement(
					cce.VNF{
						ID:          vmVNFID,
						Type:        "vm",
						Name:        "vm vnf",
						Vendor:      "smart edge",
						Description: "my vm vnf",
						Image:       "http://www.test.com/my_vm_vnf.tar.gz",
						Cores:       4,
						Memory:      1024,
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
						Image:       "http://www.test.com/my_container_vnf.tar.gz", //nolint:lll
						Cores:       4,
						Memory:      1024,
					},
				))
			},
			Entry("GET /vnfs/{id}"),
		)

		DescribeTable("404 Not Found",
			func() {
				By("Sending a GET /vnfs/{id} request")
				resp, err := http.Get(
					fmt.Sprintf("http://127.0.0.1:8080/vnfs/%s",
						uuid.New()))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
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
				req, err := http.NewRequest(
					http.MethodPatch,
					"http://127.0.0.1:8080/vnfs",
					strings.NewReader(fmt.Sprintf(reqStr, containerVNFID)))
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 204 No Content response")
				Expect(err).ToNot(HaveOccurred())
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
						"vendor": "smart edge",
						"description": "my container vnf",
						"image": "http://www.test.com/my_container_vnf.tar.gz",
						"cores": 4,
						"memory": 1024
					}
				]`,
				&cce.VNF{
					Type:        "container",
					Name:        "container vnf2",
					Vendor:      "smart edge",
					Description: "my container vnf",
					Image:       "http://www.test.com/my_container_vnf.tar.gz",
					Cores:       4,
					Memory:      1024,
				}),
			Entry("PATCH /vnfs with no description",
				`
				[
					{
						"id": "%s",
						"type": "container",
						"name": "container vnf2",
						"vendor": "smart edge",
						"image": "http://www.test.com/my_container_vnf.tar.gz",
						"cores": 4,
						"memory": 1024
					}
				]`,
				&cce.VNF{
					Type:        "container",
					Name:        "container vnf2",
					Vendor:      "smart edge",
					Description: "",
					Image:       "http://www.test.com/my_container_vnf.tar.gz",
					Cores:       4,
					Memory:      1024,
				}),
		)

		DescribeTable("400 Bad Request",
			func(reqStr string, expectedResp string) {
				By("Sending a PATCH /vnfs request")
				if strings.Contains(reqStr, "%s") {
					reqStr = fmt.Sprintf(reqStr, containerVNFID)
				}
				req, err := http.NewRequest(
					http.MethodPatch,
					"http://127.0.0.1:8080/vnfs",
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
				"PATCH /vnfs without id",
				`
				[
					{
						"type": "container",
						"name": "container vnf2",
						"vendor": "smart edge",
						"description": "my container vnf",
						"image": "http://www.test.com/my_container_vnf.tar.gz",
						"cores": 4,
						"memory": 1024
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
						"vendor": "smart edge",
						"description": "my container vnf",
						"image": "http://www.test.com/my_container_vnf.tar.gz",
						"cores": 4,
						"memory": 1024
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
						"vendor": "smart edge",
						"description": "my container vnf",
						"image": "http://www.test.com/my_container_vnf.tar.gz",
						"cores": 4,
						"memory": 1024
					}
				]`,
				"Validation failed: name cannot be empty"),
			Entry("PATCH /vnfs without vendor",
				`
				[
					{
						"id": "%s",
						"type": "container",
						"name": "container vnf2",
						"description": "my container vnf",
						"image": "http://www.test.com/my_container_vnf.tar.gz",
						"cores": 4,
						"memory": 1024
					}
				]`,
				"Validation failed: vendor cannot be empty"),
			Entry("PATCH /vnfs without image",
				`
				[
					{
						"id": "%s",
						"type": "container",
						"name": "container vnf2",
						"vendor": "smart edge",
						"description": "my container vnf",
						"cores": 4,
						"memory": 1024
					}
				]`,
				"Validation failed: image cannot be empty"),
			Entry("PATCH /vnfs with cores not in [1..8]",
				`
				[
					{
						"id": "%s",
						"type": "container",
						"name": "container vnf2",
						"vendor": "smart edge",
						"description": "my container vnf",
						"image": "http://www.test.com/my_container_vnf.tar.gz",
						"cores": 9,
						"memory": 1024
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
						"vendor": "smart edge",
						"description": "my container vnf",
						"image": "http://www.test.com/my_container_vnf.tar.gz",
						"cores": 4,
						"memory": 16385
					}
				]`,
				"Validation failed: memory must be in [1..16384]"),
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
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("http://127.0.0.1:8080/vnfs/%s",
						containerVNFID),
					nil)
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Verifying the VNF was deleted")

				By("Sending a GET /vnfs/{id} request")
				resp, err = http.Get(
					fmt.Sprintf("http://127.0.0.1:8080/vnfs/%s",
						containerVNFID))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("DELETE /vnfs/{id}"),
		)

		DescribeTable("404 Not Found",
			func(id string) {
				By("Sending a DELETE /vnfs/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("http://127.0.0.1:8080/vnfs/%s", id),
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
				"DELETE /vnfs/{id} with nonexistent ID",
				uuid.New()),
		)

		DescribeTable("422 Unprocessable Entity",
			func() {
				postDNSVNFAliases(containerVNFID)

				By("Sending a DELETE /vnfs/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf(
						"http://127.0.0.1:8080/vnfs/%s",
						containerVNFID),
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
					"cannot delete vnf_id %s: record in use in "+
						"dns_vnf_aliases",
					containerVNFID)))
			},
			Entry("DELETE /vnfs/{id} with dns_vnf_aliases record"),
		)
	})
})
