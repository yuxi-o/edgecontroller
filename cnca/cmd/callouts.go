// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package cnca

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Connectivity constants
const (
	NgcOAMServiceEndpoint      = "http://localhost:30070/ngcoam/v1/af"
	NgcAFServiceEndpoint       = "http://localhost:30050/af/v1"
	LteOAMServiceEndpoint      = "http://localhost:8082/"
	NgcOAMServiceHttp2Endpoint = "https://localhost:30070/ngcoam/v1/af"
	NgcAFServiceHttp2Endpoint  = "https://localhost:30050/af/v1"
	LteOAMServiceHttp2Endpoint = "https://localhost:8082/"
)

// HTTP client
var client http.Client

func getNgcOAMServiceUrl() string {
	if UseHttpProtocol == HTTP2 {
		return NgcOAMServiceHttp2Endpoint + "/services"
	} else {
		return NgcOAMServiceEndpoint + "/services"
	}
}

func getNgcAFServiceUrl() string {
	if UseHttpProtocol == HTTP2 {
		return NgcAFServiceHttp2Endpoint + "/subscriptions"
	} else {
		return NgcAFServiceEndpoint + "/subscriptions"
	}
}

func getNgcAFPfdServiceUrl() string {
	if UseHttpProtocol == HTTP2 {
		return NgcAFServiceHttp2Endpoint + "/pfd/transactions"
	} else {
		return NgcAFServiceEndpoint + "/pfd/transactions"
	}
}

func getLteOAMServiceUrl() string {
	if UseHttpProtocol == HTTP2 {
		return LteOAMServiceHttp2Endpoint + "/userplanes"
	} else {
		return LteOAMServiceEndpoint + "/userplanes"
	}
}

// OAM5gRegisterAFService register controller to AF services registry
func OAM5gRegisterAFService(locService []byte) (string, error) {

	url := getNgcOAMServiceUrl()

	req, err := http.NewRequest("POST", url, bytes.NewReader(locService))
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

	url := getNgcOAMServiceUrl() + "/" + serviceID

	req, err := http.NewRequest("DELETE", url, nil)
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

	url := getNgcAFServiceUrl()

	req, err := http.NewRequest("POST", url, bytes.NewReader(sub))
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

	// retrive URI of the newly created subscription from response header
	subLoc := resp.Header.Get("Location")
	if subLoc == "" {
		return "", fmt.Errorf("Empty subscription URI returned from AF")
	}
	return subLoc, nil
}

// AFPatchSubscription update an active subscription for the AF
func AFPatchSubscription(subID string, sub []byte) error {

	url := getNgcAFServiceUrl() + "/" + subID

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(sub))
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
	var url string

	if subID == "all" {
		url = getNgcAFServiceUrl()
	} else {
		url = getNgcAFServiceUrl() + "/" + subID
	}

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return sub, err
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

// AFDeleteSubscription delete an active Traffic Influence Subscription for AF
func AFDeleteSubscription(subID string) error {

	url := getNgcAFServiceUrl() + "/" + subID

	req, err := http.NewRequest("DELETE", url, nil)
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

	url := getLteOAMServiceUrl()

	req, err := http.NewRequest("POST", url, bytes.NewReader(up))
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

	url := getLteOAMServiceUrl() + "/" + upID

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(up))
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
	var url string

	if upID == "all" {
		url = getLteOAMServiceUrl()
	} else {
		url = getLteOAMServiceUrl() + "/" + upID
	}

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return up, err
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

	url := getLteOAMServiceUrl() + "/" + upID

	req, err := http.NewRequest("DELETE", url, nil)
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

// AFCreatePfdTransaction create new PFD transaction at AF
func AFCreatePfdTransaction(trans []byte) ([]byte, string, error) {

	var pfdData []byte

	url := getNgcAFPfdServiceUrl()

	req, err := http.NewRequest("POST", url, bytes.NewReader(trans))
	if err != nil {
		return nil, "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, "", fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	// retrive URI of the newly created transaction from response header
	self := resp.Header.Get("Self")
	if resp.Body != nil {
		pfdData, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, "", err
		}
	}

	return pfdData, self, nil
}

// AFGetPfdTransaction get the active PFD Transaction for the AF
func AFGetPfdTransaction(transID string) ([]byte, error) {
	var trans []byte
	var req *http.Request
	var err error
	var url string

	if transID == "all" {
		url = getNgcAFPfdServiceUrl()
	} else {
		url = getNgcAFPfdServiceUrl() + "/" + transID
	}

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return trans, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return trans, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return trans, fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	trans, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return trans, err
	}
	return trans, nil
}

// AFPatchPfdTransaction update an active PFD Transaction for the AF
func AFPatchPfdTransaction(transID string, trans []byte) error {

	url := getNgcAFPfdServiceUrl() + "/" + transID

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(trans))
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

// AFDeletePfdTransaction delete an active PFD Transaction for the AF
func AFDeletePfdTransaction(transID string) error {

	url := getNgcAFPfdServiceUrl() + "/" + transID

	req, err := http.NewRequest("DELETE", url, nil)
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

// AFGetPfdApplication get the active PFD Application for the AF
func AFGetPfdApplication(transID string, appID string) ([]byte, error) {
	var trans []byte
	var req *http.Request
	var err error

	url := getNgcAFPfdServiceUrl() + "/" + transID + "/applications/" + appID

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return trans, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return trans, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return trans, fmt.Errorf("HTTP failure: %d", resp.StatusCode)
	}

	trans, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return trans, err
	}
	return trans, nil
}

// AFPatchPfdApplication update an active PFD Application for the AF
func AFPatchPfdApplication(transID string, appID string, trans []byte) error {

	url := getNgcAFPfdServiceUrl() + "/" + transID + "/applications/" + appID

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(trans))
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

// AFDeletePfdApplication delete an active PFD Application for the AF
func AFDeletePfdApplication(transID string, appID string) error {

	url := getNgcAFPfdServiceUrl() + "/" + transID + "/applications/" + appID

	req, err := http.NewRequest("DELETE", url, nil)
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
