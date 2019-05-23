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
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	cce "github.com/smartedgemec/controller-ce"
)

type handler struct {
	model    cce.Persistable
	reqModel cce.ReqEntity

	// these funcs provide db constraint (unique/foreign key) checks
	checkDBCreate func(
		context.Context,
		cce.PersistenceService,
		cce.Persistable,
	) (statusCode int, err error)
	checkDBDelete func(
		ctx context.Context,
		ps cce.PersistenceService,
		id string,
	) (statusCode int, err error)

	// these funcs provide application logic
	handleCreate func(
		context.Context,
		cce.PersistenceService,
		cce.Persistable,
	) error
	handleGet func(
		context.Context,
		cce.PersistenceService,
		cce.Persistable,
	) (cce.RespEntity, error)
	handleUpdate func(
		context.Context,
		cce.PersistenceService,
		cce.Validatable,
	) error
	handleDelete func(
		context.Context,
		cce.PersistenceService,
		cce.Persistable,
	) error
}

func (h *handler) create(w http.ResponseWriter, r *http.Request) { //nolint:gocyclo
	ctrl := r.Context().Value(contextKey("controller")).(*cce.Controller)
	body := r.Context().Value(contextKey("body")).([]byte)

	p := reflect.New(reflect.ValueOf(h.model).Elem().Type()).Interface().(cce.Persistable)
	if err := json.Unmarshal(body, p); err != nil {
		log.Errf("Error unmarshalling json: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if p.GetID() != "" {
		w.WriteHeader(http.StatusBadRequest)

		if _, err := w.Write([]byte("Validation failed: id cannot be specified in POST request")); err != nil {
			log.Errf("Error writing response: %v", err)
		}
		return
	}

	p.SetID(uuid.NewV4().String())

	if err := p.(cce.Validatable).Validate(); err != nil {
		log.Debugf("Validation failed for %#v: %v", p, err)
		w.WriteHeader(http.StatusBadRequest)
		if _, err = w.Write([]byte(fmt.Sprintf("Validation failed: %v", err))); err != nil {
			log.Errf("Error writing response: %v", err)
		}
		return
	}

	if h.checkDBCreate != nil {
		if statusCode, err := h.checkDBCreate(r.Context(), ctrl.PersistenceService, p); err != nil {
			log.Errf("Error checking DB create: %v", err)
			w.WriteHeader(statusCode)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Errf("Error writing response: %v", err)
			}
			return
		}
	}

	if h.handleCreate != nil {
		if err := h.handleCreate(r.Context(), ctrl.PersistenceService, p); err != nil {
			log.Errf("Error handling create logic: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Errf("Error writing response: %v", err)
			}
			return
		}
	}

	if err := ctrl.PersistenceService.Create(r.Context(), p); err != nil {
		log.Errf("Error creating entity: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header()["Content-Type"] = []string{"application/json"}
	w.WriteHeader(http.StatusCreated)

	if _, err := w.Write([]byte(fmt.Sprintf(`{"id":"%s"}`, p.GetID()))); err != nil {
		log.Errf("Error writing response: %v", err)
	}
}

func (h *handler) filter(w http.ResponseWriter, r *http.Request) {
	ctrl := r.Context().Value(contextKey("controller")).(*cce.Controller)

	var filters []cce.Filter
	for k, v := range r.URL.Query() {
		filters = append(filters, cce.Filter{Field: k, Value: v[0]})
	}

	ps, err := ctrl.PersistenceService.Filter(r.Context(), h.model, filters)
	if err != nil {
		log.Errf("Error reading entities: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var res []cce.RespEntity
	for _, p := range ps {
		if h.handleGet != nil {
			re, err := h.handleGet(r.Context(), ctrl.PersistenceService, p)
			if err != nil {
				log.Errf("Error handling get logic: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				if _, err = w.Write([]byte(err.Error())); err != nil {
					log.Errf("Error writing response: %v", err)
					return
				}
			}
			res = append(res, re)
		} else {
			res = append(res, p)
		}
	}

	var bytes []byte
	bytes = append(bytes, byte('['))
	for _, re := range res {
		var appBytes []byte
		appBytes, err := json.Marshal(re)
		if err != nil {
			log.Errf("Error marshalling json: %v", err)
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
	if _, err := w.Write(bytes); err != nil {
		log.Errf("Error writing response: %v", err)
		return
	}
}

func (h *handler) getByID(w http.ResponseWriter, r *http.Request) {
	ctrl := r.Context().Value(contextKey("controller")).(*cce.Controller)

	id := mux.Vars(r)["id"]
	if id == "" {
		// TODO add test for this
		log.Debug("ID missing from request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p, err := ctrl.PersistenceService.Read(r.Context(), id, h.model)
	if err != nil {
		log.Errf("Error reading entity: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if p == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var re cce.RespEntity
	if h.handleGet != nil {
		re, err = h.handleGet(r.Context(), ctrl.PersistenceService, p)
		if err != nil {
			log.Errf("Error handling get logic: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Errf("Error writing response: %v", err)
			}
			return
		}
	} else {
		re = p
	}

	bytes, err := json.Marshal(re)
	if err != nil {
		log.Errf("Error marshalling json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header()["Content-Type"] = []string{"application/json"}
	if _, err = w.Write(bytes); err != nil {
		log.Errf("Error writing response: %v", err)
		return
	}
}

func (h *handler) bulkUpdate(w http.ResponseWriter, r *http.Request) { //nolint:gocyclo
	ctrl := r.Context().Value(contextKey("controller")).(*cce.Controller)
	body := r.Context().Value(contextKey("body")).([]byte)

	var is []interface{}
	if err := json.Unmarshal(body, &is); err != nil {
		log.Errf("Error unmarshalling json: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var ps []cce.Persistable
	for _, i := range is {
		bytes, err := json.Marshal(i)
		if err != nil {
			log.Errf("Error marshalling json: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var v cce.Validatable
		if h.reqModel != nil {
			v = reflect.New(reflect.ValueOf(h.reqModel).Elem().Type()).Interface().(cce.Validatable)
		} else {
			v = reflect.New(reflect.ValueOf(h.model).Elem().Type()).Interface().(cce.Validatable)
		}

		if err := json.Unmarshal(bytes, &v); err != nil {
			log.Errf("Error unmarshalling json: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := v.Validate(); err != nil {
			log.Debugf("Validation failed for %#v: %v", v, err)
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(fmt.Sprintf("Validation failed: %v", err)))
			if err != nil {
				log.Errf("Error writing response: %v", err)
			}
			return
		}

		if h.handleUpdate != nil {
			if err := h.handleUpdate(r.Context(), ctrl.PersistenceService, v); err != nil {
				log.Errf("Error handling update logic: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				_, err = w.Write([]byte(err.Error()))
				if err != nil {
					log.Errf("Error writing response: %v", err)
				}
				return
			}
		}

		ps = append(ps, v.(cce.Persistable))
	}

	if err := ctrl.PersistenceService.BulkUpdate(r.Context(), ps); err != nil {
		log.Errf("Error updating entities: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) delete(w http.ResponseWriter, r *http.Request) { //nolint:gocyclo
	ctrl := r.Context().Value(contextKey("controller")).(*cce.Controller)

	id := mux.Vars(r)["id"]
	if id == "" {
		// TODO add test for this
		log.Debug("ID missing from request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if h.checkDBDelete != nil {
		if statusCode, err := h.checkDBDelete(r.Context(), ctrl.PersistenceService, id); err != nil {
			log.Errf("Error running DB logic: %v", err)
			w.WriteHeader(statusCode)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Errf("Error writing response: %v", err)
			}
			return
		}
	}

	p, err := ctrl.PersistenceService.Read(r.Context(), id, h.model)
	if err != nil {
		log.Errf("Error reading entity: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Errf("Error writing response: %v", err)
		}
		return
	}

	if p == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if h.handleDelete != nil {
		if err = h.handleDelete(r.Context(), ctrl.PersistenceService, p); err != nil {
			log.Errf("Error handling delete logic: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Errf("Error writing response: %v", err)
			}
			return
		}
	}

	ok, err := ctrl.PersistenceService.Delete(r.Context(), id, h.model)
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
