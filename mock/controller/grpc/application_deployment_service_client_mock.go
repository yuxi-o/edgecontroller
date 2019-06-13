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

package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	gmock "github.com/smartedgemec/controller-ce/mock/node/grpc"
	elapb "github.com/smartedgemec/controller-ce/pb/ela"
	evapb "github.com/smartedgemec/controller-ce/pb/eva"
	"google.golang.org/grpc"
)

// MockPBApplicationDeploymentServiceClient delegates to a MockNode.
type MockPBApplicationDeploymentServiceClient struct {
	MockNode *gmock.MockNode
}

// DeployContainer delegates to a MockNode.
func (c *MockPBApplicationDeploymentServiceClient) DeployContainer(
	ctx context.Context,
	in *elapb.Application,
	opts ...grpc.CallOption,
) (*empty.Empty, error) {
	return c.MockNode.AppDeploySvc.DeployContainer(ctx, in)
}

// DeployVM delegates to a MockNode.
func (c *MockPBApplicationDeploymentServiceClient) DeployVM(
	ctx context.Context,
	in *elapb.Application,
	opts ...grpc.CallOption,
) (*empty.Empty, error) {
	return c.MockNode.AppDeploySvc.DeployVM(ctx, in)
}

// Redeploy delegates to a MockNode.
func (c *MockPBApplicationDeploymentServiceClient) Redeploy(
	ctx context.Context,
	in *elapb.Application,
	opts ...grpc.CallOption,
) (*empty.Empty, error) {
	return c.MockNode.AppDeploySvc.Redeploy(ctx, in)
}

// Undeploy delegates to a MockNode.
func (c *MockPBApplicationDeploymentServiceClient) Undeploy(
	ctx context.Context,
	in *evapb.ApplicationID,
	opts ...grpc.CallOption,
) (*empty.Empty, error) {
	return c.MockNode.AppDeploySvc.Undeploy(ctx, in)
}
