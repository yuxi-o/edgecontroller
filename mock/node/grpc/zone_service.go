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
	elapb "github.com/otcshare/edgecontroller/pb/ela"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type zoneService struct {
	zones []*elapb.NetworkZone
}

func (s *zoneService) reset() {
	s.zones = nil
}

func (s *zoneService) Create(
	ctx context.Context,
	zone *elapb.NetworkZone,
) (*empty.Empty, error) {
	s.zones = append(s.zones, zone)

	return &empty.Empty{}, nil
}

func (s *zoneService) Update(
	ctx context.Context,
	zone *elapb.NetworkZone,
) (*empty.Empty, error) {
	i := s.findIndex(zone.Id)

	if i < len(s.zones) {
		s.zones[i] = zone
		return &empty.Empty{}, nil
	}

	return nil, status.Errorf(
		codes.NotFound, "Network Zone %s not found", zone.Id)
}

func (s *zoneService) BulkUpdate(
	ctx context.Context,
	zones *elapb.NetworkZones,
) (*empty.Empty, error) {
	for _, zone := range zones.NetworkZones {
		if s.find(zone.Id) == nil {
			return nil, status.Errorf(
				codes.NotFound,
				"Network Zone %s not found", zone.Id)
		}
	}

	for _, zone := range zones.NetworkZones {
		if _, err := s.Update(ctx, zone); err != nil {
			return nil, err
		}
	}

	return &empty.Empty{}, nil
}

func (s *zoneService) GetAll(
	context.Context,
	*empty.Empty,
) (*elapb.NetworkZones, error) {
	return &elapb.NetworkZones{
		NetworkZones: s.zones,
	}, nil
}

func (s *zoneService) Get(
	ctx context.Context,
	id *elapb.ZoneID,
) (*elapb.NetworkZone, error) {
	zone := s.find(id.Id)

	if zone != nil {
		return zone, nil
	}

	return nil, status.Errorf(
		codes.NotFound, "Network Zone %s not found", id.Id)
}

func (s *zoneService) Delete(
	ctx context.Context,
	id *elapb.ZoneID,
) (*empty.Empty, error) {
	i := s.findIndex(id.Id)

	if i < len(s.zones) {
		s.delete(i)
		return &empty.Empty{}, nil
	}

	return nil, status.Errorf(
		codes.NotFound, "Network Zone %s not found", id.Id)
}

func (s *zoneService) find(id string) *elapb.NetworkZone {
	for _, zone := range s.zones {
		if zone.Id == id {
			return zone
		}
	}

	return nil
}

func (s *zoneService) findIndex(id string) int {
	for i, zone := range s.zones {
		if zone.Id == id {
			return i
		}
	}

	return len(s.zones)
}

func (s *zoneService) delete(i int) {
	copy(s.zones[i:], s.zones[i+1:])
	s.zones[len(s.zones)-1] = nil
	s.zones = s.zones[:len(s.zones)-1]
}
