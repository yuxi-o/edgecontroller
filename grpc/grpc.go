package grpc

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/smartedgemec/controller-ce/pb"
	"google.golang.org/grpc"
)

// ClientConn wraps grpc.ClientConn
type ClientConn struct {
	conn *grpc.ClientConn
}

// Dial dials the remote server.
func Dial(ctx context.Context, target string) (*ClientConn, error) {
	timeoutCtx, cancelFunc := context.WithTimeout(
		ctx, 2*time.Second)
	defer cancelFunc()

	conn, err := grpc.DialContext(
		timeoutCtx,
		target,
		grpc.WithInsecure(),
		grpc.WithBlock())
	if err != nil {
		return nil, errors.Wrapf(err, "dial %s failed", target)
	}

	return &ClientConn{conn}, nil
}

// Close wraps grpc.Close()
func (c *ClientConn) Close() error {
	return c.conn.Close()
}

// NewApplicationDeploymentServiceClient wraps the pb function.
func (c *ClientConn) NewApplicationDeploymentServiceClient() pb.ApplicationDeploymentServiceClient { //nolint:lll
	return pb.NewApplicationDeploymentServiceClient(c.conn)
}

// NewApplicationLifecycleServiceClient wraps the pb function.
func (c *ClientConn) NewApplicationLifecycleServiceClient() pb.ApplicationLifecycleServiceClient { //nolint:lll
	return pb.NewApplicationLifecycleServiceClient(c.conn)
}

// NewApplicationPolicyServiceClient wraps the pb function.
func (c *ClientConn) NewApplicationPolicyServiceClient() pb.ApplicationPolicyServiceClient { //nolint:lll
	return pb.NewApplicationPolicyServiceClient(c.conn)
}

// NewVNFDeploymentServiceClient wraps the pb function.
func (c *ClientConn) NewVNFDeploymentServiceClient() pb.VNFDeploymentServiceClient { //nolint:lll
	return pb.NewVNFDeploymentServiceClient(c.conn)
}

// NewVNFLifecycleServiceClient wraps the pb function.
func (c *ClientConn) NewVNFLifecycleServiceClient() pb.VNFLifecycleServiceClient { //nolint:lll
	return pb.NewVNFLifecycleServiceClient(c.conn)
}

// NewInterfaceServiceClient wraps the pb function.
func (c *ClientConn) NewInterfaceServiceClient() pb.InterfaceServiceClient {
	return pb.NewInterfaceServiceClient(c.conn)
}

// NewZoneServiceClient wraps the pb function.
func (c *ClientConn) NewZoneServiceClient() pb.ZoneServiceClient {
	return pb.NewZoneServiceClient(c.conn)
}

// NewInterfacePolicyServiceClient wraps the pb function.
func (c *ClientConn) NewInterfacePolicyServiceClient() pb.InterfacePolicyServiceClient { //nolint:lll
	return pb.NewInterfacePolicyServiceClient(c.conn)
}
