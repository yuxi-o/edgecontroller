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
	"runtime/debug"
	"strings"

	"github.com/gorilla/mux"
	cce "github.com/smartedgemec/controller-ce"
)

// Gorilla wraps the gorilla router and application routes.
type Gorilla struct {
	// router
	router *mux.Router

	// entity routes handlers
	nodesHandler           *handler
	appsHandler            *handler
	vnfsHandler            *handler
	trafficPoliciesHandler *handler
	dnsConfigsHandler      *handler

	// join routes handlers
	dnsConfigsAppAliasesHandler     *handler
	dnsConfigsVNFAliasesHandler     *handler
	nodesDNSConfigsHandler          *handler
	nodesAppsHandler                *handler
	nodesVNFsHandler                *handler
	nodesAppsTrafficPoliciesHandler *handler
	nodesVNFsTrafficPoliciesHandler *handler
}

// NewGorilla creates a new Gorilla.
func NewGorilla( //nolint:gocyclo
	controller *cce.Controller,
) *Gorilla {
	g := &Gorilla{
		// router
		router: mux.NewRouter(),

		// entity routes handlers
		nodesHandler: &handler{
			model: &cce.Node{},

			// TODO add checkDBDelete func + tests for nodes_apps, nodes_vnfs,
			// and nodes_dns_configs
			// checkDBDelete: checkDBDeleteNodes,

			// TODO add any application logic necessary + tests
		},
		appsHandler: &handler{
			model:         &cce.App{},
			checkDBDelete: checkDBDeleteApps,
		},
		vnfsHandler: &handler{
			model:         &cce.VNF{},
			checkDBDelete: checkDBDeleteVNFs,
		},
		trafficPoliciesHandler: &handler{
			model:         &cce.TrafficPolicy{},
			checkDBDelete: checkDBDeleteTrafficPolicies,
		},
		dnsConfigsHandler: &handler{
			model:         &cce.DNSConfig{},
			checkDBDelete: checkDBDeleteDNSConfigs,
		},

		// join routes handlers
		dnsConfigsAppAliasesHandler: &handler{
			model:         &cce.DNSConfigAppAlias{},
			checkDBCreate: checkDBCreateDNSConfigsAppAliases,
		},
		dnsConfigsVNFAliasesHandler: &handler{
			model:         &cce.DNSConfigVNFAlias{},
			checkDBCreate: checkDBCreateDNSConfigsVNFAliases,
		},
		nodesAppsHandler: &handler{
			model: &cce.NodeApp{},

			checkDBCreate: checkDBCreateNodesApps,
			// TODO add checkDBDelete func + tests for
			// nodes_apps_traffic_policies
			// checkDBDelete: checkDBDeleteNodesApps,

			handleCreate: handleCreateNodesApps,
		},
		nodesDNSConfigsHandler: &handler{
			model:         &cce.NodeDNSConfig{},
			checkDBCreate: checkDBCreateNodesDNSConfigs,
			handleCreate:  handleCreateNodesDNSConfigs,
		},
		nodesVNFsHandler: &handler{
			model: &cce.NodeVNF{},

			// TODO add checkDBDelete func + tests for
			// nodes_vnfs_traffic_policies
			// checkDBDelete: checkDBDeleteNodesVNF,
		},
		nodesAppsTrafficPoliciesHandler: &handler{
			model: &cce.NodeAppTrafficPolicy{},
		},
		nodesVNFsTrafficPoliciesHandler: &handler{
			model: &cce.NodeVNFTrafficPolicy{},
		},
	}

	routes := map[string]http.HandlerFunc{
		"POST /auth": authenticate,

		// entity routes
		"POST   /nodes":      g.nodesHandler.create,
		"GET    /nodes":      g.nodesHandler.getAll,
		"GET    /nodes/{id}": g.nodesHandler.getByID,
		"PATCH  /nodes":      g.nodesHandler.bulkUpdate,
		"DELETE /nodes/{id}": g.nodesHandler.delete,

		"POST   /apps":      g.appsHandler.create,
		"GET    /apps":      g.appsHandler.getAll,
		"GET    /apps/{id}": g.appsHandler.getByID,
		"PATCH  /apps":      g.appsHandler.bulkUpdate,
		"DELETE /apps/{id}": g.appsHandler.delete,

		"POST   /vnfs":      g.vnfsHandler.create,
		"GET    /vnfs":      g.vnfsHandler.getAll,
		"GET    /vnfs/{id}": g.vnfsHandler.getByID,
		"PATCH  /vnfs":      g.vnfsHandler.bulkUpdate,
		"DELETE /vnfs/{id}": g.vnfsHandler.delete,

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

		// join routes
		"POST   /dns_configs_app_aliases":      g.dnsConfigsAppAliasesHandler.create,  //nolint:lll
		"GET    /dns_configs_app_aliases":      g.dnsConfigsAppAliasesHandler.getAll,  //nolint:lll
		"GET    /dns_configs_app_aliases/{id}": g.dnsConfigsAppAliasesHandler.getByID, //nolint:lll
		"DELETE /dns_configs_app_aliases/{id}": g.dnsConfigsAppAliasesHandler.delete,  //nolint:lll

		"POST   /dns_configs_vnf_aliases":      g.dnsConfigsVNFAliasesHandler.create,  //nolint:lll
		"GET    /dns_configs_vnf_aliases":      g.dnsConfigsVNFAliasesHandler.getAll,  //nolint:lll
		"GET    /dns_configs_vnf_aliases/{id}": g.dnsConfigsVNFAliasesHandler.getByID, //nolint:lll
		"DELETE /dns_configs_vnf_aliases/{id}": g.dnsConfigsVNFAliasesHandler.delete,  //nolint:lll

		"POST   /nodes_dns_configs":      g.nodesDNSConfigsHandler.create,
		"GET    /nodes_dns_configs":      g.nodesDNSConfigsHandler.getAll,
		"GET    /nodes_dns_configs/{id}": g.nodesDNSConfigsHandler.getByID,
		"DELETE /nodes_dns_configs/{id}": g.nodesDNSConfigsHandler.delete,

		"POST   /nodes_apps":      g.nodesAppsHandler.create,
		"GET    /nodes_apps":      g.nodesAppsHandler.getByFilter,
		"GET    /nodes_apps/{id}": g.nodesAppsHandler.getByID,
		"PATCH  /nodes_apps":      g.nodesAppsHandler.bulkUpdate,
		"DELETE /nodes_apps/{id}": g.nodesAppsHandler.delete,

		// these endpoints still need API tests
		"POST   /nodes_vnfs":      g.nodesVNFsHandler.create,
		"GET    /nodes_vnfs":      g.nodesVNFsHandler.getByFilter,
		"DELETE /nodes_vnfs/{id}": g.nodesVNFsHandler.delete,

		"POST   /nodes_apps_traffic_policies":      g.nodesAppsTrafficPoliciesHandler.create,      //nolint:lll
		"GET    /nodes_apps_traffic_policies":      g.nodesAppsTrafficPoliciesHandler.getByFilter, //nolint:lll
		"DELETE /nodes_apps_traffic_policies/{id}": g.nodesAppsTrafficPoliciesHandler.delete,      //nolint:lll

		"POST   /nodes_vnfs_traffic_policies":      g.nodesVNFsTrafficPoliciesHandler.create,      //nolint:lll
		"GET    /nodes_vnfs_traffic_policies":      g.nodesVNFsTrafficPoliciesHandler.getByFilter, //nolint:lll
		"DELETE /nodes_vnfs_traffic_policies/{id}": g.nodesVNFsTrafficPoliciesHandler.delete,      //nolint:lll
	}

	for endpoint, handlerFunc := range routes {
		split := strings.Fields(endpoint)
		g.router.HandleFunc(split[1], handlerFunc).Methods(split[0])
	}

	// Catch panics
	g.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					log.Println("Recovered in handler func:", r)
					log.Println("Stack trace:", string(debug.Stack()))
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	})

	// Inject the controller
	g.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(
				r.Context(),
				contextKey("controller"),
				controller)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	// Require auth token for all endpoints except POST /auth
	g.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost && r.RequestURI == "/auth" {
				next.ServeHTTP(w, r)
			} else {
				requireAuthHandler(next).ServeHTTP(w, r)
			}
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
