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
	"github.com/pkg/errors"
	cce "github.com/smartedgemec/controller-ce"
	"github.com/smartedgemec/controller-ce/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Network Interface Policy Service Client", func() {
	Describe("Set", func() {
		Describe("Success", func() {
			It("Should set the traffic policy", func() {
				By("Updating if2's traffic policy")
				err := interfacePolicySvcCli.Set(
					ctx,
					"if2",
					&cce.TrafficPolicy{
						ID: uuid.New(),
						Rules: []*cce.TrafficRule{
							{
								Description: "updated_rule",
								Priority:    0,
								Source: &cce.TrafficSelector{
									Description: "updated_source",
									MACs: &cce.MACFilter{
										MACAddresses: []string{
											"updated_source_mac_0",
											"updated_source_mac_1",
										},
									},
								},
								Destination: &cce.TrafficSelector{
									Description: "updated_destination",
									MACs: &cce.MACFilter{
										MACAddresses: []string{
											"updated_dest_mac_0",
											"updated_dest_mac_1",
										},
									},
								},
								Target: &cce.TrafficTarget{
									Description: "updated_target",
									Action:      "accept",
									MAC: &cce.MACModifier{
										MACAddress: "updated_target_mac",
									},
									IP: &cce.IPModifier{
										Address: "127.0.0.1",
										Port:    9999,
									},
								},
							},
						},
					},
				)

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Describe("Errors", func() {
			It("Should return an error if the ID does not exist", func() {
				By("Passing a nonexistent ID")
				badID := uuid.New()
				err := interfacePolicySvcCli.Set(ctx, badID, &cce.TrafficPolicy{
					ID: uuid.New(),
				})

				By("Verifying a NotFound response")
				Expect(err).To(HaveOccurred())
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.NotFound,
						"Network Interface %s not found", badID)))
			})
		})
	})

	Describe("Delete", func() {
		Describe("Success", func() {
			It("Should delete the traffic policy", func() {
				By("Deleting the traffic policy")
				err := interfacePolicySvcCli.Delete(ctx, "if2")

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Describe("Errors", func() {
			It("Should return an error if the interface ID does not exist", func() {
				By("Passing a nonexistent interface ID")
				badID := uuid.New()
				err := interfacePolicySvcCli.Delete(ctx, badID)

				By("Verifying a NotFound response")
				Expect(err).To(HaveOccurred())
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.NotFound,
						"Network Interface %s not found", badID)))
			})
		})
	})
})
