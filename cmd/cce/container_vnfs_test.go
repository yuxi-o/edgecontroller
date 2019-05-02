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

var _ = Describe("/container_vnfs", func() {
	postSuccess := func() (id string) {
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

	get := func(id string) *cce.ContainerVNF {
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

	Describe("POST /container_vnfs", func() {
		DescribeTable("201 Created",
			func(req string) {
				By("Sending a POST /container_vnfs request")
				resp, err := http.Post(
					"http://127.0.0.1:8080/container_vnfs",
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
				"POST /container_vnfs",
				`
                {
                    "name": "container vnf",
                    "vendor": "smart edge",
                    "description": "my container vnf",
                    "image": "http://www.test.com/my_container_vnf.tar.gz",
                    "cores": 4,
                    "memory": 1024
                }`),
			Entry(
				"POST /container_vnfs without description",
				`
                {
                    "name": "container vnf",
                    "vendor": "smart edge",
                    "image": "http://www.test.com/my_container_vnf.tar.gz",
                    "cores": 4,
                    "memory": 1024
                }`),
		)

		DescribeTable("400 Bad Request",
			func(req, expectedResp string) {
				By("Sending a POST /container_vnfs request")
				resp, err := http.Post(
					"http://127.0.0.1:8080/container_vnfs",
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
				"POST /container_vnfs with id",
				`
                {
                    "id": "123"
                }`,
				"Validation failed: id cannot be specified in POST request"),
			Entry(
				"POST /container_vnfs without name",
				`
                {
                    "vendor": "smart edge",
                    "description": "my container vnf",
                    "image": "http://www.test.com/my_container_vnf.tar.gz",
                    "cores": 4,
                    "memory": 1024
                }`,
				"Validation failed: name cannot be empty"),
			Entry(
				"POST /container_vnfs without vendor",
				`
                {
                    "name": "container vnf",
                    "description": "my container vnf",
                    "image": "http://www.test.com/my_container_vnf.tar.gz",
                    "cores": 4,
                    "memory": 1024
                }`,
				"Validation failed: vendor cannot be empty"),
			Entry(
				"POST /container_vnfs without image",
				`
                {
                    "name": "container vnf",
                    "vendor": "smart edge",
                    "description": "my container vnf",
                    "cores": 4,
                    "memory": 1024
                }`,
				"Validation failed: image cannot be empty"),
			Entry("POST /container_vnfs with cores not in [1..8]",
				`
                {
                    "name": "container vnf",
                    "vendor": "smart edge",
                    "description": "my container vnf",
                    "image": "http://www.test.com/my_container_vnf.tar.gz",
                    "cores": 9,
                    "memory": 1024
                }`,
				"Validation failed: cores must be in [1..8]"),
			Entry("POST /container_vnfs with memory not in [1..16384]",
				`
                {
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

	Describe("GET /container_vnfs", func() {
		var (
			containerVNFID  string
			containerVNF2ID string
		)

		BeforeEach(func() {
			containerVNFID = postSuccess()
			containerVNF2ID = postSuccess()
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a GET /container_vnfs request")
				resp, err := http.Get("http://127.0.0.1:8080/container_vnfs")

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				var containerVNFs []cce.ContainerVNF

				By("Unmarshalling the response")
				Expect(json.Unmarshal(body, &containerVNFs)).To(Succeed())

				By("Verifying the 2 created container VNFs were returned")
				Expect(containerVNFs).To(ContainElement(
					cce.ContainerVNF{
						ID:          containerVNFID,
						Name:        "container vnf",
						Vendor:      "smart edge",
						Description: "my container vnf",
						Image:       "http://www.test.com/my_container_vnf.tar.gz", //nolint:lll
						Cores:       4,
						Memory:      1024,
					}))
				Expect(containerVNFs).To(ContainElement(
					cce.ContainerVNF{

						ID:          containerVNF2ID,
						Name:        "container vnf",
						Vendor:      "smart edge",
						Description: "my container vnf",
						Image:       "http://www.test.com/my_container_vnf.tar.gz", //nolint:lll
						Cores:       4,
						Memory:      1024,
					}))
			},
			Entry("GET /container_vnfs"),
		)
	})

	Describe("GET /container_vnfs/{id}", func() {
		var (
			containerVNFID string
		)

		BeforeEach(func() {
			containerVNFID = postSuccess()
		})

		DescribeTable("200 OK",
			func() {
				containerVNF := get(containerVNFID)

				By("Verifying the created container VNF was returned")
				Expect(containerVNF).To(Equal(
					&cce.ContainerVNF{
						ID:          containerVNFID,
						Name:        "container vnf",
						Vendor:      "smart edge",
						Description: "my container vnf",
						Image:       "http://www.test.com/my_container_vnf.tar.gz", //nolint:lll
						Cores:       4,
						Memory:      1024,
					},
				))
			},
			Entry("GET /container_vnfs/{id}"),
		)

		DescribeTable("404 Not Found",
			func() {
				By("Sending a GET /container_vnfs/{id} request")
				resp, err := http.Get(
					fmt.Sprintf("http://127.0.0.1:8080/container_vnfs/%s",
						uuid.New()))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("GET /container_vnfs/{id} with nonexistent ID"),
		)
	})

	Describe("PATCH /container_vnfs", func() {
		var (
			containerVNFID string
		)

		BeforeEach(func() {
			containerVNFID = postSuccess()
		})

		DescribeTable("204 No Content",
			func(reqStr string, expectedVNF *cce.ContainerVNF) {
				By("Sending a PATCH /container_vnfs request")
				req, err := http.NewRequest(
					http.MethodPatch,
					"http://127.0.0.1:8080/container_vnfs",
					strings.NewReader(fmt.Sprintf(reqStr, containerVNFID)))
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 204 No Content response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				By("Getting the updated VNF")
				updatedVNF := get(containerVNFID)

				By("Verifying the container VNF was updated")
				expectedVNF.SetID(containerVNFID)
				Expect(updatedVNF).To(Equal(expectedVNF))
			},
			Entry(
				"PATCH /container_vnfs",
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
				&cce.ContainerVNF{
					Name:        "container vnf2",
					Vendor:      "smart edge",
					Description: "my container vnf",
					Image:       "http://www.test.com/my_container_vnf.tar.gz",
					Cores:       4,
					Memory:      1024,
				}),
			Entry("PATCH /container_vnfs with no description",
				`
                [
                    {
                        "id": "%s",
                        "name": "container vnf2",
                        "vendor": "smart edge",
                        "image": "http://www.test.com/my_container_vnf.tar.gz",
                        "cores": 4,
                        "memory": 1024
                    }
                ]`,
				&cce.ContainerVNF{
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
				By("Sending a PATCH /container_vnfs request")
				if strings.Contains(reqStr, "%s") {
					reqStr = fmt.Sprintf(reqStr, containerVNFID)
				}
				req, err := http.NewRequest(
					http.MethodPatch,
					"http://127.0.0.1:8080/container_vnfs",
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
				"PATCH /container_vnfs without id",
				`
                [
                    {
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
				"PATCH /container_vnfs without name",
				`
                [
                    {
                        "id": "%s",
                        "vendor": "smart edge",
                        "description": "my container vnf",
                        "image": "http://www.test.com/my_container_vnf.tar.gz",
                        "cores": 4,
                        "memory": 1024
                    }
                ]`,
				"Validation failed: name cannot be empty"),
			Entry("PATCH /container_vnfs without vendor",
				`
                [
                    {
                        "id": "%s",
                        "name": "container vnf2",
                        "description": "my container vnf",
                        "image": "http://www.test.com/my_container_vnf.tar.gz",
                        "cores": 4,
                        "memory": 1024
                    }
                ]`,
				"Validation failed: vendor cannot be empty"),
			Entry("PATCH /container_vnfs without image",
				`
                [
                    {
                        "id": "%s",
                        "name": "container vnf2",
                        "vendor": "smart edge",
                        "description": "my container vnf",
                        "cores": 4,
                        "memory": 1024
                    }
                ]`,
				"Validation failed: image cannot be empty"),
			Entry("PATCH /container_vnfs with cores not in [1..8]",
				`
                [
                    {
                        "id": "%s",
                        "name": "container vnf2",
                        "vendor": "smart edge",
                        "description": "my container vnf",
                        "image": "http://www.test.com/my_container_vnf.tar.gz",
                        "cores": 9,
                        "memory": 1024
                    }
                ]`,
				"Validation failed: cores must be in [1..8]"),
			Entry("PATCH /container_vnfs with memory not in [1..16384]",
				`
                [
                    {
                        "id": "%s",
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

	Describe("DELETE /container_vnfs/{id}", func() {
		var (
			containerVNFID string
		)

		BeforeEach(func() {
			containerVNFID = postSuccess()
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a DELETE /container_vnfs/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("http://127.0.0.1:8080/container_vnfs/%s",
						containerVNFID),
					nil)
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Verifying the container VNF was deleted")

				By("Sending a GET /container_vnfs/{id} request")
				resp, err = http.Get(
					fmt.Sprintf("http://127.0.0.1:8080/container_vnfs/%s",
						containerVNFID))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("DELETE /container_vnfs/{id}"),
		)

		DescribeTable("404 Not Found",
			func(id string) {
				By("Sending a DELETE /container_vnfs/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("http://127.0.0.1:8080/container_vnfs/%s", id),
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
				"DELETE /container_vnfs/{id} with nonexistent ID",
				uuid.New()),
		)
	})
})
