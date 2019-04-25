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
	"errors"
	"fmt"
	"strings"

	"github.com/smartedgemec/controller-ce/uuid"
)

// NodeContainerVNF represents an association between a Node and a ContainerVNF.
type NodeContainerVNF struct {
	ID             string `json:"id"`
	NodeID         string `json:"node_id"`
	ContainerVNFID string `json:"container_vnf_id"`
}

// GetTableName returns the name of the persistence table.
func (*NodeContainerVNF) GetTableName() string {
	return "nodes_container_vnfs"
}

// GetID gets the ID.
func (n_vnf *NodeContainerVNF) GetID() string {
	return n_vnf.ID
}

// SetID sets the ID.
func (n_vnf *NodeContainerVNF) SetID(id string) {
	n_vnf.ID = id
}

// GetNodeID gets the node ID.
func (n_vnf *NodeContainerVNF) GetNodeID() string {
	return n_vnf.NodeID
}

// Validate validates the model.
func (n_vnf *NodeContainerVNF) Validate() error {
	if !uuid.IsValid(n_vnf.ID) {
		return errors.New("id not a valid uuid")
	}
	if !uuid.IsValid(n_vnf.NodeID) {
		return errors.New("node_id not a valid uuid")
	}
	if !uuid.IsValid(n_vnf.ContainerVNFID) {
		return errors.New("vnf_id not a valid uuid")
	}

	return nil
}

func (n_vnf *NodeContainerVNF) String() string {
	return fmt.Sprintf(strings.TrimSpace(`
NodeContainerVNF[
    ID: %s
    NodeID: %s
    ContainerVNFID: %s
]`),
		n_vnf.ID,
		n_vnf.NodeID,
		n_vnf.ContainerVNFID)
}
