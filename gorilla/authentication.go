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
	"net/http"
	"strings"

	cce "github.com/open-ness/edgecontroller"
)

func authenticate(w http.ResponseWriter, r *http.Request) {
	var (
		ctrl = r.Context().Value(contextKey("controller")).(*cce.Controller)
		body = r.Context().Value(contextKey("body")).([]byte)
	)

	// Extract the username and password from JSON
	var u cce.AuthCreds
	if err := json.Unmarshal(body, &u); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Verify the user name and password
	if u.Username != ctrl.AdminCreds.Username {
		log.Debugf("Unsuccessful login attempt for user '%s'", u.Username)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	if u.Password != ctrl.AdminCreds.Password {
		log.Debugf("Unsuccessful login attempt for user '%s'", u.Username)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	log.Debugf("Successfully authenticated user: %s", u.Username)

	// Create an auth token
	token, err := ctrl.TokenService.Issue()
	if err != nil {
		log.Debugf("Error signing authentication token: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Wrap auth token in JSON
	bytes, err := json.Marshal(
		struct {
			Token string `json:"token"`
		}{
			token,
		})
	if err != nil {
		log.Errf("Error marshaling authentication token: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with status code 201
	w.WriteHeader(http.StatusCreated)

	// Return JSON-encoded auth token
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(bytes); err != nil {
		log.Errf("Error writing response: %v", err)
	}
}

// requireAuthHandler is a handler that only allows HTTP requests with a valid
// JSON Web Token issued by the Controller Token Authentication service.
func requireAuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctrl := r.Context().Value(contextKey("controller")).(*cce.Controller)

		// Get the Authorization header
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Extract the auth token
		bearer := strings.Split(auth, " ")
		if len(bearer) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Validate the auth token
		err := ctrl.TokenService.Validate(bearer[1])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
