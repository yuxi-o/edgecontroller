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

// NodeVNFTrafficPolicy represents an association between a
// NodeVNF and a TrafficPolicy.
type NodeVNFTrafficPolicy struct {
	ID              string `json:"id"`
	NodeVNFID       string `json:"nodes_vnfs_id"`
	TrafficPolicyID string `json:"traffic_policy_id"`
}

// GetTableName returns the name of the persistence table.
func (*NodeVNFTrafficPolicy) GetTableName() string {
	return "nodes_vnfs_traffic_policies"
}

// GetID gets the ID.
func (n_v_tp *NodeVNFTrafficPolicy) GetID() string {
	return n_v_tp.ID
}

// SetID sets the ID.
func (n_v_tp *NodeVNFTrafficPolicy) SetID(id string) {
	n_v_tp.ID = id
}

// Validate validates the model.
func (n_v_tp *NodeVNFTrafficPolicy) Validate() error {
	if !uuid.IsValid(n_v_tp.ID) {
		return errors.New("id not a valid uuid")
	}
	if !uuid.IsValid(n_v_tp.NodeVNFID) {
		return errors.New("nodes_vnfs_id not a valid uuid")
	}
	if !uuid.IsValid(n_v_tp.TrafficPolicyID) {
		return errors.New("traffic_policy_id not a valid uuid")
	}

	return nil
}

func (n_v_tp *NodeVNFTrafficPolicy) String() string {
	return fmt.Sprintf(strings.TrimSpace(`
NodeVNFTrafficPolicy[
    ID: %s
    NodeVNFID: %s
    TrafficPolicyID: %s
]`),
		n_v_tp.ID,
		n_v_tp.NodeVNFID,
		n_v_tp.TrafficPolicyID)
}
