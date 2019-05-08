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

func checkDBDeleteContainerApps(
	ctx context.Context,
	ps cce.PersistenceService,
	id string,
) (statusCode int, err error) {
	var es []cce.Entity

	if es, err = ps.Filter(
		ctx,
		&cce.DNSContainerAppAlias{},
		[]cce.Filter{
			{
				Field: "container_app_id",
				Value: id,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) != 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"cannot delete container_app_id %s: record in use in "+
				"dns_container_app_aliases",
			id)
	}

	return 0, nil
}

func checkDBDeleteVMApps(
	ctx context.Context,
	ps cce.PersistenceService,
	id string,
) (statusCode int, err error) {
	var es []cce.Entity

	if es, err = ps.Filter(
		ctx,
		&cce.DNSVMAppAlias{},
		[]cce.Filter{
			{
				Field: "vm_app_id",
				Value: id,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) != 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"cannot delete vm_app_id %s: record in use in dns_vm_app_aliases",
			id)
	}

	return 0, nil
}

func checkDBDeleteContainerVNFs(
	ctx context.Context,
	ps cce.PersistenceService,
	id string,
) (statusCode int, err error) {
	var es []cce.Entity

	if es, err = ps.Filter(
		ctx,
		&cce.DNSContainerVNFAlias{},
		[]cce.Filter{
			{
				Field: "container_vnf_id",
				Value: id,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) != 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"cannot delete container_vnf_id %s: record in use in "+
				"dns_container_vnf_aliases",
			id)
	}

	return 0, nil
}

func checkDBDeleteVMVNFs(
	ctx context.Context,
	ps cce.PersistenceService,
	id string,
) (statusCode int, err error) {
	var es []cce.Entity

	if es, err = ps.Filter(
		ctx,
		&cce.DNSVMVNFAlias{},
		[]cce.Filter{
			{
				Field: "vm_vnf_id",
				Value: id,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) != 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"cannot delete vm_vnf_id %s: record in use in dns_vm_vnf_aliases",
			id)
	}

	return 0, nil
}

func checkDBDeleteDNSConfigs(
	ctx context.Context,
	ps cce.PersistenceService,
	id string,
) (statusCode int, err error) {
	var es []cce.Entity

	if es, err = ps.Filter(
		ctx,
		&cce.DNSConfigDNSContainerAppAlias{},
		[]cce.Filter{
			{
				Field: "dns_config_id",
				Value: id,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) != 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"cannot delete dns_config_id %s: record in use in "+
				"dns_configs_dns_container_app_aliases",
			id)
	}

	if es, err = ps.Filter(
		ctx,
		&cce.DNSConfigDNSContainerVNFAlias{},
		[]cce.Filter{
			{
				Field: "dns_config_id",
				Value: id,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) != 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"cannot delete dns_config_id %s: record in use in "+
				"dns_configs_dns_container_vnf_aliases",
			id)
	}

	if es, err = ps.Filter(
		ctx,
		&cce.DNSConfigDNSVMAppAlias{},
		[]cce.Filter{
			{
				Field: "dns_config_id",
				Value: id,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) != 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"cannot delete dns_config_id %s: record in use in "+
				"dns_configs_dns_vm_app_aliases",
			id)
	}

	if es, err = ps.Filter(
		ctx,
		&cce.DNSConfigDNSVMVNFAlias{},
		[]cce.Filter{
			{
				Field: "dns_config_id",
				Value: id,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) != 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"cannot delete dns_config_id %s: record in use in "+
				"dns_configs_dns_vm_vnf_aliases",
			id)
	}

	return 0, nil
}

func checkDBDeleteDNSContainerAppAliases(
	ctx context.Context,
	ps cce.PersistenceService,
	id string,
) (statusCode int, err error) {
	var es []cce.Entity

	if es, err = ps.Filter(
		ctx,
		&cce.DNSConfigDNSContainerAppAlias{},
		[]cce.Filter{
			{
				Field: "dns_container_app_alias_id",
				Value: id,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) != 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"cannot delete dns_container_app_alias_id %s: record in use in "+
				"dns_configs_dns_container_app_aliases",
			id)
	}

	return 0, nil
}

func checkDBDeleteDNSContainerVNFAliases(
	ctx context.Context,
	ps cce.PersistenceService,
	id string,
) (statusCode int, err error) {
	var es []cce.Entity

	if es, err = ps.Filter(
		ctx,
		&cce.DNSConfigDNSContainerVNFAlias{},
		[]cce.Filter{
			{
				Field: "dns_container_vnf_alias_id",
				Value: id,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) != 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"cannot delete dns_container_vnf_alias_id %s: record in use in "+
				"dns_configs_dns_container_vnf_aliases",
			id)
	}

	return 0, nil
}

func checkDBDeleteDNSVMAppAliases(
	ctx context.Context,
	ps cce.PersistenceService,
	id string,
) (statusCode int, err error) {
	var es []cce.Entity

	if es, err = ps.Filter(
		ctx,
		&cce.DNSConfigDNSVMAppAlias{},
		[]cce.Filter{
			{
				Field: "dns_vm_app_alias_id",
				Value: id,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) != 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"cannot delete dns_vm_app_alias_id %s: record in use in "+
				"dns_configs_dns_vm_app_aliases",
			id)
	}

	return 0, nil
}

func checkDBDeleteDNSVMVNFAliases(
	ctx context.Context,
	ps cce.PersistenceService,
	id string,
) (statusCode int, err error) {
	var es []cce.Entity

	if es, err = ps.Filter(
		ctx,
		&cce.DNSConfigDNSVMVNFAlias{},
		[]cce.Filter{
			{
				Field: "dns_vm_vnf_alias_id",
				Value: id,
			},
		},
	); err != nil {
		return http.StatusInternalServerError, err
	}

	if len(es) != 0 {
		return http.StatusUnprocessableEntity, fmt.Errorf(
			"cannot delete dns_vm_vnf_alias_id %s: record in use in "+
				"dns_configs_dns_vm_vnf_aliases",
			id)
	}

	return 0, nil
}
