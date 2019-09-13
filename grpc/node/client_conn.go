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

package node

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/otcshare/common/proxy/progutil"
	"github.com/otcshare/edgecontroller/grpc"
	gclients "github.com/otcshare/edgecontroller/grpc/clients"
	ggrpc "google.golang.org/grpc"
)

// Our network callback helper
var PrefaceLis *progutil.PrefaceListener

// ClientConn wraps a Node and provides a Connect() method to create wrapped gRPC clients.
type ClientConn struct {
	Addr string
	Port string
	TLS  *tls.Config

	conn *grpc.ClientConn

	AppDeploySvcCli   *gclients.ApplicationDeploymentServiceClient
	AppLifeSvcCli     *gclients.ApplicationLifecycleServiceClient
	AppPolicySvcCli   *gclients.ApplicationPolicyServiceClient
	IfacePolicySvcCli *gclients.InterfacePolicyServiceClient
	IfaceSvcCli       *gclients.InterfaceServiceClient
	DNSSvcCli         *gclients.DNSServiceClient
	ZoneSvcCli        *gclients.ZoneServiceClient
}

// Connect connects to a node via grpc.Dial.
func (cc *ClientConn) Connect(ctx context.Context) error {
	var err error

	endpoint := net.JoinHostPort(cc.Addr, cc.Port)

	if cc.Port == "42102" { // XXX use the actual variable with this!
		cc.conn, err = grpc.Dial(ctx, endpoint, cc.TLS,
			ggrpc.WithDialer(PrefaceLis.DialEva))

		cc.AppDeploySvcCli = gclients.NewApplicationDeploymentServiceClient(cc.conn)
		cc.AppLifeSvcCli = gclients.NewApplicationLifecycleServiceClient(cc.conn)

		return nil
	} else {
		cc.conn, err = grpc.Dial(ctx, endpoint, cc.TLS,
			ggrpc.WithDialer(PrefaceLis.DialEla))

		// ELA
		cc.AppPolicySvcCli = gclients.NewApplicationPolicyServiceClient(cc.conn)
		cc.IfacePolicySvcCli = gclients.NewInterfacePolicyServiceClient(cc.conn)
		cc.DNSSvcCli = gclients.NewDNSServiceClient(cc.conn)
		cc.IfaceSvcCli = gclients.NewInterfaceServiceClient(cc.conn)

		cc.ZoneSvcCli = gclients.NewZoneServiceClient(cc.conn) // XXX unimplemented?
	}
	if err != nil {
		return err
	}

	return nil
}

func (cc *ClientConn) Disconnect() {
	cc.conn.Close()
}
