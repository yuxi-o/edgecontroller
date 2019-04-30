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
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	cce "github.com/smartedgemec/controller-ce"
	"github.com/smartedgemec/controller-ce/grpc/node"
)

// Gorilla wraps the gorilla router and application routes.
type Gorilla struct {
	// router
	router *mux.Router

	// entity routes handlers
	nodesHandler           *handler
	containerAppsHandler   *handler
	vmAppsHandler          *handler
	containerVNFsHandler   *handler
	vmVNFsHandler          *handler
	trafficPoliciesHandler *handler

	// join routes handlers
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
		mux.NewRouter(),

		&handler{model: &cce.Node{}},
		&handler{model: &cce.ContainerApp{}},
		&handler{model: &cce.VMApp{}},
		&handler{model: &cce.ContainerVNF{}},
		&handler{model: &cce.VMVNF{}},
		&handler{model: &cce.TrafficPolicy{}},

		&handler{&cce.NodeContainerApp{}, &nodesContainerAppsBLA{}},
		&handler{model: &cce.NodeVMApp{}},
		&handler{model: &cce.NodeContainerVNF{}},
		&handler{model: &cce.NodeVMVNF{}},
		&handler{model: &cce.NodeContainerAppTrafficPolicy{}},
		&handler{model: &cce.NodeVMAppTrafficPolicy{}},
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

		// dns routes here

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

	// Initialize a map of collection name to entity model. Used in GET /nodes*
	// requests to pass the correct Entity to the persistence service for
	// deserialization from the database.
	//
	// Note: this is subject to change once the proxy service is ready
	nodeCollections := map[string]cce.Entity{
		"nodes":                                 &cce.Node{},
		"nodes_container_apps":                  &cce.NodeContainerApp{},
		"nodes_vm_apps":                         &cce.NodeVMApp{},
		"nodes_container_vnfs":                  &cce.NodeContainerVNF{},
		"nodes_vm_vnfs":                         &cce.NodeVMVNF{},
		"nodes_container_apps_traffic_policies": &cce.NodeContainerAppTrafficPolicy{}, //nolint:lll
		"nodes_vm_apps_traffic_policies":        &cce.NodeVMAppTrafficPolicy{},
	}

	// For GET|POST|PATCH /nodes* requests make necessary node connections
	//
	// Note: this is subject to change once the proxy service is ready
	g.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/nodes") {
				var nodeIDs []string

				type joinResource struct {
					NodeID string `json:"node_id"`
				}

				switch {
				case r.Method == "POST":
					body := r.Context().Value(contextKey("body")).([]byte)
					var jr joinResource
					if err := json.Unmarshal(body, &jr); err != nil {
						log.Printf("Error unmarshalling json: %v", err)
						w.WriteHeader(http.StatusBadRequest)
						return
					}

					nodeIDs = append(nodeIDs, jr.NodeID)
				case r.Method == "PATCH":
					body := r.Context().Value(contextKey("body")).([]byte)
					var jrs []joinResource
					if err := json.Unmarshal(body, &jrs); err != nil {
						log.Printf("Error unmarshalling json: %v", err)
						w.WriteHeader(http.StatusBadRequest)
						return
					}

					for _, jr := range jrs {
						nodeIDs = append(nodeIDs, jr.NodeID)
					}
				case r.Method == "GET" && mux.Vars(r)["id"] != "":
					id := mux.Vars(r)["id"]
					ctrl := r.Context().Value(
						contextKey("controller")).(*cce.Controller)
					collectionName := strings.Split(r.URL.Path, "/")[1]
					jm := nodeCollections[collectionName]
					je, err := ctrl.PersistenceService.Read(
						r.Context(), id, jm)
					if err != nil {
						log.Printf("Error in persistence service: %v", err)
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					switch e := je.(type) {
					case *cce.Node:
						nodeIDs = append(nodeIDs, e.GetID())
					case cce.JoinEntity:
						nodeIDs = append(nodeIDs, e.GetNodeID())
					}
				}

				var nodeCCs []node.ClientConn
				for _, nodeID := range nodeIDs {
					nodeCC := node.ClientConn{Node: nodeMap[nodeID]}
					// TODO refactor to use proxy service
					// if err := nodeCC.Connect(r.Context()); err != nil {
					// 	log.Printf("Could not connect to node: %v", err)
					// 	w.WriteHeader(http.StatusInternalServerError)
					// 	return
					// }

					// log.Printf("Connection to node %v established",
					// 	nodeCC.Node)
					nodeCCs = append(nodeCCs, nodeCC)
				}

				ctx := context.WithValue(
					r.Context(),
					contextKey("nodes"),
					nodeCCs)
				log.Println("Injected nodes")
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
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
