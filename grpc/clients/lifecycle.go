package clients

import (
	cce "github.com/smartedgemec/controller-ce"
	evapb "github.com/smartedgemec/controller-ce/pb/eva"
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
