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

var _ = Describe("VNF Deployment Service", func() {
	var (
		vnfID  string
		vnf2ID string
	)

	BeforeEach(func() {
		var err error
		By("Deploying a VNF")
		vnfID, err = vnfDeploySvcCli.Deploy(
			ctx,
			&pb.VNF{
				Name:        "test_vnf",
				Vendor:      "test_vendor",
				Description: "test vnf",
				Image:       "http://test.com/vnf123",
				Cores:       4,
				Memory:      4096,
			})
		Expect(err).ToNot(HaveOccurred())

		By("Deploying a second VNF")
		vnf2ID, err = vnfDeploySvcCli.Deploy(
			ctx,
			&pb.VNF{
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

	Describe("GetAll", func() {
		Describe("Success", func() {
			It("Should get all deployed VNFs", func() {
				By("Getting all VNFs")
				vnfs, err := vnfDeploySvcCli.GetAll(ctx)

				By("Verifying the response includes the deployed VNFs")
				Expect(err).ToNot(HaveOccurred())
				Expect(len(vnfs.Vnfs)).To(BeNumerically(">=", 2))
				Expect(vnfs.Vnfs).To(ContainElement(
					&pb.VNF{
						Id:                   vnfID,
						Name:                 "test_vnf",
						Vendor:               "test_vendor",
						Description:          "test vnf",
						Image:                "http://test.com/vnf123",
						Cores:                4,
						Memory:               4096,
						Status:               pb.LifecycleStatus_STOPPED,
						XXX_NoUnkeyedLiteral: *new(struct{}),
						XXX_unrecognized:     nil,
						XXX_sizecache:        0,
					},
				))
				Expect(vnfs.Vnfs).To(ContainElement(
					&pb.VNF{
						Id:                   vnf2ID,
						Name:                 "test_vnf_2",
						Vendor:               "test_vendor",
						Description:          "test vnf 2",
						Image:                "http://test.com/vnf456",
						Cores:                4,
						Memory:               4096,
						Status:               pb.LifecycleStatus_STOPPED,
						XXX_NoUnkeyedLiteral: *new(struct{}),
						XXX_unrecognized:     nil,
						XXX_sizecache:        0,
					},
				))
			})
		})

		Describe("Errors", func() {})
	})

	Describe("Get", func() {
		Describe("Success", func() {
			It("Should get VNFs", func() {
				By("Getting the first VNF")
				vnf, err := vnfDeploySvcCli.Get(ctx, vnfID)

				By("Verifying the response matches the first VNF")
				Expect(err).ToNot(HaveOccurred())
				Expect(vnf).To(Equal(
					&pb.VNF{
						Id:                   vnfID,
						Name:                 "test_vnf",
						Vendor:               "test_vendor",
						Description:          "test vnf",
						Image:                "http://test.com/vnf123",
						Cores:                4,
						Memory:               4096,
						Status:               pb.LifecycleStatus_STOPPED,
						XXX_NoUnkeyedLiteral: *new(struct{}),
						XXX_unrecognized:     nil,
						XXX_sizecache:        0,
					},
				))

				By("Getting the second VNF")
				vnf2, err := vnfDeploySvcCli.Get(ctx, vnf2ID)

				By("Verifying the response matches the second VNF")
				Expect(err).ToNot(HaveOccurred())
				Expect(vnf2).To(Equal(
					&pb.VNF{
						Id:                   vnf2ID,
						Name:                 "test_vnf_2",
						Vendor:               "test_vendor",
						Description:          "test vnf 2",
						Image:                "http://test.com/vnf456",
						Cores:                4,
						Memory:               4096,
						Status:               pb.LifecycleStatus_STOPPED,
						XXX_NoUnkeyedLiteral: *new(struct{}),
						XXX_unrecognized:     nil,
						XXX_sizecache:        0,
					},
				))
			})
		})

		Describe("Errors", func() {
			It("Should return an error if the VNF does not exist", func() {
				By("Passing a nonexistent ID")
				badID := uuid.NewV4().String()
				noApp, err := vnfDeploySvcCli.Get(ctx, badID)

				By("Verifying a NotFound response")
				Expect(err).To(HaveOccurred(),
					"Expected error but got app: %v", noApp)
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.NotFound,
						"VNF %s not found", badID)))
			})
		})
	})

	Describe("Redeploy", func() {
		Describe("Success", func() {
			It("Should redeploy VNFs", func() {
				By("Getting the VNF")
				vnf, err := vnfDeploySvcCli.Get(ctx, vnfID)
				Expect(err).ToNot(HaveOccurred())

				By("Updating the VNF")
				vnf.Cores = 8
				vnf.Memory = 8192

				By("Redeploying the updated VNF")
				err = vnfDeploySvcCli.Redeploy(ctx, vnf)

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())

				By("Getting the redeployed VNF")
				vnf, err = vnfDeploySvcCli.Get(ctx, vnfID)

				By("Verifying the response matches the updated VNF")
				Expect(err).ToNot(HaveOccurred())
				Expect(vnf.Cores).To(Equal(int32(8)))
				Expect(vnf.Memory).To(Equal(int32(8192)))
			})
		})

		Describe("Errors", func() {
			It("Should return an error if the ID does not exist", func() {
				By("Passing a nonexistent ID")
				badID := uuid.NewV4().String()
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
				err := vnfDeploySvcCli.Remove(ctx, vnfID)

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the VNF was removed")
				_, err = vnfDeploySvcCli.Get(ctx, vnfID)
				Expect(err).To(HaveOccurred())
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.NotFound,
						"VNF %s not found", vnfID)))
			})
		})

		Describe("Errors", func() {
			It("Should return an error if the ID does not exist", func() {
				By("Passing a nonexistent ID")
				badID := uuid.NewV4().String()
				err := vnfDeploySvcCli.Remove(ctx, badID)

				By("Verifying a NotFound response")
				Expect(err).To(HaveOccurred())
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.NotFound,
						"VNF %s not found", badID)))
			})
		})
	})
})
