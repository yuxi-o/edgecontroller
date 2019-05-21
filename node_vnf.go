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

// NodeVNF represents an association between a Node and a VNF.
type NodeVNF struct {
	ID     string `json:"id"`
	NodeID string `json:"node_id"`
	VNFID  string `json:"vnf_id"`
}

// NodeVNFReq is a NodeVNF request.
type NodeVNFReq struct {
	NodeVNF
	Cmd string `json:"cmd,omitempty"`
}

// NodeVNFResp is a NodeVNF response.
type NodeVNFResp struct {
	NodeVNF
	Status string `json:"status"`
}

// GetTableName returns the name of the persistence table.
func (*NodeVNF) GetTableName() string {
	return "nodes_vnfs"
}

// GetID gets the ID.
func (n_vnf *NodeVNF) GetID() string {
	return n_vnf.ID
}

// SetID sets the ID.
func (n_vnf *NodeVNF) SetID(id string) {
	n_vnf.ID = id
}

// GetNodeID gets the node ID.
func (n_vnf *NodeVNF) GetNodeID() string {
	return n_vnf.NodeID
}

// Validate validates the model.
func (n_vnf *NodeVNF) Validate() error {
	if !uuid.IsValid(n_vnf.ID) {
		return errors.New("id not a valid uuid")
	}
	if !uuid.IsValid(n_vnf.NodeID) {
		return errors.New("node_id not a valid uuid")
	}
	if !uuid.IsValid(n_vnf.VNFID) {
		return errors.New("vnf_id not a valid uuid")
	}

	return nil
}

// Validate validates the request model.
func (n_ar *NodeVNFReq) Validate() error {
	if err := n_ar.NodeVNF.Validate(); err != nil {
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
func (n_ar *NodeVNFReq) GetTableName() string {
	return n_ar.NodeVNF.GetTableName()
}

func (n_vnf *NodeVNF) String() string {
	return fmt.Sprintf(strings.TrimSpace(`
NodeVNF[
    ID: %s
    NodeID: %s
    VNFID: %s
]`),
		n_vnf.ID,
		n_vnf.NodeID,
		n_vnf.VNFID)
}
