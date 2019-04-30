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

var _ = Describe("/nodes", func() {
	postSuccess := func() (id string) {
		By("Sending a POST /nodes request")
		resp, err := http.Post(
			"http://127.0.0.1:8080/nodes",
			"application/json",
			strings.NewReader(`
                {
                    "name": "node123",
                    "location": "smart edge lab",
                    "serial": "abc123"
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

	get := func(id string) *cce.Node {
		By("Sending a GET /nodes/{id} request")
		resp, err := http.Get(
			fmt.Sprintf("http://127.0.0.1:8080/nodes/%s", id))

		By("Verifying a 200 OK response")
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		By("Reading the response body")
		body, err := ioutil.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())

		var node cce.Node

		By("Unmarshalling the response")
		Expect(json.Unmarshal(body, &node)).To(Succeed())

		return &node
	}

	Describe("POST /nodes", func() {
		DescribeTable("201 Created",
			func(req string) {
				By("Sending a POST /nodes request")
				resp, err := http.Post(
					"http://127.0.0.1:8080/nodes",
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
				"POST /nodes",
				`
                {
                    "name": "node123",
                    "location": "smart edge lab",
                    "serial": "abc123"
                }`),
		)

		DescribeTable("400 Bad Request",
			func(req, expectedResp string) {
				By("Sending a POST /nodes request")
				resp, err := http.Post(
					"http://127.0.0.1:8080/nodes",
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
				"POST /nodes without name",
				`
                {
                    "location": "smart edge lab",
                    "serial": "abc123"
                }`,
				"Validation failed: name cannot be empty"),
			Entry(
				"POST /nodes without location",
				`
                {
                    "name": "node123",
                    "serial": "abc123"
                }`,
				"Validation failed: location cannot be empty"),
			Entry(
				"POST /nodes without serial",
				`
                {
                    "name": "node123",
                    "location": "smart edge lab"
                }`,
				"Validation failed: serial cannot be empty"),
		)
	})

	Describe("GET /nodes", func() {
		var (
			nodeID  string
			node2ID string
		)

		BeforeEach(func() {
			nodeID = postSuccess()
			node2ID = postSuccess()
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a GET /nodes request")
				resp, err := http.Get("http://127.0.0.1:8080/nodes")

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				var nodes []cce.Node

				By("Unmarshalling the response")
				Expect(json.Unmarshal(body, &nodes)).To(Succeed())

				By("Verifying the 2 created nodes were returned")
				Expect(nodes).To(ContainElement(
					cce.Node{
						ID:       nodeID,
						Name:     "node123",
						Location: "smart edge lab",
						Serial:   "abc123",
					}))
				Expect(nodes).To(ContainElement(
					cce.Node{

						ID:       node2ID,
						Name:     "node123",
						Location: "smart edge lab",
						Serial:   "abc123",
					}))
			},
			Entry("GET /nodes"),
		)
	})

	Describe("GET /nodes/{id}", func() {
		var (
			nodeID string
		)

		BeforeEach(func() {
			nodeID = postSuccess()
			fmt.Println("before each: " + nodeID)
		})

		DescribeTable("200 OK",
			func() {
				fmt.Println(nodeID)
				node := get(nodeID)

				By("Verifying the created node was returned")
				Expect(node).To(Equal(
					&cce.Node{
						ID:       nodeID,
						Name:     "node123",
						Location: "smart edge lab",
						Serial:   "abc123",
					},
				))
			},
			Entry("GET /nodes/{id}"),
		)

		DescribeTable("404 Not Found",
			func() {
				By("Sending a GET /nodes/{id} request")
				resp, err := http.Get(
					fmt.Sprintf("http://127.0.0.1:8080/nodes/%s",
						uuid.New()))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("GET /nodes/{id} with nonexistent ID"),
		)
	})

	Describe("PATCH /nodes", func() {
		var (
			nodeID string
		)

		BeforeEach(func() {
			nodeID = postSuccess()
		})

		DescribeTable("204 No Content",
			func(reqStr string, expectedApp *cce.Node) {
				By("Sending a PATCH /nodes request")
				req, err := http.NewRequest(
					http.MethodPatch,
					"http://127.0.0.1:8080/nodes",
					strings.NewReader(fmt.Sprintf(reqStr, nodeID)))
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 204 No Content response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				By("Getting the updated application")
				updatedApp := get(nodeID)

				By("Verifying the node was updated")
				expectedApp.SetID(nodeID)
				Expect(updatedApp).To(Equal(expectedApp))
			},
			Entry(
				"PATCH /nodes/{id}",
				`
                [
                    {
                        "id": "%s",
                        "name": "node123456",
                        "location": "smart edge lab",
                        "serial": "abc123"
                        }
                ]`,
				&cce.Node{
					Name:     "node123456",
					Location: "smart edge lab",
					Serial:   "abc123",
				}),
		)

		DescribeTable("400 Bad Request",
			func(reqStr string, expectedResp string) {
				By("Sending a PATCH /nodes request")
				req, err := http.NewRequest(
					http.MethodPatch,
					"http://127.0.0.1:8080/nodes",
					strings.NewReader(fmt.Sprintf(reqStr, nodeID)))
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
				"PATCH /nodes without name",
				`
                [
                    {
                        "id": "%s",
                        "location": "smart edge lab",
                        "serial": "abc123"
                        }
                ]`,
				"Validation failed: name cannot be empty"),
			Entry("PATCH /nodes without location",
				`
                [
                    {
                        "id": "%s",
                        "name": "node123",
                        "serial": "abc123"
                        }
                ]`,
				"Validation failed: location cannot be empty"),
			Entry("PATCH /nodes without serial",
				`
                [
                    {
                        "id": "%s",
                        "name": "node123",
                        "location": "smart edge lab"
                        }
                ]`,
				"Validation failed: serial cannot be empty"),
		)
	})

	Describe("DELETE /nodes/{id}", func() {
		var (
			nodeID string
		)

		BeforeEach(func() {
			nodeID = postSuccess()
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a DELETE /nodes/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("http://127.0.0.1:8080/nodes/%s",
						nodeID),
					nil)
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Verifying the node was deleted")

				By("Sending a GET /nodes/{id} request")
				resp, err = http.Get(
					fmt.Sprintf("http://127.0.0.1:8080/nodes/%s",
						nodeID))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("DELETE /nodes/{id}"),
		)

		DescribeTable("404 Not Found",
			func(id string) {
				By("Sending a DELETE /nodes/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("http://127.0.0.1:8080/nodes/%s", id),
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
				"DELETE /nodes/{id} with nonexistent ID",
				uuid.New()),
		)
	})
})
