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

var _ = Describe("Application Deployment Service Client", func() {
	var (
		containerAppID string
		vmAppID        string
	)

	BeforeEach(func() {
		var err error

		By("Generating new IDs")
		containerAppID = uuid.New()
		vmAppID = uuid.New()

		By("Deploying a container application")
		err = appDeploySvcCli.Deploy(
			ctx,
			&cce.App{
				ID:          containerAppID,
				Type:        "container",
				Name:        "test_container_app",
				Vendor:      "test_vendor",
				Description: "test container app",
				Image:       "http://test.com/container_app_123",
				Cores:       4,
				Memory:      4096,
			})
		Expect(err).ToNot(HaveOccurred())

		By("Deploying a VM application")
		err = appDeploySvcCli.Deploy(
			ctx,
			&cce.App{
				ID:          vmAppID,
				Type:        "vm",
				Name:        "test_vm_app",
				Vendor:      "test_vendor",
				Description: "test vm app",
				Image:       "http://test.com/vm_app_123",
				Cores:       4,
				Memory:      4096,
			})
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("Deploy", func() {
		Describe("Success", func() {
			It("Should deploy container applications", func() {
				By("Verifying the response is an ID")
				Expect(containerAppID).ToNot(BeNil())
			})

			It("Should deploy VM applications", func() {
				By("Verifying the response is an ID")
				Expect(vmAppID).ToNot(BeNil())
			})
		})

		Describe("Errors", func() {})
	})

	Describe("GetStatus", func() {
		Describe("Success", func() {
			It("Should get container application status", func() {
				By("Getting the container application's status")
				status, err := appDeploySvcCli.GetStatus(ctx, containerAppID)

				By("Verifying the status is Deployed")
				Expect(err).ToNot(HaveOccurred())
				Expect(status).To(Equal(cce.Deployed))
			})

			It("Should get VM application status", func() {
				By("Getting the VM application's status")
				status, err := appDeploySvcCli.GetStatus(ctx, vmAppID)

				By("Verifying the status is Deployed")
				Expect(err).ToNot(HaveOccurred())
				Expect(status).To(Equal(cce.Deployed))
			})
		})

		Describe("Errors", func() {
			It("Should return an error if the application does not "+
				"exist", func() {
				By("Passing a nonexistent ID")
				badID := uuid.New()
				s, err := appDeploySvcCli.GetStatus(ctx, badID)

				By("Verifying a NotFound response")
				Expect(err).To(HaveOccurred(),
					"Expected error but got status: %v", s)
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.NotFound,
						"Application %s not found", badID)))
			})
		})
	})

	Describe("Redeploy", func() {
		Describe("Success", func() {
			It("Should redeploy container applications", func() {
				By("Redeploying the container application")
				err := appDeploySvcCli.Redeploy(
					ctx,
					&cce.App{
						ID:          containerAppID,
						Type:        "container",
						Name:        "test_container_app",
						Vendor:      "test_vendor",
						Description: "test app",
						Image:       "http://test.com/container_app_123",
						Cores:       8,
						Memory:      8192,
					})

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())

				By("Getting the redeployed application's status")
				status, err := appDeploySvcCli.GetStatus(ctx, containerAppID)

				By("Verifying the status is Deployed")
				Expect(err).ToNot(HaveOccurred())
				Expect(status).To(Equal(cce.Deployed))
			})

			It("Should redeploy VM applications", func() {
				By("Redeploying the VM application")
				err := appDeploySvcCli.Redeploy(
					ctx,
					&cce.App{
						ID:          vmAppID,
						Type:        "vm",
						Name:        "test_vm_app",
						Vendor:      "test_vendor",
						Description: "test vm app",
						Image:       "http://test.com/vm_app_123",
						Cores:       8,
						Memory:      8192,
					})

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())

				By("Getting the redeployed application's status")
				status, err := appDeploySvcCli.GetStatus(ctx, vmAppID)

				By("Verifying the status is Deployed")
				Expect(err).ToNot(HaveOccurred())
				Expect(status).To(Equal(cce.Deployed))
			})
		})

		Describe("Errors", func() {
			It("Should return an error if the ID does not exist", func() {
				By("Passing a nonexistent ID")
				badID := uuid.New()
				err := appDeploySvcCli.Redeploy(
					ctx, &cce.App{
						ID: badID,
					})

				By("Verifying a NotFound response")
				Expect(err).To(HaveOccurred())
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.NotFound,
						"Application %s not found", badID)))
			})
		})
	})

	Describe("Remove", func() {
		Describe("Success", func() {
			It("Should remove container applications", func() {
				By("Removing the container application")
				err := appDeploySvcCli.Undeploy(ctx, containerAppID)

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the container application was removed")
				_, err = appDeploySvcCli.GetStatus(ctx, containerAppID)
				Expect(err).To(HaveOccurred())
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.NotFound,
						"Application %s not found", containerAppID)))
			})

			It("Should remove VM applications", func() {
				By("Removing the VM application")
				err := appDeploySvcCli.Undeploy(ctx, vmAppID)

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the VM application was removed")
				_, err = appDeploySvcCli.GetStatus(ctx, vmAppID)
				Expect(err).To(HaveOccurred())
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.NotFound,
						"Application %s not found", vmAppID)))
			})
		})

		Describe("Errors", func() {
			It("Should return an error if the ID does not exist", func() {
				By("Passing a nonexistent ID")
				badID := uuid.New()
				err := appDeploySvcCli.Undeploy(ctx, badID)

				By("Verifying a NotFound response")
				Expect(err).To(HaveOccurred())
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.NotFound,
						"Application %s not found", badID)))
			})
		})
	})
})
