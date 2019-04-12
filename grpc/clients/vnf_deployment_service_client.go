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

// VNFDeploymentServiceClient wraps the PB client.
type VNFDeploymentServiceClient struct {
	pbCli pb.VNFDeploymentServiceClient
}

// NewVNFDeploymentServiceClient creates a new client.
func NewVNFDeploymentServiceClient(
	conn *grpc.ClientConn,
) *VNFDeploymentServiceClient {
	return &VNFDeploymentServiceClient{
		conn.NewVNFDeploymentServiceClient(),
	}
}

// Deploy deploys a VNF.
func (c *VNFDeploymentServiceClient) Deploy(
	ctx context.Context,
	vnf *pb.VNF,
) (string, error) {
	pbVNFID, err := c.pbCli.Deploy(
		ctx,
		vnf)

	if err != nil {
		return "", errors.Wrap(err, "error deploying vnf")
	}

	return pbVNFID.Id, nil
}

// GetAll retrieves all VNFs.
func (c *VNFDeploymentServiceClient) GetAll(
	ctx context.Context,
) (*pb.VNFs, error) {
	pbVNFs, err := c.pbCli.GetAll(ctx, &empty.Empty{})

	if err != nil {
		return nil, errors.Wrap(err, "error retrieving all vnfs")
	}

	return pbVNFs, nil
}

// Get retrieves a VNF.
func (c *VNFDeploymentServiceClient) Get(
	ctx context.Context,
	id string,
) (*pb.VNF, error) {
	pbApp, err := c.pbCli.Get(
		ctx,
		&pb.VNFID{
			Id: id,
		})

	if err != nil {
		return nil, errors.Wrap(err, "error retrieving vnf")
	}

	return pbApp, nil
}

// Redeploy redeploys a VNF.
func (c *VNFDeploymentServiceClient) Redeploy(
	ctx context.Context,
	vnf *pb.VNF,
) error {
	_, err := c.pbCli.Redeploy(
		ctx,
		vnf)

	if err != nil {
		return errors.Wrap(err, "error redeploying vnf")
	}

	return nil
}

// Remove removes a VNF.
func (c *VNFDeploymentServiceClient) Remove(
	ctx context.Context,
	id string,
) error {
	_, err := c.pbCli.Remove(
		ctx,
		&pb.VNFID{
			Id: id,
		})

	if err != nil {
		return errors.Wrap(err, "error removing vnf")
	}

	return nil
}
