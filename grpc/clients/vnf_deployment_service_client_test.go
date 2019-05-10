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

var _ = Describe("VNF Deployment Service Client", func() {
	var (
		containerVNFID string
		vmVNFID        string
	)

	BeforeEach(func() {
		var err error

		By("Generating new IDs")
		containerVNFID = uuid.New()
		vmVNFID = uuid.New()

		By("Deploying a container VNF")
		err = vnfDeploySvcCli.Deploy(
			ctx,
			&cce.VNF{
				ID:          containerVNFID,
				Type:        "container",
				Name:        "test_container_vnf",
				Vendor:      "test_vendor",
				Description: "test container vnf",
				Image:       "http://test.com/container_vnf_123",
				Cores:       4,
				Memory:      4096,
			})
		Expect(err).ToNot(HaveOccurred())

		By("Deploying a VM VNF")
		err = vnfDeploySvcCli.Deploy(
			ctx,
			&cce.VNF{
				ID:          vmVNFID,
				Type:        "vm",
				Name:        "test_vm_vnf",
				Vendor:      "test_vendor",
				Description: "test vm vnf",
				Image:       "http://test.com/vm_vnf_123",
				Cores:       4,
				Memory:      4096,
			})
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("Deploy", func() {
		Describe("Success", func() {
			It("Should deploy container VNFs", func() {
				By("Verifying the response is an ID")
				Expect(containerVNFID).ToNot(BeNil())
			})

			It("Should deploy VM VNFs", func() {
				By("Verifying the response is an ID")
				Expect(vmVNFID).ToNot(BeNil())
			})
		})

		Describe("Errors", func() {})
	})

	Describe("GetStatus", func() {
		Describe("Success", func() {
			It("Should get container VNF status", func() {
				By("Getting the container VNF's status")
				status, err := vnfDeploySvcCli.GetStatus(ctx, containerVNFID)

				By("Verifying the status is Deployed")
				Expect(err).ToNot(HaveOccurred())
				Expect(status).To(Equal(cce.Deployed))

			})

			It("Should get VM VNF status", func() {
				By("Getting the VM VNF's status")
				status, err := vnfDeploySvcCli.GetStatus(ctx, vmVNFID)

				By("Verifying the status is Deployed")
				Expect(err).ToNot(HaveOccurred())
				Expect(status).To(Equal(cce.Deployed))
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
			It("Should redeploy container VNFs", func() {
				By("Redeploying the container VNF")
				err := vnfDeploySvcCli.Redeploy(
					ctx,
					&cce.VNF{
						ID:          containerVNFID,
						Type:        "container",
						Name:        "test_container_vnf",
						Vendor:      "test_vendor",
						Description: "test container vnf",
						Image:       "http://test.com/container_vnf_123",
						Cores:       8,
						Memory:      8192,
					})

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())

				By("Getting the redeployed VNF's status")
				status, err := vnfDeploySvcCli.GetStatus(ctx, containerVNFID)

				By("Verifying the status is Deployed")
				Expect(err).ToNot(HaveOccurred())
				Expect(status).To(Equal(cce.Deployed))
			})

			It("Should redeploy VM VNFs", func() {
				By("Redeploying the VM VNF")
				err := vnfDeploySvcCli.Redeploy(
					ctx,
					&cce.VNF{
						ID:          vmVNFID,
						Type:        "vm",
						Name:        "test_vm_vnf",
						Vendor:      "test_vendor",
						Description: "test vm vnf",
						Image:       "http://test.com/vm_vnf_123",
						Cores:       8,
						Memory:      8192,
					})

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())

				By("Getting the redeployed VNF's status")
				status, err := vnfDeploySvcCli.GetStatus(ctx, containerVNFID)

				By("Verifying the status is Deployed")
				Expect(err).ToNot(HaveOccurred())
				Expect(status).To(Equal(cce.Deployed))
			})
		})

		Describe("Errors", func() {
			It("Should return an error if the ID does not exist", func() {
				By("Passing a nonexistent ID")
				badID := uuid.New()
				err := vnfDeploySvcCli.Redeploy(ctx, &cce.VNF{
					ID: badID,
				})

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
			It("Should remove container VNFs", func() {
				By("Removing the container VNF")
				err := vnfDeploySvcCli.Undeploy(ctx, containerVNFID)

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the VNF was removed")
				_, err = vnfDeploySvcCli.GetStatus(ctx, containerVNFID)
				Expect(err).To(HaveOccurred())
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.NotFound,
						"VNF %s not found", containerVNFID)))
			})

			It("Should remove VM VNFs", func() {
				By("Removing the VM VNF")
				err := vnfDeploySvcCli.Undeploy(ctx, vmVNFID)

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the VM VNF was removed")
				_, err = vnfDeploySvcCli.GetStatus(ctx, vmVNFID)
				Expect(err).To(HaveOccurred())
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.NotFound,
						"VNF %s not found", vmVNFID)))
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
