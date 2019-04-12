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

package main

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/smartedgemec/controller-ce/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type applicationPolicyServer struct {
	// map of application ID to traffic policy
	policies map[string]*pb.TrafficPolicy

	// reference to application server
	appServer *applicationServer
}

func newApplicationPolicyServer(
	appServer *applicationServer,
) *applicationPolicyServer {
	return &applicationPolicyServer{
		policies:  make(map[string]*pb.TrafficPolicy),
		appServer: appServer,
	}
}

func (s *applicationPolicyServer) Set(
	ctx context.Context,
	policy *pb.TrafficPolicy,
) (*empty.Empty, error) {
	if s.appServer.find(policy.Id) == nil {
		return nil, status.Errorf(
			codes.NotFound, "Application %s not found", policy.Id)
	}

	s.policies[policy.Id] = policy

	return &empty.Empty{}, nil
}

func (s *applicationPolicyServer) Get(
	ctx context.Context,
	id *pb.ApplicationID,
) (*pb.TrafficPolicy, error) {
	if s.policies[id.Id] == nil {
		return nil, status.Errorf(
			codes.NotFound, "Application %s not found", id.Id)
	}

	return s.policies[id.Id], nil
}
