package gorilla

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	cce "github.com/smartedgemec/controller-ce"
	swagger "github.com/smartedgemec/controller-ce/swagger"
	"net/http"
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
		log.Errf("Error unmarshalling json: %v", err)
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

// Used for DELETE /apps/{app_id}
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
