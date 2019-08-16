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
	"github.com/open-ness/edgecontroller/grpc"
	elapb "github.com/open-ness/edgecontroller/pb/ela"
	"github.com/pkg/errors"
)

// ZoneServiceClient wraps the PB client.
type ZoneServiceClient struct {
	PBCli elapb.ZoneServiceClient
}

// NewZoneServiceClient creates a new client.
func NewZoneServiceClient(conn *grpc.ClientConn) *ZoneServiceClient {
	return &ZoneServiceClient{
		conn.NewZoneServiceClient(),
	}
}

// Create creates a network zone.
func (c *ZoneServiceClient) Create(
	ctx context.Context,
	zone *elapb.NetworkZone,
) error {
	_, err := c.PBCli.Create(
		ctx,
		zone)

	if err != nil {
		return errors.Wrap(err, "error creating network zone")
	}

	return nil
}

// Update updates a network zone.
func (c *ZoneServiceClient) Update(
	ctx context.Context,
	ni *elapb.NetworkZone,
) error {
	_, err := c.PBCli.Update(
		ctx,
		ni)

	if err != nil {
		return errors.Wrap(err, "error updating network zone")
	}

	return nil
}

// BulkUpdate updates multiple network zones.
func (c *ZoneServiceClient) BulkUpdate(
	ctx context.Context,
	nis *elapb.NetworkZones,
) error {
	_, err := c.PBCli.BulkUpdate(
		ctx,
		nis)

	if err != nil {
		return errors.Wrap(err, "error bulk updating network zones")
	}

	return nil
}

// GetAll retrieves all network zones.
func (c *ZoneServiceClient) GetAll(
	ctx context.Context,
) (*elapb.NetworkZones, error) {
	nis, err := c.PBCli.GetAll(ctx, &empty.Empty{})

	if err != nil {
		return nil, errors.Wrap(err, "error retrieving all network zones")
	}

	return nis, nil
}

// Get retrieves a network zone.
func (c *ZoneServiceClient) Get(
	ctx context.Context,
	id string,
) (*elapb.NetworkZone, error) {
	ni, err := c.PBCli.Get(
		ctx,
		&elapb.ZoneID{
			Id: id,
		})

	if err != nil {
		return nil, errors.Wrap(err, "error retrieving network zone")
	}

	return ni, nil
}

// Delete delets a network zone.
func (c *ZoneServiceClient) Delete(
	ctx context.Context,
	id string,
) error {
	_, err := c.PBCli.Delete(
		ctx,
		&elapb.ZoneID{
			Id: id,
		})

	if err != nil {
		return errors.Wrap(err, "error deleting network zone")
	}

	return nil
}
