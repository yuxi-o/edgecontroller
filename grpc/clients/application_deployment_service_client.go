// Copyright 2019 Smart-Edge.com, Inc. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

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

// ApplicationDeploymentServiceClient wraps the PB client.
type ApplicationDeploymentServiceClient struct {
	PBCli pb.ApplicationDeploymentServiceClient
}

// NewApplicationDeploymentServiceClient creates a new client.
func NewApplicationDeploymentServiceClient(
	conn *grpc.ClientConn,
) *ApplicationDeploymentServiceClient {
	return &ApplicationDeploymentServiceClient{
		conn.NewApplicationDeploymentServiceClient(),
	}
}

// Deploy deploys an application. Depending on the type of the application,
// either DeployContainer or DeployVM is called on the gRPC service.
func (c *ApplicationDeploymentServiceClient) Deploy(
	ctx context.Context,
	app *cce.App,
) error {
	var err error
	switch app.Type {
	case "container":
		_, err = c.PBCli.DeployContainer(ctx, toPBApp(app))
	case "vm":
		_, err = c.PBCli.DeployVM(ctx, toPBApp(app))
	}

	if err != nil {
		return errors.Wrap(err, "error deploying application")
	}

	return nil
}

func toPBApp(app *cce.App) *pb.Application {
	return &pb.Application{
		Id:          app.ID,
		Name:        app.Name,
		Vendor:      app.Vendor,
		Description: app.Description,
		Version:     app.Version,
		Cores:       int32(app.Cores),
		Memory:      int32(app.Memory),
		Source: &pb.Application_HttpUri{
			HttpUri: &pb.Application_HTTPSource{
				HttpUri: app.Source,
			},
		},
	}
}

// GetStatus retrieves an application's status.
func (c *ApplicationDeploymentServiceClient) GetStatus(
	ctx context.Context,
	id string,
) (cce.LifecycleStatus, error) {
	pbStatus, err := c.PBCli.GetStatus(
		ctx,
		&pb.ApplicationID{Id: id})

	if err != nil {
		return cce.Unknown, errors.Wrap(err, "error retrieving application")
	}

	return fromPBLifecycleStatus(pbStatus), nil
}

// Redeploy redeploys an application.
func (c *ApplicationDeploymentServiceClient) Redeploy(
	ctx context.Context,
	app *cce.App,
) error {
	_, err := c.PBCli.Redeploy(ctx, toPBApp(app))

	if err != nil {
		return errors.Wrap(err, "error redeploying application")
	}

	return nil
}

// Undeploy undeploys an application.
func (c *ApplicationDeploymentServiceClient) Undeploy(
	ctx context.Context,
	id string,
) error {
	_, err := c.PBCli.Undeploy(
		ctx,
		&pb.ApplicationID{
			Id: id,
		})

	if err != nil {
		return errors.Wrap(err, "error removing application")
	}

	return nil
}
