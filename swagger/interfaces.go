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

// InterfaceSummary is a summary representation of the interface.
type InterfaceSummary struct {
	ID                string   `json:"id"`
	Description       string   `json:"description"`
	Driver            string   `json:"driver"`
	Type              string   `json:"type"`
	MACAddress        string   `json:"mac_address"`
	VLAN              int      `json:"vlan"`
	Zones             []string `json:"zones"`
	FallbackInterface string   `json:"fallback_interface"`
}

// InterfaceDetail is a detailed representation of the interface.
type InterfaceDetail struct {
	InterfaceSummary
}

// NodeAppList is a list representation of interfaces.
type InterfaceList struct {
	Interfaces []InterfaceSummary `json:"interfaces"`
}
