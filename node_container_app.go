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

// NodeContainerApp represents an association between a Node and a ContainerApp.
type NodeContainerApp struct {
	ID             string `json:"id"`
	NodeID         string `json:"node_id"`
	ContainerAppID string `json:"container_app_id"`
}

// GetTableName returns the name of the persistence table.
func (*NodeContainerApp) GetTableName() string {
	return "nodes_container_apps"
}

// GetID gets the ID.
func (n_ca *NodeContainerApp) GetID() string {
	return n_ca.ID
}

// SetID sets the ID.
func (n_ca *NodeContainerApp) SetID(id string) {
	n_ca.ID = id
}

// GetNodeID gets the node ID.
func (n_ca *NodeContainerApp) GetNodeID() string {
	return n_ca.NodeID
}

// Validate validates the model.
func (n_ca *NodeContainerApp) Validate() error {
	if !uuid.IsValid(n_ca.ID) {
		return errors.New("id not a valid uuid")
	}
	if !uuid.IsValid(n_ca.NodeID) {
		return errors.New("node_id not a valid uuid")
	}
	if !uuid.IsValid(n_ca.ContainerAppID) {
		return errors.New("container_app_id not a valid uuid")
	}

	return nil
}

func (n_ca *NodeContainerApp) String() string {
	return fmt.Sprintf(strings.TrimSpace(`
NodeContainerApp[
    ID: %s
    NodeID: %s
    ContainerAppID: %s
]`),
		n_ca.ID,
		n_ca.NodeID,
		n_ca.ContainerAppID)
}
