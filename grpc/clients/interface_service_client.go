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
	cce "github.com/smartedgemec/controller-ce"
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
	ni *cce.NetworkInterface,
) error {
	_, err := c.PBCli.Update(
		ctx,
		toPBNetworkInterface(ni))

	if err != nil {
		return errors.Wrap(err, "error updating network interface")
	}

	return nil
}

// BulkUpdate updates multiple network interfaces.
func (c *InterfaceServiceClient) BulkUpdate(
	ctx context.Context,
	nis []*cce.NetworkInterface,
) error {
	var pbNIs []*pb.NetworkInterface
	for _, ni := range nis {
		pbNIs = append(pbNIs, toPBNetworkInterface(ni))
	}

	_, err := c.PBCli.BulkUpdate(
		ctx,
		&pb.NetworkInterfaces{
			NetworkInterfaces: pbNIs,
		})

	if err != nil {
		return errors.Wrap(err, "error bulk updating network interfaces")
	}

	return nil
}

// GetAll retrieves all network interfaces.
func (c *InterfaceServiceClient) GetAll(
	ctx context.Context,
) ([]*cce.NetworkInterface, error) {
	pbNIs, err := c.PBCli.GetAll(ctx, &empty.Empty{})

	if err != nil {
		return nil, errors.Wrap(err, "error retrieving all network interfaces")
	}

	var nis []*cce.NetworkInterface
	for _, pbNI := range pbNIs.NetworkInterfaces {
		nis = append(nis, fromPBNetworkInterface(pbNI))
	}

	return nis, nil
}

// Get retrieves a network interface.
func (c *InterfaceServiceClient) Get(
	ctx context.Context,
	id string,
) (*cce.NetworkInterface, error) {
	pbNI, err := c.PBCli.Get(
		ctx,
		&pb.InterfaceID{
			Id: id,
		})

	if err != nil {
		return nil, errors.Wrap(err, "error retrieving network interface")
	}

	return fromPBNetworkInterface(pbNI), nil
}

func toPBNetworkInterface(ni *cce.NetworkInterface) *pb.NetworkInterface {
	return &pb.NetworkInterface{
		Id:                ni.ID,
		Description:       ni.Description,
		Driver:            toPBInterfaceDriver(ni.Driver),
		Type:              toPBInterfaceType(ni.Type),
		MacAddress:        ni.MACAddress,
		Vlan:              uint32(ni.VLAN),
		Zones:             ni.Zones,
		FallbackInterface: ni.FallbackInterface,
	}
}

func toPBInterfaceDriver(driver string) pb.NetworkInterface_InterfaceDriver {
	switch driver {
	case "kernel":
		return pb.NetworkInterface_KERNEL
	case "userspace":
		return pb.NetworkInterface_USERSPACE
	default:
		return 0 // this should never happen
	}
}

func toPBInterfaceType(ifType string) pb.NetworkInterface_InterfaceType {
	switch ifType {
	case "none":
		return pb.NetworkInterface_NONE
	case "upstream":
		return pb.NetworkInterface_UPSTREAM
	case "downstream":
		return pb.NetworkInterface_DOWNSTREAM
	case "bidirectional":
		return pb.NetworkInterface_BIDIRECTIONAL
	case "breakout":
		return pb.NetworkInterface_BREAKOUT
	default:
		return 0 // this should never happen
	}
}

func fromPBNetworkInterface(pbNI *pb.NetworkInterface) *cce.NetworkInterface {
	return &cce.NetworkInterface{
		ID:                pbNI.Id,
		Description:       pbNI.Description,
		Driver:            fromPBInterfaceDriver(pbNI.Driver),
		Type:              fromPBInterfaceType(pbNI.Type),
		MACAddress:        pbNI.MacAddress,
		VLAN:              int(pbNI.Vlan),
		Zones:             pbNI.Zones,
		FallbackInterface: pbNI.FallbackInterface,
	}
}

func fromPBInterfaceDriver(pbDriver pb.NetworkInterface_InterfaceDriver) string {
	switch pbDriver {
	case pb.NetworkInterface_KERNEL:
		return "kernel"
	case pb.NetworkInterface_USERSPACE:
		return "userspace"
	default:
		return "kernel" // this should never happen
	}
}

func fromPBInterfaceType(pbIfType pb.NetworkInterface_InterfaceType) string {
	switch pbIfType {
	case pb.NetworkInterface_NONE:
		return "none"
	case pb.NetworkInterface_UPSTREAM:
		return "upstream"
	case pb.NetworkInterface_DOWNSTREAM:
		return "downstream"
	case pb.NetworkInterface_BIDIRECTIONAL:
		return "bidirectional"
	case pb.NetworkInterface_BREAKOUT:
		return "breakout"
	default:
		return "none" // this should never happen
	}
}
