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

// NodeVMApp represents an association between a Node and a VMApp.
type NodeVMApp struct {
	ID      string `json:"id"`
	NodeID  string `json:"node_id"`
	VMAppID string `json:"vm_app_id"`
}

// GetTableName returns the name of the persistence table.
func (*NodeVMApp) GetTableName() string {
	return "nodes_vm_apps"
}

// GetID gets the ID.
func (n_vma *NodeVMApp) GetID() string {
	return n_vma.ID
}

// SetID sets the ID.
func (n_vma *NodeVMApp) SetID(id string) {
	n_vma.ID = id
}

// GetNodeID gets the node ID.
func (n_vma *NodeVMApp) GetNodeID() string {
	return n_vma.NodeID
}

// Validate validates the model.
func (n_vma *NodeVMApp) Validate() error {
	if !uuid.IsValid(n_vma.ID) {
		return errors.New("id not a valid uuid")
	}
	if !uuid.IsValid(n_vma.NodeID) {
		return errors.New("node_id not a valid uuid")
	}
	if !uuid.IsValid(n_vma.VMAppID) {
		return errors.New("vm_app_id not a valid uuid")
	}

	return nil
}

func (n_vma *NodeVMApp) String() string {
	return fmt.Sprintf(strings.TrimSpace(`
NodeVMApp[
    ID: %s
    NodeID: %s
    VMAppID: %s
]`),
		n_vma.ID,
		n_vma.NodeID,
		n_vma.VMAppID)
}
