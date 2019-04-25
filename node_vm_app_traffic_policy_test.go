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

var _ = Describe("Join Entities: NodeVMAppTrafficPolicy", func() {
	var (
		nvmatp *cce.NodeVMAppTrafficPolicy
	)

	BeforeEach(func() {
		nvmatp = &cce.NodeVMAppTrafficPolicy{
			ID:              "8b4fa278-9de1-4b55-8173-47d2b55c24df",
			NodeVMAppID:     "a77c4642-c6b3-4554-b793-83103a5517df",
			TrafficPolicyID: "9d740cee-035f-4076-847c-d1c80cdf19db",
		}
	})

	Describe("GetTableName", func() {
		It(`Should return "nodes_vm_apps_traffic_policies"`, func() {
			Expect(nvmatp.GetTableName()).To(Equal(
				"nodes_vm_apps_traffic_policies"))
		})
	})

	Describe("GetID", func() {
		It("Should return the ID", func() {
			Expect(nvmatp.GetID()).To(Equal(
				"8b4fa278-9de1-4b55-8173-47d2b55c24df"))
		})
	})

	Describe("SetID", func() {
		It("Should set and return the updated ID", func() {
			By("Setting the ID")
			nvmatp.SetID("456")

			By("Getting the updated ID")
			Expect(nvmatp.ID).To(Equal("456"))
		})
	})

	Describe("Validate", func() {
		It("Should return an error if ID is not a UUID", func() {
			nvmatp.ID = "123"
			Expect(nvmatp.Validate()).To(MatchError("id not a valid uuid"))
		})

		It("Should return an error if NodeVMAppID is not a UUID", func() {
			nvmatp.NodeVMAppID = "123"
			Expect(nvmatp.Validate()).To(MatchError(
				"node_vm_app_id not a valid uuid"))
		})

		It("Should return an error if TrafficPolicyID is not a UUID", func() {
			nvmatp.TrafficPolicyID = "123"
			Expect(nvmatp.Validate()).To(MatchError(
				"traffic_policy_id not a valid uuid"))
		})
	})

	Describe("String", func() {
		It("Should return the string value", func() {
			Expect(nvmatp.String()).To(Equal(strings.TrimSpace(`
NodeVMAppTrafficPolicy[
    ID: 8b4fa278-9de1-4b55-8173-47d2b55c24df
    NodeVMAppID: a77c4642-c6b3-4554-b793-83103a5517df
    TrafficPolicyID: 9d740cee-035f-4076-847c-d1c80cdf19db
]`,
			)))
		})
	})
})
