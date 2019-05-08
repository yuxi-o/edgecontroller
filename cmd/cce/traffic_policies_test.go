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

var _ = Describe("/traffic_policies", func() {
	var trafficPolicy *cce.TrafficPolicy

	BeforeEach(func() {
		trafficPolicy = &cce.TrafficPolicy{
			Rules: []*cce.TrafficRule{
				{
					Description: "test-rule-1",
					Priority:    1,
					Source: &cce.TrafficSelector{
						Description: "test-source-1",
						MACs: &cce.MACFilter{
							MACAddresses: []string{
								"F0-59-8E-7B-36-8A",
								"23-20-8E-15-89-D1",
								"35-A4-38-73-35-45",
							},
						},
						IP: &cce.IPFilter{
							Address:   "223.1.1.0",
							Mask:      16,
							BeginPort: 2000,
							EndPort:   2012,
							Protocol:  "tcp",
						},
						GTP: &cce.GTPFilter{
							Address: "10.6.7.2",
							Mask:    12,
							IMSIs: []string{
								"310150123456789",
								"310150123456790",
								"310150123456791",
							},
						},
					},
					Destination: &cce.TrafficSelector{
						Description: "test-destination-1",
						MACs: &cce.MACFilter{
							MACAddresses: []string{
								"7D-C2-3A-1C-63-D9",
								"E9-6B-D1-D2-1A-6B",
								"C8-32-A9-43-85-55",
							},
						},
						IP: &cce.IPFilter{
							Address:   "64.1.1.0",
							Mask:      16,
							BeginPort: 1000,
							EndPort:   1012,
							Protocol:  "tcp",
						},
						GTP: &cce.GTPFilter{
							Address: "108.6.7.2",
							Mask:    4,
							IMSIs: []string{
								"310150123456792",
								"310150123456793",
								"310150123456794",
							},
						},
					},
					Target: &cce.TrafficTarget{
						Description: "test-target-1",
						Action:      "accept",
						MAC: &cce.MACModifier{
							MACAddress: "C7-5A-E7-98-1B-A3",
						},
						IP: &cce.IPModifier{
							Address: "123.2.3.4",
							Port:    1600,
						},
					},
				},
			},
		}
	})

	Describe("POST /traffic_policies", func() {
		DescribeTable("201 Created",
			func(req string) {
				By("Sending a POST /traffic_policies request")
				resp, err := http.Post(
					"http://127.0.0.1:8080/traffic_policies",
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
				"POST /traffic_policies",
				`
				{
					"rules": [{
						"description": "rule1",
						"priority": 1,
						"source": {
							"description": "source1",
							"macs": {
								"mac_addresses": [
									"F0-59-8E-7B-36-8A",
									"23-20-8E-15-89-D1",
									"35-A4-38-73-35-45"
								]
							},
							"ip": {
								"address": "223.1.1.0",
								"mask": 16,
								"begin_port": 2000,
								"end_port": 2012,
								"protocol": "tcp"
							},
							"gtp": {
								"address": "10.6.7.2",
								"mask": 12,
								"imsis": [
									"310150123456789",
									"310150123456790",
									"310150123456791"
								]
							}
						},
						"destination": {
							"description": "destination1",
							"macs": {
								"mac_addresses": [
									"7D-C2-3A-1C-63-D9",
									"E9-6B-D1-D2-1A-6B",
									"C8-32-A9-43-85-55"
								]
							},
							"ip": {
								"address": "64.1.1.0",
								"mask": 16,
								"begin_port": 1000,
								"end_port": 1012,
								"protocol": "tcp"
							},
							"gtp": {
								"address": "108.6.7.2",
								"mask": 4,
								"imsis": [
									"310150123456792",
									"310150123456793",
									"310150123456794"
								]
							}
						},
						"target": {
							"description": "target1",
							"action": "accept",
							"mac": {
								"mac_address": "C7-5A-E7-98-1B-A3"
							},
							"ip": {
								"address": "123.2.3.4",
								"port": 1600
							}
						}
					}]
				}`),
		)

		DescribeTable("400 Bad Request",
			func(req, expectedResp string) {
				By("Sending a POST /traffic_policies request")
				resp, err := http.Post(
					"http://127.0.0.1:8080/traffic_policies",
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
				"POST /traffic_policies with id",
				`
				{
					"id": "123"
				}`,
				"Validation failed: id cannot be specified in POST request"),
			Entry(
				"POST /traffic_policies without rules",
				`
				{
					"rules": []
				}`,
				"Validation failed: rules cannot be empty"),
			Entry(
				"POST /traffic_policies without rules[0].description",
				`
				{
					"rules": [{}]
				}`,
				"Validation failed: rules[0].description cannot be empty"),
			Entry(
				"POST /traffic_policies with rules[0].priority not in [1..65536]", //nolint:lll
				`
				{
					"rules": [{
						"description": "rule1",
						"priority": 65537
					}]
				}`,
				"Validation failed: rules[0].priority must be in [1..65536]"),
			Entry("POST /traffic_policies without rules[0].source",
				`
				{
					"rules": [{
						"description": "rule1",
						"priority": 1
					}]
				}`,
				"Validation failed: rules[0].source cannot be empty"),
			Entry("POST /traffic_policies without rules[0].destination",
				`
				{
					"rules": [{
						"description": "rule1",
						"priority": 1,
						"source": {
							"description": "source1",
							"macs": {
								"mac_addresses": []
							}
						}
					}]
				}`,
				"Validation failed: rules[0].destination cannot be empty"),
			Entry("POST /traffic_policies without rules[0].target",
				`
				{
					"rules": [{
						"description": "rule1",
						"priority": 1,
						"source": {
							"description": "source1",
							"macs": {
								"mac_addresses": []
							}
						},
						"destination": {
							"description": "destination1",
							"macs": {
								"mac_addresses": []
							}
						}
					}]
				}`,
				"Validation failed: rules[0].target cannot be empty"),
			Entry("POST /traffic_policies without rules[0].source.description",
				`
				{
					"rules": [{
						"description": "rule1",
						"priority": 1,
						"source": {
							"macs": {
								"mac_addresses": []
							}
						}
					}]
				}`,
				"Validation failed: rules[0].source.description cannot be empty"), //nolint:lll
			Entry("POST /traffic_policies without rules[0].source.macs|ip|gtp",
				`
				{
					"rules": [{
						"description": "rule1",
						"priority": 1,
						"source": {
							"description": "source1"
						}
					}]
				}`,
				"Validation failed: rules[0].source.macs|ip|gtp cannot all be nil"), //nolint:lll
			Entry("POST /traffic_policies with invalid rules[0].source.macs.mac_addresses[0]", //nolint:lll
				`
				{
					"rules": [{
						"description": "rule1",
						"priority": 1,
						"source": {
							"description": "source1",
							"macs": {
								"mac_addresses": [
									"abc-def"
								]
							}
						}
					}]
				}`,
				"Validation failed: rules[0].source.macs.mac_addresses[0] could not be parsed (address abc-def: invalid MAC address)"), //nolint:lll
			Entry("POST /traffic_policies with invalid rules[0].source.ip.address", //nolint:lll
				`
				{
					"rules": [{
						"description": "rule1",
						"priority": 1,
						"source": {
							"description": "source1",
							"ip": {
								"address": "2234.1.1.0"
							}
						}
					}]
				}`,
				"Validation failed: rules[0].source.ip.address could not be parsed"), //nolint:lll
			Entry("POST /traffic_policies with rules[0].source.ip.mask not in [0..128]", //nolint:lll
				`
				{
					"rules": [{
						"description": "rule1",
						"priority": 1,
						"source": {
							"description": "source1",
							"ip": {
								"address": "223.1.1.0",
								"mask": 129
							}
						}
					}]
				}`,
				"Validation failed: rules[0].source.ip.mask must be in [0..128]"), //nolint:lll
			Entry("POST /traffic_policies with rules[0].source.ip.begin_port not in [1..65536]", //nolint:lll
				`
				{
					"rules": [{
						"description": "rule1",
						"priority": 1,
						"source": {
							"description": "source1",
							"ip": {
								"address": "223.1.1.0",
								"mask": 128,
								"begin_port": 65537
							}
						}
					}]
				}`,
				"Validation failed: rules[0].source.ip.begin_port must be in [1..65536]"), //nolint:lll
			Entry("POST /traffic_policies with rules[0].source.ip.end_port not in [1..65536]", //nolint:lll
				`
				{
					"rules": [{
						"description": "rule1",
						"priority": 1,
						"source": {
							"description": "source1",
							"ip": {
								"address": "223.1.1.0",
								"mask": 128,
								"begin_port": 65,
								"end_port": 65537
							}
						}
					}]
				}`,
				"Validation failed: rules[0].source.ip.end_port must be in [1..65536]"), //nolint:lll
			Entry("POST /traffic_policies with invalid rules[0].source.ip.protocol", //nolint:lll
				`
				{
					"rules": [{
						"description": "rule1",
						"priority": 1,
						"source": {
							"description": "source1",
							"ip": {
								"address": "223.1.1.0",
								"mask": 128,
								"begin_port": 65,
								"end_port": 65536,
								"protocol": "udtcp"
							}
						}
					}]
				}`,
				"Validation failed: rules[0].source.ip.protocol must be one of [tcp, udp, icmp, sctp]"), //nolint:lll
			Entry("POST /traffic_policies without rules[0].target.description", //nolint:lll
				`
				{
					"rules": [{
						"description": "rule1",
						"priority": 1,
						"source": {
							"description": "source1",
							"macs": {
								"mac_addresses": []
							}
						},
						"destination": {
							"description": "destination1",
							"macs": {
								"mac_addresses": []
							}
						},
						"target": {
							"action": "accept"
						}
					}]
				}`,
				"Validation failed: rules[0].target.description cannot be empty"), //nolint:lll
			Entry("POST /traffic_policies with invalid rules[0].target.action",
				`
				{
					"rules": [{
						"description": "rule1",
						"priority": 1,
						"source": {
							"description": "source1",
							"macs": {
								"mac_addresses": []
							}
						},
						"destination": {
							"description": "destination1",
							"macs": {
								"mac_addresses": []
							}
						},
						"target": {
							"description": "target1",
							"action": "forward"
						}
					}]
				}`,
				"Validation failed: rules[0].target.action must be one of [accept, reject, drop]"), //nolint:lll
			Entry("POST /traffic_policies without rules[0].target.mac|ip",
				`
				{
					"rules": [{
						"description": "rule1",
						"priority": 1,
						"source": {
							"description": "source1",
							"macs": {
								"mac_addresses": []
							}
						},
						"destination": {
							"description": "destination1",
							"macs": {
								"mac_addresses": []
							}
						},
						"target": {
							"description": "target1",
							"action": "accept"
						}
					}]
				}`,
				"Validation failed: rules[0].target.mac|ip cannot both be nil"),
			Entry("POST /traffic_policies with invalid rules[0].target.mac",
				`
				{
					"rules": [{
						"description": "rule1",
						"priority": 1,
						"source": {
							"description": "source1",
							"macs": {
								"mac_addresses": []
							}
						},
						"destination": {
							"description": "destination1",
							"macs": {
								"mac_addresses": []
							}
						},
						"target": {
							"description": "target1",
							"action": "accept",
							"mac": {
								"mac_address": "abc-123"
							}
						}
					}]
				}`,
				"Validation failed: rules[0].target.mac.mac_address could not be parsed (address abc-123: invalid MAC address)"), //nolint:lll
			Entry("POST /traffic_policies with invalid rules[0].target.ip.address", //nolint:lll
				`
				{
					"rules": [{
						"description": "rule1",
						"priority": 1,
						"source": {
							"description": "source1",
							"macs": {
								"mac_addresses": []
							}
						},
						"destination": {
							"description": "destination1",
							"macs": {
								"mac_addresses": []
							}
						},
						"target": {
							"description": "target1",
							"action": "accept",
							"ip": {
								"address": "424.2.2.93"
							}
						}
					}]
				}`,
				"Validation failed: rules[0].target.ip.address could not be parsed"), //nolint:lll
			Entry("POST /traffic_policies with rules[0].target.ip.port not in [1..65536]", //nolint:lll
				`
				{
					"rules": [{
						"description": "rule1",
						"priority": 1,
						"source": {
							"description": "source1",
							"macs": {
								"mac_addresses": []
							}
						},
						"destination": {
							"description": "destination1",
							"macs": {
								"mac_addresses": []
							}
						},
						"target": {
							"description": "target1",
							"action": "accept",
							"ip": {
								"address": "123.2.3.4",
								"port": 65537
							}
						}
					}]
				}`,
				"Validation failed: rules[0].target.ip.port must be in [1..65536]"), //nolint:lll
		)
	})

	Describe("GET /traffic_policies", func() {
		var (
			trafficPolicyID  string
			trafficPolicy2ID string
		)

		BeforeEach(func() {
			trafficPolicyID = postTrafficPolicies()
			trafficPolicy2ID = postTrafficPolicies()
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a GET /traffic_policies request")
				resp, err := http.Get("http://127.0.0.1:8080/traffic_policies")

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				var trafficPolicies []*cce.TrafficPolicy

				By("Unmarshalling the response")
				Expect(json.Unmarshal(body, &trafficPolicies)).To(Succeed())

				By("Verifying the 2 created traffic policies were returned")
				trafficPolicy.SetID(trafficPolicyID)
				Expect(trafficPolicies).To(ContainElement(trafficPolicy))
				tp2 := trafficPolicy
				tp2.SetID(trafficPolicy2ID)
				Expect(trafficPolicies).To(ContainElement(tp2))
			},
			Entry("GET /traffic_policies"),
		)
	})

	Describe("GET /traffic_policies/{id}", func() {
		var (
			trafficPolicyID string
		)

		BeforeEach(func() {
			trafficPolicyID = postTrafficPolicies()
		})

		DescribeTable("200 OK",
			func() {
				By("Verifying the created traffic policy was returned")
				trafficPolicy.SetID(trafficPolicyID)
				Expect(trafficPolicy).To(Equal(trafficPolicy))
			},
			Entry("GET /traffic_policies/{id}"),
		)

		DescribeTable("404 Not Found",
			func() {
				By("Sending a GET /traffic_policies/{id} request")
				resp, err := http.Get(
					fmt.Sprintf("http://127.0.0.1:8080/traffic_policies/%s",
						uuid.New()))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("GET /traffic_policies/{id} with nonexistent ID"),
		)
	})

	Describe("PATCH /traffic_policies", func() {
		var (
			trafficPolicyID string
		)

		BeforeEach(func() {
			trafficPolicyID = postTrafficPolicies()
		})

		DescribeTable("204 No Content",
			func(reqStr string) {
				By("Sending a PATCH /traffic_policies request")
				req, err := http.NewRequest(
					http.MethodPatch,
					"http://127.0.0.1:8080/traffic_policies",
					strings.NewReader(fmt.Sprintf(reqStr, trafficPolicyID)))
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 204 No Content response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				By("Getting the updated traffic policy")
				updatedPolicy := getTrafficPolicy(trafficPolicyID)

				By("Verifying the traffic policy was updated")
				expectedPolicy := trafficPolicy
				expectedPolicy.SetID(trafficPolicyID)
				expectedPolicy.Rules[0].Description = "test-rule-2"
				expectedPolicy.Rules[0].Priority = 2
				expectedPolicy.Rules[0].Source.Description = "test-source-2"
				expectedPolicy.Rules[0].Destination.Description =
					"test-destination-2"
				expectedPolicy.Rules[0].Target.Description =
					"test-target-2"
				Expect(updatedPolicy).To(Equal(expectedPolicy))
			},
			Entry(
				"PATCH /traffic_policies",
				`
				[{
					"id": "%s",
					"rules": [{
						"description": "test-rule-2",
						"priority": 2,
						"source": {
							"description": "test-source-2",
							"macs": {
								"mac_addresses": [
									"F0-59-8E-7B-36-8A",
									"23-20-8E-15-89-D1",
									"35-A4-38-73-35-45"
								]
							},
							"ip": {
								"address": "223.1.1.0",
								"mask": 16,
								"begin_port": 2000,
								"end_port": 2012,
								"protocol": "tcp"
							},
							"gtp": {
								"address": "10.6.7.2",
								"mask": 12,
								"imsis": [
									"310150123456789",
									"310150123456790",
									"310150123456791"
								]
							}
						},
						"destination": {
							"description": "test-destination-2",
							"macs": {
								"mac_addresses": [
									"7D-C2-3A-1C-63-D9",
									"E9-6B-D1-D2-1A-6B",
									"C8-32-A9-43-85-55"
								]
							},
							"ip": {
								"address": "64.1.1.0",
								"mask": 16,
								"begin_port": 1000,
								"end_port": 1012,
								"protocol": "tcp"
							},
							"gtp": {
								"address": "108.6.7.2",
								"mask": 4,
								"imsis": [
									"310150123456792",
									"310150123456793",
									"310150123456794"
								]
							}
						},
						"target": {
							"description": "test-target-2",
							"action": "accept",
							"mac": {
								"mac_address": "C7-5A-E7-98-1B-A3"
							},
							"ip": {
								"address": "123.2.3.4",
								"port": 1600
							}
						}
					}]
				}]`),
		)

		DescribeTable("400 Bad Request",
			func(reqStr string, expectedResp string) {
				By("Sending a PATCH /traffic_policies request")
				if strings.Contains(reqStr, "%s") {
					reqStr = fmt.Sprintf(reqStr, trafficPolicyID)
				}
				req, err := http.NewRequest(
					http.MethodPatch,
					"http://127.0.0.1:8080/traffic_policies",
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
			// Don't repeat all the validation testing we did in POST, just
			// one for ID and another one as a sanity check.
			Entry(
				"PATCH /nodes without id",
				`
				[{}]`,
				"Validation failed: id not a valid uuid"),
			Entry(
				"PATCH /traffic_policies without rules[0].description",
				`
				[{
					"id": "%s",
					"rules": [{}]
				}]`,
				"Validation failed: rules[0].description cannot be empty"),
		)
	})

	Describe("DELETE /traffic_policies/{id}", func() {
		var (
			trafficPolicyID string
		)

		BeforeEach(func() {
			trafficPolicyID = postTrafficPolicies()
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a DELETE /traffic_policies/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("http://127.0.0.1:8080/traffic_policies/%s",
						trafficPolicyID),
					nil)
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Verifying the traffic policy was deleted")

				By("Sending a GET /traffic_policies/{id} request")
				resp, err = http.Get(
					fmt.Sprintf("http://127.0.0.1:8080/traffic_policies/%s",
						trafficPolicyID))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("DELETE /traffic_policies/{id}"),
		)

		DescribeTable("404 Not Found",
			func(id string) {
				By("Sending a DELETE /traffic_policies/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("http://127.0.0.1:8080/traffic_policies/%s",
						id),
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
				"DELETE /traffic_policies/{id} with nonexistent ID",
				uuid.New()),
		)
	})
})
