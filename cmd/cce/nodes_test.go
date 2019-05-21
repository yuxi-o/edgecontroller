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
	Describe("POST /nodes", func() {
		DescribeTable("201 Created",
			func(req string) {
				By("Sending a POST /nodes request")
				resp, err := apiCli.Post(
					"http://127.0.0.1:8080/nodes",
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
				resp, err := apiCli.Post(
					"http://127.0.0.1:8080/nodes",
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
				"POST /nodes with id",
				`
				{
					"id": "123"
				}`,
				"Validation failed: id cannot be specified in POST request"),
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
			nodeID = postNodes()
			node2ID = postNodes()
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a GET /nodes request")
				resp, err := apiCli.Get("http://127.0.0.1:8080/nodes")
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 200 OK response")
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				var nodes []*cce.Node

				By("Unmarshalling the response")
				Expect(json.Unmarshal(body, &nodes)).To(Succeed())

				By("Verifying the 2 created nodes were returned")
				Expect(nodes).To(ContainElement(
					&cce.Node{
						ID:         nodeID,
						Name:       "Test Node 1",
						Location:   "Localhost port 8082",
						Serial:     "ABC-123",
						GRPCTarget: "127.0.0.1:8082",
					}))
				Expect(nodes).To(ContainElement(
					&cce.Node{
						ID:         node2ID,
						Name:       "Test Node 1",
						Location:   "Localhost port 8082",
						Serial:     "ABC-123",
						GRPCTarget: "127.0.0.1:8082",
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
			nodeID = postNodes()
		})

		DescribeTable("200 OK",
			func() {
				node := getNode(nodeID)

				By("Verifying the created node was returned")
				Expect(node).To(Equal(
					&cce.Node{
						ID:         nodeID,
						Name:       "Test Node 1",
						Location:   "Localhost port 8082",
						Serial:     "ABC-123",
						GRPCTarget: "127.0.0.1:8082",
					},
				))
			},
			Entry("GET /nodes/{id}"),
		)

		DescribeTable("404 Not Found",
			func() {
				By("Sending a GET /nodes/{id} request")
				resp, err := apiCli.Get(
					fmt.Sprintf("http://127.0.0.1:8080/nodes/%s",
						uuid.New()))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 404 Not Found response")
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
			nodeID = postNodes()
		})

		DescribeTable("204 No Content",
			func(reqStr string, expectedNode *cce.Node) {
				By("Sending a PATCH /nodes request")
				resp, err := apiCli.Patch(
					"http://127.0.0.1:8080/nodes",
					"application/json",
					strings.NewReader(fmt.Sprintf(reqStr, nodeID)))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 204 No Content response")
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				By("Getting the updated node")
				updatedNode := getNode(nodeID)

				By("Verifying the node was updated")
				expectedNode.SetID(nodeID)
				Expect(updatedNode).To(Equal(expectedNode))
			},
			Entry(
				"PATCH /nodes",
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
				if strings.Contains(reqStr, "%s") {
					reqStr = fmt.Sprintf(reqStr, nodeID)
				}
				resp, err := apiCli.Patch(
					"http://127.0.0.1:8080/nodes",
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
				"PATCH /nodes without id",
				`
				[
					{
						"name": "node123",
						"location": "smart edge lab",
						"serial": "abc123"
					}
				]`,
				"Validation failed: id not a valid uuid"),
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
			nodeID = postNodes()
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a DELETE /nodes/{id} request")
				resp, err := apiCli.Delete(
					fmt.Sprintf("http://127.0.0.1:8080/nodes/%s",
						nodeID))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 200 OK response")
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Verifying the node was deleted")

				By("Sending a GET /nodes/{id} request")
				resp, err = apiCli.Get(
					fmt.Sprintf("http://127.0.0.1:8080/nodes/%s",
						nodeID))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 404 Not Found response")
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("DELETE /nodes/{id}"),
		)

		DescribeTable("404 Not Found",
			func(id string) {
				By("Sending a DELETE /nodes/{id} request")
				resp, err := apiCli.Delete(
					fmt.Sprintf("http://127.0.0.1:8080/nodes/%s", id))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 404 Not Found response")
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry(
				"DELETE /nodes/{id} with nonexistent ID",
				uuid.New()),
		)
	})
})
