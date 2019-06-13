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

// MockPBZoneServiceClient delegates to a MockNode.
type MockPBZoneServiceClient struct {
	MockNode *gmock.MockNode
}

// Create delegates to a MockNode.
func (c *MockPBZoneServiceClient) Create(
	ctx context.Context,
	in *elapb.NetworkZone,
	opts ...grpc.CallOption,
) (*empty.Empty, error) {
	return c.MockNode.ZoneSvc.Create(ctx, in)
}

// Update delegates to a MockNode.
func (c *MockPBZoneServiceClient) Update(
	ctx context.Context,
	in *elapb.NetworkZone,
	opts ...grpc.CallOption,
) (*empty.Empty, error) {
	return c.MockNode.ZoneSvc.Update(ctx, in)
}

// BulkUpdate delegates to a MockNode.
func (c *MockPBZoneServiceClient) BulkUpdate(
	ctx context.Context,
	in *elapb.NetworkZones,
	opts ...grpc.CallOption,
) (*empty.Empty, error) {
	return c.MockNode.ZoneSvc.BulkUpdate(ctx, in)
}

// GetAll delegates to a MockNode.
func (c *MockPBZoneServiceClient) GetAll(
	ctx context.Context,
	in *empty.Empty,
	opts ...grpc.CallOption,
) (*elapb.NetworkZones, error) {
	return c.MockNode.ZoneSvc.GetAll(ctx, in)
}

// Get delegates to a MockNode.
func (c *MockPBZoneServiceClient) Get(
	ctx context.Context,
	in *elapb.ZoneID,
	opts ...grpc.CallOption,
) (*elapb.NetworkZone, error) {
	return c.MockNode.ZoneSvc.Get(ctx, in)
}

// Delete delegates to a MockNode.
func (c *MockPBZoneServiceClient) Delete(
	ctx context.Context,
	in *elapb.ZoneID,
	opts ...grpc.CallOption,
) (*empty.Empty, error) {
	return c.MockNode.ZoneSvc.Delete(ctx, in)
}
