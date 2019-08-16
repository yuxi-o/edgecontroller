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
	cce "github.com/open-ness/edgecontroller"
)

// PolicySummary is a summary of an application traffic policy or interface traffic policy.
type PolicySummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// PolicyDetail is a detailed representation of the traffic policy.
type PolicyDetail struct {
	PolicySummary
	Rules []*cce.TrafficRule `json:"traffic_rules"`
}

// PolicyList is a list representation of traffic policies.
type PolicyList struct {
	Policies []PolicySummary `json:"policies"`
}
