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
)

type dnsService struct {
	// map of record name to records
	records map[string]*pb.DNSARecordSet
	// map of ip address to ip address
	forwarders map[string]string
}

func newDNSService() *dnsService {
	return &dnsService{
		records:    make(map[string]*pb.DNSARecordSet),
		forwarders: make(map[string]string),
	}
}

func (s *dnsService) SetA(
	ctx context.Context,
	record *pb.DNSARecordSet,
) (*empty.Empty, error) {
	s.records[record.Name] = record

	return &empty.Empty{}, nil
}

func (s *dnsService) DeleteA(
	ctx context.Context,
	record *pb.DNSARecordSet,
) (*empty.Empty, error) {
	delete(s.records, record.Name)

	return &empty.Empty{}, nil
}

func (s *dnsService) SetForwarders(
	ctx context.Context,
	forwarders *pb.DNSForwarders,
) (*empty.Empty, error) {
	for _, forwarder := range forwarders.IpAddresses {
		s.forwarders[forwarder] = forwarder
	}

	return &empty.Empty{}, nil
}

func (s *dnsService) DeleteForwarders(
	ctx context.Context,
	forwarders *pb.DNSForwarders,
) (*empty.Empty, error) {
	for _, forwarder := range forwarders.IpAddresses {
		delete(s.forwarders, forwarder)
	}

	return &empty.Empty{}, nil
}
