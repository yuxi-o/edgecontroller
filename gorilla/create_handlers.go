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
	"log"

	"github.com/pkg/errors"
	cce "github.com/smartedgemec/controller-ce"
	"github.com/smartedgemec/controller-ce/grpc/node"
)

func connectNode(
	ctx context.Context,
	ps cce.PersistenceService,
	e cce.Entity,
) (*node.ClientConn, error) {
	n, err := ps.Read(
		ctx,
		e.(cce.JoinEntity).GetNodeID(),
		&cce.Node{})
	if err != nil {
		return nil, errors.Wrapf(err, "could not fetch node from DB")
	}
	nodeCC := node.ClientConn{Node: n.(*cce.Node)}
	if err := nodeCC.Connect(ctx); err != nil {
		log.Printf("Could not connect to node: %v", err)
		return nil, errors.Wrap(err, "could not connect to node: %v")
	}

	log.Printf("Connection to node %v established", nodeCC.Node.ID)
	log.Println(nodeCC.Node)

	return &nodeCC, nil
}

func handleCreateNodesApps(
	ctx context.Context,
	ps cce.PersistenceService,
	e cce.Entity,
) error {
	nodeCC, err := connectNode(ctx, ps, e)

	if err != nil {
		return err
	}

	// TODO add gRPC calls
	log.Println(nodeCC.Node)

	return nil
}

func handleCreateNodesDNSConfigs(
	ctx context.Context,
	ps cce.PersistenceService,
	e cce.Entity,
) error {
	nodeCC, err := connectNode(ctx, ps, e)

	if err != nil {
		return err
	}

	// TODO add gRPC calls
	log.Println(nodeCC.Node)

	return nil
}
