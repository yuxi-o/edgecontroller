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

var _ = Describe("/container_apps", func() {
	Describe("POST /container_apps", func() {
		DescribeTable("201 Created",
			func(req string) {
				By("Sending a POST /container_apps request")
				resp, err := http.Post(
					"http://127.0.0.1:8080/container_apps",
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
				"POST /container_apps",
				`
                {
                    "name": "container app",
                    "vendor": "smart edge",
                    "description": "my container app",
                    "image": "http://www.test.com/my_container_app.tar.gz",
                    "cores": 4,
                    "memory": 1024
                }`),
			Entry(
				"POST /container_apps without description",
				`
                {
                    "name": "container app",
                    "vendor": "smart edge",
                    "image": "http://www.test.com/my_container_app.tar.gz",
                    "cores": 4,
                    "memory": 1024
                }`),
		)

		DescribeTable("400 Bad Request",
			func(req, expectedResp string) {
				By("Sending a POST /container_apps request")
				resp, err := http.Post(
					"http://127.0.0.1:8080/container_apps",
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
				"POST /container_apps with id",
				`
                {
                    "id": "123"
                }`,
				"Validation failed: id cannot be specified in POST request"),
			Entry(
				"POST /container_apps without name",
				`
                {
                    "vendor": "smart edge",
                    "description": "my container app",
                    "image": "http://www.test.com/my_container_app.tar.gz",
                    "cores": 4,
                    "memory": 1024
                }`,
				"Validation failed: name cannot be empty"),
			Entry(
				"POST /container_apps without vendor",
				`
                {
                    "name": "container app",
                    "description": "my container app",
                    "image": "http://www.test.com/my_container_app.tar.gz",
                    "cores": 4,
                    "memory": 1024
                }`,
				"Validation failed: vendor cannot be empty"),
			Entry(
				"POST /container_apps without image",
				`
                {
                    "name": "container app",
                    "vendor": "smart edge",
                    "description": "my container app",
                    "cores": 4,
                    "memory": 1024
                }`,
				"Validation failed: image cannot be empty"),
			Entry("POST /container_apps with cores not in [1..8]",
				`
                {
                    "name": "container app",
                    "vendor": "smart edge",
                    "description": "my container app",
                    "image": "http://www.test.com/my_container_app.tar.gz",
                    "cores": 9,
                    "memory": 1024
                }`,
				"Validation failed: cores must be in [1..8]"),
			Entry("POST /container_apps with memory not in [1..16384]",
				`
                {
                    "name": "container app",
                    "vendor": "smart edge",
                    "description": "my container app",
                    "image": "http://www.test.com/my_container_app.tar.gz",
                    "cores": 8,
                    "memory": 16385
                }`,
				"Validation failed: memory must be in [1..16384]"),
		)
	})

	Describe("GET /container_apps", func() {
		var (
			containerAppID  string
			containerApp2ID string
		)

		BeforeEach(func() {
			containerAppID = postContainerApps()
			containerApp2ID = postContainerApps()
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a GET /container_apps request")
				resp, err := http.Get("http://127.0.0.1:8080/container_apps")

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				var containerApps []cce.ContainerApp

				By("Unmarshalling the response")
				Expect(json.Unmarshal(body, &containerApps)).To(Succeed())

				By("Verifying the 2 created container apps were returned")
				Expect(containerApps).To(ContainElement(
					cce.ContainerApp{
						ID:          containerAppID,
						Name:        "container app",
						Vendor:      "smart edge",
						Description: "my container app",
						Image:       "http://www.test.com/my_container_app.tar.gz", //nolint:lll
						Cores:       4,
						Memory:      1024,
					}))
				Expect(containerApps).To(ContainElement(
					cce.ContainerApp{

						ID:          containerApp2ID,
						Name:        "container app",
						Vendor:      "smart edge",
						Description: "my container app",
						Image:       "http://www.test.com/my_container_app.tar.gz", //nolint:lll
						Cores:       4,
						Memory:      1024,
					}))
			},
			Entry("GET /container_apps"),
		)
	})

	Describe("GET /container_apps/{id}", func() {
		var (
			containerAppID string
		)

		BeforeEach(func() {
			containerAppID = postContainerApps()
		})

		DescribeTable("200 OK",
			func() {
				containerApp := getContainerApp(containerAppID)

				By("Verifying the created container app was returned")
				Expect(containerApp).To(Equal(
					&cce.ContainerApp{
						ID:          containerAppID,
						Name:        "container app",
						Vendor:      "smart edge",
						Description: "my container app",
						Image:       "http://www.test.com/my_container_app.tar.gz", //nolint:lll
						Cores:       4,
						Memory:      1024,
					},
				))
			},
			Entry("GET /container_apps/{id}"),
		)

		DescribeTable("404 Not Found",
			func() {
				By("Sending a GET /container_apps/{id} request")
				resp, err := http.Get(
					fmt.Sprintf("http://127.0.0.1:8080/container_apps/%s",
						uuid.New()))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("GET /container_apps/{id} with nonexistent ID"),
		)
	})

	Describe("PATCH /container_apps", func() {
		var (
			containerAppID string
		)

		BeforeEach(func() {
			containerAppID = postContainerApps()
		})

		DescribeTable("204 No Content",
			func(reqStr string, expectedApp *cce.ContainerApp) {
				By("Sending a PATCH /container_apps request")
				req, err := http.NewRequest(
					http.MethodPatch,
					"http://127.0.0.1:8080/container_apps",
					strings.NewReader(fmt.Sprintf(reqStr, containerAppID)))
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 204 No Content response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				By("Getting the updated application")
				updatedApp := getContainerApp(containerAppID)

				By("Verifying the container app was updated")
				expectedApp.SetID(containerAppID)
				Expect(updatedApp).To(Equal(expectedApp))
			},
			Entry(
				"PATCH /container_apps",
				`
                [
                    {
                        "id": "%s",
                        "name": "container app2",
                        "vendor": "smart edge",
                        "description": "my container app",
                        "image": "http://www.test.com/my_container_app.tar.gz",
                        "cores": 4,
                        "memory": 1024
                    }
                ]`,
				&cce.ContainerApp{
					Name:        "container app2",
					Vendor:      "smart edge",
					Description: "my container app",
					Image:       "http://www.test.com/my_container_app.tar.gz",
					Cores:       4,
					Memory:      1024,
				}),
			Entry("PATCH /container_apps with no description",
				`
                [
                    {
                        "id": "%s",
                        "name": "container app2",
                        "vendor": "smart edge",
                        "image": "http://www.test.com/my_container_app.tar.gz",
                        "cores": 4,
                        "memory": 1024
                    }
                ]`,
				&cce.ContainerApp{
					Name:        "container app2",
					Vendor:      "smart edge",
					Description: "",
					Image:       "http://www.test.com/my_container_app.tar.gz",
					Cores:       4,
					Memory:      1024,
				}),
		)

		DescribeTable("400 Bad Request",
			func(reqStr string, expectedResp string) {
				By("Sending a PATCH /container_apps request")
				if strings.Contains(reqStr, "%s") {
					reqStr = fmt.Sprintf(reqStr, containerAppID)
				}
				req, err := http.NewRequest(
					http.MethodPatch,
					"http://127.0.0.1:8080/container_apps",
					strings.NewReader(reqStr))
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
				"PATCH /container_apps without id",
				`
                [
                    {
                        "name": "container app2",
                        "vendor": "smart edge",
                        "description": "my container app",
                        "image": "http://www.test.com/my_container_app.tar.gz",
                        "cores": 4,
                        "memory": 1024
                    }
                ]`,
				"Validation failed: id not a valid uuid"),
			Entry(
				"PATCH /container_apps without name",
				`
                [
                    {
                        "id": "%s",
                        "vendor": "smart edge",
                        "description": "my container app",
                        "image": "http://www.test.com/my_container_app.tar.gz",
                        "cores": 4,
                        "memory": 1024
                    }
                ]`,
				"Validation failed: name cannot be empty"),
			Entry("PATCH /container_apps without vendor",
				`
                [
                    {
                        "id": "%s",
                        "name": "container app2",
                        "description": "my container app",
                        "image": "http://www.test.com/my_container_app.tar.gz",
                        "cores": 4,
                        "memory": 1024
                    }
                ]`,
				"Validation failed: vendor cannot be empty"),
			Entry("PATCH /container_apps without image",
				`
                [
                    {
                        "id": "%s",
                        "name": "container app2",
                        "vendor": "smart edge",
                        "description": "my container app",
                        "cores": 4,
                        "memory": 1024
                    }
                ]`,
				"Validation failed: image cannot be empty"),
			Entry("PATCH /container_apps with cores not in [1..8]",
				`
                [
                    {
                        "id": "%s",
                        "name": "container app2",
                        "vendor": "smart edge",
                        "description": "my container app",
                        "image": "http://www.test.com/my_container_app.tar.gz",
                        "cores": 9,
                        "memory": 1024
                    }
                ]`,
				"Validation failed: cores must be in [1..8]"),
			Entry("PATCH /container_apps with memory not in [1..16384]",
				`
                [
                    {
                        "id": "%s",
                        "name": "container app2",
                        "vendor": "smart edge",
                        "description": "my container app",
                        "image": "http://www.test.com/my_container_app.tar.gz",
                        "cores": 4,
                        "memory": 16385
                    }
                ]`,
				"Validation failed: memory must be in [1..16384]"),
		)
	})

	Describe("DELETE /container_apps/{id}", func() {
		var (
			containerAppID string
		)

		BeforeEach(func() {
			containerAppID = postContainerApps()
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a DELETE /container_apps/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("http://127.0.0.1:8080/container_apps/%s",
						containerAppID),
					nil)
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Verifying the container app was deleted")

				By("Sending a GET /container_apps/{id} request")
				resp, err = http.Get(
					fmt.Sprintf("http://127.0.0.1:8080/container_apps/%s",
						containerAppID))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("DELETE /container_apps/{id}"),
		)

		DescribeTable("404 Not Found",
			func(id string) {
				By("Sending a DELETE /container_apps/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("http://127.0.0.1:8080/container_apps/%s", id),
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
				"DELETE /container_apps/{id} with nonexistent ID",
				uuid.New()),
		)
	})
})
