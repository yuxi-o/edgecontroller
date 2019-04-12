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

// ApplicationPolicyServiceClient wraps the PB client.
type ApplicationPolicyServiceClient struct {
	pbCli pb.ApplicationPolicyServiceClient
}

// NewApplicationPolicyServiceClient creates a new client.
func NewApplicationPolicyServiceClient(
	conn *grpc.ClientConn,
) *ApplicationPolicyServiceClient {
	return &ApplicationPolicyServiceClient{
		conn.NewApplicationPolicyServiceClient(),
	}
}

// Set sets the traffic policy.
func (c *ApplicationPolicyServiceClient) Set(
	ctx context.Context,
	policy *pb.TrafficPolicy,
) error {
	_, err := c.pbCli.Set(
		ctx,
		policy)

	if err != nil {
		return errors.Wrap(err, "error setting policy")
	}

	return nil
}

// Get gets the traffic policy.
func (c *ApplicationPolicyServiceClient) Get(
	ctx context.Context,
	id string,
) (*pb.TrafficPolicy, error) {
	policy, err := c.pbCli.Get(
		ctx,
		&pb.ApplicationID{Id: id})

	if err != nil {
		return nil, errors.Wrap(err, "error getting policy")
	}

	return policy, nil
}
