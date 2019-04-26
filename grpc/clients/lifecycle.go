package clients

import (
	cce "github.com/smartedgemec/controller-ce"
	"github.com/smartedgemec/controller-ce/pb"
)

func fromPBLifecycleStatus(status *pb.LifecycleStatus) cce.LifecycleStatus {
	switch status.Status {
	case pb.LifecycleStatus_UNKNOWN:
		return cce.Unknown
	case pb.LifecycleStatus_READY:
		return cce.Deployed
	case pb.LifecycleStatus_STARTING:
		return cce.Starting
	case pb.LifecycleStatus_RUNNING:
		return cce.Running
	case pb.LifecycleStatus_STOPPING:
		return cce.Stopping
	case pb.LifecycleStatus_STOPPED:
		return cce.Stopped
	case pb.LifecycleStatus_ERROR:
		return cce.Error
	default:
		return cce.Unknown
	}
}
