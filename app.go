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

// App is an application.
type App struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Vendor      string `json:"vendor"`
	Description string `json:"description"`
	Cores       int    `json:"cores"`
	Memory      int    `json:"memory"` // in MB
	Source      string `json:"source"`
}

// GetTableName returns the name of the persistence table.
func (*App) GetTableName() string {
	return "apps"
}

// GetID gets the ID.
func (app *App) GetID() string {
	return app.ID
}

// SetID sets the ID.
func (app *App) SetID(id string) {
	app.ID = id
}

// Validate validates the model.
func (app *App) Validate() error { // nolint: gocyclo
	if !uuid.IsValid(app.ID) {
		return errors.New("id not a valid uuid")
	}
	if app.Type != "container" && app.Type != "vm" {
		return errors.New(`type must be either "container" or "vm"`)
	}
	if app.Name == "" {
		return errors.New("name cannot be empty")
	}
	if app.Vendor == "" {
		return errors.New("vendor cannot be empty")
	}
	if app.Version == "" {
		return errors.New("version cannot be empty")
	}
	if app.Cores < 1 || app.Cores > MaxCores {
		return fmt.Errorf("cores must be in [1..%d]", MaxCores)
	}
	if app.Memory < 1 || app.Memory > MaxMemory {
		return fmt.Errorf("memory must be in [1..%d]", MaxMemory)
	}
	if app.Source == "" {
		return errors.New("source cannot be empty")
	}
	if _, err := url.ParseRequestURI(app.Source); err != nil {
		return errors.New("source cannot be parsed as a URI")
	}

	return nil
}

func (app *App) String() string {
	return fmt.Sprintf(strings.TrimSpace(`
App[
    ID: %s
    Name: %s
    Version: %s
    Vendor: %s
    Description: %s
    Cores: %d
    Memory: %d
    Source: %s
]`),
		app.ID,
		app.Name,
		app.Version,
		app.Vendor,
		app.Description,
		app.Cores,
		app.Memory,
		app.Source)
}
