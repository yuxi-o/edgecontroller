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
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gclients "github.com/otcshare/edgecontroller/grpc/clients"
	ctrlgmock "github.com/otcshare/edgecontroller/mock/controller/grpc"
	nodegmock "github.com/otcshare/edgecontroller/mock/node/grpc"
)

var (
	ctx             = context.Background()
	mockNode        = nodegmock.NewMockNode()
	appDeploySvcCli = &gclients.ApplicationDeploymentServiceClient{
		PBCli: &ctrlgmock.MockPBApplicationDeploymentServiceClient{
			MockNode: mockNode,
		},
	}
	appLifeSvcCli = &gclients.ApplicationLifecycleServiceClient{
		PBCli: &ctrlgmock.MockPBApplicationLifecycleServiceClient{
			MockNode: mockNode,
		},
	}
	appPolicySvcCli = &gclients.ApplicationPolicyServiceClient{
		PBCli: &ctrlgmock.MockPBApplicationPolicyServiceClient{
			MockNode: mockNode,
		},
	}
	interfaceSvcCli = &gclients.InterfaceServiceClient{
		PBCli: &ctrlgmock.MockPBInterfaceServiceClient{
			MockNode: mockNode,
		},
	}
	zoneSvcCli = &gclients.ZoneServiceClient{
		PBCli: &ctrlgmock.MockPBZoneServiceClient{
			MockNode: mockNode,
		},
	}
	interfacePolicySvcCli = &gclients.InterfacePolicyServiceClient{
		PBCli: &ctrlgmock.MockPBInterfacePolicyServiceClient{
			MockNode: mockNode,
		},
	}
	dnsSvcCli = &gclients.DNSServiceClient{
		PBCli: &ctrlgmock.MockPBDNSServiceClient{
			MockNode: mockNode,
		},
	}
)

func TestApplicationClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gRPC Clients Suite")
}
