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

var _ = Describe("/nodes_vnfs_traffic_policies", func() {
	var (
		nodeID          string
		vnfID           string
		nodeVNFID       string
		trafficPolicyID string
	)

	BeforeEach(func() {
		nodeID = postNodes()
		vnfID = postVNFs("container")
		nodeVNFID = postNodesVNFs(nodeID, vnfID)
		trafficPolicyID = postTrafficPolicies()
	})

	Describe("POST /nodes_vnfs_traffic_policies", func() {
		DescribeTable("201 Created",
			func() {
				By("Sending a POST /nodes_vnfs_traffic_policies request")
				resp, err := apiCli.Post(
					"http://127.0.0.1:8080/nodes_vnfs_traffic_policies",
					"application/json",
					strings.NewReader(fmt.Sprintf(
						`
						{
							"nodes_vnfs_id": "%s",
							"traffic_policy_id": "%s"
						}`, nodeVNFID, trafficPolicyID)))
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
				"POST /nodes_vnfs_traffic_policies"),
		)

		DescribeTable("400 Bad Request",
			func(req, expectedResp string) {
				By("Sending a POST /nodes_vnfs_traffic_policies request")
				resp, err := apiCli.Post(
					"http://127.0.0.1:8080/nodes_vnfs_traffic_policies",
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
				"POST /nodes_vnfs_traffic_policies with id",
				`
				{
					"id": "123"
				}`,
				"Validation failed: id cannot be specified in POST request"),
			Entry(
				"POST /nodes_vnfs_traffic_policies without nodes_vnfs_id",
				`
				{
				}`,
				"Validation failed: nodes_vnfs_id not a valid uuid"),
			Entry(
				"POST /nodes_vnfs_traffic_policies without traffic_policy_id",
				fmt.Sprintf(`
				{
					"nodes_vnfs_id": "%s"
				}`, uuid.New()),
				"Validation failed: traffic_policy_id not a valid uuid"),
		)

		DescribeTable("422 Unprocessable Entity",
			func() {
				var (
					resp *http.Response
					err  error
				)

				By("Sending a POST /nodes_vnfs_traffic_policies request")
				postNodesVNFsTrafficPolicies(nodeVNFID, trafficPolicyID)

				By("Repeating the first POST /nodes_vnfs_traffic_policies request")
				resp, err = apiCli.Post(
					"http://127.0.0.1:8080/nodes_vnfs_traffic_policies",
					"application/json",
					strings.NewReader(fmt.Sprintf(
						`
						{
							"nodes_vnfs_id": "%s",
							"traffic_policy_id": "%s"
						}`, nodeVNFID, trafficPolicyID)))
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
					"duplicate record detected for nodes_vnfs_id %s and "+
						"traffic_policy_id %s",
					nodeVNFID,
					trafficPolicyID)))
			},
			Entry("POST /nodes_vnfs_traffic_policies with duplicate nodes_vnfs_id and traffic_policy_id"),
		)
	})

	Describe("GET /nodes_vnfs_traffic_policies", func() {
		var (
			nodeVNFTrafficPolicyID  string
			node2ID                 string
			vnf2ID                  string
			nodeVNF2ID              string
			trafficPolicy2ID        string
			nodeVNFTrafficPolicy2ID string
		)

		BeforeEach(func() {
			nodeVNFTrafficPolicyID = postNodesVNFsTrafficPolicies(nodeVNFID, trafficPolicyID)
			node2ID = postNodes()
			vnf2ID = postVNFs("container")
			nodeVNF2ID = postNodesVNFs(node2ID, vnf2ID)
			trafficPolicy2ID = postTrafficPolicies()
			nodeVNFTrafficPolicy2ID = postNodesVNFsTrafficPolicies(nodeVNF2ID, trafficPolicy2ID)
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a GET /nodes_vnfs_traffic_policies request")
				resp, err := apiCli.Get("http://127.0.0.1:8080/nodes_vnfs_traffic_policies")
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 200 OK response")
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				var nodeVNFTrafficPolicies []*cce.NodeVNFTrafficPolicy

				By("Unmarshalling the response")
				Expect(json.Unmarshal(body, &nodeVNFTrafficPolicies)).To(Succeed())

				By("Verifying the 2 created node VNF traffic policies were returned")
				Expect(nodeVNFTrafficPolicies).To(ContainElement(
					&cce.NodeVNFTrafficPolicy{
						ID:              nodeVNFTrafficPolicyID,
						NodeVNFID:       nodeVNFID,
						TrafficPolicyID: trafficPolicyID,
					}))
				Expect(nodeVNFTrafficPolicies).To(ContainElement(
					&cce.NodeVNFTrafficPolicy{
						ID:              nodeVNFTrafficPolicy2ID,
						NodeVNFID:       nodeVNF2ID,
						TrafficPolicyID: trafficPolicy2ID,
					}))
			},
			Entry("GET /nodes_vnfs_traffic_policies"),
		)
	})

	Describe("GET /nodes_vnfs_traffic_policies/{id}", func() {
		var (
			nodeVNFTrafficPolicyID string
		)

		BeforeEach(func() {
			nodeVNFTrafficPolicyID = postNodesVNFsTrafficPolicies(nodeVNFID, trafficPolicyID)
		})

		DescribeTable("200 OK",
			func() {
				NodeVNFTrafficPolicy := getNodeVNFTrafficPolicy(nodeVNFTrafficPolicyID)

				By("Verifying the created node VNF traffic policy was returned")
				Expect(NodeVNFTrafficPolicy).To(Equal(
					&cce.NodeVNFTrafficPolicy{
						ID:              nodeVNFTrafficPolicyID,
						NodeVNFID:       nodeVNFID,
						TrafficPolicyID: trafficPolicyID,
					},
				))
			},
			Entry("GET /nodes_vnfs_traffic_policies/{id}"),
		)

		DescribeTable("404 Not Found",
			func() {
				By("Sending a GET /nodes_vnfs_traffic_policies/{id} request")
				resp, err := apiCli.Get(
					fmt.Sprintf(
						"http://127.0.0.1:8080/nodes_vnfs_traffic_policies/%s",
						uuid.New()))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("GET /nodes_vnfs_traffic_policies/{id} with nonexistent ID"),
		)
	})

	Describe("DELETE /nodes_vnfs_traffic_policies/{id}", func() {
		var (
			nodeVNFTrafficPolicyID string
		)

		BeforeEach(func() {
			nodeVNFTrafficPolicyID = postNodesVNFsTrafficPolicies(nodeVNFID, trafficPolicyID)
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a DELETE /nodes_vnfs_traffic_policies/{id} request")
				resp, err := apiCli.Delete(
					fmt.Sprintf(
						"http://127.0.0.1:8080/nodes_vnfs_traffic_policies/%s",
						nodeVNFTrafficPolicyID))

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Verifying the node VNF traffic policy was deleted")

				By("Sending a GET /nodes_vnfs_traffic_policies/{id} request")
				resp2, err := apiCli.Get(
					fmt.Sprintf(
						"http://127.0.0.1:8080/nodes_vnfs_traffic_policies/%s",
						nodeVNFTrafficPolicyID))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				defer resp2.Body.Close()
				Expect(resp2.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("DELETE /nodes_vnfs_traffic_policies/{id}"),
		)

		DescribeTable("404 Not Found",
			func(id string) {
				By("Sending a DELETE /nodes_vnfs_traffic_policies/{id} request")
				resp, err := apiCli.Delete(
					fmt.Sprintf(
						"http://127.0.0.1:8080/nodes_vnfs_traffic_policies/%s",
						id))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry(
				"DELETE /nodes_vnfs_traffic_policies/{id} with nonexistent ID",
				uuid.New()),
		)
	})
})
