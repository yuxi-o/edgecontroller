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
	cce "github.com/smartedgemec/controller-ce"
	"github.com/smartedgemec/controller-ce/grpc"
	elapb "github.com/smartedgemec/controller-ce/pb/ela"
)

// DNSServiceClient wraps the PB client.
type DNSServiceClient struct {
	PBCli elapb.DNSServiceClient
}

// NewDNSServiceClient creates a new client.
func NewDNSServiceClient(conn *grpc.ClientConn) *DNSServiceClient {
	return &DNSServiceClient{
		conn.NewDNSServiceClient(),
	}
}

// SetA sets a DNS A record.
func (c *DNSServiceClient) SetA(
	ctx context.Context,
	record *cce.DNSARecord,
) error {
	_, err := c.PBCli.SetA(
		ctx,
		&elapb.DNSARecordSet{
			Name:   record.Name,
			Values: record.IPs,
		})

	if err != nil {
		return errors.Wrap(err, "error setting A records")
	}

	return nil
}

// DeleteA deletes a DNS A record.
func (c *DNSServiceClient) DeleteA(
	ctx context.Context,
	record *cce.DNSARecord,
) error {
	_, err := c.PBCli.DeleteA(
		ctx,
		&elapb.DNSARecordSet{
			Name:   record.Name,
			Values: record.IPs,
		})

	if err != nil {
		return errors.Wrap(err, "error deleting A records")
	}

	return nil
}

// SetForwarders sets DNS forwarders.
func (c *DNSServiceClient) SetForwarders(
	ctx context.Context,
	forwarders []*cce.DNSForwarder,
) error {
	var ips []string
	for _, forwarder := range forwarders {
		ips = append(ips, forwarder.IP)
	}

	_, err := c.PBCli.SetForwarders(ctx, &elapb.DNSForwarders{
		IpAddresses: ips,
	})

	if err != nil {
		return errors.Wrap(err, "error setting forwarders")
	}

	return nil
}

// DeleteForwarders sets DNS forwarders.
func (c *DNSServiceClient) DeleteForwarders(
	ctx context.Context,
	forwarders []*cce.DNSForwarder,
) error {
	var ips []string
	for _, forwarder := range forwarders {
		ips = append(ips, forwarder.IP)
	}

	_, err := c.PBCli.DeleteForwarders(ctx, &elapb.DNSForwarders{
		IpAddresses: ips,
	})

	if err != nil {
		return errors.Wrap(err, "error deleting forwarders")
	}

	return nil
}
