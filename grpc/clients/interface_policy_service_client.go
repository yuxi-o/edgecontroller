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

	cce "github.com/open-ness/edgecontroller"
	"github.com/open-ness/edgecontroller/grpc"
	elapb "github.com/open-ness/edgecontroller/pb/ela"
	"github.com/pkg/errors"
)

// InterfacePolicyServiceClient wraps the PB client.
type InterfacePolicyServiceClient struct {
	PBCli elapb.InterfacePolicyServiceClient
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
	interfaceID string,
	interfacePolicy *cce.TrafficPolicy,
) error {
	_, err := c.PBCli.Set(
		ctx,
		toPBTrafficPolicy(interfaceID, interfacePolicy))

	if err != nil {
		return errors.Wrap(err, "error setting interface policy")
	}

	return nil
}

// Delete deletes an interface's traffic policy. This resets it to the default policy.
func (c *InterfacePolicyServiceClient) Delete(
	ctx context.Context,
	appID string,
) error {
	if err := c.Set(ctx, appID, nil); err != nil {
		return errors.Wrap(err, "error deleting interface policy")
	}

	return nil
}
