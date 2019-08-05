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

	"github.com/otcshare/edgecontroller/uuid"
)

// NodeApp represents an association between a Node and an App.
type NodeApp struct {
	ID     string `json:"id"`
	NodeID string `json:"node_id"`
	AppID  string `json:"app_id"`
}

// NodeAppReq is a NodeApp request.
// TODO add a String() method and test for this struct.
type NodeAppReq struct {
	NodeApp
	Cmd string `json:"cmd,omitempty"`
}

// NodeAppResp is a NodeApp response.
// TODO add a String() method and test for this struct.
type NodeAppResp struct {
	NodeApp
	Status string `json:"status"`
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

// FilterFields returns the filterable fields for this model.
func (*NodeApp) FilterFields() []string {
	return []string{
		"node_id",
		"app_id",
	}
}

// Validate validates the request model.
// TODO add a test for this method.
func (n_ar *NodeAppReq) Validate() error {
	if err := n_ar.NodeApp.Validate(); err != nil {
		return err
	}
	switch n_ar.Cmd {
	case "start", "stop", "restart":
		return nil
	case "":
		return errors.New("cmd missing")
	default:
		return fmt.Errorf(`cmd "%s" is invalid`, n_ar.Cmd)
	}
}

// GetTableName returns the name of the persistence table.
func (n_ar *NodeAppReq) GetTableName() string {
	return n_ar.NodeApp.GetTableName()
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
