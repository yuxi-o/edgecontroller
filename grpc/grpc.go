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
	"crypto/tls"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	logger "github.com/otcshare/common/log"
	elapb "github.com/otcshare/edgecontroller/pb/ela"
	evapb "github.com/otcshare/edgecontroller/pb/eva"
)

var log = logger.DefaultLogger.WithField("pkg", "grpc")

// ClientConn wraps grpc.ClientConn
type ClientConn struct {
	conn *grpc.ClientConn
}

// Dial dials the remote server.
func Dial(ctx context.Context, target string, conf *tls.Config, opts ...grpc.DialOption) (*ClientConn, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if conf != nil {
		opts = append(opts, grpc.WithTransportCredentials(
			credentials.NewTLS(conf)))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	conn, err := grpc.DialContext(ctx, target, opts...)
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
func (c *ClientConn) NewApplicationDeploymentServiceClient() evapb.ApplicationDeploymentServiceClient {
	return evapb.NewApplicationDeploymentServiceClient(c.conn)
}

// NewApplicationLifecycleServiceClient wraps the pb function.
func (c *ClientConn) NewApplicationLifecycleServiceClient() evapb.ApplicationLifecycleServiceClient {
	return evapb.NewApplicationLifecycleServiceClient(c.conn)
}

// NewApplicationPolicyServiceClient wraps the pb function.
func (c *ClientConn) NewApplicationPolicyServiceClient() elapb.ApplicationPolicyServiceClient {
	return elapb.NewApplicationPolicyServiceClient(c.conn)
}

// NewInterfaceServiceClient wraps the pb function.
func (c *ClientConn) NewInterfaceServiceClient() elapb.InterfaceServiceClient {
	return elapb.NewInterfaceServiceClient(c.conn)
}

// NewInterfacePolicyServiceClient wraps the pb function.
func (c *ClientConn) NewInterfacePolicyServiceClient() elapb.InterfacePolicyServiceClient {
	return elapb.NewInterfacePolicyServiceClient(c.conn)
}

// NewZoneServiceClient wraps the pb function.
func (c *ClientConn) NewZoneServiceClient() elapb.ZoneServiceClient {
	return elapb.NewZoneServiceClient(c.conn)
}

// NewDNSServiceClient wraps the pb function.
func (c *ClientConn) NewDNSServiceClient() elapb.DNSServiceClient {
	return elapb.NewDNSServiceClient(c.conn)
}
