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
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	cce "github.com/smartedgemec/controller-ce"
	"github.com/smartedgemec/controller-ce/swagger"
	"github.com/smartedgemec/controller-ce/uuid"
)

// The following handlers are compliant to our published Swagger (OpenAPI 3.0) schema.

// Used for GET /apps endpoint
func (g *Gorilla) swagGETApps(w http.ResponseWriter, r *http.Request) {
	// Load the controller to access the persistence
	ctrl := r.Context().Value(contextKey("controller")).(*cce.Controller)

	// Fetch the nodes from persistence
	persisted, err := ctrl.PersistenceService.ReadAll(r.Context(), &cce.App{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Construct the response object
	apps := swagger.AppList{Apps: []swagger.AppSummary{}}
	for _, a := range persisted {
		app := swagger.AppSummary{
			ID:          a.(*cce.App).ID,
			Type:        a.(*cce.App).Type,
			Name:        a.(*cce.App).Name,
			Version:     a.(*cce.App).Version,
			Vendor:      a.(*cce.App).Vendor,
			Description: a.(*cce.App).Description,
		}
		apps.Apps = append(apps.Apps, app)
	}

	// Marshal the response object to JSON
	appsJSON, err := json.Marshal(apps)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(appsJSON); err != nil {
		log.Errf("Error writing response: %v", err)
	}
}

// Used for POST /apps endpoint
func (g *Gorilla) swagPOSTApps(w http.ResponseWriter, r *http.Request) {
	g.appsHandler.create(w, r)
}

// Used for GET /apps/{app_id} endpoint
func (g *Gorilla) swagGETAppByID(w http.ResponseWriter, r *http.Request) {
	// Load the controller to access the persistence
	ctrl := r.Context().Value(contextKey("controller")).(*cce.Controller)

	// Fetch the entity from persistence and check if it's there
	persisted, err := ctrl.PersistenceService.Read(r.Context(), mux.Vars(r)["app_id"], &cce.App{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if persisted == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Construct the response object
	app := swagger.AppDetail{
		AppSummary: swagger.AppSummary{
			ID:          persisted.(*cce.App).ID,
			Type:        persisted.(*cce.App).Type,
			Name:        persisted.(*cce.App).Name,
			Version:     persisted.(*cce.App).Version,
			Vendor:      persisted.(*cce.App).Vendor,
			Description: persisted.(*cce.App).Description,
		},
		Cores:  persisted.(*cce.App).Cores,
		Memory: persisted.(*cce.App).Memory,
		Ports:  []swagger.PortProto{},
		Source: persisted.(*cce.App).Source,
	}
	for _, port := range persisted.(*cce.App).Ports {
		app.Ports = append(app.Ports, swagger.PortProto{PortProto: port})
	}

	// Marshal the response object to JSON
	appJSON, err := json.Marshal(app)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(appJSON); err != nil {
		log.Errf("Error writing response: %v", err)
	}
}

// Used for PATCH /apps/{app_id} endpoint
func (g *Gorilla) swagPATCHAppByID(w http.ResponseWriter, r *http.Request) {
	// Load the controller to access the persistence and the payload
	ctrl := r.Context().Value(contextKey("controller")).(*cce.Controller)
	body := r.Context().Value(contextKey("body")).([]byte)

	// Unmarshal the payload
	app := swagger.AppDetail{}
	if err := json.Unmarshal(body, &app); err != nil {
		log.Errf("Error unmarshaling json: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Convert it to a persistable object
	persisted := cce.App{
		ID:          app.ID,
		Type:        app.Type,
		Name:        app.Name,
		Version:     app.Version,
		Vendor:      app.Vendor,
		Description: app.Description,
		Cores:       app.Cores,
		Memory:      app.Memory,
		Ports:       []cce.PortProto{},
		Source:      app.Source,
	}
	for _, port := range app.Ports {
		persisted.Ports = append(persisted.Ports, cce.PortProto{Port: port.Port, Protocol: port.Protocol})
	}

	// Validate the object
	if err := persisted.Validate(); err != nil {
		log.Debugf("Validation failed for %#v: %v", persisted, err)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf("Validation failed: %v", err)))
		if err != nil {
			log.Errf("Error writing response: %v", err)
		}
		return
	}

	// Persist the object
	if err := ctrl.PersistenceService.BulkUpdate(r.Context(), []cce.Persistable{&persisted}); err != nil {
		log.Errf("Error updating entities: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Used for DELETE /apps/{app_id} endpoint
func (g *Gorilla) swagDELETEAppByID(w http.ResponseWriter, r *http.Request) {
	// Load the controller to access the persistence and the payload
	ctrl := r.Context().Value(contextKey("controller")).(*cce.Controller)

	// Check that we can delete the entity
	if statusCode, err := checkDBDeleteApps(r.Context(), ctrl.PersistenceService, mux.Vars(r)["app_id"]); err != nil {
		log.Errf("Error running DB logic: %v", err)
		w.WriteHeader(statusCode)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Errf("Error writing response: %v", err)
		}
		return
	}

	// Fetch the entity from persistence and check if it's there
	persisted, err := ctrl.PersistenceService.Read(r.Context(), mux.Vars(r)["app_id"], &cce.App{})
	if err != nil {
		log.Errf("Error reading entity: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Errf("Error writing response: %v", err)
		}
		return
	}
	if persisted == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ok, err := ctrl.PersistenceService.Delete(r.Context(), mux.Vars(r)["app_id"], &cce.App{})
	if err != nil {
		log.Errf("Error deleting entity: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// we just fetched the entity, so if !ok then something went wrong
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Used for GET /policies endpoint
func (g *Gorilla) swagGETPolicies(w http.ResponseWriter, r *http.Request) {
	// Load the controller to access the persistence
	ctrl := r.Context().Value(contextKey("controller")).(*cce.Controller)

	// Fetch the nodes from persistence
	persisted, err := ctrl.PersistenceService.ReadAll(r.Context(), &cce.TrafficPolicy{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Construct the response object
	policies := swagger.PolicyList{Policies: []swagger.PolicySummary{}}
	for _, a := range persisted {
		policy := swagger.PolicySummary{
			ID:   a.(*cce.TrafficPolicy).ID,
			Name: a.(*cce.TrafficPolicy).Name,
		}
		policies.Policies = append(policies.Policies, policy)
	}

	// Marshal the response object to JSON
	policiesJSON, err := json.Marshal(policies)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(policiesJSON); err != nil {
		log.Errf("Error writing response: %v", err)
	}
}

// Used for POST /policies endpoint
func (g *Gorilla) swagPOSTPolicies(w http.ResponseWriter, r *http.Request) {
	g.trafficPoliciesHandler.create(w, r)
}

// Used for GET /policies/{policy_id} endpoint
func (g *Gorilla) swagGETPolicyByID(w http.ResponseWriter, r *http.Request) {
	// Load the controller to access the persistence
	ctrl := r.Context().Value(contextKey("controller")).(*cce.Controller)

	// Fetch the entity from persistence and check if it's there
	persisted, err := ctrl.PersistenceService.Read(r.Context(), mux.Vars(r)["policy_id"], &cce.TrafficPolicy{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if persisted == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Construct the response object
	policy := swagger.PolicyDetail{
		PolicySummary: swagger.PolicySummary{
			ID:   persisted.(*cce.TrafficPolicy).ID,
			Name: persisted.(*cce.TrafficPolicy).Name,
		},
		Rules: persisted.(*cce.TrafficPolicy).Rules,
	}

	// Marshal the response object to JSON
	policyJSON, err := json.Marshal(policy)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(policyJSON); err != nil {
		log.Errf("Error writing response: %v", err)
	}
}

// Used for PATCH /policies/{policy_id} endpoint
func (g *Gorilla) swagPATCHPolicyByID(w http.ResponseWriter, r *http.Request) {
	// Load the controller to access the persistence and the payload
	ctrl := r.Context().Value(contextKey("controller")).(*cce.Controller)
	body := r.Context().Value(contextKey("body")).([]byte)

	// Unmarshal the payload
	policy := swagger.PolicyDetail{}
	if err := json.Unmarshal(body, &policy); err != nil {
		log.Errf("Error unmarshaling json: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Convert it to a persistable object
	persisted := cce.TrafficPolicy{
		ID:    policy.ID,
		Name:  policy.Name,
		Rules: policy.Rules,
	}

	// Validate the object
	if err := persisted.Validate(); err != nil {
		log.Debugf("Validation failed for %#v: %v", persisted, err)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf("Validation failed: %v", err)))
		if err != nil {
			log.Errf("Error writing response: %v", err)
		}
		return
	}

	// Persist the object
	if err := ctrl.PersistenceService.BulkUpdate(r.Context(), []cce.Persistable{&persisted}); err != nil {
		log.Errf("Error updating entities: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Used for DELETE /policies/{policy_id}
func (g *Gorilla) swagDELETEPolicyByID(w http.ResponseWriter, r *http.Request) {
	// Load the controller to access the persistence and the payload
	ctrl := r.Context().Value(contextKey("controller")).(*cce.Controller)

	// Check that we can delete the entity
	if statusCode, err := checkDBDeleteTrafficPolicies(
		r.Context(), ctrl.PersistenceService, mux.Vars(r)["policy_id"],
	); err != nil {
		log.Errf("Error running DB logic: %v", err)
		w.WriteHeader(statusCode)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Errf("Error writing response: %v", err)
		}
		return
	}

	// Fetch the entity from persistence and check if it's there
	persisted, err := ctrl.PersistenceService.Read(r.Context(), mux.Vars(r)["policy_id"], &cce.TrafficPolicy{})
	if err != nil {
		log.Errf("Error reading entity: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Errf("Error writing response: %v", err)
		}
		return
	}
	if persisted == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ok, err := ctrl.PersistenceService.Delete(r.Context(), mux.Vars(r)["policy_id"], &cce.TrafficPolicy{})
	if err != nil {
		log.Errf("Error deleting entity: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// we just fetched the entity, so if !ok then something went wrong
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Used for GET /nodes/{node_id}/interfaces endpoint
func (g *Gorilla) swagGETInterfaces(w http.ResponseWriter, r *http.Request) {
	// Load the controller to access the persistence
	ctrl := r.Context().Value(contextKey("controller")).(*cce.Controller)

	// Fetch the entity from persistence and check if it's there
	persisted, err := ctrl.PersistenceService.Read(r.Context(), mux.Vars(r)["node_id"], &cce.Node{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if persisted == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get the data from the remote entity
	response, err := handleGetNodes(r.Context(), ctrl.PersistenceService, persisted)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Construct the response object
	ifaces := swagger.InterfaceList{Interfaces: []swagger.InterfaceSummary{}}
	for _, res := range response.(*cce.NodeResp).NetworkInterfaces {
		iface := swagger.InterfaceSummary{
			ID:                res.ID,
			Description:       res.Description,
			Driver:            res.Driver,
			Type:              res.Type,
			MACAddress:        res.MACAddress,
			VLAN:              res.VLAN,
			Zones:             res.Zones,
			FallbackInterface: res.FallbackInterface,
		}
		ifaces.Interfaces = append(ifaces.Interfaces, iface)
	}

	// Marshal the response object to JSON
	ifacesJSON, err := json.Marshal(ifaces)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(ifacesJSON); err != nil {
		log.Errf("Error writing response: %v", err)
	}
}

// Used for PATCH /nodes/{node_id}/interfaces endpoint
func (g *Gorilla) swagPATCHInterfaces(w http.ResponseWriter, r *http.Request) {
	// Load the controller to access the persistence and the payload
	ctrl := r.Context().Value(contextKey("controller")).(*cce.Controller)
	body := r.Context().Value(contextKey("body")).([]byte)

	// Unmarshal the payload
	ifaces := swagger.InterfaceList{}
	if err := json.Unmarshal(body, &ifaces); err != nil {
		log.Errf("Error unmarshaling json: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Fetch the entity from persistence and check if it's there
	persisted, err := ctrl.PersistenceService.Read(r.Context(), mux.Vars(r)["node_id"], &cce.Node{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if persisted == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Convert it to a persistable object
	requested := cce.NodeReq{
		Node:              *persisted.(*cce.Node),
		NetworkInterfaces: []*cce.NetworkInterface{},
	}

	for _, iface := range ifaces.Interfaces {
		p := &cce.NetworkInterface{
			ID:                iface.ID,
			Description:       iface.Description,
			Driver:            iface.Driver,
			Type:              iface.Type,
			MACAddress:        iface.MACAddress,
			VLAN:              iface.VLAN,
			Zones:             iface.Zones,
			FallbackInterface: iface.FallbackInterface,
		}
		requested.NetworkInterfaces = append(requested.NetworkInterfaces, p)
	}

	// Validate the object
	if err = requested.Validate(); err != nil {
		log.Debugf("Validation failed for %#v: %v", requested, err)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf("Validation failed: %v", err)))
		if err != nil {
			log.Errf("Error writing response: %v", err)
		}
		return
	}

	code, err := handleUpdateNodes(r.Context(), ctrl.PersistenceService, &requested)
	switch {
	case code != 0:
		log.Errf("Error updating remote entities: %v", err)
		w.WriteHeader(code)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Errf("Error writing response: %v", err)
		}
		return
	}

	// Persist the object
	if err := ctrl.PersistenceService.BulkUpdate(r.Context(), []cce.Persistable{&requested}); err != nil {
		log.Errf("Error updating entities: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Used for GET /nodes/{node_id}/interfaces/{interface_id} endpoint
func (g *Gorilla) swagGETInterfaceByID(w http.ResponseWriter, r *http.Request) {
	// Load the controller to access the persistence
	ctrl := r.Context().Value(contextKey("controller")).(*cce.Controller)

	// Fetch the entity from persistence and check if it's there
	persisted, err := ctrl.PersistenceService.Read(r.Context(), mux.Vars(r)["node_id"], &cce.Node{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if persisted == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get the data from the remote entity
	response, err := handleGetNodes(r.Context(), ctrl.PersistenceService, persisted)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Construct the response object
	iface := swagger.InterfaceDetail{}
	for _, res := range response.(*cce.NodeResp).NetworkInterfaces {
		if res.ID == mux.Vars(r)["interface_id"] {
			iface = swagger.InterfaceDetail{
				InterfaceSummary: swagger.InterfaceSummary{
					ID:                res.ID,
					Description:       res.Description,
					Driver:            res.Driver,
					Type:              res.Type,
					MACAddress:        res.MACAddress,
					VLAN:              res.VLAN,
					Zones:             res.Zones,
					FallbackInterface: res.FallbackInterface,
				},
			}
		}
	}

	// If its not set, then it's not found
	if iface.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Marshal the response object to JSON
	ifaceJSON, err := json.Marshal(iface)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(ifaceJSON); err != nil {
		log.Errf("Error writing response: %v", err)
	}
}

// Used for GET /nodes/{node_id}/apps/{app_id}/policy endpoint
func (g *Gorilla) swagGETNodeAppPolicy(w http.ResponseWriter, r *http.Request) {
	// Load the controller to access the persistence
	ctrl := r.Context().Value(contextKey("controller")).(*cce.Controller)

	// Filter nodes_apps to get the node_app_id
	nodeApps, err := ctrl.PersistenceService.Filter(
		r.Context(),
		&cce.NodeApp{},
		[]cce.Filter{
			{
				Field: "node_id",
				Value: mux.Vars(r)["node_id"],
			},
			{
				Field: "app_id",
				Value: mux.Vars(r)["app_id"],
			},
		})
	if err != nil {
		log.Errf("Error filtering node_apps: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(nodeApps) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if len(nodeApps) > 1 {
		log.Errf("Filter node_apps returned %d records", len(nodeApps))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Filter nodes_apps_traffic_policies to get the traffic_policy_id
	nodeAppTrafficPolicies, err := ctrl.PersistenceService.Filter(
		r.Context(),
		&cce.NodeAppTrafficPolicy{},
		[]cce.Filter{
			{
				Field: "nodes_apps_id",
				Value: nodeApps[0].GetID(),
			},
		})
	if err != nil {
		log.Errf("Error filtering nodes_apps_traffic_policies: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(nodeAppTrafficPolicies) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if len(nodeAppTrafficPolicies) > 1 {
		log.Errf("Filter nodes_apps_traffic_policies returned %d records", len(nodeAppTrafficPolicies))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Construct the response object
	baseResource := swagger.BaseResource{
		ID: nodeAppTrafficPolicies[0].(*cce.NodeAppTrafficPolicy).TrafficPolicyID,
	}

	// Marshal the response object to JSON
	baseResourceJSON, err := json.Marshal(baseResource)
	if err != nil {
		log.Errf("Error marshaling response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(baseResourceJSON); err != nil {
		log.Errf("Error writing response: %v", err)
	}
}

// Used for PATCH /nodes/{node_id}/apps/{app_id}/policy endpoint
func (g *Gorilla) swagPATCHNodeAppPolicy(w http.ResponseWriter, r *http.Request) { //nolint:gocyclo
	// Load the controller to access the persistence and the payload
	ctrl := r.Context().Value(contextKey("controller")).(*cce.Controller)
	body := r.Context().Value(contextKey("body")).([]byte)

	// Unmarshal the payload
	var baseResource swagger.BaseResource
	if err := json.Unmarshal(body, &baseResource); err != nil {
		log.Errf("Error unmarshaling json: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Filter nodes_apps to get the node_app_id
	nodeApps, err := ctrl.PersistenceService.Filter(
		r.Context(),
		&cce.NodeApp{},
		[]cce.Filter{
			{
				Field: "node_id",
				Value: mux.Vars(r)["node_id"],
			},
			{
				Field: "app_id",
				Value: mux.Vars(r)["app_id"],
			},
		})
	if err != nil {
		log.Errf("Error filtering node_apps: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(nodeApps) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if len(nodeApps) > 1 {
		log.Errf("Filter node_apps returned %d records", len(nodeApps))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Query traffic_policies to verify the baseResourceID is valid
	policy, err := ctrl.PersistenceService.Read(r.Context(), baseResource.ID, &cce.TrafficPolicy{})
	if err != nil {
		log.Errf("Error reading traffic_policies: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if policy == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Connect to node
	nodePort := ctrl.ELAPort
	if nodePort == "" {
		nodePort = defaultELAPort
	}
	nodeCC, err := connectNode(
		r.Context(),
		ctrl.PersistenceService,
		nodeApps[0].(*cce.NodeApp),
		nodePort,
		ctrl.EdgeNodeCreds)
	if err != nil {
		log.Errf("Error connecting to node: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Make gRPC call to node to set the policy
	if err = nodeCC.AppPolicySvcCli.Set(
		r.Context(),
		nodeApps[0].(*cce.NodeApp).AppID,
		policy.(*cce.TrafficPolicy),
	); err != nil {
		log.Errf("Error setting policy: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Filter nodes_apps_traffic_policies to see if a record already exists
	nodeAppPolicies, err := ctrl.PersistenceService.Filter(
		r.Context(),
		&cce.NodeAppTrafficPolicy{},
		[]cce.Filter{
			{
				Field: "nodes_apps_id",
				Value: nodeApps[0].GetID(),
			},
		})
	if err != nil {
		log.Errf("Error reading nodes_apps_traffic_policies: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// If it exists, delete it
	if len(nodeAppPolicies) == 1 {
		ok, err := ctrl.PersistenceService.Delete(r.Context(), nodeAppPolicies[0].GetID(), &cce.NodeAppTrafficPolicy{})
		if err != nil {
			log.Errf("Error deleting from nodes_apps_traffic_policies: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !ok {
			log.Err("Did not delete 1 record from nodes_apps_traffic_policies")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// Convert the base resource to a persistable object
	persisted := cce.NodeAppTrafficPolicy{
		ID:              uuid.New(),
		NodeAppID:       nodeApps[0].GetID(),
		TrafficPolicyID: baseResource.ID,
	}

	// Persist the object
	if err := ctrl.PersistenceService.Create(r.Context(), &persisted); err != nil {
		log.Errf("Error creating entity: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Used for DELETE /nodes/{node_id}/apps/{app_id}/policy endpoint
func (g *Gorilla) swagDELETENodeAppPolicy(w http.ResponseWriter, r *http.Request) {
	// Load the controller to access the persistence and the payload
	ctrl := r.Context().Value(contextKey("controller")).(*cce.Controller)

	// Filter nodes_apps to get the node_app_id
	nodeApps, err := ctrl.PersistenceService.Filter(
		r.Context(),
		&cce.NodeApp{},
		[]cce.Filter{
			{
				Field: "node_id",
				Value: mux.Vars(r)["node_id"],
			},
			{
				Field: "app_id",
				Value: mux.Vars(r)["app_id"],
			},
		})
	if err != nil {
		log.Errf("Error filtering node_apps: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(nodeApps) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if len(nodeApps) > 1 {
		log.Errf("Filter node_apps returned %d records", len(nodeApps))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Filter nodes_apps_traffic_policies to get the ID
	nodeAppPolicies, err := ctrl.PersistenceService.Filter(
		r.Context(),
		&cce.NodeAppTrafficPolicy{},
		[]cce.Filter{
			{
				Field: "nodes_apps_id",
				Value: nodeApps[0].GetID(),
			},
		})
	if err != nil {
		log.Errf("Error reading nodes_apps_traffic_policies: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Connect to node
	nodePort := ctrl.ELAPort
	if nodePort == "" {
		nodePort = defaultELAPort
	}
	nodeCC, err := connectNode(
		r.Context(),
		ctrl.PersistenceService,
		nodeApps[0].(*cce.NodeApp),
		nodePort,
		ctrl.EdgeNodeCreds)
	if err != nil {
		log.Errf("Error connecting to node: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Make gRPC call to node to delete the policy
	if err = nodeCC.AppPolicySvcCli.Delete(
		r.Context(),
		nodeApps[0].(*cce.NodeApp).AppID,
	); err != nil {
		log.Errf("Error deleting policy: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Delete the resource
	ok, err := ctrl.PersistenceService.Delete(r.Context(), nodeAppPolicies[0].GetID(), &cce.NodeAppTrafficPolicy{})
	if err != nil {
		log.Errf("Error deleting from nodes_apps_traffic_policies: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		log.Err("Did not delete 1 record from nodes_apps_traffic_policies")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
