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
	NGCOAMServer = "http://localhost:8081"
	NGCAFServer  = "http://localhost:8080"
	LteOAMServer = "http://localhost:8082"
)

// HTTP client
var client = &http.Client{
	Timeout: 10 * time.Second,
}

// OAM5gRegisterAFService register controller to AF services registry
func OAM5gRegisterAFService(service []byte) (string, error) {
	var afService string
	req, err := http.NewRequest("POST",
		NGCOAMServer+"/oam/v1/af/services",
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
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return afService, err
		}
		afService = string(b)
	} else {
		return afService, fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	return afService, nil
}

// AFCreateSubscription create new Traffic Influence Subscription at AF
func AFCreateSubscription(sub []byte) (string, error) {
	var subID string
	req, err := http.NewRequest("POST",
		NGCAFServer+"/CNCA/1.0.1/subscriptions",
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
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return subID, err
		}
		subID = string(b)
	} else {
		return subID, fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	return subID, nil
}

// AFPatchSubscription update an active subscription for the AF
func AFPatchSubscription(subID string, sub []byte) error {

	req, err := http.NewRequest("PATCH",
		NGCAFServer+"/CNCA/1.0.1/subscriptions/"+subID,
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
			NGCAFServer+"/CNCA/1.0.1/subscriptions", nil)
		if err != nil {
			return sub, err
		}
	} else {
		req, err = http.NewRequest("GET",
			NGCAFServer+"/CNCA/1.0.1/subscriptions/"+subID, nil)
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
		NGCAFServer+"/CNCA/1.0.1/subscriptions/"+subID, nil)
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

// LteCreateUserplane create new LTE userplane
func LteCreateUserplane(up []byte) (string, error) {
	var ID string
	req, err := http.NewRequest("POST",
		LteOAMServer+"/userplanes",
		bytes.NewReader(up))
	if err != nil {
		return ID, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return ID, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return ID, err
		}
		ID = string(b)
	} else {
		return ID, fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	return ID, nil
}

// LtePatchUserplane update an active LTE CUPS userplane
func LtePatchUserplane(upID string, up []byte) error {

	req, err := http.NewRequest("PATCH",
		LteOAMServer+"/userplanes/"+upID,
		bytes.NewReader(up))
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

// LteGetUserplane get the active CUPS userplane
func LteGetUserplane(upID string) ([]byte, error) {
	var up []byte
	var req *http.Request
	var err error

	if upID == "all" {
		req, err = http.NewRequest("GET",
			LteOAMServer+"/userplanes", nil)
		if err != nil {
			return up, err
		}
	} else {
		req, err = http.NewRequest("GET",
			LteOAMServer+"/userplanes/"+upID, nil)
		if err != nil {
			return up, err
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return up, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		up, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return up, err
		}
		return up, nil
	}
	return up, fmt.Errorf("HTTP failure: %d", resp.StatusCode)
}

// LteDeleteUserplane delete an active LTE CUPS userplane
func LteDeleteUserplane(upID string) error {

	req, err := http.NewRequest("DELETE",
		LteOAMServer+"/userplanes/"+upID, nil)
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
