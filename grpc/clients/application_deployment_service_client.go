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

	cce "github.com/open-ness/edgecontroller"
	"github.com/open-ness/edgecontroller/grpc"
	evapb "github.com/open-ness/edgecontroller/pb/eva"
	"github.com/pkg/errors"
)

// ApplicationDeploymentServiceClient wraps the PB client.
type ApplicationDeploymentServiceClient struct {
	PBCli evapb.ApplicationDeploymentServiceClient
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

func toPBApp(app *cce.App) *evapb.Application {
	var ports []*evapb.PortProto
	for _, pp := range app.Ports {
		// If the protocol is "all", make it empty in the protobuf
		// since there's currently no support for this on the other
		// end
		// TODO: remove this logic when other end supports the "all"
		// protocol value.
		protocol := ""
		switch {
		case pp.Protocol != "all":
			protocol = pp.Protocol
		}
		ports = append(ports, &evapb.PortProto{Port: pp.Port, Protocol: protocol})
	}

	return &evapb.Application{
		Id:          app.ID,
		Name:        app.Name,
		Vendor:      app.Vendor,
		Description: app.Description,
		Version:     app.Version,
		Cores:       int32(app.Cores),
		Memory:      int32(app.Memory),
		Ports:       ports,
		Source: &evapb.Application_HttpUri{
			HttpUri: &evapb.Application_HTTPSource{
				HttpUri: app.Source,
			},
		},
	}
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
		&evapb.ApplicationID{
			Id: id,
		})

	if err != nil {
		return errors.Wrap(err, "error removing application")
	}

	return nil
}
