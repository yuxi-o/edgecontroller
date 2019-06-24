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

// NodeAppSummary is a summary representation of the node app.
type NodeAppSummary struct {
	ID string `json:"id"`
}

// NodeAppDetail is a detailed representation of the node app.
type NodeAppDetail struct {
	NodeAppSummary
	Status  string `json:"status"`
	Command string `json:"command"`
}

// NodeAppList is a list representation of node apps.
type NodeAppList struct {
	NodeApps []NodeAppSummary `json:"apps"`
}
