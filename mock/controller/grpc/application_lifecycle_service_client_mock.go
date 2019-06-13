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

// MockPBApplicationLifecycleServiceClient delegates to a MockNode.
type MockPBApplicationLifecycleServiceClient struct {
	MockNode *gmock.MockNode
}

// Start delegates to a MockNode.
func (c *MockPBApplicationLifecycleServiceClient) Start(
	ctx context.Context,
	in *evapb.LifecycleCommand,
	opts ...grpc.CallOption,
) (*empty.Empty, error) {
	return c.MockNode.AppLifeSvc.Start(ctx, in)
}

// Stop delegates to a MockNode.
func (c *MockPBApplicationLifecycleServiceClient) Stop(
	ctx context.Context,
	in *evapb.LifecycleCommand,
	opts ...grpc.CallOption,
) (*empty.Empty, error) {
	return c.MockNode.AppLifeSvc.Stop(ctx, in)
}

// Restart delegates to a MockNode.
func (c *MockPBApplicationLifecycleServiceClient) Restart(
	ctx context.Context,
	in *evapb.LifecycleCommand,
	opts ...grpc.CallOption,
) (*empty.Empty, error) {
	return c.MockNode.AppLifeSvc.Restart(ctx, in)
}

// GetStatus delegates to a MockNode.
func (c *MockPBApplicationLifecycleServiceClient) GetStatus(
	ctx context.Context,
	in *evapb.ApplicationID,
	opts ...grpc.CallOption,
) (*elapb.LifecycleStatus, error) {
	return c.MockNode.AppLifeSvc.GetStatus(ctx, in)
}
