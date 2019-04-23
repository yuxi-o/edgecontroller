// Copyright 2019 Smart-Edge.com, Inc. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grpc

import (
	"github.com/smartedgemec/controller-ce/pb"
)

// MockNode provides a mock node gRPC server.
type MockNode struct {
	AppDeploySvc pb.ApplicationDeploymentServiceServer
	AppLifeSvc   pb.ApplicationLifecycleServiceServer
	AppPolicySvc pb.ApplicationPolicyServiceServer
	VNFDeploySvc pb.VNFDeploymentServiceServer
	VNFLifeSvc   pb.VNFLifecycleServiceServer
	InterfaceSvc pb.InterfaceServiceServer
	IfPolicySvc  pb.InterfacePolicyServiceServer
	ZoneSvc      pb.ZoneServiceServer
}

// NewMockNode creates a new MockNode with node services initialized.
// AppDeploySvc and AppLifeSvc are combined into appDeployLifeService;
// VNFDeploySvc and VNFLifeSvc are combined into vnfDeployLifeService.
func NewMockNode() *MockNode {
	var (
		appDeployLifeSvc = newAppDeployLifeService()
		appPolicySvc     = newApplicationPolicyService(appDeployLifeSvc)
		vnfDeployLifeSvc = &vnfDeployLifeService{}
		interfaceSvc     = newInterfaceService()
		ifPolicySvc      = newInterfacePolicyService(interfaceSvc)
		zoneSvc          = &zoneService{}
	)

	appDeployLifeSvc.policyService = appPolicySvc

	return &MockNode{
		AppDeploySvc: appDeployLifeSvc,
		AppLifeSvc:   appDeployLifeSvc,
		AppPolicySvc: appPolicySvc,
		VNFDeploySvc: vnfDeployLifeSvc,
		VNFLifeSvc:   vnfDeployLifeSvc,
		InterfaceSvc: interfaceSvc,
		IfPolicySvc:  ifPolicySvc,
		ZoneSvc:      zoneSvc,
	}
}
