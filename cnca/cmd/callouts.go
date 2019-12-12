// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package cnca

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Connectivity constants
const (
	NgcOAMServiceEndpoint = "http://localhost:30070/ngcoam/v1/af"
	NgcAFServiceEndpoint  = "http://localhost:30050/af/v1"
	LteOAMServiceEndpoint = "http://localhost:8082/"
)

// HTTP client
var client = &http.Client{
	Timeout: 10 * time.Second,
}

// OAM5gRegisterAFService register controller to AF services registry
func OAM5gRegisterAFService(locService []byte) (string, error) {

	req, err := http.NewRequest("POST",
		NgcOAMServiceEndpoint+"/services",
		bytes.NewReader(locService))
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		var s AFServiceID
		err = json.Unmarshal(b, &s)
		if err != nil {
			return "", err
		}
		return s.AFServiceID, nil
	} else {
		return "", fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}
}

// OAM5gUnregisterAFService unregister controller from AF services registry
func OAM5gUnregisterAFService(serviceID string) error {

	req, err := http.NewRequest("DELETE",
		NgcOAMServiceEndpoint+"/services/"+serviceID, nil)
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

// AFCreateSubscription create new Traffic Influence Subscription at AF
func AFCreateSubscription(sub []byte) (string, error) {

	req, err := http.NewRequest("POST",
		NgcAFServiceEndpoint+"/subscriptions",
		bytes.NewReader(sub))
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var s SubscriptionID
	err = json.Unmarshal(b, &s)
	if err != nil {
		return "", err
	}
	return s.ID, nil
}

// AFPatchSubscription update an active subscription for the AF
func AFPatchSubscription(subID string, sub []byte) error {

	req, err := http.NewRequest("PATCH",
		NgcAFServiceEndpoint+"/subscriptions/"+subID,
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
			NgcAFServiceEndpoint+"/subscriptions", nil)
		if err != nil {
			return sub, err
		}
	} else {
		req, err = http.NewRequest("GET",
			NgcAFServiceEndpoint+"/subscriptions/"+subID, nil)
		if err != nil {
			return sub, err
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return sub, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return sub, fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	sub, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return sub, err
	}
	return sub, nil
}

// AFDeleteSubscription delete an active Traffic Influence Subscription for the AF
func AFDeleteSubscription(subID string) error {

	req, err := http.NewRequest("DELETE",
		NgcAFServiceEndpoint+"/subscriptions/"+subID, nil)
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
	req, err := http.NewRequest("POST",
		LteOAMServiceEndpoint+"/userplanes",
		bytes.NewReader(up))
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var u CupsUserplaneID
	err = json.Unmarshal(b, &u)
	if err != nil {
		return "", err
	}

	return u.ID, nil
}

// LtePatchUserplane update an active LTE CUPS userplane
func LtePatchUserplane(upID string, up []byte) error {

	req, err := http.NewRequest("PATCH",
		LteOAMServiceEndpoint+"/userplanes/"+upID,
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
			LteOAMServiceEndpoint+"/userplanes", nil)
		if err != nil {
			return up, err
		}
	} else {
		req, err = http.NewRequest("GET",
			LteOAMServiceEndpoint+"/userplanes/"+upID, nil)
		if err != nil {
			return up, err
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return up, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return up, fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	up, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return up, err
	}
	return up, nil
}

// LteDeleteUserplane delete an active LTE CUPS userplane
func LteDeleteUserplane(upID string) error {

	req, err := http.NewRequest("DELETE",
		LteOAMServiceEndpoint+"/userplanes/"+upID, nil)
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
