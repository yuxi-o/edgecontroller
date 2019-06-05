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

var _ = Describe("/nodes_apps", func() {
	var (
		appID string
	)

	BeforeEach(func() {
		clearGRPCTargetsTable()
		appID = postApps("container")
	})

	Describe("POST /nodes_apps", func() {
		DescribeTable("201 Created",
			func() {
				nodeCfg := createAndRegisterNode()

				By("Sending a POST /nodes_apps request")
				resp, err := apiCli.Post(
					"http://127.0.0.1:8080/nodes_apps",
					"application/json",
					strings.NewReader(fmt.Sprintf(
						`
						{
							"node_id": "%s",
							"app_id": "%s"
						}`, nodeCfg.nodeID, appID)))
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
				"POST /nodes_apps"),
		)

		DescribeTable("400 Bad Request",
			func(req, expectedResp string) {
				By("Sending a POST /nodes_apps request")
				resp, err := apiCli.Post(
					"http://127.0.0.1:8080/nodes_apps",
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
				"POST /nodes_apps with id",
				`
				{
					"id": "123"
				}`,
				"Validation failed: id cannot be specified in POST request"),
			Entry(
				"POST /nodes_apps without node_id",
				`
				{
				}`,
				"Validation failed: node_id not a valid uuid"),
			Entry(
				"POST /nodes_apps without app_id",
				fmt.Sprintf(`
				{
					"node_id": "%s"
				}`, uuid.New()),
				"Validation failed: app_id not a valid uuid"),
		)

		DescribeTable("422 Unprocessable Entity",
			func() {
				nodeCfg := createAndRegisterNode()

				By("Sending a POST /nodes_apps request")
				postNodesApps(nodeCfg.nodeID, appID)

				By("Repeating the first POST /nodes_apps request")
				resp, err := apiCli.Post(
					"http://127.0.0.1:8080/nodes_apps",
					"application/json",
					strings.NewReader(fmt.Sprintf(
						`
						{
							"node_id": "%s",
							"app_id": "%s"
						}`, nodeCfg.nodeID, appID)))
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
					"duplicate record in nodes_apps detected for node_id %s and "+
						"app_id %s",
					nodeCfg.nodeID,
					appID)))
			},
			Entry("POST /nodes_apps with duplicate node_id and app_id"),
		)
	})

	Describe("GET /nodes_apps", func() {
		var (
			nodeCfg    *nodeConfig
			nodeAppID  string
			app2ID     string
			nodeApp2ID string
		)

		BeforeEach(func() {
			nodeCfg = createAndRegisterNode()
			nodeAppID = postNodesApps(nodeCfg.nodeID, appID)
			app2ID = postApps("container")
			nodeApp2ID = postNodesApps(nodeCfg.nodeID, app2ID)
		})

		DescribeTable("200 OK",
			func() {
				nodeApps := getNodeApps(nodeCfg.nodeID)

				By("Verifying the 2 created node apps were returned")
				Expect(nodeApps).To(ContainElement(
					&cce.NodeApp{
						ID:     nodeAppID,
						NodeID: nodeCfg.nodeID,
						AppID:  appID,
					}))
				Expect(nodeApps).To(ContainElement(
					&cce.NodeApp{
						ID:     nodeApp2ID,
						NodeID: nodeCfg.nodeID,
						AppID:  app2ID,
					}))
			},
			Entry("GET /nodes_apps"),
		)
	})

	Describe("GET /nodes_apps/{id}", func() {
		DescribeTable("200 OK",
			func() {
				nodeCfg := createAndRegisterNode()
				nodeAppID := postNodesApps(nodeCfg.nodeID, appID)
				nodeAppResp := getNodeApp(nodeAppID)

				By("Verifying the created node app was returned")
				Expect(nodeAppResp).To(Equal(
					&cce.NodeAppResp{
						NodeApp: cce.NodeApp{
							ID:     nodeAppID,
							NodeID: nodeCfg.nodeID,
							AppID:  appID,
						},
						Status: "deployed",
					},
				))
			},
			Entry("GET /nodes_apps/{id}"),
		)

		DescribeTable("404 Not Found",
			func() {
				By("Sending a GET /nodes_apps/{id} request")
				resp, err := apiCli.Get(
					fmt.Sprintf(
						"http://127.0.0.1:8080/nodes_apps/%s",
						uuid.New()))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("GET /nodes_apps/{id} with nonexistent ID"),
		)
	})

	Describe("PATCH /nodes_apps", func() {
		DescribeTable("204 No Content",
			func(reqStr string, expectedNodeAppResp *cce.NodeAppResp) {
				nodeCfg := createAndRegisterNode()
				nodeAppID := postNodesApps(nodeCfg.nodeID, appID)

				By("Sending a PATCH /nodes_apps request")
				resp, err := apiCli.Patch(
					"http://127.0.0.1:8080/nodes_apps",
					"application/json",
					strings.NewReader(
						fmt.Sprintf(reqStr, nodeAppID, nodeCfg.nodeID, appID)))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 204 No Content response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				By("Getting the updated node")
				updatedNodeAppResp := getNodeApp(nodeAppID)

				By("Verifying the node was updated")
				expectedNodeAppResp.NodeApp = cce.NodeApp{
					ID:     nodeAppID,
					NodeID: nodeCfg.nodeID,
					AppID:  appID,
				}
				Expect(updatedNodeAppResp).To(Equal(expectedNodeAppResp))
			},
			Entry(
				"PATCH /nodes_apps",
				`
				[
					{
						"id": "%s",
						"node_id": "%s",
						"app_id": "%s",
						"cmd": "start"
					}
				]`,
				&cce.NodeAppResp{
					Status: "running",
				}),
		)

		DescribeTable("400 Bad Request",
			func(reqStr string, expectedResp string) {
				nodeCfg := createAndRegisterNode()
				nodeAppID := postNodesApps(nodeCfg.nodeID, appID)

				By("Sending a PATCH /nodes_apps request")
				switch strings.Count(reqStr, "%s") {
				case 2:
					reqStr = fmt.Sprintf(reqStr, nodeAppID, nodeCfg.nodeID)
				case 3:
					reqStr = fmt.Sprintf(reqStr, nodeAppID, nodeCfg.nodeID, appID)
				}
				resp, err := apiCli.Patch(
					"http://127.0.0.1:8080/nodes_apps",
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
				"PATCH /nodes_apps without id",
				`
				[
					{
						"node_id": "123",
						"app_id": "456"
					}
				]`,
				"Validation failed: id not a valid uuid"),
			Entry(
				"PATCH /nodes_apps without node_id",
				`
				[
					{
						"id": "%s",
						"app_id": "%s"
					}
				]`,
				"Validation failed: node_id not a valid uuid"),
			Entry("PATCH /nodes_apps without app_id",
				`
				[
					{
						"id": "%s",
						"node_id": "%s"
					}
				]`,
				"Validation failed: app_id not a valid uuid"),
			Entry("PATCH /nodes_apps without cmd",
				`
				[
					{
						"id": "%s",
						"node_id": "%s",
						"app_id": "%s"
					}
				]`,
				"Validation failed: cmd missing"),
			Entry("PATCH /nodes_apps with invalid cmd",
				`
				[
					{
						"id": "%s",
						"node_id": "%s",
						"app_id": "%s",
						"cmd": "abc"
					}
				]`,
				`Validation failed: cmd "abc" is invalid`),
		)
	})

	Describe("DELETE /nodes_apps/{id}", func() {
		DescribeTable("200 OK",
			func() {
				nodeCfg := createAndRegisterNode()
				nodeAppID := postNodesApps(nodeCfg.nodeID, appID)

				By("Sending a DELETE /nodes_apps/{id} request")
				resp, err := apiCli.Delete(
					fmt.Sprintf(
						"http://127.0.0.1:8080/nodes_apps/%s",
						nodeAppID))

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Verifying the node app was deleted")

				By("Sending a GET /nodes_apps/{id} request")
				resp2, err := apiCli.Get(
					fmt.Sprintf(
						"http://127.0.0.1:8080/nodes_apps/%s",
						nodeAppID))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				defer resp2.Body.Close()
				Expect(resp2.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("DELETE /nodes_apps/{id}"),
		)

		DescribeTable("404 Not Found",
			func(id string) {
				By("Sending a DELETE /nodes_apps/{id} request")
				resp, err := apiCli.Delete(
					fmt.Sprintf(
						"http://127.0.0.1:8080/nodes_apps/%s",
						id))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry(
				"DELETE /nodes_apps/{id} with nonexistent ID",
				uuid.New()),
		)

		DescribeTable("422 Unprocessable Entity",
			func(resource, expectedResp string) {
				nodeCfg := createAndRegisterNode()
				nodeAppID := postNodesApps(nodeCfg.nodeID, appID)

				switch resource {
				case "nodes_apps_traffic_policies":
					postNodesAppsTrafficPolicies(
						nodeAppID,
						postTrafficPolicies())
				}

				By("Sending a DELETE /nodes_apps/{id} request")
				resp, err := apiCli.Delete(
					fmt.Sprintf("http://127.0.0.1:8080/nodes_apps/%s",
						nodeAppID))
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
					fmt.Sprintf(expectedResp, nodeAppID)))
			},
			Entry(
				"DELETE /nodes_apps/{id} with nodes_apps_traffic_policies record",
				"nodes_apps_traffic_policies",
				"cannot delete node_app_id %s: record in use in nodes_apps_traffic_policies",
			),
		)
	})
})
