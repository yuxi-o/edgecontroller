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

var _ = Describe("Join Entities: NodeContainerVNF", func() {
	var (
		nvnf *cce.NodeContainerVNF
	)

	BeforeEach(func() {
		nvnf = &cce.NodeContainerVNF{
			ID:             "f01e0777-1368-4688-8745-30da374cbc4a",
			NodeID:         "48606c73-3905-47e0-864f-14bc7466f5bb",
			ContainerVNFID: "28bbfdb2-dace-421d-a680-9ae893a95d37",
		}
	})

	Describe("GetTableName", func() {
		It(`Should return "nodes_container_vnfs"`, func() {
			Expect(nvnf.GetTableName()).To(Equal("nodes_container_vnfs"))
		})
	})

	Describe("GetID", func() {
		It("Should return the ID", func() {
			Expect(nvnf.GetID()).To(Equal(
				"f01e0777-1368-4688-8745-30da374cbc4a"))
		})
	})

	Describe("SetID", func() {
		It("Should set and return the updated ID", func() {
			By("Setting the ID")
			nvnf.SetID("456")

			By("Getting the updated ID")
			Expect(nvnf.ID).To(Equal("456"))
		})
	})

	Describe("GetNodeID", func() {
		It("Should get the node ID", func() {
			Expect(nvnf.GetNodeID()).To(Equal(
				"48606c73-3905-47e0-864f-14bc7466f5bb"))
		})
	})

	Describe("Validate", func() {
		It("Should return an error if ID is not a UUID", func() {
			nvnf.ID = "123"
			Expect(nvnf.Validate()).To(MatchError("id not a valid uuid"))
		})

		It("Should return an error if NodeID is not a UUID", func() {
			nvnf.NodeID = "123"
			Expect(nvnf.Validate()).To(MatchError("node_id not a valid uuid"))
		})

		It("Should return an error if ContainerVNFID is not a UUID", func() {
			nvnf.ContainerVNFID = "123"
			Expect(nvnf.Validate()).To(MatchError(
				"vnf_id not a valid uuid"))
		})
	})

	Describe("String", func() {
		It("Should return the string value", func() {
			Expect(nvnf.String()).To(Equal(strings.TrimSpace(`
NodeContainerVNF[
    ID: f01e0777-1368-4688-8745-30da374cbc4a
    NodeID: 48606c73-3905-47e0-864f-14bc7466f5bb
    ContainerVNFID: 28bbfdb2-dace-421d-a680-9ae893a95d37
]`,
			)))
		})
	})
})
