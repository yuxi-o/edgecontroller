// Copyright 2019 Intel Corporation. All rights reserved
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

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"k8s.io/klog"
	"time"
)

// Connectivity constants
const (
	AFServer  = "http://localhost:80/"
	OAMServer = "http://localhost:80/"
)

// HTTP client
var client = &http.Client{
	Timeout: 10 * time.Second,
}

// AFCreateSubscription create new Traffic Influence Subscription at AF
func AFCreateSubscription(sub TrafficInfluSub) error {

	subBytes, err := json.Marshal(sub)
	if err != nil {
		klog.Error("Failed to marshal TrafficInfluSub:", err)
		return err
	}

	req, err := http.NewRequest("POST",
					AFServer + "AFTransactions",
					bytes.NewReader(subBytes))
	if err != nil {
		klog.Error("Create request failed:", err)
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		klog.Error("Create request transmission failed:", err)
		return err
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}
