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

var _ = Describe("Entities: DNSAppAlias", func() {
	var (
		alias *cce.DNSAppAlias
	)

	BeforeEach(func() {
		alias = &cce.DNSAppAlias{
			ID:          "a48145cc-87de-4aa9-814d-51d23a47eccd",
			Name:        "patient-checkin.choc.org",
			Description: "Patient Check-in Dashboard",
			AppID:       "efcece3c-6b58-4993-8d45-bde6239d4baa",
		}
	})

	Describe("GetTableName", func() {
		It(`Should return "dns_app_aliases"`, func() {
			Expect(alias.GetTableName()).To(Equal("dns_app_aliases"))
		})
	})

	Describe("GetID", func() {
		It("Should return the ID", func() {
			Expect(alias.GetID()).To(Equal(
				"a48145cc-87de-4aa9-814d-51d23a47eccd"))
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

		It("Should return an error if AppID is not a UUID", func() {
			alias.AppID = "123"
			Expect(alias.Validate()).To(MatchError(
				"app_id not a valid uuid"))
		})
	})

	Describe("String", func() {
		It("Should return the string value", func() {
			Expect(alias.String()).To(Equal(strings.TrimSpace(`
DNSAppAlias[
    ID: a48145cc-87de-4aa9-814d-51d23a47eccd
    Name: patient-checkin.choc.org
    Description: Patient Check-in Dashboard
    AppID: efcece3c-6b58-4993-8d45-bde6239d4baa
]`,
			)))
		})
	})
})
