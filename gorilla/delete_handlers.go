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

	cce "github.com/smartedgemec/controller-ce"
)

func handleDeleteNodesApps(ctx context.Context, ps cce.PersistenceService, e cce.Persistable) error {
	app, err := ps.Read(
		ctx,
		e.(*cce.NodeApp).AppID,
		&cce.App{})
	if err != nil {
		return err
	}

	nodeCC, err := connectNode(ctx, ps, e.(*cce.NodeApp))
	if err != nil {
		return err
	}
	log.Debug(nodeCC.Node)

	// if kubernetes un-deploy application
	ctrl := getController(ctx)

	if ctrl.OrchestrationMode == cce.OrchestrationModeKubernetes {
		if err = ctrl.KubernetesClient.Undeploy(
			ctx,
			e.(*cce.NodeApp).NodeID,
			e.(*cce.NodeApp).AppID,
		); err != nil {
			return err
		}
	}

	if err = nodeCC.AppDeploySvcCli.Undeploy(ctx, app.GetID()); err != nil {
		return err
	}
	return nil
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

	nodeCC, err := connectNode(ctx, ps, e.(*cce.NodeDNSConfig))
	if err != nil {
		return err
	}

	log.Debug(nodeCC.Node)

	for _, aRecord := range dnsConfig.(*cce.DNSConfig).ARecords {
		if err := nodeCC.DNSSvcCli.DeleteA(ctx, aRecord); err != nil {
			return err
		}
	}

	if err := nodeCC.DNSSvcCli.DeleteForwarders(ctx, dnsConfig.(*cce.DNSConfig).Forwarders); err != nil {
		return err
	}

	return nil
}

func handleDeleteNodesAppsTrafficPolicies(
	ctx context.Context,
	ps cce.PersistenceService,
	e cce.Persistable,
) error {
	nodeApp, err := ps.Read(ctx, e.(*cce.NodeAppTrafficPolicy).NodeAppID, &cce.NodeApp{})
	if err != nil {
		return err
	}
	log.Debugf("Loaded node app %s: %+v", nodeApp.GetID(), nodeApp)

	nodeCC, err := connectNode(ctx, ps, nodeApp.(*cce.NodeApp))
	if err != nil {
		return err
	}

	log.Debugf("Connection to node established: %+v", nodeCC.Node)

	if err := nodeCC.AppPolicySvcCli.Delete(ctx, nodeApp.(*cce.NodeApp).AppID); err != nil {
		return err
	}

	return nil
}
