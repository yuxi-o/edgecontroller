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

var _ = Describe("Entities: VNF", func() {
	var (
		vnf *cce.VNF
	)

	BeforeEach(func() {
		vnf = &cce.VNF{
			ID:          "28bbfdb2-dace-421d-a680-9ae893a95d37",
			Type:        "container",
			Name:        "test-container-vnf",
			Vendor:      "test-vendor",
			Description: "test-description",
			Image:       "test-image",
			Cores:       4,
			Memory:      1024,
		}
	})

	Describe("GetTableName", func() {
		It(`Should return "vnfs"`, func() {
			Expect(vnf.GetTableName()).To(Equal("vnfs"))
		})
	})

	Describe("GetID", func() {
		It("Should return the ID", func() {
			Expect(vnf.GetID()).To(Equal(
				"28bbfdb2-dace-421d-a680-9ae893a95d37"))
		})
	})

	Describe("SetID", func() {
		It("Should set and return the updated ID", func() {
			By("Setting the ID")
			vnf.SetID("456")

			By("Getting the updated ID")
			Expect(vnf.ID).To(Equal("456"))
		})
	})

	Describe("Validate", func() {
		It("Should return an error if ID is not a UUID", func() {
			vnf.ID = "123"
			Expect(vnf.Validate()).To(MatchError("id not a valid uuid"))
		})

		It(`Should return an error if Type is not "container" or "vm"`, func() {
			vnf.Type = "foo"
			Expect(vnf.Validate()).To(MatchError(
				`type must be either "container" or "vm"`))
		})

		It("Should return an error if Name is empty", func() {
			vnf.Name = ""
			Expect(vnf.Validate()).To(MatchError("name cannot be empty"))
		})

		It("Should return an error if Vendor is empty", func() {
			vnf.Vendor = ""
			Expect(vnf.Validate()).To(MatchError("vendor cannot be empty"))
		})

		It("Should return an error if Image is empty", func() {
			vnf.Image = ""
			Expect(vnf.Validate()).To(MatchError("image cannot be empty"))
		})

		It("Should return an error if Cores is < 1", func() {
			vnf.Cores = 0
			Expect(vnf.Validate()).To(MatchError("cores must be in [1..8]"))
		})

		It("Should return an error if Cores is > 8", func() {
			vnf.Cores = 9
			Expect(vnf.Validate()).To(MatchError("cores must be in [1..8]"))
		})

		It("Should return an error if Memory is < 1", func() {
			vnf.Memory = 0
			Expect(vnf.Validate()).To(MatchError(
				"memory must be in [1..16384]"))
		})

		It("Should return an error if Memory is > 16384", func() {
			vnf.Memory = 16385
			Expect(vnf.Validate()).To(MatchError(
				"memory must be in [1..16384]"))
		})
	})

	Describe("String", func() {
		It("Should return the string value", func() {
			Expect(vnf.String()).To(Equal(strings.TrimSpace(`
VNF[
    ID: 28bbfdb2-dace-421d-a680-9ae893a95d37
    Name: test-container-vnf
    Vendor: test-vendor
    Description: test-description
    Image: test-image
    Cores: 4
    Memory: 1024
]`,
			)))
		})
	})
})
