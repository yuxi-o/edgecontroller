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

var _ = Describe("Entities: Node", func() {
	var (
		node *cce.Node
	)

	BeforeEach(func() {
		node = &cce.Node{
			ID:       "48606c73-3905-47e0-864f-14bc7466f5bb",
			Name:     "test-node",
			Location: "test-location",
			Serial:   "test-serial",
		}
	})

	Describe("GetTableName", func() {
		It(`Should return "nodes"`, func() {
			Expect(node.GetTableName()).To(Equal("nodes"))
		})
	})

	Describe("GetID", func() {
		It("Should return the ID", func() {
			Expect(node.GetID()).To(Equal(
				"48606c73-3905-47e0-864f-14bc7466f5bb"))
		})
	})

	Describe("SetID", func() {
		It("Should set and return the updated ID", func() {
			By("Setting the ID")
			node.SetID("456")

			By("Getting the updated ID")
			Expect(node.ID).To(Equal("456"))
		})
	})

	Describe("GetNodeID", func() {
		It("Should return the node ID", func() {
			Expect(node.GetNodeID()).To(Equal(
				"48606c73-3905-47e0-864f-14bc7466f5bb"))
		})
	})

	Describe("Validate", func() {
		It("Should return an error if ID is not a UUID", func() {
			node.ID = "123"
			Expect(node.Validate()).To(MatchError("id not a valid uuid"))
		})

		It("Should return an error if Name is empty", func() {
			node.Name = ""
			Expect(node.Validate()).To(MatchError("name cannot be empty"))
		})

		It("Should return an error if Location is empty", func() {
			node.Location = ""
			Expect(node.Validate()).To(MatchError("location cannot be empty"))
		})

		It("Should return an error if Serial is empty", func() {
			node.Serial = ""
			Expect(node.Validate()).To(MatchError("serial cannot be empty"))
		})
	})

	Describe("String", func() {
		It("Should return the string value", func() {
			Expect(node.String()).To(Equal(strings.TrimSpace(`
Node[
    ID: 48606c73-3905-47e0-864f-14bc7466f5bb
    Name: test-node
    Location: test-location
    Serial: test-serial
]`,
			)))
		})
	})
})
