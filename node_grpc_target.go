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

package cce

import (
	"fmt"
	"strings"
)

// NodeGRPCTarget is a node's GRPC target.
type NodeGRPCTarget struct {
	ID         string `json:"id"`
	NodeID     string `json:"node_id"`
	GRPCTarget string `json:"grpc_target"`
}

// GetTableName returns the name of the persistence table.
func (*NodeGRPCTarget) GetTableName() string {
	return "node_grpc_targets"
}

// GetID gets the ID.
func (t *NodeGRPCTarget) GetID() string {
	return t.ID
}

// SetID sets the ID.
func (t *NodeGRPCTarget) SetID(id string) {
	t.ID = id
}

// GetNodeID gets the node ID.
func (t *NodeGRPCTarget) GetNodeID() string {
	return t.NodeID
}

func (t *NodeGRPCTarget) String() string {
	return fmt.Sprintf(strings.TrimSpace(`
NodeGRPCTarget[
    ID: %s
    NodeID: %s
    GRPCTarget: %s
]`),
		t.ID,
		t.NodeID,
		t.GRPCTarget)
}
