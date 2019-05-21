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
	"log"

	cce "github.com/smartedgemec/controller-ce"
)

func handleDeleteNodesApps(
	ctx context.Context,
	ps cce.PersistenceService,
	e cce.Persistable,
) error {
	nodeCC, err := connectNode(ctx, ps, e.(*cce.NodeApp))
	if err != nil {
		return err
	}

	log.Println(nodeCC.Node)

	if err := nodeCC.AppDeploySvcCli.Undeploy(ctx, e.(*cce.NodeApp).AppID); err != nil {
		return err
	}

	return nil
}

func handleDeleteNodesVNFs(
	ctx context.Context,
	ps cce.PersistenceService,
	e cce.Persistable,
) error {
	nodeCC, err := connectNode(ctx, ps, e.(*cce.NodeVNF))
	if err != nil {
		return err
	}

	log.Println(nodeCC.Node)

	if err := nodeCC.VNFDeploySvcCli.Undeploy(ctx, e.(*cce.NodeVNF).VNFID); err != nil {
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
	log.Printf("Loaded DNS Config %s", dnsConfig.GetID())
	log.Println(dnsConfig)

	nodeCC, err := connectNode(ctx, ps, e.(*cce.NodeDNSConfig))
	if err != nil {
		return err
	}

	log.Println(nodeCC.Node)

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
	log.Printf("Loaded node app %s", nodeApp.GetID())
	log.Println(nodeApp)

	nodeCC, err := connectNode(ctx, ps, nodeApp.(*cce.NodeApp))
	if err != nil {
		return err
	}

	log.Println("Connection to node established:", nodeCC.Node)

	if err := nodeCC.AppPolicySvcCli.Delete(ctx, nodeApp.(*cce.NodeApp).AppID); err != nil {
		return err
	}

	return nil
}

func handleDeleteNodesVNFsTrafficPolicies(
	ctx context.Context,
	ps cce.PersistenceService,
	e cce.Persistable,
) error {
	nodeVNF, err := ps.Read(ctx, e.(*cce.NodeVNFTrafficPolicy).NodeVNFID, &cce.NodeVNF{})
	if err != nil {
		return err
	}
	log.Printf("Loaded node VNF %s", nodeVNF.GetID())
	log.Println(nodeVNF)

	nodeCC, err := connectNode(ctx, ps, nodeVNF.(*cce.NodeVNF))
	if err != nil {
		return err
	}

	log.Println("Connection to node established:", nodeCC.Node)

	// TODO there is currently no VNFPolicyService in https://github.com/smartedgemec/schema/blob/master/pb/ela.proto
	// if err := nodeCC.VNFPolicySvcCli.Delete(ctx, nodeVNF.(*cce.NodeVNF).VNFID); err != nil {
	// 	return err
	// }

	return nil
}
