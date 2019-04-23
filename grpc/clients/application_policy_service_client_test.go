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
	"github.com/smartedgemec/controller-ce/pb"
	"github.com/smartedgemec/controller-ce/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Application Policy Service Client", func() {
	var (
		appID string
	)

	BeforeEach(func() {
		var err error

		By("Generating new IDs")
		appID = uuid.New()

		By("Deploying an application")
		err = appDeploySvcCli.DeployContainer(
			ctx,
			&pb.Application{
				Id:          appID,
				Name:        "test_container_app",
				Vendor:      "test_vendor",
				Description: "test container app",
				Image:       "http://test.com/container123",
				Cores:       4,
				Memory:      4096,
			})
		Expect(err).ToNot(HaveOccurred())
		Expect(appID).ToNot(BeNil())
	})

	Describe("Set", func() {
		Describe("Success", func() {
			It("Should set the traffic policy", func() {
				By("Updating the traffic policy")
				err := appPolicySvcCli.Set(
					ctx,
					&pb.TrafficPolicy{
						Id: appID,
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
			})
		})

		Describe("Errors", func() {
			It("Should return an error if the ID does not exist", func() {
				By("Passing a nonexistent ID")
				badID := uuid.New()
				err := appPolicySvcCli.Set(ctx, &pb.TrafficPolicy{Id: badID})

				By("Verifying a NotFound response")
				Expect(err).To(HaveOccurred())
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.NotFound,
						"Application %s not found", badID)))
			})
		})
	})
})
