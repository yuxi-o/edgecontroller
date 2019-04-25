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

// NodeVMAppTrafficPolicy represents an association between a Node, a VMApp, and
// a TrafficPolicy.
type NodeVMAppTrafficPolicy struct {
	ID              string `json:"id"`
	NodeVMAppID     string `json:"vm_app_id"`
	TrafficPolicyID string `json:"traffic_policy_id"`
}

// GetTableName returns the name of the persistence table.
func (*NodeVMAppTrafficPolicy) GetTableName() string {
	return "nodes_vm_apps_traffic_policies"
}

// GetID gets the ID.
func (n_vma_tp *NodeVMAppTrafficPolicy) GetID() string {
	return n_vma_tp.ID
}

// SetID sets the ID.
func (n_vma_tp *NodeVMAppTrafficPolicy) SetID(id string) {
	n_vma_tp.ID = id
}

// Validate validates the model.
func (n_vma_tp *NodeVMAppTrafficPolicy) Validate() error {
	if !uuid.IsValid(n_vma_tp.ID) {
		return errors.New("id not a valid uuid")
	}
	if !uuid.IsValid(n_vma_tp.NodeVMAppID) {
		return errors.New("node_vm_app_id not a valid uuid")
	}
	if !uuid.IsValid(n_vma_tp.TrafficPolicyID) {
		return errors.New("traffic_policy_id not a valid uuid")
	}

	return nil
}

func (n_vma_tp *NodeVMAppTrafficPolicy) String() string {
	return fmt.Sprintf(strings.TrimSpace(`
NodeVMAppTrafficPolicy[
    ID: %s
    NodeVMAppID: %s
    TrafficPolicyID: %s
]`),
		n_vma_tp.ID,
		n_vma_tp.NodeVMAppID,
		n_vma_tp.TrafficPolicyID)
}
