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
)

// ContainerVNF is a containerized VNF.
type ContainerVNF struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Vendor      string `json:"vendor"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Cores       int    `json:"cores"`
	Memory      int    `json:"memory"`
}

// GetTableName returns the name of the persistence table.
func (*ContainerVNF) GetTableName() string {
	return "container_vnfs"
}

// GetID gets the ID.
func (vnf *ContainerVNF) GetID() string {
	return vnf.ID
}

// SetID sets the ID.
func (vnf *ContainerVNF) SetID(id string) {
	vnf.ID = id
}

// Validate validates the model.
func (vnf *ContainerVNF) Validate() error {
	if vnf.ID == "" {
		return errors.New("id cannot be empty")
	}
	if vnf.Name == "" {
		return errors.New("name cannot be empty")
	}
	if vnf.Vendor == "" {
		return errors.New("vendor cannot be empty")
	}
	if vnf.Image == "" {
		return errors.New("image cannot be empty")
	}
	if vnf.Cores < 1 || vnf.Cores > MaxCores {
		return fmt.Errorf("cores must be in [1..%d]", MaxCores)
	}
	if vnf.Memory < 1 || vnf.Memory > MaxMemory {
		return fmt.Errorf("memory must be in [1..%d]", MaxMemory)
	}

	return nil
}

func (vnf *ContainerVNF) String() string {
	return fmt.Sprintf(strings.TrimSpace(`
ContainerVNF[
    ID: %s
    Name: %s
    Vendor: %s
    Description: %s
    Image: %s
    Cores: %d
    Memory: %d
]`),
		vnf.ID,
		vnf.Name,
		vnf.Vendor,
		vnf.Description,
		vnf.Image,
		vnf.Cores,
		vnf.Memory)
}
