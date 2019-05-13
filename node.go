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

// Node is a node (aka appliance or device).
type Node struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Location   string `json:"location"`
	Serial     string `json:"serial"`
	GRPCTarget string `json:"grpc_target"`
	// TODO figure out interface model
}

// GetTableName returns the name of the persistence table.
func (*Node) GetTableName() string {
	return "nodes"
}

// GetID gets the ID.
func (n *Node) GetID() string {
	return n.ID
}

// SetID sets the ID.
func (n *Node) SetID(id string) {
	n.ID = id
}

// Validate validates the model.
func (n *Node) Validate() error {
	if !uuid.IsValid(n.ID) {
		return errors.New("id not a valid uuid")
	}
	if n.Name == "" {
		return errors.New("name cannot be empty")
	}
	if n.Location == "" {
		return errors.New("location cannot be empty")
	}
	if n.Serial == "" {
		return errors.New("serial cannot be empty")
	}

	return nil
}

func (n *Node) String() string {
	return fmt.Sprintf(strings.TrimSpace(`
Node[
    ID: %s
    Name: %s
    Location: %s
    Serial: %s
]`),
		n.ID,
		n.Name,
		n.Location,
		n.Serial)
}
