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

// MockPBDNSServiceClient delegates to a MockNode.
type MockPBDNSServiceClient struct {
	MockNode *gmock.MockNode
}

// SetA delegates to a MockNode.
func (c *MockPBDNSServiceClient) SetA(
	ctx context.Context,
	in *elapb.DNSARecordSet,
	opts ...grpc.CallOption,
) (*empty.Empty, error) {
	return c.MockNode.DNSSvc.SetA(ctx, in)
}

// DeleteA delegates to a MockNode.
func (c *MockPBDNSServiceClient) DeleteA(
	ctx context.Context,
	in *elapb.DNSARecordSet,
	opts ...grpc.CallOption,
) (*empty.Empty, error) {
	return c.MockNode.DNSSvc.DeleteA(ctx, in)
}

// SetForwarders delegates to a MockNode.
func (c *MockPBDNSServiceClient) SetForwarders(
	ctx context.Context,
	in *elapb.DNSForwarders,
	opts ...grpc.CallOption,
) (*empty.Empty, error) {
	return c.MockNode.DNSSvc.SetForwarders(ctx, in)
}

// DeleteForwarders delegates to a MockNode.
func (c *MockPBDNSServiceClient) DeleteForwarders(
	ctx context.Context,
	in *elapb.DNSForwarders,
	opts ...grpc.CallOption,
) (*empty.Empty, error) {
	return c.MockNode.DNSSvc.DeleteForwarders(ctx, in)
}
