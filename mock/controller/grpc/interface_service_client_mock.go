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
	"google.golang.org/grpc"
)

// MockPBInterfaceServiceClient delegates to a MockNode.
type MockPBInterfaceServiceClient struct {
	MockNode *gmock.MockNode
}

// Update delegates to a MockNode.
func (c *MockPBInterfaceServiceClient) Update(
	ctx context.Context,
	in *elapb.NetworkInterface,
	opts ...grpc.CallOption,
) (*empty.Empty, error) {
	return c.MockNode.InterfaceSvc.Update(ctx, in)
}

// BulkUpdate delegates to a MockNode.
func (c *MockPBInterfaceServiceClient) BulkUpdate(
	ctx context.Context,
	in *elapb.NetworkInterfaces,
	opts ...grpc.CallOption,
) (*empty.Empty, error) {
	return c.MockNode.InterfaceSvc.BulkUpdate(ctx, in)
}

// GetAll delegates to a MockNode.
func (c *MockPBInterfaceServiceClient) GetAll(
	ctx context.Context,
	in *empty.Empty,
	opts ...grpc.CallOption,
) (*elapb.NetworkInterfaces, error) {
	return c.MockNode.InterfaceSvc.GetAll(ctx, in)
}

// Get delegates to a MockNode.
func (c *MockPBInterfaceServiceClient) Get(
	ctx context.Context,
	in *elapb.InterfaceID,
	opts ...grpc.CallOption,
) (*elapb.NetworkInterface, error) {
	return c.MockNode.InterfaceSvc.Get(ctx, in)
}
