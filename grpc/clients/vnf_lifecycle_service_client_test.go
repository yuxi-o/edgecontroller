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

var _ = Describe("VNF Lifecycle Service Client", func() {
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

	Describe("Start", func() {
		Describe("Success", func() {
			It("Should start container VNFs", func() {
				By("Starting the container VNF")
				err := vnfLifeSvcCli.Start(ctx, containerVNFID)

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the container VNF is started")
				status, err := vnfDeploySvcCli.GetStatus(ctx, containerVNFID)
				Expect(err).ToNot(HaveOccurred())
				Expect(status).To(Equal(cce.Running))
			})

			It("Should start VM VNFs", func() {
				By("Starting the VM VNF")
				err := vnfLifeSvcCli.Start(ctx, vmVNFID)

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the VM VNF is started")
				status, err := vnfDeploySvcCli.GetStatus(ctx, vmVNFID)
				Expect(err).ToNot(HaveOccurred())
				Expect(status).To(Equal(cce.Running))
			})
		})

		Describe("Errors", func() {
			It("Should return an error if the VNF is already "+
				"running", func() {
				By("Starting the container VNF")
				err := vnfLifeSvcCli.Start(ctx, containerVNFID)
				Expect(err).ToNot(HaveOccurred())

				By("Attempting to start the container VNF again")
				err = vnfLifeSvcCli.Start(ctx, containerVNFID)

				By("Verifying a FailedPrecondition response")
				Expect(err).To(HaveOccurred())
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.FailedPrecondition,
						"VNF %s not stopped or ready", containerVNFID)))
			})
		})
	})

	Describe("Restart", func() {
		Describe("Success", func() {
			It("Should restart container VNFs", func() {
				By("Starting the container VNF")
				err := vnfLifeSvcCli.Start(ctx, containerVNFID)
				Expect(err).ToNot(HaveOccurred())

				By("Restarting the container VNF")
				err = vnfLifeSvcCli.Restart(ctx, containerVNFID)

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the container VNF is restarted")
				status, err := vnfDeploySvcCli.GetStatus(ctx, containerVNFID)
				Expect(err).ToNot(HaveOccurred())
				Expect(status).To(Equal(cce.Running))
			})

			It("Should restart VM VNFs", func() {
				By("Starting the VM VNF")
				err := vnfLifeSvcCli.Start(ctx, vmVNFID)
				Expect(err).ToNot(HaveOccurred())

				By("Restarting the VM VNF")
				err = vnfLifeSvcCli.Restart(ctx, vmVNFID)

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the VM VNF is restarted")
				status, err := vnfDeploySvcCli.GetStatus(ctx, vmVNFID)
				Expect(err).ToNot(HaveOccurred())
				Expect(status).To(Equal(cce.Running))
			})
		})

		Describe("Errors", func() {
			It("Should return an error if the VNF is not running", func() {
				By("Attempting to restart the container VNF")
				err := vnfLifeSvcCli.Restart(ctx, containerVNFID)

				By("Verifying a FailedPrecondition response")
				Expect(err).To(HaveOccurred())
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.FailedPrecondition,
						"VNF %s not running", containerVNFID)))
			})
		})
	})

	Describe("Stop", func() {
		Describe("Success", func() {
			It("Should stop container VNFs", func() {
				By("Starting the container VNF")
				err := vnfLifeSvcCli.Start(ctx, containerVNFID)
				Expect(err).ToNot(HaveOccurred())

				By("Stopping the container VNF")
				err = vnfLifeSvcCli.Stop(ctx, containerVNFID)

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the container VNF is stopped")
				status, err := vnfDeploySvcCli.GetStatus(ctx, containerVNFID)
				Expect(err).ToNot(HaveOccurred())
				Expect(status).To(Equal(cce.Stopped))
			})

			It("Should stop VM VNFs", func() {
				By("Starting the VM VNF")
				err := vnfLifeSvcCli.Start(ctx, vmVNFID)
				Expect(err).ToNot(HaveOccurred())

				By("Stopping the VM VNF")
				err = vnfLifeSvcCli.Stop(ctx, vmVNFID)

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the VM VNF is stopped")
				status, err := vnfDeploySvcCli.GetStatus(ctx, vmVNFID)
				Expect(err).ToNot(HaveOccurred())
				Expect(status).To(Equal(cce.Stopped))
			})
		})

		Describe("Errors", func() {
			It("Should return an error if the VNF is already stopped", func() {
				By("Starting the container VNF")
				err := vnfLifeSvcCli.Start(ctx, containerVNFID)
				Expect(err).ToNot(HaveOccurred())

				By("Stopping the container VNF")
				err = vnfLifeSvcCli.Stop(ctx, containerVNFID)
				Expect(err).ToNot(HaveOccurred())

				By("Attempting to stop the container VNF again")
				err = vnfLifeSvcCli.Stop(ctx, containerVNFID)

				By("Verifying a FailedPrecondition response")
				Expect(err).To(HaveOccurred())
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.FailedPrecondition,
						"VNF %s not running", containerVNFID)))

			})
		})
	})
})
