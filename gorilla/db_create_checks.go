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

func checkDBCreateDNSConfigsAppAliases(
	ctx context.Context,
	ps cce.PersistenceService,
	e cce.Entity,
) (statusCode int, err error) {
	var es []cce.Entity

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
			"duplicate record detected for dns_config_id %s and app_id %s",
			e.(*cce.DNSConfigAppAlias).DNSConfigID,
			e.(*cce.DNSConfigAppAlias).AppID)
	}

	return 0, nil
}

func checkDBCreateDNSConfigsVNFAliases(
	ctx context.Context,
	ps cce.PersistenceService,
	e cce.Entity,
) (statusCode int, err error) {
	var es []cce.Entity

	if es, err = ps.Filter(
		ctx,
		&cce.DNSConfigVNFAlias{},
		[]cce.Filter{
			{
				Field: "dns_config_id",
				Value: e.(*cce.DNSConfigVNFAlias).DNSConfigID,
			},
			{
				Field: "vnf_id",
				Value: e.(*cce.DNSConfigVNFAlias).VNFID,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) != 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"duplicate record detected for dns_config_id %s and vnf_id %s",
			e.(*cce.DNSConfigVNFAlias).DNSConfigID,
			e.(*cce.DNSConfigVNFAlias).VNFID)
	}

	return 0, nil
}

func checkDBCreateNodeDNSConfigs(
	ctx context.Context,
	ps cce.PersistenceService,
	e cce.Entity,
) (statusCode int, err error) {
	var es []cce.Entity

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
			"duplicate record detected for node_id %s and "+
				"dns_config_id %s",
			e.(*cce.NodeDNSConfig).NodeID,
			e.(*cce.NodeDNSConfig).DNSConfigID)
	}

	return 0, nil
}
