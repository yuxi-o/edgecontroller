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

package cce_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cce "github.com/smartedgemec/controller-ce"
)

var _ = Describe("Join Entities: NodeVNFTrafficPolicy", func() {
	var (
		nvtp *cce.NodeVNFTrafficPolicy
	)

	BeforeEach(func() {
		nvtp = &cce.NodeVNFTrafficPolicy{
			ID:              "a2243693-4fcb-4b80-a914-3c3662424abd",
			NodeVNFID:       "f01e0777-1368-4688-8745-30da374cbc4a",
			TrafficPolicyID: "9d740cee-035f-4076-847c-d1c80cdf19db",
		}
	})

	Describe("GetTableName", func() {
		It(`Should return "nodes_vnfs_traffic_policies"`, func() {
			Expect(nvtp.GetTableName()).To(Equal(
				"nodes_vnfs_traffic_policies"))
		})
	})

	Describe("GetID", func() {
		It("Should return the ID", func() {
			Expect(nvtp.GetID()).To(Equal(
				"a2243693-4fcb-4b80-a914-3c3662424abd"))
		})
	})

	Describe("SetID", func() {
		It("Should set and return the updated ID", func() {
			By("Setting the ID")
			nvtp.SetID("456")

			By("Getting the updated ID")
			Expect(nvtp.ID).To(Equal("456"))
		})
	})

	Describe("Validate", func() {
		It("Should return an error if ID is not a UUID", func() {
			nvtp.ID = "123"
			Expect(nvtp.Validate()).To(MatchError("id not a valid uuid"))
		})

		It("Should return an error if NodeVNFID is not a "+
			"UUID", func() {
			nvtp.NodeVNFID = "123"
			Expect(nvtp.Validate()).To(MatchError(
				"nodes_vnfs_id not a valid uuid"))
		})

		It("Should return an error if TrafficPolicyID is not a UUID", func() {
			nvtp.TrafficPolicyID = "123"
			Expect(nvtp.Validate()).To(MatchError(
				"traffic_policy_id not a valid uuid"))
		})
	})

	Describe("String", func() {
		It("Should return the string value", func() {
			Expect(nvtp.String()).To(Equal(strings.TrimSpace(`
NodeVNFTrafficPolicy[
    ID: a2243693-4fcb-4b80-a914-3c3662424abd
    NodeVNFID: f01e0777-1368-4688-8745-30da374cbc4a
    TrafficPolicyID: 9d740cee-035f-4076-847c-d1c80cdf19db
]`,
			)))
		})
	})
})
