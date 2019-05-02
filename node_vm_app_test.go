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

var _ = Describe("Join Entities: NodeVMApp", func() {
	var (
		nvma *cce.NodeVMApp
	)

	BeforeEach(func() {
		nvma = &cce.NodeVMApp{
			ID:      "a77c4642-c6b3-4554-b793-83103a5517df",
			NodeID:  "48606c73-3905-47e0-864f-14bc7466f5bb",
			VMAppID: "c53ca266-6678-439c-be4e-f37b49e10c37",
		}
	})

	Describe("GetTableName", func() {
		It(`Should return "nodes_vm_apps"`, func() {
			Expect(nvma.GetTableName()).To(Equal("nodes_vm_apps"))
		})
	})

	Describe("GetID", func() {
		It("Should return the ID", func() {
			Expect(nvma.GetID()).To(Equal(
				"a77c4642-c6b3-4554-b793-83103a5517df"))
		})
	})

	Describe("SetID", func() {
		It("Should set and return the updated ID", func() {
			By("Setting the ID")
			nvma.SetID("456")

			By("Getting the updated ID")
			Expect(nvma.ID).To(Equal("456"))
		})
	})

	Describe("GetNodeID", func() {
		It("Should return the node ID", func() {
			Expect(nvma.GetNodeID()).To(Equal(
				"48606c73-3905-47e0-864f-14bc7466f5bb"))
		})
	})

	Describe("Validate", func() {
		It("Should return an error if ID is not a UUID", func() {
			nvma.ID = "123"
			Expect(nvma.Validate()).To(MatchError("id not a valid uuid"))
		})

		It("Should return an error if NodeID is not a UUID", func() {
			nvma.NodeID = "123"
			Expect(nvma.Validate()).To(MatchError("node_id not a valid uuid"))
		})

		It("Should return an error if VMAppID is not a UUID", func() {
			nvma.VMAppID = "123"
			Expect(nvma.Validate()).To(MatchError(
				"vm_app_id not a valid uuid"))
		})
	})

	Describe("String", func() {
		It("Should return the string value", func() {
			Expect(nvma.String()).To(Equal(strings.TrimSpace(`
NodeVMApp[
    ID: a77c4642-c6b3-4554-b793-83103a5517df
    NodeID: 48606c73-3905-47e0-864f-14bc7466f5bb
    VMAppID: c53ca266-6678-439c-be4e-f37b49e10c37
]`,
			)))
		})
	})
})
