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

var _ = Describe("Entities: DNSConfigDNSVMAppAlias", func() {
	var (
		cfgAlias *cce.DNSConfigDNSVMAppAlias
	)

	BeforeEach(func() {
		cfgAlias = &cce.DNSConfigDNSVMAppAlias{
			ID:              "c5beb302-f2bf-4cba-ae53-136bb4204624",
			DNSConfigID:     "84c1f7b9-53e7-408e-9223-deab73befc54",
			DNSVMAppAliasID: "9de5767a-97b2-44ba-9d1f-faf44c10ad27",
		}
	})

	Describe("GetTableName", func() {
		It(`Should return "dns_configs_dns_vm_app_aliases"`, func() {
			Expect(cfgAlias.GetTableName()).To(Equal(
				"dns_configs_dns_vm_app_aliases"))
		})
	})

	Describe("GetID", func() {
		It("Should return the ID", func() {
			Expect(cfgAlias.GetID()).To(Equal(
				"c5beb302-f2bf-4cba-ae53-136bb4204624"))
		})
	})

	Describe("SetID", func() {
		It("Should set and return the updated ID", func() {
			By("Setting the ID")
			cfgAlias.SetID("456")

			By("Getting the updated ID")
			Expect(cfgAlias.ID).To(Equal("456"))
		})
	})

	Describe("Validate", func() {
		It("Should return an error if ID is not a UUID", func() {
			cfgAlias.ID = "123"
			Expect(cfgAlias.Validate()).To(MatchError("id not a valid uuid"))
		})
	})

	Describe("String", func() {
		It("Should return the string value", func() {
			Expect(cfgAlias.String()).To(Equal(strings.TrimSpace(`
DNSConfigDNSVMAppAlias[
    ID: c5beb302-f2bf-4cba-ae53-136bb4204624
    DNSConfigID: 84c1f7b9-53e7-408e-9223-deab73befc54
    DNSVMAppAliasID: 9de5767a-97b2-44ba-9d1f-faf44c10ad27
]`,
			)))
		})
	})
})
