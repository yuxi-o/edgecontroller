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
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	cce "github.com/smartedgemec/controller-ce"
)

type handler struct {
	model cce.Entity

	// these funcs provide db checks
	checkDBCreate func(
		context.Context,
		cce.PersistenceService,
		cce.Entity,
	) (statusCode int, err error)
	checkDBDelete func(
		ctx context.Context,
		ps cce.PersistenceService,
		id string,
	) (statusCode int, err error)

	// these funcs handle node logic
	handleCreate func(context.Context) error
}

func (h *handler) create(w http.ResponseWriter, r *http.Request) { //nolint:gocyclo,lll
	var (
		ctrl = r.Context().Value(contextKey("controller")).(*cce.Controller)
		body = r.Context().Value(contextKey("body")).([]byte)
		e    cce.Entity
		err  error
	)

	e = reflect.New(
		reflect.ValueOf(h.model).Elem().Type()).Interface().(cce.Entity)
	if err = json.Unmarshal(body, e); err != nil {
		log.Printf("Error unmarshalling json: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if e.GetID() != "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(
			"Validation failed: id cannot be specified in POST request"))
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
		return
	}

	e.SetID(uuid.NewV4().String())

	if err = e.Validate(); err != nil {
		log.Printf("Validation failed for %#v: %v", e, err)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf("Validation failed: %v", err)))
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
		return
	}

	if h.checkDBCreate != nil {
		var statusCode int
		if statusCode, err = h.checkDBCreate(
			r.Context(),
			ctrl.PersistenceService,
			e,
		); err != nil {
			log.Printf("Error running DB logic: %v", err)
			w.WriteHeader(statusCode)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Printf("Error writing response: %v", err)
			}
			return
		}
	}

	if err = ctrl.PersistenceService.Create(r.Context(), e); err != nil {
		log.Printf("Error creating entity: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header()["Content-Type"] = []string{"application/json"}
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(fmt.Sprintf(`{"id":"%s"}`, e.GetID())))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (h *handler) getAll(w http.ResponseWriter, r *http.Request) {
	var (
		ctrl  = r.Context().Value(contextKey("controller")).(*cce.Controller)
		es    []cce.Entity
		e     cce.Entity
		err   error
		bytes []byte
	)

	if es, err = ctrl.PersistenceService.ReadAll(
		r.Context(),
		h.model,
	); err != nil {
		log.Printf("Error reading entities: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytes = append(bytes, byte('['))
	for _, e = range es {
		var appBytes []byte
		if appBytes, err = json.Marshal(e); err != nil {
			log.Printf("Error marshalling json: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		bytes = append(bytes, appBytes...)
		bytes = append(bytes, byte(','))
	}

	if len(bytes) > 1 {
		bytes = bytes[:len(bytes)-1]
	}
	bytes = append(bytes, byte(']'))

	w.Header()["Content-Type"] = []string{"application/json"}
	if _, err = w.Write(bytes); err != nil {
		log.Printf("Error writing response: %v", err)
		return
	}
}

func (h *handler) getByFilter(w http.ResponseWriter, r *http.Request) {
	var (
		ctrl  = r.Context().Value(contextKey("controller")).(*cce.Controller)
		es    []cce.Entity
		e     cce.Entity
		err   error
		bytes []byte
	)

	// TODO parse from request
	if es, err = ctrl.PersistenceService.Filter(
		r.Context(),
		h.model,
		[]cce.Filter{
			{Field: "node_id", Value: "9112538c-4df3-4a7a-a7e6-5db9ec203d03"},
			{Field: "node_id", Value: "9112538c-4df3-4a7a-a7e6-5db9ec203d03"},
		},
	); err != nil {
		log.Printf("Error getting entityRoutes: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytes = append(bytes, byte('['))
	for _, e = range es {
		var appBytes []byte
		appBytes, err = json.Marshal(e)
		if err != nil {
			log.Printf("Error marshalling json: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		bytes = append(bytes, appBytes...)
		bytes = append(bytes, byte(','))
	}

	if len(bytes) > 1 {
		bytes = bytes[:len(bytes)-1]
	}
	bytes = append(bytes, byte(']'))

	w.Header()["Content-Type"] = []string{"application/json"}
	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("Error writing response: %v", err)
		return
	}
}

func (h *handler) getByID(w http.ResponseWriter, r *http.Request) {
	var (
		ctrl  = r.Context().Value(contextKey("controller")).(*cce.Controller)
		id    = mux.Vars(r)["id"]
		e     cce.Entity
		err   error
		bytes []byte
	)

	if e, err = ctrl.PersistenceService.Read(
		r.Context(),
		id,
		h.model,
	); err != nil {
		log.Printf("Error reading entity: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if e == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if bytes, err = json.Marshal(e); err != nil {
		log.Printf("Error marshalling json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header()["Content-Type"] = []string{"application/json"}
	if _, err = w.Write(bytes); err != nil {
		log.Printf("Error writing response: %v", err)
		return
	}
}

func (h *handler) bulkUpdate(
	w http.ResponseWriter,
	r *http.Request,
) {
	var (
		ctrl = r.Context().Value(contextKey("controller")).(*cce.Controller)
		body = r.Context().Value(contextKey("body")).([]byte)
		ies  []interface{}
		ie   interface{}
		e    cce.Entity
		err  error
	)

	if err = json.Unmarshal(body, &ies); err != nil {
		log.Printf("Error unmarshalling json: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var es []cce.Entity
	for _, ie = range ies {
		var bytes []byte
		if bytes, err = json.Marshal(ie); err != nil {
			log.Printf("Error marshalling json: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		e = reflect.New(reflect.ValueOf(h.model).Elem().Type()).
			Interface().(cce.Entity)
		if err = json.Unmarshal(bytes, e); err != nil {
			log.Printf("Error unmarshalling json: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err = e.Validate(); err != nil {
			log.Printf("Validation failed for %#v: %v", e, err)
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(fmt.Sprintf("Validation failed: %v", err)))
			if err != nil {
				log.Printf("Error writing response: %v", err)
			}
			return
		}

		es = append(es, e)
	}

	if err = ctrl.PersistenceService.BulkUpdate(r.Context(), es); err != nil {
		log.Printf("Error updating entities: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) delete(w http.ResponseWriter, r *http.Request) {
	var (
		ctrl = r.Context().Value(contextKey("controller")).(*cce.Controller)
		id   = mux.Vars(r)["id"]
		ok   bool
		err  error
	)

	if h.checkDBDelete != nil {
		var statusCode int
		if statusCode, err = h.checkDBDelete(
			r.Context(),
			ctrl.PersistenceService,
			id,
		); err != nil {
			log.Printf("Error running DB logic: %v", err)
			w.WriteHeader(statusCode)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Printf("Error writing response: %v", err)
			}
			return
		}
	}

	if ok, err = ctrl.PersistenceService.Delete(
		r.Context(),
		id,
		h.model,
	); err != nil {
		log.Printf("Error deleting entity: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}