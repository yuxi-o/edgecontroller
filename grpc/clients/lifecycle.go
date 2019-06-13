package clients

import (
	cce "github.com/smartedgemec/controller-ce"
	elapb "github.com/smartedgemec/controller-ce/pb/ela"
)

func fromPBLifecycleStatus(status *elapb.LifecycleStatus) cce.LifecycleStatus {
	switch status.Status {
	case elapb.LifecycleStatus_UNKNOWN:
		return cce.Unknown
	case elapb.LifecycleStatus_READY:
		return cce.Deployed
	case elapb.LifecycleStatus_STARTING:
		return cce.Starting
	case elapb.LifecycleStatus_RUNNING:
		return cce.Running
	case elapb.LifecycleStatus_STOPPING:
		return cce.Stopping
	case elapb.LifecycleStatus_STOPPED:
		return cce.Stopped
	case elapb.LifecycleStatus_ERROR:
		return cce.Error
	default:
		return cce.Unknown
	}
}
