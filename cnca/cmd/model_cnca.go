// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package cnca

// Header holds version & type of YAML configuration
type Header struct {
	Version string `yaml:"apiVersion"`
	//  ngc: 5G Traffic Influence Subscription, or
	//  lte: LTE CUPS Userplane
	Kind string `yaml:"kind"`
}

// AFTrafficInfluSub describes NGC AF Traffic Influence Subscription
type AFTrafficInfluSub struct {
	H      Header
	Policy struct {
		// Identifies a service on behalf of which the AF is issuing the request.
		AfServiceID string `yaml:"afServiceId,omitempty"`
		// Identifies an application.
		AfAppID string `yaml:"afAppId,omitempty"`
		// Identifies an NEF Northbound interface transaction, generated by the AF.
		AfTransID string `yaml:"afTransId,omitempty"`
		// Identifies whether an application can be relocated once a location of the application has been selected.
		AppReloInd bool `yaml:"appReloInd,omitempty"`
		// Identifies data network name
		DNN string `yaml:"dnn,omitempty"`
		// Snssai
		SNSSAI struct {
			SST int32  `yaml:"sst"`
			SD  string `yaml:"sd,omitempty"`
		} `yaml:"snssai,omitempty"`
		// string containing a local identifier followed by \"@\" and a domain identifier. Both the local identifier and the domain identifier shall be encoded as strings that do not contain any \"@\" characters. See Clauses 4.6.2 and 4.6.3 of 3GPP TS 23.682 for more information.
		ExternalGroupID string `yaml:"externalGroupId,omitempty"`
		// Identifies whether the AF request applies to any UE.
		AnyUeInd bool `yaml:"anyUeInd,omitempty"`
		// Identifies the requirement to be notified of the event(s):
		// - UP_PATH_CHANGE
		SubscribedEvents []string `yaml:"subscribedEvents,omitempty"`
		// Gpsi
		GPSI string `yaml:"gpsi,omitempty"`
		// string identifying a Ipv4 address formatted in the \"dotted decimal\" notation as defined in IETF RFC 1166.
		IPv4Addr string `yaml:"ipv4Addr,omitempty"`
		// string identifying a Ipv6 address formatted according to clause 4 in IETF RFC 5952.
		IPv6Addr string `yaml:"ipv6Addr,omitempty"`
		// string identifying MAC Address
		MACAddr string `yaml:"macAddr,omitempty"`
		// Identifies the type of notification regarding UP path management event. Possible values are EARLY - early notification of UP path reconfiguration. EARLY_LATE - early and late notification of UP path reconfiguration. This value shall only be present in the subscription to the DNAI change event. LATE - late notification of UP path reconfiguration.
		DNAIChgType string `yaml:"dnaiChgType,omitempty"`
		// string formatted according to IETF RFC 3986 identifying a referenced resource.
		NotificationDestination string `yaml:"notificationDestination,omitempty"`
		// Set to true by the AF to request the NEF to send a test notification. Set to false or omitted otherwise.
		RequestTestNotification bool `yaml:"requestTestNotification,omitempty"`
		// WebsockNotifConfig
		WebsockNotifConfig struct {
			WebsocketURI        string `yaml:"websocketUri,omitempty"`
			RequestWebsocketURI bool   `yaml:"requestWebsocketUri,omitempty"`
		} `yaml:"websockNotifConfig,omitempty"`
		// Identifies IP packet filters.
		TrafficFilters []struct {
			// Indicates the IP flow.
			FlowID int32 `yaml:"flowId"`
			// Indicates the packet filters of the IP flow. Refer to subclause 5.3.8 of 3GPP TS 29.214 for encoding. It shall contain UL and/or DL IP flow description.
			FlowDescriptions []string `yaml:"flowDescriptions,omitempty"`
		} `yaml:"trafficFilters,omitempty"`
		// Identifies Ethernet packet filters.
		EthTrafficFilters []struct {
			DestMACAddr string `yaml:"destMacAddr,omitempty"`
			EthType     string `yaml:"ethType"`
			// Defines a packet filter of an IP flow.
			FDesc string `yaml:"fDesc,omitempty"`
			// Possible values are :
			// - DOWNLINK - The corresponding filter applies for traffic to the UE.
			// - UPLINK - The corresponding filter applies for traffic from the UE.
			// - BIDIRECTIONAL The corresponding filter applies for traffic both to and from the UE.
			// UNSPECIFIED - The corresponding filter applies for traffic to the UE (downlink), but has no specific direction declared. The service data flow detection shall apply the filter for uplink traffic as if the filter was bidirectional.
			FDir          string   `yaml:"fDir,omitempty"`
			SourceMACAddr string   `yaml:"sourceMacAddr,omitempty"`
			VLANTags      []string `yaml:"vlanTags,omitempty"`
		} `yaml:"ethTrafficFilters,omitempty"`
		// Identifies the N6 traffic routing requirement.
		TrafficRoutes []struct {
			DNAI      string `yaml:"dnai"`
			RouteInfo struct {
				IPv4Addr   string `yaml:"ipv4Addr,omitempty"`
				IPv6Addr   string `yaml:"ipv6Addr,omitempty"`
				PortNumber int32  `yaml:"portNumber"`
			} `yaml:"routeInfo,omitempty"`
			RouteProfID string `yaml:"routeProfId,omitempty"`
		} `yaml:"trafficRoutes,omitempty"`
		// Indicates the time interval(s) during which the AF request is to be applied
		TempValidities []struct {
			// string with format \"date-time\" as defined in OpenAPI.
			StartTime string `yaml:"startTime,omitempty"`
			// string with format \"date-time\" as defined in OpenAPI.
			StopTime string `yaml:"stopTime,omitempty"`
		} `yaml:"tempValidities,omitempty"`
		// Identifies a geographic zone that the AF request applies only to the traffic of UE(s) located in this specific zone.
		ValidGeoZoneIds []string `yaml:"validGeoZoneIds,omitempty"`
	} `yaml:"policy"`
}

// LTEUserplane describes LTE Userplane Configuration
type LTEUserplane struct {
	H      Header
	Policy struct {
		ID       string `yaml:"id,omitempty"`
		UUID     string `yaml:"uuid,omitempty"`
		Function string `yaml:"function,omitempty"`
		Config   struct {
			Sxa      LteConfigInfoCpup `yaml:"sxa,omitempty"`
			Sxb      LteConfigInfoCpup `yaml:"sxb,omitempty"`
			S1u      LteConfigInfoUp   `yaml:"s1u,omitempty"`
			S5uSGW   LteConfigInfoUp   `yaml:"s5u_sgw,omitempty"`
			S5uPGW   LteConfigInfoUp   `yaml:"s5u_pgw,omitempty"`
			SGi      LteConfigInfoUp   `yaml:"sgi,omitempty"`
			Breakout []LteConfigInfoUp `yaml:"breakout,omitempty"`
			DNS      []LteConfigInfoUp `yaml:"dns,omitempty"`
		} `yaml:"config,omitempty"`

		Selectors []struct {
			ID      string `yaml:"id,omitempty"`
			Network struct {
				MCC string `yaml:"mcc,omitempty"`
				MNC string `yaml:"mnc,omitempty"`
			} `yaml:"network,omitempty"`
			ULI struct {
				TAI struct {
					// Tracking area code (TAC), which is typically an unsigned integer
					// from 1 to 2^16, inclusive.
					TAC int64 `yaml:"tac,omitempty"`
				} `yaml:"tai,omitempty"`
				ECGI struct {
					// E-UTRAN cell identifier (ECI), which is typically an unsigned integer
					// from 1 to 2^32, inclusive.
					ECI int64 `yaml:"eci,omitempty"`
				} `yaml:"ecgi,omitempty"`
			} `yaml:"uli,omitempty"`
			PDN struct {
				APNs []string `yaml:"apns,omitempty"`
			} `yaml:"pdn,omitempty"`
		} `yaml:"selectors,omitempty"`

		// The UEs that should be entitled to access privileged networks via this
		// userplane.  Note: UEs not in this list will still be able to get a bearer
		// via the userplane. The UEs in this list are just for entitlement
		// purposes. (optional)
		Entitlements []struct {
			ID    string   `yaml:"id,omitempty"`
			APNs  []string `yaml:"apns,omitempty"`
			IMSIs []struct {
				Begin string `yaml:"begin,omitempty"`
				End   string `yaml:"end,omitempty"`
			} `yaml:"imsis,omitempty"`
		} `yaml:"entitlements,omitempty"`
	}
}

// LteConfigInfoCpup Information that the userplane should configure, which
// relates to the control plane (CP) side and the user plane (UP) side.
type LteConfigInfoCpup struct {
	CpIPAddress string `yaml:"cp_ip_address,omitempty"`
	UpIPAddress string `yaml:"up_ip_address,omitempty"`
}

// LteConfigInfoUp Information that the userplane should configure, which
// relates to the user plane (UP) side only.
type LteConfigInfoUp struct {
	UpIPAddress string `yaml:"up_ip_address,omitempty"`
}

// AFPfdManagement describes NGC AF Pfd Transaction
type AFPfdManagement struct {
	H      Header
	Policy struct {
		// Link to the resource "Individual PFD Management Transaction".
		// This parameter shall be supplied by the AF in HTTP responses.
		Self string `yaml:"self,omitempty"`
		// String identifying supported optional features of PFD Management
		// This attribute shall be provided in the POST request and in the
		// response of successful resource creation.
		SuppFeat *string `yaml:"suppFeat,omitempty"`
		// Each element uniquely identifies the PFDs for an external application
		// identifier. Each element is identified in the map via an external
		// application identifier as key. The response shall include successfully
		// provisioned PFD data of application(s).
		PfdDatas []struct {
			// Each element uniquely identifies external application identifier
			ExternalAppID string `yaml:"externalAppID"`
			// Link to the resource. This parameter shall be supplied by the AF in
			// HTTP responses that include an object of PfdData type
			Self string `yaml:"self,omitempty"`
			// Contains the PFDs of the external application identifier. Each PFD is
			// identified in the map via a key containing the PFD identifier.
			Pfds []struct {
				// Identifies a PDF of an application identifier.
				PfdID string `yaml:"pfdID"`
				// Represents a 3-tuple with protocol, server ip and server port for UL/DL
				// application traffic. The content of the string has the same encoding as
				// the IPFilterRule AVP value as defined in IETFÂ RFCÂ 6733.
				FlowDescriptions []string `yaml:"flowDescriptions,omitempty"`
				// Indicates a URL or a regular expression which is used to match the
				// significant parts of the URL.
				Urls []string `yaml:"urls,omitempty"`
				// Indicates an FQDN or a regular expression as a domain name matching
				// criteria.
				DomainNames []string `yaml:"domainNames,omitempty"`
			} `yaml:"pfds"`
			// Indicates that the list of PFDs in this request should be deployed
			// within the time interval indicated by the Allowed Delay
			AllowedDelay *uint64 `yaml:"allowedDelay,omitempty"`
			// SCEF supplied property, inclusion of this property means the allowed
			// delayed cannot be satisfied, i.e. it is smaller than the caching time,
			// but the PFD data is still stored.
			CachingTime *uint64 `yaml:"cachingTime,omitempty"`
		} `yaml:"pfdDatas"`
		// Supplied by the AF and contains the external application identifiers
		// for which PFD(s) are not added or modified successfully. The failure
		// reason is also included. Each element provides the related information
		// for one or more external application identifier(s) and is identified in
		// the map via the failure identifier as key.
		PfdReports map[string]PfdReport `yaml:"pfdReports,omitempty"`
	} `yaml:"policy"`
}

// AFPfdData describes NGC AF Pfd Application
type AFPfdData struct {
	H      Header
	Policy struct {
		// Each element uniquely identifies external application identifier
		ExternalAppID string `yaml:"externalAppID"`
		// Link to the resource. This parameter shall be supplied by the AF in
		// HTTP responses that include an object of PfdData type
		Self string `yaml:"self,omitempty"`
		// Contains the PFDs of the external application identifier. Each PFD is
		// identified in the map via a key containing the PFD identifier.
		Pfds []struct {
			// Identifies a PDF of an application identifier.
			PfdID string `yaml:"pfdID"`
			// Represents a 3-tuple with protocol, server ip and server port for UL/DL
			// application traffic. The content of the string has the same encoding as
			// the IPFilterRule AVP value as defined in IETFÂ RFCÂ 6733.
			FlowDescriptions []string `yaml:"flowDescriptions,omitempty"`
			// Indicates a URL or a regular expression which is used to match the
			// significant parts of the URL.
			Urls []string `yaml:"urls,omitempty"`
			// Indicates an FQDN or a regular expression as a domain name matching
			// criteria.
			DomainNames []string `yaml:"domainNames,omitempty"`
		} `yaml:"pfds"`
		// Indicates that the list of PFDs in this request should be deployed
		// within the time interval indicated by the Allowed Delay
		AllowedDelay *uint64 `yaml:"allowedDelay,omitempty"`
		// SCEF supplied property, inclusion of this property means the allowed
		// delayed cannot be satisfied, i.e. it is smaller than the caching time,
		// but the PFD data is still stored.
		CachingTime *uint64 `yaml:"cachingTime,omitempty"`
	} `yaml:"policy"`
}
