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
	"net/url"
	"strings"

	"github.com/smartedgemec/controller-ce/uuid"
)

// VNF is a Virtual Network Function.
type VNF struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Vendor      string `json:"vendor"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Cores       int    `json:"cores"`
	Memory      int    `json:"memory"`
	Source      string `json:"source"`
}

// GetTableName returns the name of the persistence table.
func (*VNF) GetTableName() string {
	return "vnfs"
}

// GetID gets the ID.
func (vnf *VNF) GetID() string {
	return vnf.ID
}

// SetID sets the ID.
func (vnf *VNF) SetID(id string) {
	vnf.ID = id
}

// Validate validates the model.
func (vnf *VNF) Validate() error { // nolint: gocyclo
	if !uuid.IsValid(vnf.ID) {
		return errors.New("id not a valid uuid")
	}
	if vnf.Type != "container" && vnf.Type != "vm" {
		return errors.New(`type must be either "container" or "vm"`)
	}
	if vnf.Name == "" {
		return errors.New("name cannot be empty")
	}
	if vnf.Vendor == "" {
		return errors.New("vendor cannot be empty")
	}
	if vnf.Version == "" {
		return errors.New("version cannot be empty")
	}
	if vnf.Cores < 1 || vnf.Cores > MaxCores {
		return fmt.Errorf("cores must be in [1..%d]", MaxCores)
	}
	if vnf.Memory < 1 || vnf.Memory > MaxMemory {
		return fmt.Errorf("memory must be in [1..%d]", MaxMemory)
	}
	if vnf.Source == "" {
		return errors.New("source cannot be empty")
	}
	if _, err := url.ParseRequestURI(vnf.Source); err != nil {
		return errors.New("source cannot be parsed as a URI")
	}

	return nil
}

func (vnf *VNF) String() string {
	return fmt.Sprintf(strings.TrimSpace(`
VNF[
    ID: %s
    Name: %s
    Version: %s
    Vendor: %s
    Description: %s
    Cores: %d
    Memory: %d
    Source: %s
]`),
		vnf.ID,
		vnf.Name,
		vnf.Version,
		vnf.Vendor,
		vnf.Description,
		vnf.Cores,
		vnf.Memory,
		vnf.Source)
}
