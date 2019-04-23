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

package clients

import (
	"context"

	"github.com/pkg/errors"
	"github.com/smartedgemec/controller-ce/grpc"
	"github.com/smartedgemec/controller-ce/pb"
)

// InterfacePolicyServiceClient wraps the PB client.
type InterfacePolicyServiceClient struct {
	PBCli pb.InterfacePolicyServiceClient
}

// NewInterfacePolicyServiceClient creates a new client.
func NewInterfacePolicyServiceClient(
	conn *grpc.ClientConn,
) *InterfacePolicyServiceClient {
	return &InterfacePolicyServiceClient{
		conn.NewInterfacePolicyServiceClient(),
	}
}

// Set sets the traffic policy.
func (c *InterfacePolicyServiceClient) Set(
	ctx context.Context,
	interfacePolicy *pb.TrafficPolicy,
) error {
	_, err := c.PBCli.Set(
		ctx,
		interfacePolicy)

	if err != nil {
		return errors.Wrap(err, "error setting interface policy")
	}

	return nil
}
