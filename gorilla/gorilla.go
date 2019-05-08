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
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	cce "github.com/smartedgemec/controller-ce"
)

// Gorilla wraps the gorilla router and application routes.
type Gorilla struct {
	// router
	router *mux.Router

	// entity routes handlers
	nodesHandler                  *handler
	containerAppsHandler          *handler
	vmAppsHandler                 *handler
	containerVNFsHandler          *handler
	vmVNFsHandler                 *handler
	trafficPoliciesHandler        *handler
	dnsConfigsHandler             *handler
	dnsContainerAppAliasesHandler *handler
	dnsContainerVNFAliasesHandler *handler
	dnsVMAppAliasesHandler        *handler
	dnsVMVNFAliasesHandler        *handler

	// join routes handlers
	dnsConfigsDNSContainerAppAliasesHandler  *handler
	dnsConfigsDNSContainerVNFAliasesHandler  *handler
	dnsConfigsDNSVMAppAliasesHandler         *handler
	dnsConfigsDNSVMVNFAliasesHandler         *handler
	nodesDNSConfigsHandler                   *handler
	nodesContainerAppsHandler                *handler
	nodesVMAppsHandler                       *handler
	nodesContainerVNFsHandler                *handler
	nodesVMVNFsHandler                       *handler
	nodesContainerAppsTrafficPoliciesHandler *handler
	nodesVMAppsTrafficPoliciesHandler        *handler
}

// NewGorilla creates a new Gorilla.
func NewGorilla( //nolint:gocyclo
	controller *cce.Controller,
	nodeMap map[string]*cce.Node,
) *Gorilla {
	g := &Gorilla{
		router: mux.NewRouter(),

		nodesHandler: &handler{model: &cce.Node{}},
		containerAppsHandler: &handler{
			model:         &cce.ContainerApp{},
			checkDBDelete: checkDBDeleteContainerApps,
		},
		vmAppsHandler: &handler{
			model:         &cce.VMApp{},
			checkDBDelete: checkDBDeleteVMApps,
		},
		containerVNFsHandler: &handler{
			model:         &cce.ContainerVNF{},
			checkDBDelete: checkDBDeleteContainerVNFs,
		},
		vmVNFsHandler: &handler{
			model:         &cce.VMVNF{},
			checkDBDelete: checkDBDeleteVMVNFs,
		},
		trafficPoliciesHandler: &handler{model: &cce.TrafficPolicy{}},
		dnsConfigsHandler: &handler{
			model:         &cce.DNSConfig{},
			checkDBDelete: checkDBDeleteDNSConfigs,
		},
		dnsContainerAppAliasesHandler: &handler{
			model:         &cce.DNSContainerAppAlias{},
			checkDBDelete: checkDBDeleteDNSContainerAppAliases,
		},
		dnsContainerVNFAliasesHandler: &handler{
			model:         &cce.DNSContainerVNFAlias{},
			checkDBDelete: checkDBDeleteDNSContainerVNFAliases,
		},
		dnsVMAppAliasesHandler: &handler{
			model:         &cce.DNSVMAppAlias{},
			checkDBDelete: checkDBDeleteDNSVMAppAliases,
		},
		dnsVMVNFAliasesHandler: &handler{
			model:         &cce.DNSVMVNFAlias{},
			checkDBDelete: checkDBDeleteDNSVMVNFAliases,
		},

		dnsConfigsDNSContainerAppAliasesHandler: &handler{
			model:         &cce.DNSConfigDNSContainerAppAlias{},
			checkDBCreate: checkDBCreateDNSConfigsDNSContainerAppAliases,
		},
		dnsConfigsDNSContainerVNFAliasesHandler: &handler{
			model:         &cce.DNSConfigDNSContainerVNFAlias{},
			checkDBCreate: checkDBCreateDNSConfigsDNSContainerVNFAliases,
		},
		dnsConfigsDNSVMAppAliasesHandler: &handler{
			model:         &cce.DNSConfigDNSVMAppAlias{},
			checkDBCreate: checkDBCreateDNSConfigsDNSVMAppAliases,
		},
		dnsConfigsDNSVMVNFAliasesHandler: &handler{
			model:         &cce.DNSConfigDNSVMVNFAlias{},
			checkDBCreate: checkDBCreateDNSConfigsDNSVMVNFAliases,
		},
		nodesDNSConfigsHandler: &handler{
			model:         &cce.NodeDNSConfig{},
			checkDBCreate: checkDBCreateNodeDNSConfigs,
			// TODO add logic to apply DNS config to node, and tests
		},
		nodesContainerAppsHandler:                &handler{model: &cce.NodeContainerApp{}},              //nolint:lll
		nodesVMAppsHandler:                       &handler{model: &cce.NodeVMApp{}},                     //nolint:lll
		nodesContainerVNFsHandler:                &handler{model: &cce.NodeContainerVNF{}},              //nolint:lll
		nodesVMVNFsHandler:                       &handler{model: &cce.NodeVMVNF{}},                     //nolint:lll
		nodesContainerAppsTrafficPoliciesHandler: &handler{model: &cce.NodeContainerAppTrafficPolicy{}}, //nolint:lll
		nodesVMAppsTrafficPoliciesHandler:        &handler{model: &cce.NodeVMAppTrafficPolicy{}},        //nolint:lll
	}

	routes := map[string]http.HandlerFunc{
		"POST   /nodes":      g.nodesHandler.create,
		"GET    /nodes":      g.nodesHandler.getAll,
		"GET    /nodes/{id}": g.nodesHandler.getByID,
		"PATCH  /nodes":      g.nodesHandler.bulkUpdate,
		"DELETE /nodes/{id}": g.nodesHandler.delete,

		"POST   /container_apps":      g.containerAppsHandler.create,
		"GET    /container_apps":      g.containerAppsHandler.getAll,
		"GET    /container_apps/{id}": g.containerAppsHandler.getByID,
		"PATCH  /container_apps":      g.containerAppsHandler.bulkUpdate,
		"DELETE /container_apps/{id}": g.containerAppsHandler.delete,

		"POST   /vm_apps":      g.vmAppsHandler.create,
		"GET    /vm_apps":      g.vmAppsHandler.getAll,
		"GET    /vm_apps/{id}": g.vmAppsHandler.getByID,
		"PATCH  /vm_apps":      g.vmAppsHandler.bulkUpdate,
		"DELETE /vm_apps/{id}": g.vmAppsHandler.delete,

		"POST   /container_vnfs":      g.containerVNFsHandler.create,
		"GET    /container_vnfs":      g.containerVNFsHandler.getAll,
		"GET    /container_vnfs/{id}": g.containerVNFsHandler.getByID,
		"PATCH  /container_vnfs":      g.containerVNFsHandler.bulkUpdate,
		"DELETE /container_vnfs/{id}": g.containerVNFsHandler.delete,

		"POST   /vm_vnfs":      g.vmVNFsHandler.create,
		"GET    /vm_vnfs":      g.vmVNFsHandler.getAll,
		"GET    /vm_vnfs/{id}": g.vmVNFsHandler.getByID,
		"PATCH  /vm_vnfs":      g.vmVNFsHandler.bulkUpdate,
		"DELETE /vm_vnfs/{id}": g.vmVNFsHandler.delete,

		"POST   /traffic_policies":      g.trafficPoliciesHandler.create,
		"GET    /traffic_policies":      g.trafficPoliciesHandler.getAll,
		"GET    /traffic_policies/{id}": g.trafficPoliciesHandler.getByID,
		"PATCH  /traffic_policies":      g.trafficPoliciesHandler.bulkUpdate,
		"DELETE /traffic_policies/{id}": g.trafficPoliciesHandler.delete,

		"POST   /dns_configs":      g.dnsConfigsHandler.create,
		"GET    /dns_configs":      g.dnsConfigsHandler.getAll,
		"GET    /dns_configs/{id}": g.dnsConfigsHandler.getByID,
		"PATCH  /dns_configs":      g.dnsConfigsHandler.bulkUpdate,
		"DELETE /dns_configs/{id}": g.dnsConfigsHandler.delete,

		"POST   /dns_container_app_aliases":      g.dnsContainerAppAliasesHandler.create,     //nolint:lll
		"GET    /dns_container_app_aliases":      g.dnsContainerAppAliasesHandler.getAll,     //nolint:lll
		"GET    /dns_container_app_aliases/{id}": g.dnsContainerAppAliasesHandler.getByID,    //nolint:lll
		"PATCH  /dns_container_app_aliases":      g.dnsContainerAppAliasesHandler.bulkUpdate, //nolint:lll
		"DELETE /dns_container_app_aliases/{id}": g.dnsContainerAppAliasesHandler.delete,     //nolint:lll

		"POST   /dns_container_vnf_aliases":      g.dnsContainerVNFAliasesHandler.create,     //nolint:lll
		"GET    /dns_container_vnf_aliases":      g.dnsContainerVNFAliasesHandler.getAll,     //nolint:lll
		"GET    /dns_container_vnf_aliases/{id}": g.dnsContainerVNFAliasesHandler.getByID,    //nolint:lll
		"PATCH  /dns_container_vnf_aliases":      g.dnsContainerVNFAliasesHandler.bulkUpdate, //nolint:lll
		"DELETE /dns_container_vnf_aliases/{id}": g.dnsContainerVNFAliasesHandler.delete,     //nolint:lll

		"POST   /dns_vm_app_aliases":      g.dnsVMAppAliasesHandler.create,     //nolint:lll
		"GET    /dns_vm_app_aliases":      g.dnsVMAppAliasesHandler.getAll,     //nolint:lll
		"GET    /dns_vm_app_aliases/{id}": g.dnsVMAppAliasesHandler.getByID,    //nolint:lll
		"PATCH  /dns_vm_app_aliases":      g.dnsVMAppAliasesHandler.bulkUpdate, //nolint:lll
		"DELETE /dns_vm_app_aliases/{id}": g.dnsVMAppAliasesHandler.delete,     //nolint:lll

		"POST   /dns_vm_vnf_aliases":      g.dnsVMVNFAliasesHandler.create,     //nolint:lll
		"GET    /dns_vm_vnf_aliases":      g.dnsVMVNFAliasesHandler.getAll,     //nolint:lll
		"GET    /dns_vm_vnf_aliases/{id}": g.dnsVMVNFAliasesHandler.getByID,    //nolint:lll
		"PATCH  /dns_vm_vnf_aliases":      g.dnsVMVNFAliasesHandler.bulkUpdate, //nolint:lll
		"DELETE /dns_vm_vnf_aliases/{id}": g.dnsVMVNFAliasesHandler.delete,     //nolint:lll

		"POST   /dns_configs_dns_container_app_aliases":      g.dnsConfigsDNSContainerAppAliasesHandler.create,  //nolint:lll
		"GET    /dns_configs_dns_container_app_aliases":      g.dnsConfigsDNSContainerAppAliasesHandler.getAll,  //nolint:lll
		"GET    /dns_configs_dns_container_app_aliases/{id}": g.dnsConfigsDNSContainerAppAliasesHandler.getByID, //nolint:lll
		"DELETE /dns_configs_dns_container_app_aliases/{id}": g.dnsConfigsDNSContainerAppAliasesHandler.delete,  //nolint:lll

		"POST   /dns_configs_dns_container_vnf_aliases":      g.dnsConfigsDNSContainerVNFAliasesHandler.create,  //nolint:lll
		"GET    /dns_configs_dns_container_vnf_aliases":      g.dnsConfigsDNSContainerVNFAliasesHandler.getAll,  //nolint:lll
		"GET    /dns_configs_dns_container_vnf_aliases/{id}": g.dnsConfigsDNSContainerVNFAliasesHandler.getByID, //nolint:lll
		"DELETE /dns_configs_dns_container_vnf_aliases/{id}": g.dnsConfigsDNSContainerVNFAliasesHandler.delete,  //nolint:lll

		"POST   /dns_configs_dns_vm_app_aliases":      g.dnsConfigsDNSVMAppAliasesHandler.create,  //nolint:lll
		"GET    /dns_configs_dns_vm_app_aliases":      g.dnsConfigsDNSVMAppAliasesHandler.getAll,  //nolint:lll
		"GET    /dns_configs_dns_vm_app_aliases/{id}": g.dnsConfigsDNSVMAppAliasesHandler.getByID, //nolint:lll
		"DELETE /dns_configs_dns_vm_app_aliases/{id}": g.dnsConfigsDNSVMAppAliasesHandler.delete,  //nolint:lll

		"POST   /dns_configs_dns_vm_vnf_aliases":      g.dnsConfigsDNSVMVNFAliasesHandler.create,  //nolint:lll
		"GET    /dns_configs_dns_vm_vnf_aliases":      g.dnsConfigsDNSVMVNFAliasesHandler.getAll,  //nolint:lll
		"GET    /dns_configs_dns_vm_vnf_aliases/{id}": g.dnsConfigsDNSVMVNFAliasesHandler.getByID, //nolint:lll
		"DELETE /dns_configs_dns_vm_vnf_aliases/{id}": g.dnsConfigsDNSVMVNFAliasesHandler.delete,  //nolint:lll

		"POST   /nodes_dns_configs":      g.nodesDNSConfigsHandler.create,
		"GET    /nodes_dns_configs":      g.nodesDNSConfigsHandler.getAll,
		"GET    /nodes_dns_configs/{id}": g.nodesDNSConfigsHandler.getByID,
		"PATCH  /nodes_dns_configs":      g.nodesDNSConfigsHandler.bulkUpdate,
		"DELETE /nodes_dns_configs/{id}": g.nodesDNSConfigsHandler.delete,

		"POST   /nodes_container_apps":      g.nodesContainerAppsHandler.create,      //nolint:lll
		"GET    /nodes_container_apps":      g.nodesContainerAppsHandler.getByFilter, //nolint:lll
		"GET    /nodes_container_apps/{id}": g.nodesContainerAppsHandler.getByID,     //nolint:lll
		"PATCH  /nodes_container_apps":      g.nodesContainerAppsHandler.bulkUpdate,  //nolint:lll
		"DELETE /nodes_container_apps/{id}": g.nodesContainerAppsHandler.delete,      //nolint:lll

		"POST   /nodes_vm_apps":      g.nodesVMAppsHandler.create,
		"GET    /nodes_vm_apps":      g.nodesVMAppsHandler.getByFilter,
		"DELETE /nodes_vm_apps/{id}": g.nodesVMAppsHandler.delete,

		"POST   /nodes_container_vnfs":      g.nodesContainerVNFsHandler.create,
		"GET    /nodes_container_vnfs":      g.nodesContainerVNFsHandler.getByFilter, //nolint:lll
		"DELETE /nodes_container_vnfs/{id}": g.nodesContainerVNFsHandler.delete,

		"POST   /nodes_vm_vnfs":      g.nodesVMVNFsHandler.create,
		"GET    /nodes_vm_vnfs":      g.nodesVMVNFsHandler.getByFilter,
		"DELETE /nodes_vm_vnfs/{id}": g.nodesVMVNFsHandler.delete,

		"POST   /nodes_container_apps_traffic_policies":      g.nodesContainerAppsTrafficPoliciesHandler.create,      //nolint:lll
		"GET    /nodes_container_apps_traffic_policies":      g.nodesContainerAppsTrafficPoliciesHandler.getByFilter, //nolint:lll
		"DELETE /nodes_container_apps_traffic_policies/{id}": g.nodesContainerAppsTrafficPoliciesHandler.delete,      //nolint:lll

		"POST   /nodes_vm_apps_traffic_policies":      g.nodesVMAppsTrafficPoliciesHandler.create,      //nolint:lll
		"GET    /nodes_vm_apps_traffic_policies":      g.nodesVMAppsTrafficPoliciesHandler.getByFilter, //nolint:lll
		"DELETE /nodes_vm_apps_traffic_policies/{id}": g.nodesVMAppsTrafficPoliciesHandler.delete,      //nolint:lll
	}

	for endpoint, handlerFunc := range routes {
		split := strings.Fields(endpoint)
		g.router.HandleFunc(split[1], handlerFunc).Methods(split[0])
	}

	// Inject the controller
	g.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(
				r.Context(),
				contextKey("controller"),
				controller)
			log.Printf("Injected controller %#v", controller)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	// Read and inject the body for POST and PATCH requests
	g.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "POST", "PATCH":
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					log.Printf("Error reading body: %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				ctx := context.WithValue(r.Context(), contextKey("body"), body)
				log.Println("Injected body", string(body))
				next.ServeHTTP(w, r.WithContext(ctx))
			default:
				next.ServeHTTP(w, r)
			}
		})
	})

	return g
}

type contextKey string

func (c contextKey) String() string {
	return "controller-ce context key " + string(c)
}

// ServeHTTP wraps mux.ServeHTTP.
func (g *Gorilla) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	g.router.ServeHTTP(w, req)
}
