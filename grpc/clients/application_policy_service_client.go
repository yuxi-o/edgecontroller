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

	cce "github.com/otcshare/edgecontroller"
	"github.com/otcshare/edgecontroller/grpc"
	elapb "github.com/otcshare/edgecontroller/pb/ela"
	"github.com/pkg/errors"
)

// ApplicationPolicyServiceClient wraps the PB client.
type ApplicationPolicyServiceClient struct {
	PBCli elapb.ApplicationPolicyServiceClient
}

// NewApplicationPolicyServiceClient creates a new client.
func NewApplicationPolicyServiceClient(
	conn *grpc.ClientConn,
) *ApplicationPolicyServiceClient {
	return &ApplicationPolicyServiceClient{
		conn.NewApplicationPolicyServiceClient(),
	}
}

// Set sets an app's traffic policy.
func (c *ApplicationPolicyServiceClient) Set(
	ctx context.Context,
	appID string,
	policy *cce.TrafficPolicy,
) error {
	_, err := c.PBCli.Set(
		ctx,
		toPBTrafficPolicy(appID, policy))

	if err != nil {
		return errors.Wrap(err, "error setting application policy")
	}

	return nil
}

// Delete deletes an app's traffic policy. This resets it to the default policy.
func (c *ApplicationPolicyServiceClient) Delete(
	ctx context.Context,
	appID string,
) error {
	if err := c.Set(ctx, appID, nil); err != nil {
		return errors.Wrap(err, "error deleting application policy")
	}

	return nil
}
