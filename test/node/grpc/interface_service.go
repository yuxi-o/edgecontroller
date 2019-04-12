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

type interfaceService struct {
	nis []*pb.NetworkInterface

	// reference to policy server
	policyService *interfacePolicyService
}

func newInterfaceService() *interfaceService {
	return &interfaceService{
		nis: []*pb.NetworkInterface{
			{
				Id:          "if0",
				Description: "interface0",
				Driver:      pb.NetworkInterface_KERNEL,
				Type:        pb.NetworkInterface_NONE,
				MacAddress:  "mac0",
				Vlan:        0,
				Zones:       nil,
			},
			{
				Id:          "if1",
				Description: "interface1",
				Driver:      pb.NetworkInterface_KERNEL,
				Type:        pb.NetworkInterface_NONE,
				MacAddress:  "mac1",
				Vlan:        1,
				Zones:       nil,
			},
			{
				Id:          "if2",
				Description: "interface2",
				Driver:      pb.NetworkInterface_KERNEL,
				Type:        pb.NetworkInterface_NONE,
				MacAddress:  "mac2",
				Vlan:        2,
				Zones:       nil,
			},
			{
				Id:          "if3",
				Description: "interface3",
				Driver:      pb.NetworkInterface_KERNEL,
				Type:        pb.NetworkInterface_NONE,
				MacAddress:  "mac3",
				Vlan:        3,
				Zones:       nil,
			},
		},
	}
}

func (s *interfaceService) init(policyService *interfacePolicyService) {
	s.policyService = policyService

	s.policyService.policies["if0"] = defaultPolicy("if0")
	s.policyService.policies["if1"] = defaultPolicy("if1")
	s.policyService.policies["if2"] = defaultPolicy("if2")
	s.policyService.policies["if3"] = defaultPolicy("if3")
}

func (s *interfaceService) Update(
	ctx context.Context,
	ni *pb.NetworkInterface,
) (*empty.Empty, error) {
	i := s.findIndex(ni.Id)

	if i < len(s.nis) {
		s.nis[i] = ni
		return &empty.Empty{}, nil
	}

	return nil, status.Errorf(
		codes.NotFound, "Network Interface %s not found", ni.Id)
}

func (s *interfaceService) BulkUpdate(
	ctx context.Context,
	nis *pb.NetworkInterfaces,
) (*empty.Empty, error) {
	for _, ni := range nis.NetworkInterfaces {
		if s.find(ni.Id) == nil {
			return nil, status.Errorf(
				codes.NotFound,
				"Network Interface %s not found", ni.Id)
		}
	}

	for _, ni := range nis.NetworkInterfaces {
		if _, err := s.Update(ctx, ni); err != nil {
			return nil, err
		}
	}

	return &empty.Empty{}, nil
}

func (s *interfaceService) GetAll(
	context.Context,
	*empty.Empty,
) (*pb.NetworkInterfaces, error) {
	return &pb.NetworkInterfaces{
		NetworkInterfaces: s.nis,
	}, nil
}

func (s *interfaceService) Get(
	ctx context.Context,
	id *pb.InterfaceID,
) (*pb.NetworkInterface, error) {
	ni := s.find(id.Id)

	if ni != nil {
		return ni, nil
	}

	return nil, status.Errorf(
		codes.NotFound, "Network Interface %s not found", id.Id)
}

func (s *interfaceService) find(id string) *pb.NetworkInterface {
	for _, ni := range s.nis {
		if ni.Id == id {
			return ni
		}
	}

	return nil
}

func (s *interfaceService) findIndex(id string) int {
	for i, ni := range s.nis {
		if ni.Id == id {
			return i
		}
	}

	return len(s.nis)
}
