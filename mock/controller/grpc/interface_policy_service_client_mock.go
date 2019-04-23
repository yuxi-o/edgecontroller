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
	"github.com/smartedgemec/controller-ce/pb"
	"google.golang.org/grpc"
)

// MockPBInterfacePolicyServiceClient delegates to a MockNode.
type MockPBInterfacePolicyServiceClient struct {
	MockNode *gmock.MockNode
}

// Set delegates to a MockNode.
func (c *MockPBInterfacePolicyServiceClient) Set(
	ctx context.Context,
	in *pb.TrafficPolicy,
	opts ...grpc.CallOption,
) (*empty.Empty, error) {
	return c.MockNode.IfPolicySvc.Set(ctx, in)
}
