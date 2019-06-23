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
	"github.com/smartedgemec/controller-ce/swagger"
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

				By("Unmarshaling the response")
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
			nodeCfg *nodeConfig
		)

		BeforeEach(func() {
			clearGRPCTargetsTable()
			nodeCfg = createAndRegisterNode()
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

				By("Unmarshaling the response")
				Expect(json.Unmarshal(body, &nodes)).To(Succeed())

				By("Verifying the 2 created nodes were returned")
				Expect(nodes).To(ContainElement(
					&cce.Node{
						ID:       nodeCfg.nodeID,
						Name:     "Test Node 1",
						Location: "Localhost port 42101",
						Serial:   nodeCfg.serial,
					}))
			},
			Entry("GET /nodes"),
		)

		DescribeTable("400 Bad Request",
			func(field, value string) {
				By("Sending a GET /nodes request with a disallowed filter")
				resp, err := apiCli.Get(fmt.Sprintf(
					"http://127.0.0.1:8080/nodes?%s=%s",
					field, value,
				))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 400 Bad Request response")
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(body)).To(Equal(
					fmt.Sprintf("disallowed filter field %q\n", field),
				))
			},
			Entry("GET /nodes", "location", "usa"),
		)
	})

	Describe("GET /nodes/{id}", func() {
		DescribeTable("200 OK",
			func() {
				clearGRPCTargetsTable()
				nodeCfg := createAndRegisterNode()
				node := getNode(nodeCfg.nodeID)

				By("Verifying the created node was returned")
				Expect(node).To(Equal(
					&cce.NodeResp{
						Node: cce.Node{
							ID:       nodeCfg.nodeID,
							Name:     "Test Node 1",
							Location: "Localhost port 42101",
							Serial:   nodeCfg.serial,
						},
						NetworkInterfaces: []*cce.NetworkInterface{
							{
								ID:                "if0",
								Description:       "interface0",
								Driver:            "kernel",
								Type:              "none",
								MACAddress:        "mac0",
								VLAN:              0,
								Zones:             nil,
								FallbackInterface: "",
							},
							{
								ID:                "if1",
								Description:       "interface1",
								Driver:            "kernel",
								Type:              "none",
								MACAddress:        "mac1",
								VLAN:              1,
								Zones:             nil,
								FallbackInterface: "",
							},
							{
								ID:                "if2",
								Description:       "interface2",
								Driver:            "kernel",
								Type:              "none",
								MACAddress:        "mac2",
								VLAN:              2,
								Zones:             nil,
								FallbackInterface: "",
							},
							{
								ID:                "if3",
								Description:       "interface3",
								Driver:            "kernel",
								Type:              "none",
								MACAddress:        "mac3",
								VLAN:              3,
								Zones:             nil,
								FallbackInterface: "",
							},
						},
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
			nodeCfg *nodeConfig
		)

		BeforeEach(func() {
			clearGRPCTargetsTable()
			nodeCfg = createAndRegisterNode()
		})

		DescribeTable("204 No Content",
			func(reqStr string, expectedNodeResp *cce.NodeResp) {
				By("Sending a PATCH /nodes request")
				switch strings.Count(reqStr, "%s") {
				case 1:
					reqStr = fmt.Sprintf(reqStr, nodeCfg.nodeID)
				case 5:
					trafficPolicyID := postPolicies()
					reqStr = fmt.Sprintf(reqStr, nodeCfg.nodeID, trafficPolicyID, trafficPolicyID, trafficPolicyID,
						trafficPolicyID)
				}
				resp, err := apiCli.Patch(
					"http://127.0.0.1:8080/nodes",
					"application/json",
					strings.NewReader(reqStr))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 204 No Content response")
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				By("Getting the updated node")
				updatedNodeResp := getNode(nodeCfg.nodeID)

				By("Verifying the node was updated")
				expectedNodeResp.SetID(nodeCfg.nodeID)
				Expect(updatedNodeResp).To(Equal(expectedNodeResp))
			},
			Entry(
				"PATCH /nodes with network interfaces",
				`
				[
					{
						"id": "%s",
						"name": "node123456",
						"location": "smart edge lab",
						"serial": "abc123",
						"network_interfaces": [
							{
								"id": "if0",
								"description": "interface0",
								"driver": "userspace",
								"type": "upstream",
								"mac_address": "mac0",
								"vlan": 50,
								"zones": null,
								"fallback_interface": ""
							},
							{
								"id": "if1",
								"description": "interface1",
								"driver": "kernel",
								"type": "none",
								"mac_address": "mac1",
								"vlan": 1,
								"zones": null,
								"fallback_interface": ""
							},
							{
								"id": "if2",
								"description": "interface2",
								"driver": "kernel",
								"type": "none",
								"mac_address": "mac2",
								"vlan": 2,
								"zones": null,
								"fallback_interface": ""
							},
							{
								"id": "if3",
								"description": "interface3",
								"driver": "kernel",
								"type": "none",
								"mac_address": "mac3",
								"vlan": 3,
								"zones": null,
								"fallback_interface": ""
							}
						]
					}
				]`,
				&cce.NodeResp{
					Node: cce.Node{
						Name:     "node123456",
						Location: "smart edge lab",
						Serial:   "abc123",
					},
					NetworkInterfaces: []*cce.NetworkInterface{
						{
							ID:                "if0",
							Description:       "interface0",
							Driver:            "userspace",
							Type:              "upstream",
							MACAddress:        "mac0",
							VLAN:              50,
							Zones:             nil,
							FallbackInterface: "",
						},
						{
							ID:                "if1",
							Description:       "interface1",
							Driver:            "kernel",
							Type:              "none",
							MACAddress:        "mac1",
							VLAN:              1,
							Zones:             nil,
							FallbackInterface: "",
						},
						{
							ID:                "if2",
							Description:       "interface2",
							Driver:            "kernel",
							Type:              "none",
							MACAddress:        "mac2",
							VLAN:              2,
							Zones:             nil,
							FallbackInterface: "",
						},
						{
							ID:                "if3",
							Description:       "interface3",
							Driver:            "kernel",
							Type:              "none",
							MACAddress:        "mac3",
							VLAN:              3,
							Zones:             nil,
							FallbackInterface: "",
						},
					},
				}),
			Entry(
				"PATCH /nodes without network interfaces",
				`
				[
					{
						"id": "%s",
						"name": "node123456",
						"location": "smart edge lab",
						"serial": "abc123"
					}
				]`,
				&cce.NodeResp{
					Node: cce.Node{
						Name:     "node123456",
						Location: "smart edge lab",
						Serial:   "abc123",
					},
					NetworkInterfaces: []*cce.NetworkInterface{
						{
							ID:                "if0",
							Description:       "interface0",
							Driver:            "kernel",
							Type:              "none",
							MACAddress:        "mac0",
							VLAN:              0,
							Zones:             nil,
							FallbackInterface: "",
						},
						{
							ID:                "if1",
							Description:       "interface1",
							Driver:            "kernel",
							Type:              "none",
							MACAddress:        "mac1",
							VLAN:              1,
							Zones:             nil,
							FallbackInterface: "",
						},
						{
							ID:                "if2",
							Description:       "interface2",
							Driver:            "kernel",
							Type:              "none",
							MACAddress:        "mac2",
							VLAN:              2,
							Zones:             nil,
							FallbackInterface: "",
						},
						{
							ID:                "if3",
							Description:       "interface3",
							Driver:            "kernel",
							Type:              "none",
							MACAddress:        "mac3",
							VLAN:              3,
							Zones:             nil,
							FallbackInterface: "",
						},
					},
				}),
			Entry(
				"PATCH /nodes with traffic policies",
				`
				[
					{
						"id": "%s",
						"name": "node123456",
						"location": "smart edge lab",
						"serial": "abc123",
						"traffic_policies": [
							{
								"network_interface_id": "if0",
								"traffic_policy_id": "%s"
							},
							{
								"network_interface_id": "if1",
								"traffic_policy_id": "%s"
							},
							{
								"network_interface_id": "if2",
								"traffic_policy_id": "%s"
							},
							{
								"network_interface_id": "if3",
								"traffic_policy_id": "%s"
							}
						]
					}
				]`,
				&cce.NodeResp{
					Node: cce.Node{
						Name:     "node123456",
						Location: "smart edge lab",
						Serial:   "abc123",
					},
					NetworkInterfaces: []*cce.NetworkInterface{
						{
							ID:                "if0",
							Description:       "interface0",
							Driver:            "kernel",
							Type:              "none",
							MACAddress:        "mac0",
							VLAN:              0,
							Zones:             nil,
							FallbackInterface: "",
						},
						{
							ID:                "if1",
							Description:       "interface1",
							Driver:            "kernel",
							Type:              "none",
							MACAddress:        "mac1",
							VLAN:              1,
							Zones:             nil,
							FallbackInterface: "",
						},
						{
							ID:                "if2",
							Description:       "interface2",
							Driver:            "kernel",
							Type:              "none",
							MACAddress:        "mac2",
							VLAN:              2,
							Zones:             nil,
							FallbackInterface: "",
						},
						{
							ID:                "if3",
							Description:       "interface3",
							Driver:            "kernel",
							Type:              "none",
							MACAddress:        "mac3",
							VLAN:              3,
							Zones:             nil,
							FallbackInterface: "",
						},
					},
				}),
		)

		DescribeTable("400 Bad Request",
			func(reqStr string, expectedResp string) {
				By("Sending a PATCH /nodes request")
				if strings.Contains(reqStr, "%s") {
					reqStr = fmt.Sprintf(reqStr, nodeCfg.nodeID)
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

		DescribeTable("404 Not Found",
			func(reqStr string, expectedResp string) {
				By("Sending a PATCH /nodes request")
				resp, err := apiCli.Patch(
					"http://127.0.0.1:8080/nodes",
					"application/json",
					strings.NewReader(fmt.Sprintf(reqStr, nodeCfg.nodeID)))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 404 Not Found")
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the response body")
				Expect(string(body)).To(Equal(expectedResp))
			},
			Entry("PATCH /nodes with traffic policies and invalid traffic_policy_id",
				`
				[
					{
						"id": "%s",
						"name": "node123",
						"location": "smart edge lab",
						"serial": "abc123",
						"traffic_policies": [
							{
								"network_interface_id": "if0",
								"traffic_policy_id": "2886fc50-58a0-4dad-9853-5e0a5310a294"
							}
						]
					}
				]`,
				"traffic policy 2886fc50-58a0-4dad-9853-5e0a5310a294 not found"),
			Entry(
				"PATCH /nodes with network interfaces",
				`
				[
					{
						"id": "%s",
						"name": "node123456",
						"location": "smart edge lab",
						"serial": "abc123",
						"network_interfaces": [
							{
								"id": "if03",
								"description": "interface0",
								"driver": "userspace",
								"type": "upstream",
								"mac_address": "mac0",
								"vlan": 50,
								"zones": null,
								"fallback_interface": ""
							},
							{
								"id": "if1",
								"description": "interface1",
								"driver": "kernel",
								"type": "none",
								"mac_address": "mac1",
								"vlan": 1,
								"zones": null,
								"fallback_interface": ""
							},
							{
								"id": "if2",
								"description": "interface2",
								"driver": "kernel",
								"type": "none",
								"mac_address": "mac2",
								"vlan": 2,
								"zones": null,
								"fallback_interface": ""
							},
							{
								"id": "if3",
								"description": "interface3",
								"driver": "kernel",
								"type": "none",
								"mac_address": "mac3",
								"vlan": 3,
								"zones": null,
								"fallback_interface": ""
							}
						]
					}
				]`,
				"Network Interface if03 not found"),
		)

		DescribeTable("500 Internal server error",
			func(reqStr string, expectedResp string) {
				By("Sending a PATCH /nodes request")
				if strings.Contains(reqStr, "%s") {
					reqStr = fmt.Sprintf(reqStr, nodeCfg.nodeID)
				}
				resp, err := apiCli.Patch(
					"http://127.0.0.1:8080/nodes",
					"application/json",
					strings.NewReader(reqStr))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 400 Bad Request")
				Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the response body")
				Expect(string(body)).To(Equal(expectedResp))
			},
			Entry("PATCH /nodes without all interfaces",
				`
				[
					{
						"id": "%s",
						"name": "node123456",
						"location": "smart edge lab",
						"serial": "abc123",
						"network_interfaces": [
							{
								"id": "if0",
								"description": "interface0",
								"driver": "userspace",
								"type": "upstream",
								"mac_address": "mac0",
								"vlan": 50,
								"zones": null,
								"fallback_interface": ""
							}
						]
					}
				]`,
				"error bulk updating network interfaces: rpc error: code = FailedPrecondition desc = Network Interface if1 missing from request"), //nolint:lll
		)
	})

	Describe("DELETE /nodes/{id}", func() {
		DescribeTable("200 OK",
			func() {
				nodeID := postNodesSerial("abc-123")
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

		DescribeTable("422 Unprocessable Entity",
			func(resource, expectedResp string) {
				// we need a new nodeCfg because postNodesDNSConfigs has a duplicate check on node_id
				clearGRPCTargetsTable()
				nodeCfg := createAndRegisterNode()
				switch resource {
				case "nodes_apps":
					postNodesApps(
						nodeCfg.nodeID,
						postApps("container"))
				case "nodes_dns_configs":
					postNodesDNSConfigs(
						nodeCfg.nodeID,
						postDNSConfigs())
				}

				By("Sending a DELETE /nodes/{id} request")
				resp, err := apiCli.Delete(
					fmt.Sprintf("http://127.0.0.1:8080/nodes/%s",
						nodeCfg.nodeID))
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
					fmt.Sprintf(expectedResp, nodeCfg.nodeID)))
			},
			Entry(
				"DELETE /nodes/{id} with nodes_apps record",
				"nodes_apps",
				"cannot delete node_id %s: record in use in nodes_apps",
			),
			Entry(
				"DELETE /nodes/{id} with nodes_dns_configs record",
				"nodes_dns_configs",
				"cannot delete node_id %s: record in use in nodes_dns_configs",
			),
		)
	})

	Describe("GET /nodes/{node_id}/interfaces", func() {
		DescribeTable("200 OK",
			func() {
				clearGRPCTargetsTable()
				nodeCfg := createAndRegisterNode()
				By("Sending a GET /nodes/{node_id}/interfaces request")
				resp, err := apiCli.Get(fmt.Sprintf("http://127.0.0.1:8080/nodes/%s/interfaces", nodeCfg.nodeID))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 200 OK response")
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				var ifaces swagger.InterfaceList

				By("Unmarshaling the response")
				Expect(json.Unmarshal(body, &ifaces)).To(Succeed())

				By("Verifying the created node was returned")
				Expect(ifaces).To(Equal(
					swagger.InterfaceList{
						Interfaces: []swagger.InterfaceSummary{
							{
								ID:                "if0",
								Description:       "interface0",
								Driver:            "kernel",
								Type:              "none",
								MACAddress:        "mac0",
								VLAN:              0,
								Zones:             nil,
								FallbackInterface: "",
							},
							{
								ID:                "if1",
								Description:       "interface1",
								Driver:            "kernel",
								Type:              "none",
								MACAddress:        "mac1",
								VLAN:              1,
								Zones:             nil,
								FallbackInterface: "",
							},
							{
								ID:                "if2",
								Description:       "interface2",
								Driver:            "kernel",
								Type:              "none",
								MACAddress:        "mac2",
								VLAN:              2,
								Zones:             nil,
								FallbackInterface: "",
							},
							{
								ID:                "if3",
								Description:       "interface3",
								Driver:            "kernel",
								Type:              "none",
								MACAddress:        "mac3",
								VLAN:              3,
								Zones:             nil,
								FallbackInterface: "",
							},
						},
					},
				))
			},
			Entry("GET /nodes/{node_id}/interfaces"),
		)

		DescribeTable("404 Not Found",
			func() {
				By("Sending a GET /nodes/{node_id}/interfaces/{interface_id} request")
				clearGRPCTargetsTable()
				nodeCfg := createAndRegisterNode()
				resp, err := apiCli.Get(
					fmt.Sprintf("http://127.0.0.1:8080/nodes/%s/interfaces/%s",
						nodeCfg.nodeID,
						uuid.New()))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 404 Not Found response")
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("GET /nodes/{node_id}/interfaces/{interface_id} with nonexistent ID"),
		)
	})

	Describe("PATCH /nodes/{node_id}/interfaces", func() {
		var (
			nodeCfg *nodeConfig
		)

		BeforeEach(func() {
			clearGRPCTargetsTable()
			nodeCfg = createAndRegisterNode()
		})

		DescribeTable("200 OK",
			func(reqStr string, expectedNodeResp *cce.NodeResp) {
				By("Sending a PATCH /nodes/{node_id}/interfaces request")
				switch strings.Count(reqStr, "%s") {
				case 1:
					reqStr = fmt.Sprintf(reqStr, nodeCfg.nodeID)
				case 5:
					trafficPolicyID := postPolicies()
					reqStr = fmt.Sprintf(reqStr, nodeCfg.nodeID, trafficPolicyID, trafficPolicyID, trafficPolicyID,
						trafficPolicyID)
				}
				resp, err := apiCli.Patch(
					fmt.Sprintf("http://127.0.0.1:8080/nodes/%s/interfaces", nodeCfg.nodeID),
					"application/json",
					strings.NewReader(reqStr))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 200 OK response")
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Getting the updated node")
				updatedNodeResp := getNode(nodeCfg.nodeID)

				By("Verifying the node was updated")
				expectedNodeResp.SetID(nodeCfg.nodeID)
				Expect(updatedNodeResp.NetworkInterfaces).To(Equal(expectedNodeResp.NetworkInterfaces))
			},
			Entry(
				"PATCH /nodes/{node_id}/interfaces with network interfaces",
				`
					{
						"interfaces": [
							{
								"id": "if0",
								"description": "CHANGED IN PATCH",
								"driver": "userspace",
								"type": "upstream",
								"mac_address": "mac0",
								"vlan": 50,
								"zones": null,
								"fallback_interface": ""
							},
							{
								"id": "if1",
								"description": "interface1",
								"driver": "kernel",
								"type": "none",
								"mac_address": "mac1",
								"vlan": 1,
								"zones": null,
								"fallback_interface": ""
							},
							{
								"id": "if2",
								"description": "interface2",
								"driver": "kernel",
								"type": "none",
								"mac_address": "mac2",
								"vlan": 2,
								"zones": null,
								"fallback_interface": ""
							},
							{
								"id": "if3",
								"description": "interface3",
								"driver": "kernel",
								"type": "none",
								"mac_address": "mac3",
								"vlan": 3,
								"zones": null,
								"fallback_interface": ""
							}
						]
					}
				`,
				&cce.NodeResp{
					NetworkInterfaces: []*cce.NetworkInterface{
						{
							ID:                "if0",
							Description:       "CHANGED IN PATCH",
							Driver:            "userspace",
							Type:              "upstream",
							MACAddress:        "mac0",
							VLAN:              50,
							Zones:             nil,
							FallbackInterface: "",
						},
						{
							ID:                "if1",
							Description:       "interface1",
							Driver:            "kernel",
							Type:              "none",
							MACAddress:        "mac1",
							VLAN:              1,
							Zones:             nil,
							FallbackInterface: "",
						},
						{
							ID:                "if2",
							Description:       "interface2",
							Driver:            "kernel",
							Type:              "none",
							MACAddress:        "mac2",
							VLAN:              2,
							Zones:             nil,
							FallbackInterface: "",
						},
						{
							ID:                "if3",
							Description:       "interface3",
							Driver:            "kernel",
							Type:              "none",
							MACAddress:        "mac3",
							VLAN:              3,
							Zones:             nil,
							FallbackInterface: "",
						},
					},
				}),
		)

		DescribeTable("400 Bad Request",
			func(reqStr string, expectedResp string) {
				By("Sending a PATCH /nodes/{node_id}/interfaces request")
				if strings.Contains(reqStr, "%s") {
					reqStr = fmt.Sprintf(reqStr, nodeCfg.nodeID)
				}
				resp, err := apiCli.Patch(
					fmt.Sprintf("http://127.0.0.1:8080/nodes/%s/interfaces", nodeCfg.nodeID),
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
				"PATCH /nodes/{node_id}/interfaces by leaving the driver of an interface empty",
				`
				{
					"interfaces": [
						{
							"id": "if0",
							"description": "",
							"driver": "",
							"type": "",
							"mac_address": "mac0",
							"vlan": 50,
							"zones": null,
							"fallback_interface": ""
						}
					]
				}
				`,
				"Validation failed: network_interfaces[0].driver must be one of [kernel, userspace]"),
			Entry(
				"PATCH /nodes/{node_id}/interfaces by leaving the type of an interface empty",
				`
					{
						"interfaces": [
							{
								"id": "if0",
								"description": "",
								"driver": "kernel",
								"type": "",
								"mac_address": "mac0",
								"vlan": 50,
								"zones": null,
								"fallback_interface": ""
							}
						]
					}
					`,
				"Validation failed: network_interfaces[0].type must be one of [none, upstream, "+
					"downstream, bidirectional, breakout]"),
		)

		DescribeTable("404 Not Found",
			func(reqStr string, expectedResp string) {
				By("Sending a PATCH /nodes/{node_id}/interfaces request")
				resp, err := apiCli.Patch(
					fmt.Sprintf("http://127.0.0.1:8080/nodes/%s/interfaces", nodeCfg.nodeID),
					"application/json",
					strings.NewReader(reqStr))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 404 Not Found")
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the response body")
				Expect(string(body)).To(Equal(expectedResp))
			},
			Entry(
				"PATCH /nodes with network interfaces",
				`
				{
					"interfaces": [
						{
							"id": "if03",
							"description": "interface0",
							"driver": "userspace",
							"type": "upstream",
							"mac_address": "mac0",
							"vlan": 50,
							"zones": null,
							"fallback_interface": ""
						},
						{
							"id": "if1",
							"description": "interface1",
							"driver": "kernel",
							"type": "none",
							"mac_address": "mac1",
							"vlan": 1,
							"zones": null,
							"fallback_interface": ""
						},
						{
							"id": "if2",
							"description": "interface2",
							"driver": "kernel",
							"type": "none",
							"mac_address": "mac2",
							"vlan": 2,
							"zones": null,
							"fallback_interface": ""
						},
						{
							"id": "if3",
							"description": "interface3",
							"driver": "kernel",
							"type": "none",
							"mac_address": "mac3",
							"vlan": 3,
							"zones": null,
							"fallback_interface": ""
						}
					]
				}
				`,
				"Network Interface if03 not found"),
		)

		DescribeTable("500 Internal server error",
			func(reqStr string, expectedResp string) {
				By("Sending a PATCH /nodes request")
				if strings.Contains(reqStr, "%s") {
					reqStr = fmt.Sprintf(reqStr, nodeCfg.nodeID)
				}
				resp, err := apiCli.Patch(
					fmt.Sprintf("http://127.0.0.1:8080/nodes/%s/interfaces", nodeCfg.nodeID),
					"application/json",
					strings.NewReader(reqStr))
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()

				By("Verifying a 500 Internal server error")
				Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the response body")
				Expect(string(body)).To(Equal(expectedResp))
			},
			Entry("PATCH /nodes/{node_id}/interfaces without all interfaces",
				`
				{
					"interfaces": [
						{
							"id": "if0",
							"description": "interface0",
							"driver": "userspace",
							"type": "upstream",
							"mac_address": "mac0",
							"vlan": 50,
							"zones": null,
							"fallback_interface": ""
						}
					]
				}
				`,
				"error bulk updating network interfaces: rpc error: code = FailedPrecondition desc = Network Interface if1 missing from request"), //nolint:lll
		)
	})
})
