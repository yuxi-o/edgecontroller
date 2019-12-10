// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package clients

import (
	"context"
	"encoding/json"

	cce "github.com/otcshare/edgecontroller"
	"github.com/otcshare/edgecontroller/grpc"
	evapb "github.com/otcshare/edgecontroller/pb/eva"
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

	tmp, err := json.Marshal(app.EPAFeatures)
	if err != nil {
		return nil
	}

	pb := evapb.Application{
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
		EACJsonBlob: string(tmp),
	}

	return &pb
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
