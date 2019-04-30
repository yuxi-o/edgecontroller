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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	cce "github.com/smartedgemec/controller-ce"
	"github.com/smartedgemec/controller-ce/uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("/vm_apps", func() {
	postSuccess := func() (id string) {
		By("Sending a POST /vm_apps request")
		resp, err := http.Post(
			"http://127.0.0.1:8080/vm_apps",
			"application/json",
			strings.NewReader(`
                {
                    "name": "vm app",
                    "vendor": "smart edge",
                    "description": "my vm app",
                    "image": "http://www.test.com/my_vm_app.tar.gz",
                    "cores": 4,
                    "memory": 1024
                }`))

		By("Verifying a 201 Created response")
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusCreated))

		By("Reading the response body")
		body, err := ioutil.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())

		var respBody struct {
			ID string
		}

		By("Unmarshalling the response")
		Expect(json.Unmarshal(body, &respBody)).To(Succeed())

		return respBody.ID
	}

	get := func(id string) *cce.VMApp {
		By("Sending a GET /vm_apps/{id} request")
		resp, err := http.Get(
			fmt.Sprintf("http://127.0.0.1:8080/vm_apps/%s", id))

		By("Verifying a 200 OK response")
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		By("Reading the response body")
		body, err := ioutil.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())

		var vmApp cce.VMApp

		By("Unmarshalling the response")
		Expect(json.Unmarshal(body, &vmApp)).To(Succeed())

		return &vmApp
	}

	Describe("POST /vm_apps", func() {
		DescribeTable("201 Created",
			func(req string) {
				By("Sending a POST /vm_apps request")
				resp, err := http.Post(
					"http://127.0.0.1:8080/vm_apps",
					"application/json",
					strings.NewReader(req))
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 201 response")
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				var respBody struct {
					ID string
				}

				By("Unmarshalling the response")
				Expect(json.Unmarshal(body, &respBody)).To(Succeed())

				By("Verifying a UUID was returned")
				Expect(uuid.IsValid(respBody.ID)).To(BeTrue())
			},
			Entry(
				"POST /vm_apps",
				`
                {
                    "name": "vm app",
                    "vendor": "smart edge",
                    "description": "my vm app",
                    "image": "http://www.test.com/my_vm_app.tar.gz",
                    "cores": 4,
                    "memory": 1024
                }`),
			Entry(
				"POST /vm_apps without description",
				`
                {
                    "name": "vm app",
                    "vendor": "smart edge",
                    "image": "http://www.test.com/my_vm_app.tar.gz",
                    "cores": 4,
                    "memory": 1024
                }`),
		)

		DescribeTable("400 Bad Request",
			func(req, expectedResp string) {
				By("Sending a POST /vm_apps request")
				resp, err := http.Post(
					"http://127.0.0.1:8080/vm_apps",
					"application/json",
					strings.NewReader(req))
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 400 Bad Request response")
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the response body")
				Expect(string(body)).To(Equal(expectedResp))
			},
			Entry(
				"POST /vm_apps without name",
				`
                {
                    "vendor": "smart edge",
                    "description": "my vm app",
                    "image": "http://www.test.com/my_vm_app.tar.gz",
                    "cores": 4,
                    "memory": 1024
                }`,
				"Validation failed: name cannot be empty"),
			Entry(
				"POST /vm_apps without vendor",
				`
                {
                    "name": "vm app",
                    "description": "my vm app",
                    "image": "http://www.test.com/my_vm_app.tar.gz",
                    "cores": 4,
                    "memory": 1024
                }`,
				"Validation failed: vendor cannot be empty"),
			Entry(
				"POST /vm_apps without image",
				`
                {
                    "name": "vm app",
                    "vendor": "smart edge",
                    "description": "my vm app",
                    "cores": 4,
                    "memory": 1024
                }`,
				"Validation failed: image cannot be empty"),
			Entry("POST /vm_apps with cores not in [1..8]",
				`
                {
                    "name": "vm app",
                    "vendor": "smart edge",
                    "description": "my vm app",
                    "image": "http://www.test.com/my_vm_app.tar.gz",
                    "cores": 9,
                    "memory": 1024
                }`,
				"Validation failed: cores must be in [1..8]"),
			Entry("POST /vm_apps with memory not in [1..16384]",
				`
                {
                    "name": "vm app",
                    "vendor": "smart edge",
                    "description": "my vm app",
                    "image": "http://www.test.com/my_vm_app.tar.gz",
                    "cores": 8,
                    "memory": 16385
                }`,
				"Validation failed: memory must be in [1..16384]"),
		)
	})

	Describe("GET /vm_apps", func() {
		var (
			vmAppID  string
			vmApp2ID string
		)

		BeforeEach(func() {
			vmAppID = postSuccess()
			vmApp2ID = postSuccess()
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a GET /vm_apps request")
				resp, err := http.Get("http://127.0.0.1:8080/vm_apps")

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				var vmApps []cce.VMApp

				By("Unmarshalling the response")
				Expect(json.Unmarshal(body, &vmApps)).To(Succeed())

				By("Verifying the 2 created vm apps were returned")
				Expect(vmApps).To(ContainElement(
					cce.VMApp{
						ID:          vmAppID,
						Name:        "vm app",
						Vendor:      "smart edge",
						Description: "my vm app",
						Image:       "http://www.test.com/my_vm_app.tar.gz",
						Cores:       4,
						Memory:      1024,
					}))
				Expect(vmApps).To(ContainElement(
					cce.VMApp{

						ID:          vmApp2ID,
						Name:        "vm app",
						Vendor:      "smart edge",
						Description: "my vm app",
						Image:       "http://www.test.com/my_vm_app.tar.gz",
						Cores:       4,
						Memory:      1024,
					}))
			},
			Entry("GET /vm_apps"),
		)
	})

	Describe("GET /vm_apps/{id}", func() {
		var (
			vmAppID string
		)

		BeforeEach(func() {
			vmAppID = postSuccess()
		})

		DescribeTable("200 OK",
			func() {
				vmApp := get(vmAppID)

				By("Verifying the created vm app was returned")
				Expect(vmApp).To(Equal(
					&cce.VMApp{
						ID:          vmAppID,
						Name:        "vm app",
						Vendor:      "smart edge",
						Description: "my vm app",
						Image:       "http://www.test.com/my_vm_app.tar.gz",
						Cores:       4,
						Memory:      1024,
					},
				))
			},
			Entry("GET /vm_apps/{id}"),
		)

		DescribeTable("404 Not Found",
			func() {
				By("Sending a GET /vm_apps/{id} request")
				resp, err := http.Get(
					fmt.Sprintf("http://127.0.0.1:8080/vm_apps/%s",
						uuid.New()))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("GET /vm_apps/{id} with nonexistent ID"),
		)
	})

	Describe("PATCH /vm_apps", func() {
		var (
			vmAppID string
		)

		BeforeEach(func() {
			vmAppID = postSuccess()
		})

		DescribeTable("204 No Content",
			func(reqStr string, expectedApp *cce.VMApp) {
				By("Sending a PATCH /vm_apps request")
				req, err := http.NewRequest(
					http.MethodPatch,
					"http://127.0.0.1:8080/vm_apps",
					strings.NewReader(fmt.Sprintf(reqStr, vmAppID)))
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 204 No Content response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				By("Getting the updated application")
				updatedApp := get(vmAppID)

				By("Verifying the vm app was updated")
				expectedApp.SetID(vmAppID)
				Expect(updatedApp).To(Equal(expectedApp))
			},
			Entry(
				"PATCH /vm_apps/{id}",
				`
                [
                    {
                        "id": "%s",
                        "name": "vm app2",
                        "vendor": "smart edge",
                        "description": "my vm app",
                        "image": "http://www.test.com/my_vm_app.tar.gz",
                        "cores": 4,
                        "memory": 1024
                    }
                ]`,
				&cce.VMApp{
					Name:        "vm app2",
					Vendor:      "smart edge",
					Description: "my vm app",
					Image:       "http://www.test.com/my_vm_app.tar.gz",
					Cores:       4,
					Memory:      1024,
				}),
			Entry("PATCH /vm_apps/{id} with no description",
				`
                [
                    {
                        "id": "%s",
                        "name": "vm app2",
                        "vendor": "smart edge",
                        "image": "http://www.test.com/my_vm_app.tar.gz",
                        "cores": 4,
                        "memory": 1024
                    }
                ]`,
				&cce.VMApp{
					Name:        "vm app2",
					Vendor:      "smart edge",
					Description: "",
					Image:       "http://www.test.com/my_vm_app.tar.gz",
					Cores:       4,
					Memory:      1024,
				}),
		)

		DescribeTable("400 Bad Request",
			func(reqStr string, expectedResp string) {
				By("Sending a PATCH /vm_apps request")
				req, err := http.NewRequest(
					http.MethodPatch,
					"http://127.0.0.1:8080/vm_apps",
					strings.NewReader(fmt.Sprintf(reqStr, vmAppID)))
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)

				By("Verifying a 400 Bad Request")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the response body")
				Expect(string(body)).To(Equal(expectedResp))
			},
			Entry(
				"PATCH /vm_apps without name",
				`
                [
                    {
                        "id": "%s",
                        "vendor": "smart edge",
                        "description": "my vm app",
                        "image": "http://www.test.com/my_vm_app.tar.gz",
                        "cores": 4,
                        "memory": 1024
                    }
                ]`,
				"Validation failed: name cannot be empty"),
			Entry("PATCH /vm_apps without vendor",
				`
                [
                    {
                        "id": "%s",
                        "name": "vm app2",
                        "description": "my vm app",
                        "image": "http://www.test.com/my_vm_app.tar.gz",
                        "cores": 4,
                        "memory": 1024
                    }
                ]`,
				"Validation failed: vendor cannot be empty"),
			Entry("PATCH /vm_apps without image",
				`
                [
                    {
                        "id": "%s",
                        "name": "vm app2",
                        "vendor": "smart edge",
                        "description": "my vm app",
                        "cores": 4,
                        "memory": 1024
                    }
                ]`,
				"Validation failed: image cannot be empty"),
			Entry("PATCH /vm_apps with cores not in [1..8]",
				`
                [
                    {
                        "id": "%s",
                        "name": "vm app2",
                        "vendor": "smart edge",
                        "description": "my vm app",
                        "image": "http://www.test.com/my_vm_app.tar.gz",
                        "cores": 9,
                        "memory": 1024
                    }
                ]`,
				"Validation failed: cores must be in [1..8]"),
			Entry("PATCH /vm_apps with memory not in [1..16384]",
				`
                [
                    {
                        "id": "%s",
                        "name": "vm app2",
                        "vendor": "smart edge",
                        "description": "my vm app",
                        "image": "http://www.test.com/my_vm_app.tar.gz",
                        "cores": 4,
                        "memory": 16385
                    }
                ]`,
				"Validation failed: memory must be in [1..16384]"),
		)
	})

	Describe("DELETE /vm_apps/{id}", func() {
		var (
			vmAppID string
		)

		BeforeEach(func() {
			vmAppID = postSuccess()
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a DELETE /vm_apps/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("http://127.0.0.1:8080/vm_apps/%s",
						vmAppID),
					nil)
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Verifying the vm app was deleted")

				By("Sending a GET /vm_apps/{id} request")
				resp, err = http.Get(
					fmt.Sprintf("http://127.0.0.1:8080/vm_apps/%s",
						vmAppID))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("DELETE /vm_apps/{id}"),
		)

		DescribeTable("404 Not Found",
			func(id string) {
				By("Sending a DELETE /vm_apps/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("http://127.0.0.1:8080/vm_apps/%s", id),
					nil)
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry(
				"DELETE /vm_apps/{id} with nonexistent ID",
				uuid.New()),
		)
	})
})
