// Copyright 2019 Smart-Edge.com, Inc. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/smartedgemec/controller-ce/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type appDeployLifeService struct {
	// maps of application ID to application
	containerApps map[string]*pb.Application
	vmApps        map[string]*pb.Application

	// reference to policy server
	appPolicyService *appPolicyService
}

func newAppDeployLifeService() *appDeployLifeService {
	return &appDeployLifeService{
		containerApps: make(map[string]*pb.Application),
		vmApps:        make(map[string]*pb.Application),
	}
}

func (s *appDeployLifeService) DeployContainer(
	ctx context.Context,
	containerApp *pb.Application,
) (*empty.Empty, error) {
	id := containerApp.Id
	s.containerApps[id] = containerApp
	containerApp.Status = pb.LifecycleStatus_READY

	return &empty.Empty{}, nil
}

func (s *appDeployLifeService) DeployVM(
	ctx context.Context,
	vmApp *pb.Application,
) (*empty.Empty, error) {
	id := vmApp.Id
	s.vmApps[id] = vmApp
	vmApp.Status = pb.LifecycleStatus_READY

	return &empty.Empty{}, nil
}

func (s *appDeployLifeService) GetStatus(
	ctx context.Context,
	id *pb.ApplicationID,
) (*pb.LifecycleStatus, error) {
	if containerApp, ok := s.containerApps[id.Id]; ok {
		return &pb.LifecycleStatus{
			Status: containerApp.Status,
		}, nil
	}

	if vmApp, ok := s.vmApps[id.Id]; ok {
		return &pb.LifecycleStatus{
			Status: vmApp.Status,
		}, nil
	}

	return nil, status.Errorf(codes.NotFound, "Application %s not found", id.Id)
}

func (s *appDeployLifeService) Redeploy(
	ctx context.Context,
	app *pb.Application,
) (*empty.Empty, error) {
	if oldApp, ok := s.containerApps[app.Id]; ok {
		app.Status = oldApp.Status
		s.containerApps[app.Id] = app
		return &empty.Empty{}, nil
	}

	if oldApp, ok := s.vmApps[app.Id]; ok {
		app.Status = oldApp.Status
		s.vmApps[app.Id] = app
		return &empty.Empty{}, nil
	}

	return nil, status.Errorf(
		codes.NotFound, "Application %s not found", app.Id)
}

func (s *appDeployLifeService) Undeploy(
	ctx context.Context,
	id *pb.ApplicationID,
) (*empty.Empty, error) {
	var policyExists bool
	if _, policyExists = s.appPolicyService.policies[id.Id]; policyExists {
		delete(s.appPolicyService.policies, id.Id)
		policyExists = true
	}

	if _, ok := s.containerApps[id.Id]; ok {
		delete(s.containerApps, id.Id)
		return &empty.Empty{}, nil
	}

	if _, ok := s.vmApps[id.Id]; ok {
		delete(s.vmApps, id.Id)
		return &empty.Empty{}, nil
	}

	if policyExists {
		return nil, status.Errorf(codes.DataLoss,
			"Application %s not found but had a policy!", id.Id)
	}

	return nil, status.Errorf(codes.NotFound, "Application %s not found", id.Id)
}

func (s *appDeployLifeService) Start(
	ctx context.Context,
	cmd *pb.LifecycleCommand,
) (*empty.Empty, error) {
	app := s.find(cmd.Id)

	if app != nil {
		switch app.Status {
		case pb.LifecycleStatus_READY:
		case pb.LifecycleStatus_STOPPED:
		default:
			return nil, status.Errorf(
				codes.FailedPrecondition, "Application %s not stopped or ready",
				cmd.Id)
		}

		app.Status = pb.LifecycleStatus_RUNNING
		return &empty.Empty{}, nil
	}

	return nil, status.Errorf(
		codes.NotFound, "Application %s not found", cmd.Id)
}

func (s *appDeployLifeService) Stop(
	ctx context.Context,
	cmd *pb.LifecycleCommand,
) (*empty.Empty, error) {
	app := s.find(cmd.Id)

	if app != nil {
		if app.Status != pb.LifecycleStatus_RUNNING {
			return nil, status.Errorf(
				codes.FailedPrecondition, "Application %s not running", cmd.Id)
		}

		app.Status = pb.LifecycleStatus_STOPPED
		return &empty.Empty{}, nil
	}

	return nil, status.Errorf(
		codes.NotFound, "Application %s not found", cmd.Id)
}

func (s *appDeployLifeService) Restart(
	ctx context.Context,
	cmd *pb.LifecycleCommand,
) (*empty.Empty, error) {
	app := s.find(cmd.Id)

	if app != nil {
		if app.Status != pb.LifecycleStatus_RUNNING {
			return nil, status.Errorf(
				codes.FailedPrecondition, "Application %s not running", cmd.Id)
		}

		return &empty.Empty{}, nil
	}

	return nil, status.Errorf(
		codes.NotFound, "Application %s not found", cmd.Id)
}

func (s *appDeployLifeService) find(id string) *pb.Application {
	if containerApp, ok := s.containerApps[id]; ok {
		return containerApp
	}

	if vmApp, ok := s.vmApps[id]; ok {
		return vmApp
	}

	return nil
}
