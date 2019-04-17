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

var _ = Describe("VNF Lifecycle Service", func() {
	var (
		vnfID string
	)

	BeforeEach(func() {
		var err error

		By("Generating new IDs")
		vnfID = uuid.New()

		By("Deploying a VNF")
		err = vnfDeploySvcCli.Deploy(
			ctx,
			&pb.VNF{
				Id:          vnfID,
				Name:        "test_vnf",
				Vendor:      "test_vendor",
				Description: "test vnf",
				Image:       "http://test.com/vnf123",
				Cores:       4,
				Memory:      4096,
			})
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("Start", func() {
		Describe("Success", func() {
			It("Should start VNFs", func() {
				By("Starting the first VNF")
				err := vnfLifeSvcCli.Start(ctx, vnfID)

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the first VNF is started")
				status, err := vnfDeploySvcCli.GetStatus(ctx, vnfID)
				Expect(err).ToNot(HaveOccurred())
				Expect(status).To(Equal(
					&pb.LifecycleStatus{
						Status: pb.LifecycleStatus_RUNNING,
					},
				))
			})
		})

		Describe("Errors", func() {
			It("Should return an error if the VNF is already "+
				"running", func() {
				By("Starting the first VNF")
				err := vnfLifeSvcCli.Start(ctx, vnfID)
				Expect(err).ToNot(HaveOccurred())

				By("Attempting to start the first VNF again")
				err = vnfLifeSvcCli.Start(ctx, vnfID)

				By("Verifying a FailedPrecondition response")
				Expect(err).To(HaveOccurred())
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.FailedPrecondition,
						"VNF %s not stopped or ready", vnfID)))
			})
		})
	})

	Describe("Restart", func() {
		Describe("Success", func() {
			It("Should restart VNFs", func() {
				By("Starting the first VNF")
				err := vnfLifeSvcCli.Start(ctx, vnfID)
				Expect(err).ToNot(HaveOccurred())

				By("Restarting the first VNF")
				err = vnfLifeSvcCli.Restart(ctx, vnfID)

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the first VNF is restarted")
				status, err := vnfDeploySvcCli.GetStatus(ctx, vnfID)
				Expect(err).ToNot(HaveOccurred())
				Expect(status).To(Equal(
					&pb.LifecycleStatus{
						Status: pb.LifecycleStatus_RUNNING,
					},
				))
			})
		})

		Describe("Errors", func() {
			It("Should return an error if the VNF is not running", func() {
				By("Attempting to restart the first VNF")
				err := vnfLifeSvcCli.Restart(ctx, vnfID)

				By("Verifying a FailedPrecondition response")
				Expect(err).To(HaveOccurred())
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.FailedPrecondition,
						"VNF %s not running", vnfID)))
			})
		})
	})

	Describe("Stop", func() {
		Describe("Success", func() {
			It("Should stop VNFs", func() {
				By("Starting the first VNF")
				err := vnfLifeSvcCli.Start(ctx, vnfID)
				Expect(err).ToNot(HaveOccurred())

				By("Stopping the first VNF")
				err = vnfLifeSvcCli.Stop(ctx, vnfID)

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the first VNF is stopped")
				status, err := vnfDeploySvcCli.GetStatus(ctx, vnfID)
				Expect(err).ToNot(HaveOccurred())
				Expect(status).To(Equal(
					&pb.LifecycleStatus{
						Status: pb.LifecycleStatus_STOPPED,
					},
				))
			})
		})

		Describe("Errors", func() {
			It("Should return an error if the first VNF is already "+
				"stopped", func() {
				By("Starting the first VNF")
				err := vnfLifeSvcCli.Start(ctx, vnfID)
				Expect(err).ToNot(HaveOccurred())

				By("Stopping the first VNF")
				err = vnfLifeSvcCli.Stop(ctx, vnfID)
				Expect(err).ToNot(HaveOccurred())

				By("Attempting to stop the first VNF again")
				err = vnfLifeSvcCli.Stop(ctx, vnfID)

				By("Verifying a FailedPrecondition response")
				Expect(err).To(HaveOccurred())
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.FailedPrecondition,
						"VNF %s not running", vnfID)))

			})
		})
	})
})
