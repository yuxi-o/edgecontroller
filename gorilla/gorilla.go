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
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gorilla/mux"
	cce "github.com/smartedgemec/controller-ce"
	logger "github.com/smartedgemec/log"
)

var log = logger.DefaultLogger.WithField("pkg", "gorilla")

// Gorilla wraps the gorilla router and application routes.
type Gorilla struct {
	// router
	router *mux.Router

	// entity routes handlers
	nodesHandler           *handler
	appsHandler            *handler
	trafficPoliciesHandler *handler
	dnsConfigsHandler      *handler

	// join routes handlers
	dnsConfigsAppAliasesHandler     *handler
	nodesDNSConfigsHandler          *handler
	nodesAppsHandler                *handler
	nodesAppsTrafficPoliciesHandler *handler
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
			model:    &cce.Node{},
			reqModel: &cce.NodeReq{},

			checkDBDelete: checkDBDeleteNodes,

			handleGet:    handleGetNodes,
			handleUpdate: handleUpdateNodes,
		},
		appsHandler: &handler{
			model:         &cce.App{},
			checkDBDelete: checkDBDeleteApps,
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
		nodesAppsHandler: &handler{
			model:    &cce.NodeApp{},
			reqModel: &cce.NodeAppReq{},

			checkDBCreate: checkDBCreateNodesApps,
			checkDBDelete: checkDBDeleteNodesApps,

			handleCreate: handleCreateNodesApps,
			handleGet:    handleGetNodesApps,
			handleUpdate: handleUpdateNodesApps,
			handleDelete: handleDeleteNodesApps,
		},
		nodesDNSConfigsHandler: &handler{
			model: &cce.NodeDNSConfig{},

			checkDBCreate: checkDBCreateNodesDNSConfigs,

			handleCreate: handleCreateNodesDNSConfigs,
			handleDelete: handleDeleteNodesDNSConfigs,
		},
		nodesAppsTrafficPoliciesHandler: &handler{
			model: &cce.NodeAppTrafficPolicy{},

			checkDBCreate: checkDBCreateNodesAppsTrafficPolicies,

			handleCreate: handleCreateNodesAppsTrafficPolicies,
			handleDelete: handleDeleteNodesAppsTrafficPolicies,
		},
	}

	routes := map[string]http.HandlerFunc{
		"POST /auth": authenticate,

		// entity routes
		"POST   /nodes":      g.nodesHandler.create,
		"GET    /nodes":      g.nodesHandler.filter,
		"GET    /nodes/{id}": g.nodesHandler.getByID,
		"PATCH  /nodes":      g.nodesHandler.bulkUpdate,
		"DELETE /nodes/{id}": g.nodesHandler.delete,

		"POST   /traffic_policies":      g.trafficPoliciesHandler.create,
		"GET    /traffic_policies":      g.trafficPoliciesHandler.filter,
		"GET    /traffic_policies/{id}": g.trafficPoliciesHandler.getByID,
		"PATCH  /traffic_policies":      g.trafficPoliciesHandler.bulkUpdate,
		"DELETE /traffic_policies/{id}": g.trafficPoliciesHandler.delete,

		"POST   /dns_configs":      g.dnsConfigsHandler.create,
		"GET    /dns_configs":      g.dnsConfigsHandler.filter,
		"GET    /dns_configs/{id}": g.dnsConfigsHandler.getByID,
		"PATCH  /dns_configs":      g.dnsConfigsHandler.bulkUpdate,
		"DELETE /dns_configs/{id}": g.dnsConfigsHandler.delete,

		// non-node join routes
		"POST   /dns_configs_app_aliases":      g.dnsConfigsAppAliasesHandler.create,
		"GET    /dns_configs_app_aliases":      g.dnsConfigsAppAliasesHandler.filter,
		"GET    /dns_configs_app_aliases/{id}": g.dnsConfigsAppAliasesHandler.getByID,
		"DELETE /dns_configs_app_aliases/{id}": g.dnsConfigsAppAliasesHandler.delete,

		// node join routes
		"POST   /nodes_apps":      g.nodesAppsHandler.create,
		"GET    /nodes_apps":      g.nodesAppsHandler.filter,
		"GET    /nodes_apps/{id}": g.nodesAppsHandler.getByID,
		"PATCH  /nodes_apps":      g.nodesAppsHandler.bulkUpdate,
		"DELETE /nodes_apps/{id}": g.nodesAppsHandler.delete,

		"POST   /nodes_dns_configs":      g.nodesDNSConfigsHandler.create,
		"GET    /nodes_dns_configs":      g.nodesDNSConfigsHandler.filter,
		"GET    /nodes_dns_configs/{id}": g.nodesDNSConfigsHandler.getByID,
		"DELETE /nodes_dns_configs/{id}": g.nodesDNSConfigsHandler.delete,

		"POST   /nodes_apps_traffic_policies":      g.nodesAppsTrafficPoliciesHandler.create,
		"GET    /nodes_apps_traffic_policies":      g.nodesAppsTrafficPoliciesHandler.filter,
		"GET    /nodes_apps_traffic_policies/{id}": g.nodesAppsTrafficPoliciesHandler.getByID,
		"DELETE /nodes_apps_traffic_policies/{id}": g.nodesAppsTrafficPoliciesHandler.delete,

		// The following endpoints are compliant with the Swagger / OpenAPI 3.0 schema. All references
		// to the `v2` prefix in code should be removed once all endpoints are implemented. In conjunction
		// with this, the old endpoints and their tests should be dropped.
		"GET      /v2/nodes":           nil, // TODO
		"POST     /v2/nodes":           nil, // TODO
		"GET      /v2/nodes/{node_id}": nil, // TODO
		"PATCH    /v2/nodes/{node_id}": nil, // TODO
		"DELETE   /v2/nodes/{node_id}": nil, // TODO

		"GET      /apps":          g.swagGETApps,
		"POST     /apps":          g.swagPOSTApps,
		"GET      /apps/{app_id}": g.swagGETAppByID,
		"PATCH    /apps/{app_id}": g.swagPATCHAppByID,
		"DELETE   /apps/{app_id}": g.swagDELETEAppByID,

		"GET      /v2/policies":             nil, // TODO
		"POST     /v2/policies":             nil, // TODO
		"GET      /v2/policies/{policy_id}": nil, // TODO
		"PATCH    /v2/policies/{policy_id}": nil, // TODO
		"DELETE   /v2/policies/{policy_id}": nil, // TODO

		"GET      /v2/nodes/{node_id}/dns": nil, // TODO
		"PATCH    /v2/nodes/{node_id}/dns": nil, // TODO
		"DELETE   /v2/nodes/{node_id}/dns": nil, // TODO

		"GET      /v2/nodes/{node_id}/interfaces":                nil, // TODO
		"PATCH    /v2/nodes/{node_id}/interfaces":                nil, // TODO
		"GET      /v2/nodes/{node_id}/interfaces/{interface_id}": nil, // TODO

		"GET      /v2/nodes/{node_id}/interfaces/{interface_id}/policy": nil, // TODO
		"PATCH    /v2/nodes/{node_id}/interfaces/{interface_id}/policy": nil, // TODO
		"DELETE   /v2/nodes/{node_id}/interfaces/{interface_id}/policy": nil, // TODO

		"GET      /v2/nodes/{node_id}/apps":                 nil, // TODO
		"POST     /v2/nodes/{node_id}/apps":                 nil, // TODO
		"GET      /v2/nodes/{node_id}/apps/{app_id}":        nil, // TODO
		"PATCH    /v2/nodes/{node_id}/apps/{app_id}":        nil, // TODO
		"DELETE   /v2/nodes/{node_id}/apps/{app_id}":        nil, // TODO
		"GET      /v2/nodes/{node_id}/apps/{app_id}/policy": nil, // TODO
		"PATCH    /v2/nodes/{node_id}/apps/{app_id}/policy": nil, // TODO
		"DELETE   /v2/nodes/{node_id}/apps/{app_id}/policy": nil, // TODO
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
					log.Critf("Recovered in handler func: %q\nStack trace:\n%s",
						r, string(debug.Stack()))
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
					log.Errf("Error reading body: %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				ctx := context.WithValue(r.Context(), contextKey("body"), body)

				// Scrub for the body payload for potentially sensitive authentication data
				// (this only affects logging, not the actual request body)
				// TODO: Log the JSON payload here but with the password field scrubbed
				if r.URL.Path == "/auth" {
					body = []byte("***** REDACTED *****")
				}

				log.Debugf("Injected body: %s", string(body))
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
