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

var _ = Describe("Join Entities: NodeContainerApp", func() {
	var (
		nca *cce.NodeContainerApp
	)

	BeforeEach(func() {
		nca = &cce.NodeContainerApp{
			ID:             "7a41f67a-086a-4ec2-a980-5db97d9c9f4e",
			NodeID:         "48606c73-3905-47e0-864f-14bc7466f5bb",
			ContainerAppID: "efcece3c-6b58-4993-8d45-bde6239d4baa",
		}
	})

	Describe("GetTableName", func() {
		It(`Should return "nodes_container_apps"`, func() {
			Expect(nca.GetTableName()).To(Equal("nodes_container_apps"))
		})
	})

	Describe("GetID", func() {
		It("Should return the ID", func() {
			Expect(nca.GetID()).To(Equal(
				"7a41f67a-086a-4ec2-a980-5db97d9c9f4e"))
		})
	})

	Describe("SetID", func() {
		It("Should set and return the updated ID", func() {
			By("Setting the ID")
			nca.SetID("456")

			By("Getting the updated ID")
			Expect(nca.ID).To(Equal("456"))
		})
	})

	Describe("GetNodeID", func() {
		It("Should return the node ID", func() {
			Expect(nca.GetNodeID()).To(Equal(
				"48606c73-3905-47e0-864f-14bc7466f5bb"))
		})
	})

	Describe("Validate", func() {
		It("Should return an error if ID is not a UUID", func() {
			nca.ID = "123"
			Expect(nca.Validate()).To(MatchError("id not a valid uuid"))
		})

		It("Should return an error if NodeID is not a UUID", func() {
			nca.NodeID = "123"
			Expect(nca.Validate()).To(MatchError("node_id not a valid uuid"))
		})

		It("Should return an error if ContainerAppID is not a UUID", func() {
			nca.ContainerAppID = "123"
			Expect(nca.Validate()).To(MatchError(
				"container_app_id not a valid uuid"))
		})
	})

	Describe("String", func() {
		It("Should return the string value", func() {
			Expect(nca.String()).To(Equal(strings.TrimSpace(`
NodeContainerApp[
    ID: 7a41f67a-086a-4ec2-a980-5db97d9c9f4e
    NodeID: 48606c73-3905-47e0-864f-14bc7466f5bb
    ContainerAppID: efcece3c-6b58-4993-8d45-bde6239d4baa
]`,
			)))
		})
	})
})
