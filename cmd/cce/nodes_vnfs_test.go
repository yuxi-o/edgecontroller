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

var _ = Describe("/nodes_vnfs", func() {
	var (
		nodeID string
		vnfID  string
	)

	BeforeEach(func() {
		nodeID = postNodes()
		vnfID = postVNFs("container")
	})

	Describe("POST /nodes_vnfs", func() {
		DescribeTable("201 Created",
			func() {
				By("Sending a POST /nodes_vnfs request")
				resp, err := apiCli.Post(
					"http://127.0.0.1:8080/nodes_vnfs",
					"application/json",
					strings.NewReader(fmt.Sprintf(
						`
						{
							"node_id": "%s",
							"vnf_id": "%s"
						}`, nodeID, vnfID)))
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
				"POST /nodes_vnfs"),
		)

		DescribeTable("400 Bad Request",
			func(req, expectedResp string) {
				By("Sending a POST /nodes_vnfs request")
				resp, err := apiCli.Post(
					"http://127.0.0.1:8080/nodes_vnfs",
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
				"POST /nodes_vnfs with id",
				`
				{
					"id": "123"
				}`,
				"Validation failed: id cannot be specified in POST request"),
			Entry(
				"POST /nodes_vnfs without node_id",
				`
				{
				}`,
				"Validation failed: node_id not a valid uuid"),
			Entry(
				"POST /nodes_vnfs without vnf_id",
				fmt.Sprintf(`
				{
					"node_id": "%s"
				}`, uuid.New()),
				"Validation failed: vnf_id not a valid uuid"),
		)

		DescribeTable("422 Unprocessable Entity",
			func() {
				var (
					resp *http.Response
					err  error
				)

				By("Sending a POST /nodes_vnfs request")
				postNodesVNFs(nodeID, vnfID)

				By("Repeating the first POST /nodes_vnfs request")
				resp, err = apiCli.Post(
					"http://127.0.0.1:8080/nodes_vnfs",
					"application/json",
					strings.NewReader(fmt.Sprintf(
						`
						{
							"node_id": "%s",
							"vnf_id": "%s"
						}`, nodeID, vnfID)))
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
					"duplicate record detected for node_id %s and "+
						"vnf_id %s",
					nodeID,
					vnfID)))
			},
			Entry("POST /nodes_vnfs with duplicate node_id and vnf_id"),
		)
	})

	Describe("GET /nodes_vnfs", func() {
		var (
			nodeVNFID  string
			vnf2ID     string
			nodeVNF2ID string
		)

		BeforeEach(func() {
			nodeVNFID = postNodesVNFs(nodeID, vnfID)
			vnf2ID = postVNFs("container")
			nodeVNF2ID = postNodesVNFs(nodeID, vnf2ID)
		})

		DescribeTable("200 OK",
			func() {
				nodeVNFsResp := getNodeVNFs(nodeID)

				By("Verifying the 2 created node VNFs were returned")
				Expect(nodeVNFsResp).To(ContainElement(
					&cce.NodeVNFResp{
						NodeVNF: cce.NodeVNF{
							ID:     nodeVNFID,
							NodeID: nodeID,
							VNFID:  vnfID,
						},
						Status: "deployed",
					}))
				Expect(nodeVNFsResp).To(ContainElement(
					&cce.NodeVNFResp{
						NodeVNF: cce.NodeVNF{
							ID:     nodeVNF2ID,
							NodeID: nodeID,
							VNFID:  vnf2ID,
						},
						Status: "deployed",
					}))
			},
			Entry("GET /nodes_vnfs"),
		)
	})

	Describe("GET /nodes_vnfs/{id}", func() {
		var (
			nodeVNFID string
		)

		BeforeEach(func() {
			nodeVNFID = postNodesVNFs(nodeID, vnfID)
		})

		DescribeTable("200 OK",
			func() {
				nodeVNFResp := getNodeVNF(nodeVNFID)

				By("Verifying the created node VNF was returned")
				Expect(nodeVNFResp).To(Equal(
					&cce.NodeVNFResp{
						NodeVNF: cce.NodeVNF{
							ID:     nodeVNFID,
							NodeID: nodeID,
							VNFID:  vnfID,
						},
						Status: "deployed",
					},
				))
			},
			Entry("GET /nodes_vnfs/{id}"),
		)

		DescribeTable("404 Not Found",
			func() {
				By("Sending a GET /nodes_vnfs/{id} request")
				resp, err := apiCli.Get(
					fmt.Sprintf(
						"http://127.0.0.1:8080/nodes_vnfs/%s",
						uuid.New()))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("GET /nodes_vnfs/{id} with nonexistent ID"),
		)
	})

	Describe("PATCH /nodes_vnfs", func() {
		var (
			nodeVNFID string
		)

		BeforeEach(func() {
			nodeVNFID = postNodesVNFs(nodeID, vnfID)
		})

		DescribeTable("204 No Content",
			func(reqStr string, expectedNodeVNFResp *cce.NodeVNFResp) {
				By("Sending a PATCH /nodes_vnfs request")
				resp, err := apiCli.Patch(
					"http://127.0.0.1:8080/nodes_vnfs",
					"application/json",
					strings.NewReader(
						fmt.Sprintf(reqStr, nodeVNFID, nodeID, vnfID)))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 204 No Content response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				By("Getting the updated node")
				updatedNodeVNFResp := getNodeVNF(nodeVNFID)

				By("Verifying the node was updated")
				expectedNodeVNFResp.NodeVNF = cce.NodeVNF{
					ID:     nodeVNFID,
					NodeID: nodeID,
					VNFID:  vnfID,
				}
				Expect(updatedNodeVNFResp).To(Equal(expectedNodeVNFResp))
			},
			Entry(
				"PATCH /nodes_vnfs",
				`
				[
					{
						"id": "%s",
						"node_id": "%s",
						"vnf_id": "%s",
						"cmd": "start"
					}
				]`,
				&cce.NodeVNFResp{
					Status: "running",
				}),
		)

		DescribeTable("400 Bad Request",
			func(reqStr string, expectedResp string) {
				By("Sending a PATCH /nodes_vnfs request")
				switch strings.Count(reqStr, "%s") {
				case 2:
					reqStr = fmt.Sprintf(reqStr, nodeVNFID, nodeID)
				case 3:
					reqStr = fmt.Sprintf(reqStr, nodeVNFID, nodeID, vnfID)
				}
				resp, err := apiCli.Patch(
					"http://127.0.0.1:8080/nodes_vnfs",
					"application/json",
					strings.NewReader(reqStr))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

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
				"PATCH /nodes_vnfs without id",
				`
				[
					{
						"node_id": "123",
						"vnf_id": "456"
					}
				]`,
				"Validation failed: id not a valid uuid"),
			Entry(
				"PATCH /nodes_vnfs without node_id",
				`
				[
					{
						"id": "%s",
						"vnf_id": "%s"
					}
				]`,
				"Validation failed: node_id not a valid uuid"),
			Entry("PATCH /nodes_vnfs without vnf_id",
				`
				[
					{
						"id": "%s",
						"node_id": "%s"
					}
				]`,
				"Validation failed: vnf_id not a valid uuid"),
			Entry("PATCH /nodes_vnfs without cmd",
				`
				[
					{
						"id": "%s",
						"node_id": "%s",
						"vnf_id": "%s"
					}
				]`,
				"Validation failed: cmd missing"),
			Entry("PATCH /nodes_vnfs with invalid cmd",
				`
				[
					{
						"id": "%s",
						"node_id": "%s",
						"vnf_id": "%s",
						"cmd": "abc"
					}
				]`,
				`Validation failed: cmd "abc" is invalid`),
		)
	})

	Describe("DELETE /nodes_vnfs/{id}", func() {
		var (
			nodeVNFID string
		)

		BeforeEach(func() {
			nodeVNFID = postNodesVNFs(nodeID, vnfID)
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a DELETE /nodes_vnfs/{id} request")
				resp, err := apiCli.Delete(
					fmt.Sprintf(
						"http://127.0.0.1:8080/nodes_vnfs/%s",
						nodeVNFID))

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Verifying the node VNF was deleted")

				By("Sending a GET /nodes_vnfs/{id} request")
				resp2, err := apiCli.Get(
					fmt.Sprintf(
						"http://127.0.0.1:8080/nodes_vnfs/%s",
						nodeVNFID))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				defer resp2.Body.Close()
				Expect(resp2.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("DELETE /nodes_vnfs/{id}"),
		)

		DescribeTable("404 Not Found",
			func(id string) {
				By("Sending a DELETE /nodes_vnfs/{id} request")
				resp, err := apiCli.Delete(
					fmt.Sprintf(
						"http://127.0.0.1:8080/nodes_vnfs/%s",
						id))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry(
				"DELETE /nodes_vnfs/{id} with nonexistent ID",
				uuid.New()),
		)
	})
})
