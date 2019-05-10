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
	vnf *cce.VNF,
) error {
	_, err := c.PBCli.Deploy(ctx, toPBVNF(vnf))

	if err != nil {
		return errors.Wrap(err, "error deploying vnf")
	}

	return nil
}

func toPBVNF(vnf *cce.VNF) *pb.VNF {
	return &pb.VNF{
		Id:          vnf.ID,
		Name:        vnf.Name,
		Vendor:      vnf.Vendor,
		Description: vnf.Description,
		Image:       vnf.Image,
		Cores:       int32(vnf.Cores),
		Memory:      int32(vnf.Memory),
	}
}

// GetStatus retrieves a VNF's status.
func (c *VNFDeploymentServiceClient) GetStatus(
	ctx context.Context,
	id string,
) (cce.LifecycleStatus, error) {
	pbStatus, err := c.PBCli.GetStatus(
		ctx,
		&pb.VNFID{
			Id: id,
		})

	if err != nil {
		return cce.Unknown, errors.Wrap(err, "error retrieving vnf")
	}

	return fromPBLifecycleStatus(pbStatus), nil
}

// Redeploy redeploys a VNF.
func (c *VNFDeploymentServiceClient) Redeploy(
	ctx context.Context,
	vnf *cce.VNF,
) error {
	_, err := c.PBCli.Redeploy(ctx, toPBVNF(vnf))

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
