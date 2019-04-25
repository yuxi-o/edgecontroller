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

var _ = Describe("Entities: DNSVMVNFAlias", func() {
	var (
		alias *cce.DNSVMVNFAlias
	)

	BeforeEach(func() {
		alias = &cce.DNSVMVNFAlias{
			ID:          "5fc828ec-3412-4265-bfd0-e4bc3cc51d8c",
			Name:        "patient-checkin.choc.org",
			Description: "Patient Check-in Dashboard",
			VMVNFID:     "8555211d-f572-45e4-bf54-4606481c84eb",
		}
	})

	Describe("GetTableName", func() {
		It(`Should return "dns_vm_vnf_aliases"`, func() {
			Expect(alias.GetTableName()).To(Equal("dns_vm_vnf_aliases"))
		})
	})

	Describe("GetID", func() {
		It("Should return the ID", func() {
			Expect(alias.GetID()).To(Equal(
				"5fc828ec-3412-4265-bfd0-e4bc3cc51d8c"))
		})
	})

	Describe("SetID", func() {
		It("Should set and return the updated ID", func() {
			By("Setting the ID")
			alias.SetID("456")

			By("Getting the updated ID")
			Expect(alias.ID).To(Equal("456"))
		})
	})

	Describe("Validate", func() {
		It("Should return an error if ID is not a UUID", func() {
			alias.ID = "123"
			Expect(alias.Validate()).To(MatchError("id not a valid uuid"))
		})

		It("Should return an error if Name is empty", func() {
			alias.Name = ""
			Expect(alias.Validate()).To(MatchError("name cannot be empty"))
		})

		It("Should return an error if Description is empty", func() {
			alias.Description = ""
			Expect(alias.Validate()).To(MatchError(
				"description cannot be empty"))
		})

		It("Should return an error if VMVNFID is not a UUID", func() {
			alias.VMVNFID = "123"
			Expect(alias.Validate()).To(MatchError(
				"vm_vnf_id not a valid uuid"))
		})
	})

	Describe("String", func() {
		It("Should return the string value", func() {
			Expect(alias.String()).To(Equal(strings.TrimSpace(`
DNSVMVNFAlias[
    ID: 5fc828ec-3412-4265-bfd0-e4bc3cc51d8c
    Name: patient-checkin.choc.org
    Description: Patient Check-in Dashboard
    VMVNFID: 8555211d-f572-45e4-bf54-4606481c84eb
]`,
			)))
		})
	})
})
