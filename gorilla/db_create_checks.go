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
	"fmt"
	"net/http"

	cce "github.com/smartedgemec/controller-ce"
)

func checkDBCreateNodesApps(
	ctx context.Context,
	ps cce.PersistenceService,
	e cce.Persistable,
) (statusCode int, err error) {
	var es []cce.Persistable

	if es, err = ps.Filter(
		ctx,
		&cce.NodeApp{},
		[]cce.Filter{
			{
				Field: "node_id",
				Value: e.(*cce.NodeApp).NodeID,
			},
			{
				Field: "app_id",
				Value: e.(*cce.NodeApp).AppID,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) != 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"duplicate record in %s detected for node_id %s and app_id %s",
			e.(*cce.NodeApp).GetTableName(),
			e.(*cce.NodeApp).NodeID,
			e.(*cce.NodeApp).AppID)
	}

	return 0, nil
}

func checkDBCreateDNSConfigsAppAliases(
	ctx context.Context,
	ps cce.PersistenceService,
	e cce.Persistable,
) (statusCode int, err error) {
	var es []cce.Persistable

	if es, err = ps.Filter(
		ctx,
		&cce.DNSConfigAppAlias{},
		[]cce.Filter{
			{
				Field: "dns_config_id",
				Value: e.(*cce.DNSConfigAppAlias).DNSConfigID,
			},
			{
				Field: "app_id",
				Value: e.(*cce.DNSConfigAppAlias).AppID,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) != 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"duplicate record in %s detected for dns_config_id %s and app_id %s",
			e.(*cce.DNSConfigAppAlias).GetTableName(),
			e.(*cce.DNSConfigAppAlias).DNSConfigID,
			e.(*cce.DNSConfigAppAlias).AppID)
	}

	return 0, nil
}

func checkDBCreateNodesDNSConfigs(
	ctx context.Context,
	ps cce.PersistenceService,
	e cce.Persistable,
) (statusCode int, err error) {
	var es []cce.Persistable

	// the nodes_dns_configs table has a unique constraint on node_id so we don't filter on dns_config_id
	if es, err = ps.Filter(
		ctx,
		&cce.NodeDNSConfig{},
		[]cce.Filter{
			{
				Field: "node_id",
				Value: e.(*cce.NodeDNSConfig).NodeID,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) != 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"duplicate record in %s detected for node_id %s",
			e.(*cce.NodeDNSConfig).GetTableName(),
			e.(*cce.NodeDNSConfig).NodeID)
	}

	return 0, nil
}

func checkDBCreateNodesAppsTrafficPolicies(
	ctx context.Context,
	ps cce.PersistenceService,
	e cce.Persistable,
) (statusCode int, err error) {
	var es []cce.Persistable

	// the nodes_dns_configs table has a unique constraint on node_id so we don't filter on dns_config_id
	if es, err = ps.Filter(
		ctx,
		&cce.NodeAppTrafficPolicy{},
		[]cce.Filter{
			{
				Field: "nodes_apps_id",
				Value: e.(*cce.NodeAppTrafficPolicy).NodeAppID,
			},
			{
				Field: "traffic_policy_id",
				Value: e.(*cce.NodeAppTrafficPolicy).TrafficPolicyID,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) != 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"duplicate record in %s detected for nodes_apps_id %s and "+
				"traffic_policy_id %s",
			e.(*cce.NodeAppTrafficPolicy).GetTableName(),
			e.(*cce.NodeAppTrafficPolicy).NodeAppID,
			e.(*cce.NodeAppTrafficPolicy).TrafficPolicyID)
	}

	return 0, nil
}
