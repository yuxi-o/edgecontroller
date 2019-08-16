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
	elapb "github.com/open-ness/edgecontroller/pb/ela"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type interfacePolicyService struct {
	// map of interface ID to traffic policy
	policies map[string]*elapb.TrafficPolicy

	// reference to interface server
	interfaceService *interfaceService
}

func newInterfacePolicyService(
	interfaceService *interfaceService,
) *interfacePolicyService {
	return &interfacePolicyService{
		policies:         make(map[string]*elapb.TrafficPolicy),
		interfaceService: interfaceService,
	}
}

func (s *interfacePolicyService) reset() {
	s.policies = make(map[string]*elapb.TrafficPolicy)
}

func (s *interfacePolicyService) Set(
	ctx context.Context,
	policy *elapb.TrafficPolicy,
) (*empty.Empty, error) {
	if s.interfaceService.find(policy.Id) == nil {
		return nil, status.Errorf(
			codes.NotFound, "Network Interface %s not found", policy.Id)
	}

	s.policies[policy.Id] = policy

	return &empty.Empty{}, nil
}

func (s *interfacePolicyService) Get(
	ctx context.Context,
	id *elapb.InterfaceID,
) (*elapb.TrafficPolicy, error) {
	if s.policies[id.Id] == nil {
		return nil, status.Errorf(
			codes.NotFound, "Network Interface %s not found", id.Id)
	}

	return s.policies[id.Id], nil
}
