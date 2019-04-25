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

var _ = Describe("Entities: VMApp", func() {
	var (
		app *cce.VMApp
	)

	BeforeEach(func() {
		app = &cce.VMApp{
			ID:          "c53ca266-6678-439c-be4e-f37b49e10c37",
			Name:        "test-vm-app",
			Vendor:      "test-vendor",
			Description: "test-description",
			Image:       "test-image",
			Cores:       4,
			Memory:      1024,
		}
	})

	Describe("GetTableName", func() {
		It(`Should return "vm_apps"`, func() {
			Expect(app.GetTableName()).To(Equal("vm_apps"))
		})
	})

	Describe("GetID", func() {
		It("Should return the ID", func() {
			Expect(app.GetID()).To(Equal(
				"c53ca266-6678-439c-be4e-f37b49e10c37"))
		})
	})

	Describe("SetID", func() {
		It("Should set and return the updated ID", func() {
			By("Setting the ID")
			app.SetID("456")

			By("Getting the updated ID")
			Expect(app.ID).To(Equal("456"))
		})
	})

	Describe("Validate", func() {
		It("Should return an error if ID is not a UUID", func() {
			app.ID = "123"
			Expect(app.Validate()).To(MatchError("id not a valid uuid"))
		})

		It("Should return an error if Name is empty", func() {
			app.Name = ""
			Expect(app.Validate()).To(MatchError("name cannot be empty"))
		})

		It("Should return an error if Vendor is empty", func() {
			app.Vendor = ""
			Expect(app.Validate()).To(MatchError("vendor cannot be empty"))
		})

		It("Should return an error if Image is empty", func() {
			app.Image = ""
			Expect(app.Validate()).To(MatchError("image cannot be empty"))
		})

		It("Should return an error if Cores is < 1", func() {
			app.Cores = 0
			Expect(app.Validate()).To(MatchError("cores must be in [1..8]"))
		})

		It("Should return an error if Cores is > 8", func() {
			app.Cores = 9
			Expect(app.Validate()).To(MatchError("cores must be in [1..8]"))
		})

		It("Should return an error if Memory is < 1", func() {
			app.Memory = 0
			Expect(app.Validate()).To(MatchError(
				"memory must be in [1..16384]"))
		})

		It("Should return an error if Memory is > 16384", func() {
			app.Memory = 16385
			Expect(app.Validate()).To(MatchError(
				"memory must be in [1..16384]"))
		})
	})

	Describe("String", func() {
		It("Should return the string value", func() {
			Expect(app.String()).To(Equal(strings.TrimSpace(`
VMApp[
    ID: c53ca266-6678-439c-be4e-f37b49e10c37
    Name: test-vm-app
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
