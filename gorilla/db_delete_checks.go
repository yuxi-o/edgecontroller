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

func checkDBDeleteNodes(
	ctx context.Context,
	ps cce.PersistenceService,
	id string,
) (statusCode int, err error) {
	var es []cce.Persistable

	if es, err = ps.Filter(
		ctx,
		&cce.NodeApp{},
		[]cce.Filter{
			{
				Field: "node_id",
				Value: id,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) > 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"cannot delete node_id %s: record in use in nodes_apps",
			id)
	}

	if es, err = ps.Filter(
		ctx,
		&cce.NodeDNSConfig{},
		[]cce.Filter{
			{
				Field: "node_id",
				Value: id,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) > 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"cannot delete node_id %s: record in use in nodes_dns_configs",
			id)
	}

	return 0, nil
}

func checkDBDeleteApps(
	ctx context.Context,
	ps cce.PersistenceService,
	id string,
) (statusCode int, err error) {
	var es []cce.Persistable

	if es, err = ps.Filter(
		ctx,
		&cce.DNSConfigAppAlias{},
		[]cce.Filter{
			{
				Field: "app_id",
				Value: id,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) > 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"cannot delete app_id %s: record in use in dns_configs_app_aliases",
			id)
	}

	if es, err = ps.Filter(
		ctx,
		&cce.NodeApp{},
		[]cce.Filter{
			{
				Field: "app_id",
				Value: id,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) > 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"cannot delete app_id %s: record in use in nodes_apps",
			id)
	}

	return 0, nil
}

func checkDBDeleteNodesApps(
	ctx context.Context,
	ps cce.PersistenceService,
	id string,
) (statusCode int, err error) {
	var es []cce.Persistable

	if es, err = ps.Filter(
		ctx,
		&cce.NodeAppTrafficPolicy{},
		[]cce.Filter{
			{
				Field: "nodes_apps_id",
				Value: id,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) > 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"cannot delete node_app_id %s: record in use in nodes_apps_traffic_policies",
			id)
	}

	return 0, nil
}

func checkDBDeleteTrafficPolicies(
	ctx context.Context,
	ps cce.PersistenceService,
	id string,
) (statusCode int, err error) {
	var es []cce.Persistable

	if es, err = ps.Filter(
		ctx,
		&cce.NodeAppTrafficPolicy{},
		[]cce.Filter{
			{
				Field: "traffic_policy_id",
				Value: id,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) > 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"cannot delete traffic_policy_id %s: record in use in "+
				"nodes_apps_traffic_policies",
			id)
	}

	return 0, nil
}

func checkDBDeleteDNSConfigs(
	ctx context.Context,
	ps cce.PersistenceService,
	id string,
) (statusCode int, err error) {
	var es []cce.Persistable

	if es, err = ps.Filter(
		ctx,
		&cce.DNSConfigAppAlias{},
		[]cce.Filter{
			{
				Field: "dns_config_id",
				Value: id,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) > 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"cannot delete dns_config_id %s: record in use in "+
				"dns_configs_app_aliases",
			id)
	}

	if es, err = ps.Filter(
		ctx,
		&cce.NodeDNSConfig{},
		[]cce.Filter{
			{
				Field: "dns_config_id",
				Value: id,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) > 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"cannot delete dns_config_id %s: record in use in "+
				"nodes_dns_configs",
			id)
	}

	return 0, nil
}
