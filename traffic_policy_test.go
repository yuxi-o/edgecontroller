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

package cce_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cce "github.com/smartedgemec/controller-ce"
)

var _ = Describe("Entities: TrafficPolicy", func() {
	var (
		tp *cce.TrafficPolicy
	)

	BeforeEach(func() {
		tp = &cce.TrafficPolicy{
			ID: "9d740cee-035f-4076-847c-d1c80cdf19db",
			Rules: []*cce.TrafficRule{
				{
					Description: "test-rule-1",
					Priority:    1,
					Source: &cce.TrafficSelector{
						Description: "source-source-1",
						MACs: &cce.MACFilter{
							MACAddresses: []string{
								"F0-59-8E-7B-36-8A",
								"23-20-8E-15-89-D1",
								"35-A4-38-73-35-45",
							},
						},
						IP: &cce.IPFilter{
							Address:   "223.1.1.0",
							Mask:      16,
							BeginPort: 2000,
							EndPort:   2012,
							Protocol:  "tcp",
						},
						GTP: &cce.GTPFilter{
							Address: "10.6.7.2",
							Mask:    12,
							IMSIs: []string{
								"310150123456789",
								"310150123456790",
								"310150123456791",
							},
						},
					},
					Destination: &cce.TrafficSelector{
						Description: "test-destination-1",
						MACs: &cce.MACFilter{
							MACAddresses: []string{
								"7D-C2-3A-1C-63-D9",
								"E9-6B-D1-D2-1A-6B",
								"C8-32-A9-43-85-55",
							},
						},
						IP: &cce.IPFilter{
							Address:   "64.1.1.0",
							Mask:      16,
							BeginPort: 1000,
							EndPort:   1012,
							Protocol:  "tcp",
						},
						GTP: &cce.GTPFilter{
							Address: "108.6.7.2",
							Mask:    4,
							IMSIs: []string{
								"310150123456792",
								"310150123456793",
								"310150123456794",
							},
						},
					},
					Target: &cce.TrafficTarget{
						Description: "test-target-1",
						Action:      "accept",
						MAC: &cce.MACModifier{
							MACAddress: "C7-5A-E7-98-1B-A3",
						},
						IP: &cce.IPModifier{
							Address: "123.2.3.4",
							Port:    1600,
						},
					},
				},
				{
					Description: "test-rule-2",
					Priority:    2,
					Source: &cce.TrafficSelector{
						Description: "source-source-2",
						MACs: &cce.MACFilter{
							MACAddresses: []string{
								"43-78-01-EE-B5-8D",
								"DB-C6-F2-CC-0B-31",
								"66-69-C2-D8-78-83",
							},
						},
						IP: &cce.IPFilter{
							Address:   "12.1.1.0",
							Mask:      6,
							BeginPort: 5000,
							EndPort:   5012,
							Protocol:  "udp",
						},
						GTP: &cce.GTPFilter{
							Address: "10.66.7.2",
							Mask:    4,
							IMSIs: []string{
								"310150123456795",
								"310150123456796",
								"310150123456797",
							},
						},
					},
					Destination: &cce.TrafficSelector{
						Description: "test-destination-2",
						MACs: &cce.MACFilter{
							MACAddresses: []string{
								"30-50-D5-50-47-90",
								"14-7C-F2-7B-DC-73",
								"40-D2-CA-98-CA-CE",
							},
						},
						IP: &cce.IPFilter{
							Address:   "132.15.16.20",
							Mask:      3,
							BeginPort: 55000,
							EndPort:   55012,
							Protocol:  "udp",
						},
						GTP: &cce.GTPFilter{
							Address: "70.6.17.2",
							Mask:    4,
							IMSIs: []string{
								"310150123456798",
								"310150123456799",
								"310150123456800",
							},
						},
					},
					Target: &cce.TrafficTarget{
						Description: "test-target-2",
						Action:      "reject",
						MAC: &cce.MACModifier{
							MACAddress: "28-3F-D2-2C-47-1F",
						},
						IP: &cce.IPModifier{
							Address: "242.25.31.14",
							Port:    2600,
						},
					},
				},
			},
		}
	})

	Describe("GetTableName", func() {
		It(`Should return "traffic_policies"`, func() {
			Expect(tp.GetTableName()).To(Equal("traffic_policies"))
		})
	})

	Describe("GetID", func() {
		It("Should return the ID", func() {
			Expect(tp.GetID()).To(Equal("9d740cee-035f-4076-847c-d1c80cdf19db"))
		})
	})

	Describe("SetID", func() {
		It("Should set and return the updated ID", func() {
			By("Setting the ID")
			tp.SetID("456")

			By("Getting the updated ID")
			Expect(tp.ID).To(Equal("456"))
		})
	})

	Describe("Validate", func() {
		It("Should return an error if ID is not a UUID", func() {
			tp.ID = "123"
			Expect(tp.Validate()).To(MatchError("id not a valid uuid"))
		})

		It("Should return an error if Rules has zero length", func() {
			tp.Rules = nil
			Expect(tp.Validate()).To(MatchError("rules cannot be empty"))
		})

		It("Should return an error if Rules.Description is empty", func() {
			tp.Rules[0].Description = ""
			Expect(tp.Validate()).To(MatchError(
				"rules[0].description cannot be empty"))
		})

		It("Should return an error if Rules.Priority is < 1", func() {
			tp.Rules[0].Priority = 0
			Expect(tp.Validate()).To(MatchError(
				"rules[0].priority must be in [1..65536]"))
		})

		It("Should return an error if Rules.Priority is > 65536", func() {
			tp.Rules[0].Priority = 65537
			Expect(tp.Validate()).To(MatchError(
				"rules[0].priority must be in [1..65536]"))
		})

		It("Should return an error if Rules.Source is nil", func() {
			tp.Rules[0].Source = nil
			Expect(tp.Validate()).To(MatchError(
				"rules[0].source cannot be empty"))
		})

		It("Should return an error if Rules.Destination is nil", func() {
			tp.Rules[0].Destination = nil
			Expect(tp.Validate()).To(MatchError(
				"rules[0].destination cannot be empty"))
		})

		It("Should return an error if Rules.Target is nil", func() {
			tp.Rules[0].Target = nil
			Expect(tp.Validate()).To(MatchError(
				"rules[0].target cannot be empty"))
		})

		It("Should return an error if Rules.Source.Description is "+
			"empty", func() {
			tp.Rules[0].Source.Description = ""
			Expect(tp.Validate()).To(MatchError(
				"rules[0].source.description cannot be empty"))
		})

		It("Should return an error if Rules.Source.MACs, Rules.Source.IP, "+
			"and Rules.Source.GTP are all nil", func() {
			tp.Rules[0].Source.MACs = nil
			tp.Rules[0].Source.IP = nil
			tp.Rules[0].Source.GTP = nil
			Expect(tp.Validate()).To(MatchError(
				"rules[0].source.macs|ip|gtp cannot all be nil"))
		})

		It("Should return an error if Rules.Source.MACs.MACAddresses "+
			"contains an invalid MAC address", func() {
			tp.Rules[0].Source.MACs.MACAddresses[0] = "abc-def"
			Expect(tp.Validate()).To(MatchError(
				"rules[0].source.macs.mac_addresses[0] could not be parsed " +
					"(address abc-def: invalid MAC address)"))
		})

		It("Should return an error if Rules.Source.IP.Address is "+
			"invalid", func() {
			tp.Rules[0].Source.IP.Address = "987.0.3.4"
			Expect(tp.Validate()).To(MatchError(
				"rules[0].source.ip.address could not be parsed"))
		})

		It("Should return an error if Rules.Source.IP.Mask is < 0", func() {
			tp.Rules[0].Source.IP.Mask = -1
			Expect(tp.Validate()).To(MatchError(
				"rules[0].source.ip.mask must be in [0..128]"))
		})

		It("Should return an error if Rules.Source.IP.Mask is > 128", func() {
			tp.Rules[0].Source.IP.Mask = 129
			Expect(tp.Validate()).To(MatchError(
				"rules[0].source.ip.mask must be in [0..128]"))
		})

		It("Should return an error if Rules.Source.IP.BeginPort "+
			"is < 1", func() {
			tp.Rules[0].Source.IP.BeginPort = 0
			Expect(tp.Validate()).To(MatchError(
				"rules[0].source.ip.begin_port must be in [1..65536]"))
		})

		It("Should return an error if Rules.Source.IP.BeginPort "+
			"is > 65536", func() {
			tp.Rules[0].Source.IP.BeginPort = 65537
			Expect(tp.Validate()).To(MatchError(
				"rules[0].source.ip.begin_port must be in [1..65536]"))
		})

		It("Should return an error if Rules.Source.IP.EndPort is < 1", func() {
			tp.Rules[0].Source.IP.EndPort = 0
			Expect(tp.Validate()).To(MatchError(
				"rules[0].source.ip.end_port must be in [1..65536]"))
		})

		It("Should return an error if Rules.Source.IP.EndPort "+
			"is > 65536", func() {
			tp.Rules[0].Source.IP.EndPort = 65537
			Expect(tp.Validate()).To(MatchError(
				"rules[0].source.ip.end_port must be in [1..65536]"))
		})

		It("Should return an error if Rules.Source.IP.BeginPort "+
			"is > EndPort", func() {
			tp.Rules[0].Source.IP.BeginPort = 1024
			tp.Rules[0].Source.IP.EndPort = 1023
			Expect(tp.Validate()).To(MatchError(
				"rules[0].source.ip.begin_port must be <= end_port"))
		})

		It("Should return an error if Rules.Source.IP.Protocol is not one of"+
			"[tcp, udp, icmp, sctp]", func() {
			tp.Rules[0].Source.IP.Protocol = "abc"
			Expect(tp.Validate()).To(MatchError(
				"rules[0].source.ip.protocol must be one of " +
					"[tcp, udp, icmp, sctp]"))
		})

		It("Should return an error if Rules.Source.GTP.Address is "+
			"empty", func() {
			tp.Rules[0].Source.GTP.Address = ""
			Expect(tp.Validate()).To(MatchError(
				"rules[0].source.gtp.address cannot be empty"))
		})

		It("Should return an error if Rules.Source.GTP.Address is "+
			"invalid", func() {
			tp.Rules[0].Source.GTP.Address = "555.3.2.9"
			Expect(tp.Validate()).To(MatchError(
				"rules[0].source.gtp.address could not be parsed"))
		})

		It("Should return an error if Rules.Source.GTP.Mask is < 0", func() {
			tp.Rules[0].Source.GTP.Mask = -1
			Expect(tp.Validate()).To(MatchError(
				"rules[0].source.gtp.mask must be in [0..128]"))
		})

		It("Should return an error if Rules.Source.GTP.Mask is > 128", func() {
			tp.Rules[0].Source.GTP.Mask = 129
			Expect(tp.Validate()).To(MatchError(
				"rules[0].source.gtp.mask must be in [0..128]"))
		})

		It("Should return an error if Rules.Source.GTP.IMSIs contains a value "+
			"that is not numeric", func() {
			tp.Rules[0].Source.GTP.IMSIs[0] = "abcdef"
			Expect(tp.Validate()).To(MatchError(
				"rules[0].source.gtp.imsis[0] must be 14 or 15 digits"))
		})

		It("Should return an error if Rules.Source.GTP.IMSIs contains a value "+
			"that is < 14 digits", func() {
			tp.Rules[0].Source.GTP.IMSIs[0] = "1234567890123"
			Expect(tp.Validate()).To(MatchError(
				"rules[0].source.gtp.imsis[0] must be 14 or 15 digits"))
		})

		It("Should return an error if Rules.Source.GTP.IMSIs contains a value "+
			"that is > 15 digits", func() {
			tp.Rules[0].Source.GTP.IMSIs[0] = "1234567890123456"
			Expect(tp.Validate()).To(MatchError(
				"rules[0].source.gtp.imsis[0] must be 14 or 15 digits"))
		})

		It("Should return an error if Rules.Target.Description is "+
			"empty", func() {
			tp.Rules[0].Target.Description = ""
			Expect(tp.Validate()).To(MatchError(
				"rules[0].target.description cannot be empty"))
		})

		It("Should return an error if Rules.Target.Action is not one of "+
			"[accept, reject, drop]", func() {
			tp.Rules[0].Target.Action = "abc"
			Expect(tp.Validate()).To(MatchError(
				"rules[0].target.action must be one of [accept, reject, drop]"))
		})

		It("Should return an error if Rules.Target.MAC.MACAddress is "+
			"invalid", func() {
			tp.Rules[0].Target.MAC.MACAddress = "abc-98-deg"
			Expect(tp.Validate()).To(MatchError(
				"rules[0].target.mac.mac_address could not be parsed " +
					"(address abc-98-deg: invalid MAC address)"))
		})

		It("Should return an error if Rules.Target.IP.Address is "+
			"invalid", func() {
			tp.Rules[0].Target.IP.Address = "123"
			Expect(tp.Validate()).To(MatchError(
				"rules[0].target.ip.address could not be parsed"))
		})

		It("Should return an error if Rules.Target.IP.Port is < 1", func() {
			tp.Rules[0].Target.IP.Port = 0
			Expect(tp.Validate()).To(MatchError(
				"rules[0].target.ip.port must be in [1..65536]"))
		})

		It("Should return an error if Rules.Target.IP.Port is > 65536", func() {
			tp.Rules[0].Target.IP.Port = 65537
			Expect(tp.Validate()).To(MatchError(
				"rules[0].target.ip.port must be in [1..65536]"))
		})
	})

	Describe("String", func() {
		It("Should return the string value", func() {
			Expect(tp.String()).To(Equal(strings.TrimSpace(`
TrafficPolicy[
    ID: 9d740cee-035f-4076-847c-d1c80cdf19db,
    Rules: [
        TrafficRule[
            Description: test-rule-1
            Priority: 1
            Source: TrafficSelector[
                Description: source-source-1
                MACs: MACFilter[
                    MACAddresses: [
                        F0-59-8E-7B-36-8A
                        23-20-8E-15-89-D1
                        35-A4-38-73-35-45
                    ]
                ]
                IP: IPFilter[
                    Address: 223.1.1.0
                    Mask: 16
                    BeginPort: 2000
                    EndPort: 2012
                    Protocol: tcp
                ]
                GTP: GTPFilter[
                    Address: 10.6.7.2
                    Mask: 12
                    IMSIs: [
                        310150123456789
                        310150123456790
                        310150123456791
                    ]
                ]
            ]
            Destination: TrafficSelector[
                Description: test-destination-1
                MACs: MACFilter[
                    MACAddresses: [
                        7D-C2-3A-1C-63-D9
                        E9-6B-D1-D2-1A-6B
                        C8-32-A9-43-85-55
                    ]
                ]
                IP: IPFilter[
                    Address: 64.1.1.0
                    Mask: 16
                    BeginPort: 1000
                    EndPort: 1012
                    Protocol: tcp
                ]
                GTP: GTPFilter[
                    Address: 108.6.7.2
                    Mask: 4
                    IMSIs: [
                        310150123456792
                        310150123456793
                        310150123456794
                    ]
                ]
            ]
            Target: TrafficTarget[
                Description: test-target-1
                Action: accept
                MAC: MACModifier[
                    MACAddress: C7-5A-E7-98-1B-A3
                ]
                IP: IPModifier[
                    Address: 123.2.3.4
                    Port: 1600
                ]
            ]
        ]
        TrafficRule[
            Description: test-rule-2
            Priority: 2
            Source: TrafficSelector[
                Description: source-source-2
                MACs: MACFilter[
                    MACAddresses: [
                        43-78-01-EE-B5-8D
                        DB-C6-F2-CC-0B-31
                        66-69-C2-D8-78-83
                    ]
                ]
                IP: IPFilter[
                    Address: 12.1.1.0
                    Mask: 6
                    BeginPort: 5000
                    EndPort: 5012
                    Protocol: udp
                ]
                GTP: GTPFilter[
                    Address: 10.66.7.2
                    Mask: 4
                    IMSIs: [
                        310150123456795
                        310150123456796
                        310150123456797
                    ]
                ]
            ]
            Destination: TrafficSelector[
                Description: test-destination-2
                MACs: MACFilter[
                    MACAddresses: [
                        30-50-D5-50-47-90
                        14-7C-F2-7B-DC-73
                        40-D2-CA-98-CA-CE
                    ]
                ]
                IP: IPFilter[
                    Address: 132.15.16.20
                    Mask: 3
                    BeginPort: 55000
                    EndPort: 55012
                    Protocol: udp
                ]
                GTP: GTPFilter[
                    Address: 70.6.17.2
                    Mask: 4
                    IMSIs: [
                        310150123456798
                        310150123456799
                        310150123456800
                    ]
                ]
            ]
            Target: TrafficTarget[
                Description: test-target-2
                Action: reject
                MAC: MACModifier[
                    MACAddress: 28-3F-D2-2C-47-1F
                ]
                IP: IPModifier[
                    Address: 242.25.31.14
                    Port: 2600
                ]
            ]
        ]
    ]
]`,
			)))
		})
	})
})
