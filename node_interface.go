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

	"github.com/open-ness/edgecontroller/uuid"
)

// NodeInterface represents an association between a Node and an Interface.
type NodeInterface struct {
	ID          string `json:"id"`
	NodeID      string `json:"node_id"`
	InterfaceID string `json:"interface_id"`
}

// NodeInterfaceReq is a NodeInterface request.
// TODO add a String() method and test for this struct.
type NodeInterfaceReq struct {
	NodeInterface
	Cmd string `json:"cmd,omitempty"`
}

// NodeInterfaceResp is a NodeInterface response.
// TODO add a String() method and test for this struct.
type NodeInterfaceResp struct {
	NodeInterface
	Status string `json:"status"`
}

// GetTableName returns the name of the persistence table.
func (*NodeInterface) GetTableName() string {
	return "nodes_interfaces"
}

// GetID gets the ID.
func (n_i *NodeInterface) GetID() string {
	return n_i.ID
}

// SetID sets the ID.
func (n_i *NodeInterface) SetID(id string) {
	n_i.ID = id
}

// GetNodeID gets the node ID.
func (n_i *NodeInterface) GetNodeID() string {
	return n_i.NodeID
}

// Validate validates the model.
func (n_i *NodeInterface) Validate() error {
	if !uuid.IsValid(n_i.ID) {
		return errors.New("id not a valid uuid")
	}
	if !uuid.IsValid(n_i.NodeID) {
		return errors.New("node_id not a valid uuid")
	}
	if !uuid.IsValid(n_i.InterfaceID) {
		return errors.New("interface_id not a valid uuid")
	}

	return nil
}

// FilterFields returns the filterable fields for this model.
func (*NodeInterface) FilterFields() []string {
	return []string{
		"node_id",
		"interface_id",
	}
}

// Validate validates the request model.
// TODO add a test for this method.
func (n_ir *NodeInterfaceReq) Validate() error {
	if err := n_ir.NodeInterface.Validate(); err != nil {
		return err
	}
	switch n_ir.Cmd {
	case "start", "stop", "restart":
		return nil
	case "":
		return errors.New("cmd missing")
	default:
		return fmt.Errorf(`cmd "%s" is invalid`, n_ir.Cmd)
	}
}

// GetTableName returns the name of the persistence table.
func (n_ir *NodeInterfaceReq) GetTableName() string {
	return n_ir.NodeInterface.GetTableName()
}

func (n_i *NodeInterface) String() string {
	return fmt.Sprintf(strings.TrimSpace(`
NodeInterface[
    ID: %s
    NodeID: %s
    InterfaceID: %s
]`),
		n_i.ID,
		n_i.NodeID,
		n_i.InterfaceID)
}
