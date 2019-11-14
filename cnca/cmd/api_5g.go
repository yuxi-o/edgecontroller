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

package cnca

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Connectivity constants
const (
	AFServer    = "http://localhost:8080"
	OAM5gServer = "http://localhost:8081"
)

// HTTP client
var client = &http.Client{
	Timeout: 10 * time.Second,
}

// OAM5gRegisterAFService register controller to AF services registry
func OAM5gRegisterAFService(service []byte) (string, error) {

	var afService string

	req, err := http.NewRequest("POST",
		OAM5gServer + "/oam/v1/af/services",
		bytes.NewReader(service))
	if err != nil {
		return afService, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return afService, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		bodyB, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return afService, err
		}
		afService = string(bodyB)
	} else {
		return afService, fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	return afService, nil
}

// AFCreateSubscription create new Traffic Influence Subscription at AF
func AFCreateSubscription(sub []byte) (string, error) {

	var subID string

	req, err := http.NewRequest("POST",
		AFServer + "/CNCA/1.0.1/subscriptions",
		bytes.NewReader(sub))
	if err != nil {
		return subID, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return subID, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		bodyB, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return subID, err
		}
		subID = string(bodyB)
	} else {
		return subID, fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	return subID, nil
}

// AFPatchSubscription update an active subscription for the AF
func AFPatchSubscription(subID string, sub []byte) error {

	req, err := http.NewRequest("PATCH",
		AFServer + "/CNCA/1.0.1/subscriptions/" + subID,
		bytes.NewReader(sub))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	return nil
}

// AFGetSubscription get the active Traffic Influence Subscription for the AF
func AFGetSubscription(subID string) ([]byte, error) {
	var sub []byte
	var req *http.Request
	var err error

	if subID == "all" {
		req, err = http.NewRequest("GET",
		AFServer + "/CNCA/1.0.1/subscriptions/", nil)
		if err != nil {
			return sub, err
		}
	} else {
		req, err = http.NewRequest("GET",
		AFServer + "/CNCA/1.0.1/subscriptions/" + subID, nil)
		if err != nil {
			return sub, err
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return sub, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		sub, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return sub, err
		}
		return sub, nil
	}
	return sub, fmt.Errorf("HTTP failure: %d", resp.StatusCode)
}

// AFDeleteSubscription delete an active Traffic Influence Subscription for the AF
func AFDeleteSubscription(subID string) error {

	req, err := http.NewRequest("DELETE",
		AFServer + "/CNCA/1.0.1/subscriptions/" + subID, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	return nil
}
