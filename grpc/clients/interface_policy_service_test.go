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
	"github.com/satori/go.uuid"
	"github.com/smartedgemec/controller-ce/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Network Interface Policy Service", func() {
	Describe("Get", func() {
		Describe("Success", func() {
			It("Should get the default policy", func() {
				By("Getting the default policy for the first interface")
				policy, err := interfacePolicySvcCli.Get(ctx, "if0")
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the response")
				Expect(policy).To(Equal(
					&pb.TrafficPolicy{
						Id: policy.Id,
						TrafficRules: []*pb.TrafficRule{
							{
								Description: "default_rule",
								Priority:    0,
								Source: &pb.TrafficSelector{
									Description: "default_source",
									Macs: &pb.MACFilter{
										MacAddresses: []string{
											"default_source_mac_0",
											"default_source_mac_1",
										},
									},
								},
								Destination: &pb.TrafficSelector{
									Description: "default_destination",
									Macs: &pb.MACFilter{
										MacAddresses: []string{
											"default_dest_mac_0",
											"default_dest_mac_1",
										},
									},
								},
								Target: &pb.TrafficTarget{
									Description: "default_target",
									Action:      pb.TrafficTarget_ACCEPT,
									Mac: &pb.MACModifier{
										MacAddress: "default_target_mac",
									},
									Ip: &pb.IPModifier{
										Address: "127.0.0.1",
										Port:    9999,
									},
								},
							},
						},
					},
				))
			})
		})

		Describe("Errors", func() {
			It("Should return an error if the ID does not exist", func() {
				By("Passing a nonexistent ID")
				badID := uuid.NewV4().String()
				_, err := interfacePolicySvcCli.Get(ctx, badID)

				By("Verifying a NotFound response")
				Expect(err).To(HaveOccurred())
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.NotFound,
						"Network Interface %s not found", badID)))
			})
		})
	})

	Describe("Set", func() {
		Describe("Success", func() {
			It("Should set the traffic policy", func() {
				By("Updating if2's traffic policy")
				err := interfacePolicySvcCli.Set(
					ctx,
					&pb.TrafficPolicy{
						Id: "if2",
						TrafficRules: []*pb.TrafficRule{
							{
								Description: "updated_rule",
								Priority:    0,
								Source: &pb.TrafficSelector{
									Description: "updated_source",
									Macs: &pb.MACFilter{
										MacAddresses: []string{
											"updated_source_mac_0",
											"updated_source_mac_1",
										},
									},
								},
								Destination: &pb.TrafficSelector{
									Description: "updated_destination",
									Macs: &pb.MACFilter{
										MacAddresses: []string{
											"updated_dest_mac_0",
											"updated_dest_mac_1",
										},
									},
								},
								Target: &pb.TrafficTarget{
									Description: "updated_target",
									Action:      pb.TrafficTarget_ACCEPT,
									Mac: &pb.MACModifier{
										MacAddress: "updated_target_mac",
									},
									Ip: &pb.IPModifier{
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

				By("Getting the updated policy")
				policy, err := interfacePolicySvcCli.Get(ctx, "if2")
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the response")
				Expect(policy).To(Equal(
					&pb.TrafficPolicy{
						Id: "if2",
						TrafficRules: []*pb.TrafficRule{
							{
								Description: "updated_rule",
								Priority:    0,
								Source: &pb.TrafficSelector{
									Description: "updated_source",
									Macs: &pb.MACFilter{
										MacAddresses: []string{
											"updated_source_mac_0",
											"updated_source_mac_1",
										},
									},
								},
								Destination: &pb.TrafficSelector{
									Description: "updated_destination",
									Macs: &pb.MACFilter{
										MacAddresses: []string{
											"updated_dest_mac_0",
											"updated_dest_mac_1",
										},
									},
								},
								Target: &pb.TrafficTarget{
									Description: "updated_target",
									Action:      pb.TrafficTarget_ACCEPT,
									Mac: &pb.MACModifier{
										MacAddress: "updated_target_mac",
									},
									Ip: &pb.IPModifier{
										Address: "127.0.0.1",
										Port:    9999,
									},
								},
							},
						},
					},
				))
			})
		})

		Describe("Errors", func() {
			It("Should return an error if the ID does not exist", func() {
				By("Passing a nonexistent ID")
				badID := uuid.NewV4().String()
				err := interfacePolicySvcCli.Set(ctx,
					&pb.TrafficPolicy{Id: badID})

				By("Verifying a NotFound response")
				Expect(err).To(HaveOccurred())
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.NotFound,
						"Network Interface %s not found", badID)))
			})
		})
	})
})
