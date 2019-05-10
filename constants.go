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

// MaxCores is the maximum number of cores that an application or VNF can use.
const MaxCores = 8

// MaxMemory is the maximum memory (in MB) that an application or VNF can use.
const MaxMemory = 16 * 1024

// LifecycleStatus is an application or VNF's status.
type LifecycleStatus int

const (
	// Unknown is an unknown lifecycle status
	Unknown LifecycleStatus = iota
	// Deploying is deploying to a node
	Deploying
	// Deployed is deployed to a node
	Deployed
	// Starting is starting
	Starting
	// Running is running
	Running
	// Stopping is stopping
	Stopping
	// Stopped is stopped
	Stopped
	// Error is an error status
	Error
)

func (s LifecycleStatus) String() string {
	switch s {
	case Deploying:
		return "deploying"
	case Deployed:
		return "deployed"
	case Starting:
		return "starting"
	case Running:
		return "running"
	case Stopping:
		return "stopping"
	case Stopped:
		return "stopped"
	case Error:
		return "error"
	case Unknown:
		fallthrough
	default:
		return "unknown"
	}
}
