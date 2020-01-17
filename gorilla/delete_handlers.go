// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

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
	defer disconnectNode(nodeCC)

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

	err = nodeCC.AppDeploySvcCli.Undeploy(ctx, app.GetID())

	return err
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
	defer disconnectNode(nodeCC)

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
