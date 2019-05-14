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
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/smartedgemec/controller-ce/pb"
)

// ClientConn wraps grpc.ClientConn
type ClientConn struct {
	conn *grpc.ClientConn
}

// Dial dials the remote server.
func Dial(ctx context.Context, target string) (*ClientConn, error) {
	timeoutCtx, cancel := context.WithTimeout(
		ctx, 2*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		timeoutCtx,
		target,
		grpc.WithInsecure(),
		grpc.WithBlock())
	if err != nil {
		return nil, errors.Wrapf(err, "dial %s failed", target)
	}

	return &ClientConn{conn}, nil
}

// Close wraps grpc.Close()
func (c *ClientConn) Close() error {
	return c.conn.Close()
}

// NewApplicationDeploymentServiceClient wraps the pb function.
func (c *ClientConn) NewApplicationDeploymentServiceClient() pb.ApplicationDeploymentServiceClient { //nolint:lll
	return pb.NewApplicationDeploymentServiceClient(c.conn)
}

// NewApplicationLifecycleServiceClient wraps the pb function.
func (c *ClientConn) NewApplicationLifecycleServiceClient() pb.ApplicationLifecycleServiceClient { //nolint:lll
	return pb.NewApplicationLifecycleServiceClient(c.conn)
}

// NewApplicationPolicyServiceClient wraps the pb function.
func (c *ClientConn) NewApplicationPolicyServiceClient() pb.ApplicationPolicyServiceClient { //nolint:lll
	return pb.NewApplicationPolicyServiceClient(c.conn)
}

// NewVNFDeploymentServiceClient wraps the pb function.
func (c *ClientConn) NewVNFDeploymentServiceClient() pb.VNFDeploymentServiceClient { //nolint:lll
	return pb.NewVNFDeploymentServiceClient(c.conn)
}

// NewVNFLifecycleServiceClient wraps the pb function.
func (c *ClientConn) NewVNFLifecycleServiceClient() pb.VNFLifecycleServiceClient { //nolint:lll
	return pb.NewVNFLifecycleServiceClient(c.conn)
}

// NewInterfaceServiceClient wraps the pb function.
func (c *ClientConn) NewInterfaceServiceClient() pb.InterfaceServiceClient {
	return pb.NewInterfaceServiceClient(c.conn)
}

// NewInterfacePolicyServiceClient wraps the pb function.
func (c *ClientConn) NewInterfacePolicyServiceClient() pb.InterfacePolicyServiceClient { //nolint:lll
	return pb.NewInterfacePolicyServiceClient(c.conn)
}

// NewZoneServiceClient wraps the pb function.
func (c *ClientConn) NewZoneServiceClient() pb.ZoneServiceClient {
	return pb.NewZoneServiceClient(c.conn)
}

// NewDNSServiceClient wraps the pb function.
func (c *ClientConn) NewDNSServiceClient() pb.DNSServiceClient {
	return pb.NewDNSServiceClient(c.conn)
}