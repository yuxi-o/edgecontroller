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

// DNSContainerAppAlias is a DNS container app alias.
type DNSContainerAppAlias struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	ContainerAppID string `json:"container_app_id"`
}

// GetTableName returns the name of the persistence table.
func (*DNSContainerAppAlias) GetTableName() string {
	return "dns_container_app_aliases"
}

// GetID gets the ID.
func (alias *DNSContainerAppAlias) GetID() string {
	return alias.ID
}

// SetID sets the ID.
func (alias *DNSContainerAppAlias) SetID(id string) {
	alias.ID = id
}

// Validate validates the model.
func (alias *DNSContainerAppAlias) Validate() error {
	if !uuid.IsValid(alias.ID) {
		return errors.New("id not a valid uuid")
	}
	if alias.Name == "" {
		return errors.New("name cannot be empty")
	}
	if alias.Description == "" {
		return errors.New("description cannot be empty")
	}
	if !uuid.IsValid(alias.ContainerAppID) {
		return errors.New("container_app_id not a valid uuid")
	}

	return nil
}

func (alias *DNSContainerAppAlias) String() string {
	return fmt.Sprintf(strings.TrimSpace(`
DNSContainerAppAlias[
    ID: %s
    Name: %s
    Description: %s
    ContainerAppID: %s
]`),
		alias.ID,
		alias.Name,
		alias.Description,
		alias.ContainerAppID)
}
