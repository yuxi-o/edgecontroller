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

// NodeApp represents an association between a Node and an App.
type NodeApp struct {
	ID     string `json:"id"`
	NodeID string `json:"node_id"`
	AppID  string `json:"app_id"`
}

// GetTableName returns the name of the persistence table.
func (*NodeApp) GetTableName() string {
	return "nodes_apps"
}

// GetID gets the ID.
func (n_a *NodeApp) GetID() string {
	return n_a.ID
}

// SetID sets the ID.
func (n_a *NodeApp) SetID(id string) {
	n_a.ID = id
}

// GetNodeID gets the node ID.
func (n_a *NodeApp) GetNodeID() string {
	return n_a.NodeID
}

// Validate validates the model.
func (n_a *NodeApp) Validate() error {
	if !uuid.IsValid(n_a.ID) {
		return errors.New("id not a valid uuid")
	}
	if !uuid.IsValid(n_a.NodeID) {
		return errors.New("node_id not a valid uuid")
	}
	if !uuid.IsValid(n_a.AppID) {
		return errors.New("app_id not a valid uuid")
	}

	return nil
}

func (n_a *NodeApp) String() string {
	return fmt.Sprintf(strings.TrimSpace(`
NodeApp[
    ID: %s
    NodeID: %s
    AppID: %s
]`),
		n_a.ID,
		n_a.NodeID,
		n_a.AppID)
}
