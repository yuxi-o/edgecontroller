// Copyright 2019 Intel Corporation and Smart-Edge.com, Inc. All rights reserved
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

package clients

import (
	cce "github.com/open-ness/edgecontroller"
	evapb "github.com/open-ness/edgecontroller/pb/eva"
)

func fromPBLifecycleStatus(status *evapb.LifecycleStatus) cce.LifecycleStatus {
	switch status.Status {
	case evapb.LifecycleStatus_UNKNOWN:
		return cce.Unknown
	case evapb.LifecycleStatus_READY:
		return cce.Deployed
	case evapb.LifecycleStatus_STARTING:
		return cce.Starting
	case evapb.LifecycleStatus_RUNNING:
		return cce.Running
	case evapb.LifecycleStatus_STOPPING:
		return cce.Stopping
	case evapb.LifecycleStatus_STOPPED:
		return cce.Stopped
	case evapb.LifecycleStatus_ERROR:
		return cce.Error
	default:
		return cce.Unknown
	}
}
