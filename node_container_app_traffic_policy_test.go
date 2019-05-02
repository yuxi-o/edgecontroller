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

var _ = Describe("Join Entities: NodeContainerAppTrafficPolicy", func() {
	var (
		ncatp *cce.NodeContainerAppTrafficPolicy
	)

	BeforeEach(func() {
		ncatp = &cce.NodeContainerAppTrafficPolicy{
			ID:                 "a2243693-4fcb-4b80-a914-3c3662424abd",
			NodeContainerAppID: "7a41f67a-086a-4ec2-a980-5db97d9c9f4e",
			TrafficPolicyID:    "9d740cee-035f-4076-847c-d1c80cdf19db",
		}
	})

	Describe("GetTableName", func() {
		It(`Should return "nodes_container_apps_traffic_policies"`, func() {
			Expect(ncatp.GetTableName()).To(Equal(
				"nodes_container_apps_traffic_policies"))
		})
	})

	Describe("GetID", func() {
		It("Should return the ID", func() {
			Expect(ncatp.GetID()).To(Equal(
				"a2243693-4fcb-4b80-a914-3c3662424abd"))
		})
	})

	Describe("SetID", func() {
		It("Should set and return the updated ID", func() {
			By("Setting the ID")
			ncatp.SetID("456")

			By("Getting the updated ID")
			Expect(ncatp.ID).To(Equal("456"))
		})
	})

	Describe("Validate", func() {
		It("Should return an error if ID is not a UUID", func() {
			ncatp.ID = "123"
			Expect(ncatp.Validate()).To(MatchError("id not a valid uuid"))
		})

		It("Should return an error if NodeContainerAppID is not a "+
			"UUID", func() {
			ncatp.NodeContainerAppID = "123"
			Expect(ncatp.Validate()).To(MatchError(
				"nodes_container_apps_id not a valid uuid"))
		})

		It("Should return an error if TrafficPolicyID is not a UUID", func() {
			ncatp.TrafficPolicyID = "123"
			Expect(ncatp.Validate()).To(MatchError(
				"traffic_policy_id not a valid uuid"))
		})
	})

	Describe("String", func() {
		It("Should return the string value", func() {
			Expect(ncatp.String()).To(Equal(strings.TrimSpace(`
NodeContainerAppTrafficPolicy[
    ID: a2243693-4fcb-4b80-a914-3c3662424abd
    NodeContainerAppID: 7a41f67a-086a-4ec2-a980-5db97d9c9f4e
    TrafficPolicyID: 9d740cee-035f-4076-847c-d1c80cdf19db
]`,
			)))
		})
	})
})
