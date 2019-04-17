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

var _ = Describe("VNF Deployment Service", func() {
	var (
		vnfID  string
		vnf2ID string
	)

	BeforeEach(func() {
		var err error

		By("Generating new IDs")
		vnfID = uuid.New()
		vnf2ID = uuid.New()

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

		By("Deploying a second VNF")
		err = vnfDeploySvcCli.Deploy(
			ctx,
			&pb.VNF{
				Id:          vnf2ID,
				Name:        "test_vnf_2",
				Vendor:      "test_vendor",
				Description: "test vnf 2",
				Image:       "http://test.com/vnf456",
				Cores:       4,
				Memory:      4096,
			})
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("Deploy", func() {
		Describe("Success", func() {
			It("Should deploy VNFs", func() {
				By("Verifying the responses are IDs")
				Expect(vnfID).ToNot(BeNil())
				Expect(vnf2ID).ToNot(BeNil())
			})
		})
		Describe("Errors", func() {})
	})

	Describe("GetStatus", func() {
		Describe("Success", func() {
			It("Should get VNF status", func() {
				By("Getting the first VNF's status")
				status, err := vnfDeploySvcCli.GetStatus(ctx, vnfID)

				By("Verifying the status is Ready")
				Expect(err).ToNot(HaveOccurred())
				Expect(status).To(Equal(
					&pb.LifecycleStatus{
						Status: pb.LifecycleStatus_READY,
					},
				))

				By("Getting the second VNF's status")
				status2, err := vnfDeploySvcCli.GetStatus(ctx, vnf2ID)

				By("Verifying the status is Ready")
				Expect(err).ToNot(HaveOccurred())
				Expect(status2).To(Equal(
					&pb.LifecycleStatus{
						Status: pb.LifecycleStatus_READY,
					},
				))
			})
		})

		Describe("Errors", func() {
			It("Should return an error if the VNF does not exist", func() {
				By("Passing a nonexistent ID")
				badID := uuid.New()
				s, err := vnfDeploySvcCli.GetStatus(ctx, badID)

				By("Verifying a NotFound response")
				Expect(err).To(HaveOccurred(),
					"Expected error but got app: %v", s)
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.NotFound,
						"VNF %s not found", badID)))
			})
		})
	})

	Describe("Redeploy", func() {
		Describe("Success", func() {
			It("Should redeploy VNFs", func() {
				By("Redeploying the VNF")
				err := vnfDeploySvcCli.Redeploy(
					ctx,
					&pb.VNF{
						Id:          vnfID,
						Name:        "test_vnf",
						Vendor:      "test_vendor",
						Description: "test vnf",
						Image:       "http://test.com/vnf123",
						Cores:       8,
						Memory:      8192,
					})

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())

				By("Getting the redeployed VNF's status")
				status, err := vnfDeploySvcCli.GetStatus(ctx, vnfID)

				By("Verifying the status is Ready")
				Expect(err).ToNot(HaveOccurred())
				Expect(status).To(Equal(
					&pb.LifecycleStatus{
						Status: pb.LifecycleStatus_READY,
					},
				))
			})
		})

		Describe("Errors", func() {
			It("Should return an error if the ID does not exist", func() {
				By("Passing a nonexistent ID")
				badID := uuid.New()
				err := vnfDeploySvcCli.Redeploy(ctx, &pb.VNF{Id: badID})

				By("Verifying a NotFound response")
				Expect(err).To(HaveOccurred())
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.NotFound,
						"VNF %s not found", badID)))
			})
		})
	})

	Describe("Remove", func() {
		Describe("Success", func() {
			It("Should remove VNFs", func() {
				By("Removing the first VNF")
				err := vnfDeploySvcCli.Undeploy(ctx, vnfID)

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the VNF was removed")
				_, err = vnfDeploySvcCli.GetStatus(ctx, vnfID)
				Expect(err).To(HaveOccurred())
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.NotFound,
						"VNF %s not found", vnfID)))
			})
		})

		Describe("Errors", func() {
			It("Should return an error if the ID does not exist", func() {
				By("Passing a nonexistent ID")
				badID := uuid.New()
				err := vnfDeploySvcCli.Undeploy(ctx, badID)

				By("Verifying a NotFound response")
				Expect(err).To(HaveOccurred())
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.NotFound,
						"VNF %s not found", badID)))
			})
		})
	})
})
