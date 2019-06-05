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

package gorilla

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	cce "github.com/smartedgemec/controller-ce"
	"github.com/smartedgemec/controller-ce/grpc/node"
	"github.com/smartedgemec/controller-ce/k8s"
)

func connectNode(
	ctx context.Context,
	ps cce.PersistenceService,
	e cce.NodeEntity,
) (*node.ClientConn, error) {
	targets, err := ps.Filter(
		ctx,
		&cce.NodeGRPCTarget{},
		[]cce.Filter{
			{
				Field: "node_id",
				Value: e.GetNodeID(),
			},
		})
	if err != nil {
		return nil, errors.Wrapf(err, "could not fetch gRPC target from DB")
	}
	// sanity check since we are about to access targets[0]
	if len(targets) != 1 {
		return nil, fmt.Errorf("filter returned %v", targets)
	}

	nodeCC := node.ClientConn{NodeGRPCTarget: targets[0].(*cce.NodeGRPCTarget)}
	if err := nodeCC.Connect(ctx); err != nil {
		log.Noticef("Could not connect to node: %v", err)
		return nil, errors.Wrap(err, "could not connect to node")
	}

	log.Debugf("Connection to node %s established: %s",
		e.GetNodeID(), nodeCC.NodeGRPCTarget.GRPCTarget)
	return &nodeCC, nil
}

func getController(ctx context.Context) *cce.Controller {
	return ctx.Value(contextKey("controller")).(*cce.Controller)
}

func toK8SApp(app *cce.App) k8s.App {
	var ports []*k8s.PortProto
	for _, port := range app.Ports {
		ports = append(ports, &k8s.PortProto{
			Port:     int32(port.Port),
			Protocol: port.Protocol,
		})
	}

	return k8s.App{
		ID:     app.ID,
		Image:  app.ID + ":latest",
		Cores:  app.Cores,
		Memory: app.Memory,
		Ports:  ports,
	}
}
