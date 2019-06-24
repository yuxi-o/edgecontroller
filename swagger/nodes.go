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

// NodeSummary is a summary representation of the node.
type NodeSummary struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Serial   string `json:"serial"`
}

// NodeDetail is a detailed representation of the node.
type NodeDetail struct {
	NodeSummary
}

// NodeList is a list representation of nodes.
type NodeList struct {
	Nodes []NodeSummary `json:"nodes"`
}
