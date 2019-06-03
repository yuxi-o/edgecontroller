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

package k8s_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	cce "github.com/smartedgemec/controller-ce"
	"github.com/smartedgemec/controller-ce/uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("/nodes_apps for k8s", func() {
	var (
		nodeID string
		appID  string
	)

	BeforeEach(func() {
		// clean
		nodeID = postNodes()

		// label node with correct id
		Expect(exec.Command("kubectl",
			"label", "nodes", "minikube", fmt.Sprintf("node-id=%s", nodeID)).Run()).To(Succeed())

		appID = postApps("container")

		// tag docker with app id
		cmd := exec.Command("docker", "tag", "nginx:1.12", fmt.Sprintf("%s:%s", appID, "latest"))
		Expect(cmd.Run()).To(Succeed())
	})

	AfterEach(func() {
		// un-label node
		Expect(exec.Command("kubectl", "label", "nodes", "minikube", "node-id-").Run()).To(Succeed())

		// clean up all k8s deployments
		cmd := exec.Command("kubectl", "delete", "--all", "deployments,pods", "--namespace=default")
		Expect(cmd.Run()).To(Succeed())
		// remove tagged docker image
		cmd = exec.Command("docker", "rmi", fmt.Sprintf("%s:%s", appID, "latest"))
		Expect(cmd.Run()).To(Succeed())
	})

	Describe("POST /nodes_apps", func() {
		DescribeTable("201 Created",
			func() {
				By("Sending a POST /nodes_apps request")
				resp, err := apiCli.Post(
					"http://127.0.0.1:8080/nodes_apps",
					"application/json",
					strings.NewReader(fmt.Sprintf(
						`
						{
							"node_id": "%s",
							"app_id": "%s"
						}`, nodeID, appID)))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 201 response")
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				By("Unmarshalling the response")
				var respBody struct {
					ID string
				}
				Expect(json.Unmarshal(body, &respBody)).To(Succeed())

				By("Verifying a UUID was returned")
				Expect(uuid.IsValid(respBody.ID)).To(BeTrue())
			},
			Entry("POST /nodes_apps"),
		)
	})

	Describe("GET /nodes_apps", func() {
		var (
			nodeAppID  string
			app2ID     string
			nodeApp2ID string
		)

		BeforeEach(func() {
			nodeAppID = postNodesApps(nodeID, appID)

			app2ID = postApps("container")

			cmd := exec.Command("docker", "tag", "nginx:1.12", fmt.Sprintf("%s:%s", app2ID, "latest"))
			Expect(cmd.Run()).To(Succeed())

			nodeApp2ID = postNodesApps(nodeID, app2ID)
		})

		AfterEach(func() {
			cmd := exec.Command("docker", "rmi", fmt.Sprintf("%s:%s", app2ID, "latest"))
			Expect(cmd.Run()).To(Succeed())
		})

		DescribeTable("200 OK",
			func() {
				By("Verifying the 2 created node apps were deployed")
				count := 0
				Eventually(func() []*cce.NodeAppResp {
					count++
					By(fmt.Sprintf("Attempt #%d: Verifying if k8s deployment status is deployed", count))
					nodeAppsResp := getNodeApps(nodeID)
					return nodeAppsResp
				}, 15*time.Second, 1*time.Second).Should(ContainElement(
					&cce.NodeAppResp{
						NodeApp: cce.NodeApp{
							ID:     nodeAppID,
							NodeID: nodeID,
							AppID:  appID,
						},
						Status: "deployed",
					},
				))

				count = 0
				Eventually(func() []*cce.NodeAppResp {
					count++
					By(fmt.Sprintf("Attempt #%d: Verifying if k8s deployment status is deployed", count))
					nodeAppsResp := getNodeApps(nodeID)
					return nodeAppsResp
				}, 15*time.Second, 1*time.Second).Should(ContainElement(
					&cce.NodeAppResp{
						NodeApp: cce.NodeApp{
							ID:     nodeApp2ID,
							NodeID: nodeID,
							AppID:  app2ID,
						},
						Status: "deployed",
					},
				))
			},
			Entry("GET /nodes_apps"),
		)
	})

	Describe("GET /nodes_apps/{id}", func() {
		var (
			nodeAppID string
		)

		BeforeEach(func() {
			nodeAppID = postNodesApps(nodeID, appID)
		})

		DescribeTable("200 OK",
			func() {
				nodeAppResp := getNodeApp(nodeAppID)
				By("Verifying the created node app was deployed")
				Expect(nodeAppResp).To(Equal(
					&cce.NodeAppResp{
						NodeApp: cce.NodeApp{
							ID:     nodeAppID,
							NodeID: nodeID,
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
		var (
			nodeAppID string
		)

		BeforeEach(func() {
			nodeAppID = postNodesApps(nodeID, appID)
		})

		DescribeTable("204 No Content",
			func(reqStr string, expectedNodeAppFull *cce.NodeAppResp) {
				By("Sending a PATCH /nodes_apps request")
				resp, err := apiCli.Patch(
					"http://127.0.0.1:8080/nodes_apps",
					"application/json",
					strings.NewReader(
						fmt.Sprintf(reqStr, nodeAppID, nodeID, appID)))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 204 No Content response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				By("Verifying the node was updated")
				expectedNodeAppFull.NodeApp = cce.NodeApp{
					ID:     nodeAppID,
					NodeID: nodeID,
					AppID:  appID,
				}
				count := 0
				Eventually(func() *cce.NodeAppResp {
					count++
					By(fmt.Sprintf("Attempt #%d: Verifying if k8s deployment status is %s",
						count, expectedNodeAppFull.Status))
					return getNodeApp(nodeAppID)
				}, 30*time.Second, 1*time.Second).Should(Equal(expectedNodeAppFull))
			},
			Entry(
				"PATCH /nodes_apps (cmd='start')",
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
					// due to limitations of minikube support other than linux
					Status: status(),
				},
			),
			Entry(
				"PATCH /nodes_apps (cmd='stop')",
				`
					[
						{
							"id": "%s",
							"node_id": "%s",
							"app_id": "%s",
							"cmd": "stop"
						}
					]`,
				&cce.NodeAppResp{
					Status: "deployed",
				},
			),
			Entry(
				"PATCH /nodes_apps (cmd='restart')",
				`
					[
						{
							"id": "%s",
							"node_id": "%s",
							"app_id": "%s",
							"cmd": "restart"
						}
					]`,
				&cce.NodeAppResp{
					// due to limitations of minikube support other than linux
					Status: status(),
				},
			),
		)
	})

	Describe("DELETE /nodes_apps/{id}", func() {
		var (
			nodeAppID string
		)

		BeforeEach(func() {
			nodeAppID = postNodesApps(nodeID, appID)
		})

		DescribeTable("200 OK",
			func() {
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
	})
})

func status() string {
	if runtime.GOOS == "linux" {
		return "deployed"
	}
	return "deploying"
}
