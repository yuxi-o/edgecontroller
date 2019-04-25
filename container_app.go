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

// ContainerApp is a containerized application.
type ContainerApp struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Vendor      string `json:"vendor"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Cores       int    `json:"cores"`
	Memory      int    `json:"memory"` // in MB
}

// GetTableName returns the name of the persistence table.
func (*ContainerApp) GetTableName() string {
	return "container_apps"
}

// GetID gets the ID.
func (app *ContainerApp) GetID() string {
	return app.ID
}

// SetID sets the ID.
func (app *ContainerApp) SetID(id string) {
	app.ID = id
}

// Validate validates the model.
func (app *ContainerApp) Validate() error {
	if !uuid.IsValid(app.ID) {
		return errors.New("id not a valid uuid")
	}
	if app.Name == "" {
		return errors.New("name cannot be empty")
	}
	if app.Vendor == "" {
		return errors.New("vendor cannot be empty")
	}
	if app.Image == "" {
		return errors.New("image cannot be empty")
	}
	if app.Cores < 1 || app.Cores > MaxCores {
		return fmt.Errorf("cores must be in [1..%d]", MaxCores)
	}
	if app.Memory < 1 || app.Memory > MaxMemory {
		return fmt.Errorf("memory must be in [1..%d]", MaxMemory)
	}

	return nil
}

func (app *ContainerApp) String() string {
	return fmt.Sprintf(strings.TrimSpace(`
ContainerApp[
    ID: %s
    Name: %s
    Vendor: %s
    Description: %s
    Image: %s
    Cores: %d
    Memory: %d
]`),
		app.ID,
		app.Name,
		app.Vendor,
		app.Description,
		app.Image,
		app.Cores,
		app.Memory)
}
