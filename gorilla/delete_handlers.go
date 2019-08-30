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

package gorilla

import (
	"context"

	cce "github.com/otcshare/edgecontroller"
)

func handleDeleteNodesApps(ctx context.Context, ps cce.PersistenceService, e cce.Persistable) error {
	app, err := ps.Read(
		ctx,
		e.(*cce.NodeApp).AppID,
		&cce.App{})
	if err != nil {
		return err
	}

	ctrl := getController(ctx)
	nodePort := ctrl.EVAPort
	if nodePort == "" {
		nodePort = defaultEVAPort
	}
	nodeCC, err := connectNode(ctx, ps, e.(*cce.NodeApp), nodePort, ctrl.EdgeNodeCreds)
	if err != nil {
		return err
	}

	// if kubernetes un-deploy application
	if ctrl.OrchestrationMode == cce.OrchestrationModeKubernetes ||
		ctrl.OrchestrationMode == cce.OrchestrationModeKubernetesOVN {
		if err = ctrl.KubernetesClient.Undeploy(
			ctx,
			e.(*cce.NodeApp).NodeID,
			e.(*cce.NodeApp).AppID,
		); err != nil {
			return err
		}
	}

	return nodeCC.AppDeploySvcCli.Undeploy(ctx, app.GetID())
}

func handleDeleteNodesDNSConfigs(
	ctx context.Context,
	ps cce.PersistenceService,
	e cce.Persistable,
) error {
	dnsConfig, err := ps.Read(ctx, e.(*cce.NodeDNSConfig).DNSConfigID, &cce.DNSConfig{})
	if err != nil {
		return err
	}
	log.Debugf("Loaded DNS Config %s\n%+v", dnsConfig.GetID(), dnsConfig)

	ctrl := getController(ctx)
	nodePort := ctrl.ELAPort
	if nodePort == "" {
		nodePort = defaultELAPort
	}
	nodeCC, err := connectNode(ctx, ps, e.(*cce.NodeDNSConfig), nodePort, ctrl.EdgeNodeCreds)
	if err != nil {
		return err
	}

	for _, aRecord := range dnsConfig.(*cce.DNSConfig).ARecords {
		if err := nodeCC.DNSSvcCli.DeleteA(ctx, aRecord); err != nil {
			return err
		}
	}

	return nodeCC.DNSSvcCli.DeleteForwarders(ctx, dnsConfig.(*cce.DNSConfig).Forwarders)
}

func handleDeleteNodesDNSConfigsWithAliases(
	ctx context.Context,
	ps cce.PersistenceService,
	nodeDNS cce.Persistable,
	dnsConfig cce.Persistable,
	dnsAliases []cce.Persistable,
) error {
	ctrl := getController(ctx)
	nodePort := ctrl.ELAPort
	if nodePort == "" {
		nodePort = defaultELAPort
	}
	nodeCC, err := connectNode(ctx, ps, nodeDNS.(*cce.NodeDNSConfig), nodePort, ctrl.EdgeNodeCreds)
	if err != nil {
		return err
	}

	for _, alias := range dnsAliases {
		record := &cce.DNSARecord{
			Name:        alias.(*cce.DNSConfigAppAlias).AppID,
			Description: alias.(*cce.DNSConfigAppAlias).Description,
			IPs:         []string{alias.(*cce.DNSConfigAppAlias).AppID},
		}

		if err := nodeCC.DNSSvcCli.DeleteA(ctx, record); err != nil {
			return err
		}
	}

	for _, aRecord := range dnsConfig.(*cce.DNSConfig).ARecords {
		if err := nodeCC.DNSSvcCli.DeleteA(ctx, aRecord); err != nil {
			return err
		}
	}

	if len(dnsConfig.(*cce.DNSConfig).Forwarders) != 0 {
		if err := nodeCC.DNSSvcCli.DeleteForwarders(ctx, dnsConfig.(*cce.DNSConfig).Forwarders); err != nil {
			return err
		}
	}

	return nil
}
