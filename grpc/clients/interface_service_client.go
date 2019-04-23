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

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/smartedgemec/controller-ce/grpc"
	"github.com/smartedgemec/controller-ce/pb"
)

// InterfaceServiceClient wraps the PB client.
type InterfaceServiceClient struct {
	PBCli pb.InterfaceServiceClient
}

// NewInterfaceServiceClient creates a new client.
func NewInterfaceServiceClient(conn *grpc.ClientConn) *InterfaceServiceClient {
	return &InterfaceServiceClient{
		conn.NewInterfaceServiceClient(),
	}
}

// Update updates a network interface.
func (c *InterfaceServiceClient) Update(
	ctx context.Context,
	ni *pb.NetworkInterface,
) error {
	_, err := c.PBCli.Update(
		ctx,
		ni)

	if err != nil {
		return errors.Wrap(err, "error updating network interface")
	}

	return nil
}

// BulkUpdate updates multiple network interfaces.
func (c *InterfaceServiceClient) BulkUpdate(
	ctx context.Context,
	nis *pb.NetworkInterfaces,
) error {
	_, err := c.PBCli.BulkUpdate(
		ctx,
		nis)

	if err != nil {
		return errors.Wrap(err, "error bulk updating network interfaces")
	}

	return nil
}

// GetAll retrieves all network interfaces.
func (c *InterfaceServiceClient) GetAll(
	ctx context.Context,
) (*pb.NetworkInterfaces, error) {
	nis, err := c.PBCli.GetAll(ctx, &empty.Empty{})

	if err != nil {
		return nil, errors.Wrap(err, "error retrieving all network interfaces")
	}

	return nis, nil
}

// Get retrieves a network interface.
func (c *InterfaceServiceClient) Get(
	ctx context.Context,
	id string,
) (*pb.NetworkInterface, error) {
	ni, err := c.PBCli.Get(
		ctx,
		&pb.InterfaceID{
			Id: id,
		})

	if err != nil {
		return nil, errors.Wrap(err, "error retrieving network interface")
	}

	return ni, nil
}
