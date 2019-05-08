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

func checkDBCreateDNSConfigsDNSContainerAppAliases(
	ctx context.Context,
	ps cce.PersistenceService,
	e cce.Entity,
) (statusCode int, err error) {
	var es []cce.Entity

	if es, err = ps.Filter(
		ctx,
		&cce.DNSConfigDNSContainerAppAlias{},
		[]cce.Filter{
			{
				Field: "dns_config_id",
				Value: e.(*cce.DNSConfigDNSContainerAppAlias).DNSConfigID,
			},
			{
				Field: "dns_container_app_alias_id",
				Value: e.(*cce.DNSConfigDNSContainerAppAlias).DNSContainerAppAliasID, //nolint:lll
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) != 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"duplicate record detected for dns_config_id %s and "+
				"dns_container_app_alias_id %s",
			e.(*cce.DNSConfigDNSContainerAppAlias).DNSConfigID,
			e.(*cce.DNSConfigDNSContainerAppAlias).DNSContainerAppAliasID)
	}

	return 0, nil
}

func checkDBCreateDNSConfigsDNSContainerVNFAliases(
	ctx context.Context,
	ps cce.PersistenceService,
	e cce.Entity,
) (statusCode int, err error) {
	var es []cce.Entity

	if es, err = ps.Filter(
		ctx,
		&cce.DNSConfigDNSContainerVNFAlias{},
		[]cce.Filter{
			{
				Field: "dns_config_id",
				Value: e.(*cce.DNSConfigDNSContainerVNFAlias).DNSConfigID,
			},
			{
				Field: "dns_container_vnf_alias_id",
				Value: e.(*cce.DNSConfigDNSContainerVNFAlias).DNSContainerVNFAliasID, //nolint:lll
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) != 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"duplicate record detected for dns_config_id %s and "+
				"dns_container_vnf_alias_id %s",
			e.(*cce.DNSConfigDNSContainerVNFAlias).DNSConfigID,
			e.(*cce.DNSConfigDNSContainerVNFAlias).DNSContainerVNFAliasID)
	}

	return 0, nil
}

func checkDBCreateDNSConfigsDNSVMAppAliases(
	ctx context.Context,
	ps cce.PersistenceService,
	e cce.Entity,
) (statusCode int, err error) {
	var es []cce.Entity

	if es, err = ps.Filter(
		ctx,
		&cce.DNSConfigDNSVMAppAlias{},
		[]cce.Filter{
			{
				Field: "dns_config_id",
				Value: e.(*cce.DNSConfigDNSVMAppAlias).DNSConfigID,
			},
			{
				Field: "dns_vm_app_alias_id",
				Value: e.(*cce.DNSConfigDNSVMAppAlias).DNSVMAppAliasID,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) != 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"duplicate record detected for dns_config_id %s and "+
				"dns_vm_app_alias_id %s",
			e.(*cce.DNSConfigDNSVMAppAlias).DNSConfigID,
			e.(*cce.DNSConfigDNSVMAppAlias).DNSVMAppAliasID)
	}

	return 0, nil
}

func checkDBCreateDNSConfigsDNSVMVNFAliases(
	ctx context.Context,
	ps cce.PersistenceService,
	e cce.Entity,
) (statusCode int, err error) {
	var es []cce.Entity

	if es, err = ps.Filter(
		ctx,
		&cce.DNSConfigDNSVMVNFAlias{},
		[]cce.Filter{
			{
				Field: "dns_config_id",
				Value: e.(*cce.DNSConfigDNSVMVNFAlias).DNSConfigID,
			},
			{
				Field: "dns_vm_vnf_alias_id",
				Value: e.(*cce.DNSConfigDNSVMVNFAlias).DNSVMVNFAliasID,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) != 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"duplicate record detected for dns_config_id %s and "+
				"dns_vm_vnf_alias_id %s",
			e.(*cce.DNSConfigDNSVMVNFAlias).DNSConfigID,
			e.(*cce.DNSConfigDNSVMVNFAlias).DNSVMVNFAliasID)
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
