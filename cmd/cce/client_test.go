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

package main_test

import (
	"fmt"
	"io"
	"net/http"
)

// apiClient is a Controller API client that automatically injects an auth token
// into the request headers of each HTTP request in the OAuth 2.0 Bearer Token
// standard format (https://tools.ietf.org/html/rfc6750).
type apiClient struct {
	// Token is a JSON Web Token.
	Token string
}

// Get sends a HTTP GET request with a token and returns an HTTP response.
func (cli apiClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return new(http.Client).Do(cli.injectToken(req))
}

// Post sends a HTTP POST request with a token and returns an HTTP response.
func (cli apiClient) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	return new(http.Client).Do(cli.injectToken(req))
}

// Patch sends a HTTP PATCH request with a token and returns an HTTP response.
func (cli apiClient) Patch(url string, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return nil, err
	}

	return new(http.Client).Do(cli.injectToken(req))
}

// Delete sends a HTTP DELETE request with a token and returns an HTTP response.
func (cli apiClient) Delete(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	return new(http.Client).Do(cli.injectToken(req))
}

func (cli apiClient) injectToken(r *http.Request) *http.Request {
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cli.Token))
	return r
}
