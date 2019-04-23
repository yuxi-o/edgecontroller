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

// VNFDeploymentServiceClient wraps the PB client.
type VNFDeploymentServiceClient struct {
	PBCli pb.VNFDeploymentServiceClient
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
) error {
	_, err := c.PBCli.Deploy(
		ctx,
		vnf)

	if err != nil {
		return errors.Wrap(err, "error deploying vnf")
	}

	return nil
}

// GetStatus retrieves a VNF's status.
func (c *VNFDeploymentServiceClient) GetStatus(
	ctx context.Context,
	id string,
) (*pb.LifecycleStatus, error) {
	status, err := c.PBCli.GetStatus(
		ctx,
		&pb.VNFID{
			Id: id,
		})

	if err != nil {
		return nil, errors.Wrap(err, "error retrieving vnf")
	}

	return status, nil
}

// Redeploy redeploys a VNF.
func (c *VNFDeploymentServiceClient) Redeploy(
	ctx context.Context,
	vnf *pb.VNF,
) error {
	_, err := c.PBCli.Redeploy(
		ctx,
		vnf)

	if err != nil {
		return errors.Wrap(err, "error redeploying vnf")
	}

	return nil
}

// Undeploy undeploys a VNF.
func (c *VNFDeploymentServiceClient) Undeploy(
	ctx context.Context,
	id string,
) error {
	_, err := c.PBCli.Undeploy(
		ctx,
		&pb.VNFID{
			Id: id,
		})

	if err != nil {
		return errors.Wrap(err, "error removing vnf")
	}

	return nil
}
