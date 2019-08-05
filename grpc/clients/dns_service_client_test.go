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

package clients_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cce "github.com/otcshare/edgecontroller"
)

var _ = Describe("DNS Service Client", func() {
	Describe("SetA", func() {
		Describe("Success", func() {
			It("Should set A records", func() {
				By("Setting A records")
				Expect(dnsSvcCli.SetA(ctx, &cce.DNSARecord{
					Name:        "patient-checkin.choc.org",
					Description: "Patient Check-in Dashboard",
					IPs: []string{
						"app-70233300-C0B4-4315-BAAF-8C4BAF953611",
					},
				})).To(Succeed())
			})
		})

		Describe("Errors", func() {})
	})

	Describe("DeleteA", func() {
		Describe("Success", func() {
			It("Should delete A records", func() {
				By("Deleting A records")
				Expect(dnsSvcCli.DeleteA(ctx, &cce.DNSARecord{
					Name:        "patient-checkin.choc.org",
					Description: "Patient Check-in Dashboard",
					IPs: []string{
						"app-70233300-C0B4-4315-BAAF-8C4BAF953611",
					},
				})).To(Succeed())
			})
		})

		Describe("Errors", func() {})
	})

	Describe("SetForwarders", func() {
		Describe("Success", func() {
			It("Should set forwarders", func() {
				By("Setting forwarders")
				Expect(dnsSvcCli.SetForwarders(ctx, []*cce.DNSForwarder{
					{
						Name:        "Google DNS #1",
						Description: "Google's DNS servers (primary)",
						IP:          "8.8.8.8",
					},
					{
						Name:        "Cloudflare DNS #1",
						Description: "Cloudflare's DNS servers (backup)",
						IP:          "1.1.1.1",
					},
				})).To(Succeed())
			})
		})

		Describe("Errors", func() {})
	})

	Describe("DeleteForwarders", func() {
		Describe("Success", func() {
			It("Should delete forwarders", func() {
				By("Deleting forwarders")
				Expect(dnsSvcCli.DeleteForwarders(ctx, []*cce.DNSForwarder{
					{
						Name:        "Google DNS #1",
						Description: "Google's DNS servers (primary)",
						IP:          "8.8.8.8",
					},
					{
						Name:        "Cloudflare DNS #1",
						Description: "Cloudflare's DNS servers (backup)",
						IP:          "1.1.1.1",
					},
				})).To(Succeed())
			})
		})

		Describe("Errors", func() {})
	})
})
