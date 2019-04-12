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

var _ = Describe("Application Deployment Service", func() {
	var (
		containerAppID string
		vmAppID        string
	)

	BeforeEach(func() {
		var err error

		By("Deploying a container application")
		containerAppID, err = appDeploySvcCli.DeployContainer(
			ctx,
			&pb.Application{
				Name:        "test_container_app",
				Vendor:      "test_vendor",
				Description: "test container app",
				Image:       "http://test.com/container123",
				Cores:       4,
				Memory:      4096,
			})
		Expect(err).ToNot(HaveOccurred())

		By("Deploying a VM application")
		vmAppID, err = appDeploySvcCli.DeployVM(
			ctx,
			&pb.Application{
				Name:        "test_vm_app",
				Vendor:      "test_vendor",
				Description: "test vm app",
				Image:       "http://test.com/vm123",
				Cores:       4,
				Memory:      4096,
			})
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("DeployContainer", func() {
		Describe("Success", func() {
			It("Should deploy container applications", func() {
				By("Verifying the response is an ID")
				Expect(containerAppID).ToNot(BeNil())
			})
		})

		Describe("Errors", func() {})
	})

	Describe("DeployVM", func() {
		Describe("Success", func() {
			It("Should deploy VM applications", func() {
				By("Verifying the response is an ID")
				Expect(vmAppID).ToNot(BeNil())
			})
		})

		Describe("Errors", func() {})
	})

	Describe("GetAll", func() {
		Describe("Success", func() {
			It("Should get all deployed applications", func() {
				By("Getting all applications")
				apps, err := appDeploySvcCli.GetAll(ctx)

				By("Verifying the response includes the container and VM " +
					"applications")
				Expect(err).ToNot(HaveOccurred())
				Expect(len(apps.Applications)).To(BeNumerically(">=", 2))
				img := "http://test.com/container123"
				Expect(apps.Applications).To(ContainElement(
					&pb.Application{
						Id:                   containerAppID,
						Name:                 "test_container_app",
						Vendor:               "test_vendor",
						Description:          "test container app",
						Image:                img,
						Cores:                4,
						Memory:               4096,
						Status:               pb.LifecycleStatus_STOPPED,
						XXX_NoUnkeyedLiteral: *new(struct{}),
						XXX_unrecognized:     nil,
						XXX_sizecache:        0,
					},
				))
				Expect(apps.Applications).To(ContainElement(
					&pb.Application{
						Id:                   vmAppID,
						Name:                 "test_vm_app",
						Vendor:               "test_vendor",
						Description:          "test vm app",
						Image:                "http://test.com/vm123",
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
			It("Should get container applications", func() {
				By("Getting the container application")
				containerApp, err := appDeploySvcCli.Get(ctx, containerAppID)

				By("Verifying the response matches the container " +
					"application")
				Expect(err).ToNot(HaveOccurred())
				img := "http://test.com/container123"
				Expect(containerApp).To(Equal(
					&pb.Application{
						Id:                   containerAppID,
						Name:                 "test_container_app",
						Vendor:               "test_vendor",
						Description:          "test container app",
						Image:                img,
						Cores:                4,
						Memory:               4096,
						Status:               pb.LifecycleStatus_STOPPED,
						XXX_NoUnkeyedLiteral: *new(struct{}),
						XXX_unrecognized:     nil,
						XXX_sizecache:        0,
					},
				))
			})

			It("Should get VM applications", func() {
				By("Getting the VM application")
				vmApp, err := appDeploySvcCli.Get(ctx, vmAppID)

				By("Verifying the response matches the VM application")
				Expect(err).ToNot(HaveOccurred())
				Expect(vmApp).To(Equal(
					&pb.Application{
						Id:                   vmAppID,
						Name:                 "test_vm_app",
						Vendor:               "test_vendor",
						Description:          "test vm app",
						Image:                "http://test.com/vm123",
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
			It("Should return an error if the application does not "+
				"exist", func() {
				By("Passing a nonexistent ID")
				badID := uuid.NewV4().String()
				noApp, err := appDeploySvcCli.Get(ctx, badID)

				By("Verifying a NotFound response")
				Expect(err).To(HaveOccurred(),
					"Expected error but got app: %v", noApp)
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.NotFound,
						"Application %s not found", badID)))
			})
		})
	})

	Describe("Redeploy", func() {
		Describe("Success", func() {
			It("Should redeploy container applications", func() {
				By("Getting the container application")
				containerApp, err := appDeploySvcCli.Get(ctx, containerAppID)
				Expect(err).ToNot(HaveOccurred())

				By("Updating the container application")
				containerApp.Cores = 8
				containerApp.Memory = 8192

				By("Redeploying the updated container application")
				err = appDeploySvcCli.Redeploy(ctx, containerApp)

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())

				By("Getting the redeployed container application")
				containerApp, err = appDeploySvcCli.Get(ctx, containerAppID)

				By("Verifying the response matches the updated container " +
					"application")
				Expect(err).ToNot(HaveOccurred())
				Expect(containerApp.Cores).To(Equal(int32(8)))
				Expect(containerApp.Memory).To(Equal(int32(8192)))
			})

			It("Should redeploy VM applications", func() {
				By("Getting the VM application")
				vmApp, err := appDeploySvcCli.Get(ctx, vmAppID)
				Expect(err).ToNot(HaveOccurred())

				By("Updating the VM application")
				vmApp.Cores = 8
				vmApp.Memory = 8192

				By("Redeploying the updated VM application")
				err = appDeploySvcCli.Redeploy(ctx, vmApp)

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())

				By("Getting the redeployed VM application")
				vmApp, err = appDeploySvcCli.Get(ctx, vmAppID)

				By("Verifying the response matches the updated VM " +
					"application")
				Expect(err).ToNot(HaveOccurred())
				Expect(vmApp.Cores).To(Equal(int32(8)))
				Expect(vmApp.Memory).To(Equal(int32(8192)))
			})
		})

		Describe("Errors", func() {
			It("Should return an error if the ID does not exist", func() {
				By("Passing a nonexistent ID")
				badID := uuid.NewV4().String()
				err := appDeploySvcCli.Redeploy(ctx, &pb.Application{Id: badID})

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
				err := appDeploySvcCli.Remove(ctx, containerAppID)

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the container application was removed")
				_, err = appDeploySvcCli.Get(ctx, containerAppID)
				Expect(err).To(HaveOccurred())
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.NotFound,
						"Application %s not found", containerAppID)))
			})

			It("Should remove VM applications", func() {
				By("Removing the VM application")
				err := appDeploySvcCli.Remove(ctx, vmAppID)

				By("Verifying a success response")
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the VM application was removed")
				_, err = appDeploySvcCli.Get(ctx, vmAppID)
				Expect(err).To(HaveOccurred())
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.NotFound,
						"Application %s not found", vmAppID)))
			})
		})

		Describe("Errors", func() {
			It("Should return an error if the ID does not exist", func() {
				By("Passing a nonexistent ID")
				badID := uuid.NewV4().String()
				err := appDeploySvcCli.Remove(ctx, badID)

				By("Verifying a NotFound response")
				Expect(err).To(HaveOccurred())
				Expect(errors.Cause(err)).To(Equal(
					status.Errorf(codes.NotFound,
						"Application %s not found", badID)))
			})
		})
	})
})
