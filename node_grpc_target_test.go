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

var _ = Describe("Entities: NodeGRPCTarget", func() {
	var (
		target *cce.NodeGRPCTarget
	)

	BeforeEach(func() {
		target = &cce.NodeGRPCTarget{
			ID:         "ca0fa495-1020-405b-a78c-9a1884349078",
			NodeID:     "48606c73-3905-47e0-864f-14bc7466f5bb",
			GRPCTarget: "127.0.0.1:8082",
		}
	})

	Describe("GetTableName", func() {
		It(`Should return "node_grpc_targets"`, func() {
			Expect(target.GetTableName()).To(Equal("node_grpc_targets"))
		})
	})

	Describe("GetID", func() {
		It("Should return the ID", func() {
			Expect(target.GetID()).To(Equal(
				"ca0fa495-1020-405b-a78c-9a1884349078"))
		})
	})

	Describe("SetID", func() {
		It("Should set and return the updated ID", func() {
			By("Setting the ID")
			target.SetID("456")

			By("Getting the updated ID")
			Expect(target.ID).To(Equal("456"))
		})
	})

	Describe("GetNodeID", func() {
		It("Should return the target ID", func() {
			Expect(target.GetNodeID()).To(Equal(
				"48606c73-3905-47e0-864f-14bc7466f5bb"))
		})
	})

	Describe("FilterFields", func() {
		It("Should return the filterable fields", func() {
			Expect(target.FilterFields()).To(Equal([]string{
				"node_id",
			}))
		})
	})

	Describe("String", func() {
		It("Should return the string value", func() {
			Expect(target.String()).To(Equal(strings.TrimSpace(`
NodeGRPCTarget[
    ID: ca0fa495-1020-405b-a78c-9a1884349078
    NodeID: 48606c73-3905-47e0-864f-14bc7466f5bb
    GRPCTarget: 127.0.0.1:8082
]`,
			)))
		})
	})
})
