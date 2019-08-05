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

package swagger

import (
	cce "github.com/otcshare/edgecontroller"
)

// AppSummary is a summary representation of the app.
type AppSummary struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Vendor      string `json:"vendor"`
	Description string `json:"description"`
}

// AppDetail is a detailed representation of the app.
type AppDetail struct {
	AppSummary
	Cores  int         `json:"cores"`
	Memory int         `json:"memory"`
	Ports  []PortProto `json:"ports"`
	Source string      `json:"source"`
}

// PortProto is a port and protocol combination.
type PortProto struct {
	cce.PortProto
}

// AppList is a list representation of apps.
type AppList struct {
	Apps []AppSummary `json:"apps"`
}
