// Copyright 2019 Smart-Edge.com, Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clients

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	cce "github.com/smartedgemec/controller-ce"
	"github.com/smartedgemec/controller-ce/grpc"
	"github.com/smartedgemec/controller-ce/pb"
)

// ApplicationPolicyServiceClient wraps the PB client.
type ApplicationPolicyServiceClient struct {
	PBCli pb.ApplicationPolicyServiceClient
}

// NewApplicationPolicyServiceClient creates a new client.
func NewApplicationPolicyServiceClient(
	conn *grpc.ClientConn,
) *ApplicationPolicyServiceClient {
	return &ApplicationPolicyServiceClient{
		conn.NewApplicationPolicyServiceClient(),
	}
}

// Set sets the traffic policy.
func (c *ApplicationPolicyServiceClient) Set(
	ctx context.Context,
	policy *cce.TrafficPolicy,
) error {
	pbPolicy := &pb.TrafficPolicy{
		Id: policy.ID,
	}

	for _, rule := range policy.Rules {
		pbPolicy.TrafficRules = append(
			pbPolicy.TrafficRules, toPBTrafficRule(rule))
	}

	_, err := c.PBCli.Set(
		ctx,
		pbPolicy)

	if err != nil {
		return errors.Wrap(err, "error setting policy")
	}

	return nil
}

func toPBTrafficRule(tr *cce.TrafficRule) *pb.TrafficRule {
	return &pb.TrafficRule{
		Description: tr.Description,
		Priority:    uint32(tr.Priority),
		Source:      toPBTrafficSelector(tr.Source),
		Destination: toPBTrafficSelector(tr.Destination),
		Target:      toPBTrafficTarget(tr.Target),
	}
}

func toPBTrafficSelector(ts *cce.TrafficSelector) *pb.TrafficSelector {
	if ts == nil {
		return nil
	}

	return &pb.TrafficSelector{
		Description: ts.Description,
		Macs:        toPBMACFilter(ts.MACs),
		Ip:          toPBIPFilter(ts.IP),
		Gtp:         toPBGTPFilter(ts.GTP),
	}
}

func toPBMACFilter(macf *cce.MACFilter) *pb.MACFilter {
	if macf == nil {
		return nil
	}

	return &pb.MACFilter{
		MacAddresses: macf.MACAddresses,
	}
}

func toPBIPFilter(ipf *cce.IPFilter) *pb.IPFilter {
	if ipf == nil {
		return nil
	}

	return &pb.IPFilter{
		Address:   ipf.Address,
		Mask:      uint32(ipf.Mask),
		BeginPort: uint32(ipf.BeginPort),
		EndPort:   uint32(ipf.EndPort),
		Protocol:  ipf.Protocol,
	}
}

func toPBGTPFilter(gtpf *cce.GTPFilter) *pb.GTPFilter {
	if gtpf == nil {
		return nil
	}

	return &pb.GTPFilter{
		Address: gtpf.Address,
		Mask:    uint32(gtpf.Mask),
		Imsis:   gtpf.IMSIs,
	}
}

func toPBTrafficTarget(target *cce.TrafficTarget) *pb.TrafficTarget {
	return &pb.TrafficTarget{
		Description: target.Description,
		Action:      toPBTargetAction(target.Action),
		Mac:         toPBMACModifier(target.MAC),
		Ip:          toPBIPModifier(target.IP),
	}
}

func toPBTargetAction(action string) pb.TrafficTarget_TargetAction {
	switch action {
	case "accept":
		return pb.TrafficTarget_ACCEPT
	case "reject":
		return pb.TrafficTarget_REJECT
	case "drop":
		return pb.TrafficTarget_DROP
	default:
		panic(fmt.Sprintf("invalid target action %s", action))
	}
}

func toPBMACModifier(macMod *cce.MACModifier) *pb.MACModifier {
	if macMod == nil {
		return nil
	}

	return &pb.MACModifier{
		MacAddress: macMod.MACAddress,
	}
}

func toPBIPModifier(ipMod *cce.IPModifier) *pb.IPModifier {
	if ipMod == nil {
		return nil
	}

	return &pb.IPModifier{
		Address: ipMod.Address,
		Port:    uint32(ipMod.Port),
	}
}
