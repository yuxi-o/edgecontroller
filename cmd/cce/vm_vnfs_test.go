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

var _ = Describe("/vm_vnfs", func() {
	Describe("POST /vm_vnfs", func() {
		DescribeTable("201 Created",
			func(req string) {
				By("Sending a POST /vm_vnfs request")
				resp, err := http.Post(
					"http://127.0.0.1:8080/vm_vnfs",
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
				"POST /vm_vnfs",
				`
                {
                    "name": "vm vnf",
                    "vendor": "smart edge",
                    "description": "my vm vnf",
                    "image": "http://www.test.com/my_vm_vnf.tar.gz",
                    "cores": 4,
                    "memory": 1024
                }`),
			Entry(
				"POST /vm_vnfs without description",
				`
                {
                    "name": "vm vnf",
                    "vendor": "smart edge",
                    "image": "http://www.test.com/my_vm_vnf.tar.gz",
                    "cores": 4,
                    "memory": 1024
                }`),
		)

		DescribeTable("400 Bad Request",
			func(req, expectedResp string) {
				By("Sending a POST /vm_vnfs request")
				resp, err := http.Post(
					"http://127.0.0.1:8080/vm_vnfs",
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
				"POST /vm_vnfs with id",
				`
                {
                    "id": "123"
                }`,
				"Validation failed: id cannot be specified in POST request"),
			Entry(
				"POST /vm_vnfs without name",
				`
                {
                    "vendor": "smart edge",
                    "description": "my vm vnf",
                    "image": "http://www.test.com/my_vm_vnf.tar.gz",
                    "cores": 4,
                    "memory": 1024
                }`,
				"Validation failed: name cannot be empty"),
			Entry(
				"POST /vm_vnfs without vendor",
				`
                {
                    "name": "vm vnf",
                    "description": "my vm vnf",
                    "image": "http://www.test.com/my_vm_vnf.tar.gz",
                    "cores": 4,
                    "memory": 1024
                }`,
				"Validation failed: vendor cannot be empty"),
			Entry(
				"POST /vm_vnfs without image",
				`
                {
                    "name": "vm vnf",
                    "vendor": "smart edge",
                    "description": "my vm vnf",
                    "cores": 4,
                    "memory": 1024
                }`,
				"Validation failed: image cannot be empty"),
			Entry("POST /vm_vnfs with cores not in [1..8]",
				`
                {
                    "name": "vm vnf",
                    "vendor": "smart edge",
                    "description": "my vm vnf",
                    "image": "http://www.test.com/my_vm_vnf.tar.gz",
                    "cores": 9,
                    "memory": 1024
                }`,
				"Validation failed: cores must be in [1..8]"),
			Entry("POST /vm_vnfs with memory not in [1..16384]",
				`
                {
                    "name": "vm vnf",
                    "vendor": "smart edge",
                    "description": "my vm vnf",
                    "image": "http://www.test.com/my_vm_vnf.tar.gz",
                    "cores": 8,
                    "memory": 16385
                }`,
				"Validation failed: memory must be in [1..16384]"),
		)
	})

	Describe("GET /vm_vnfs", func() {
		var (
			vmVNFID  string
			vmVNF2ID string
		)

		BeforeEach(func() {
			vmVNFID = postVMVNFs()
			vmVNF2ID = postVMVNFs()
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a GET /vm_vnfs request")
				resp, err := http.Get("http://127.0.0.1:8080/vm_vnfs")

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Reading the response body")
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				var vmVNFs []cce.VMVNF

				By("Unmarshalling the response")
				Expect(json.Unmarshal(body, &vmVNFs)).To(Succeed())

				By("Verifying the 2 created VM VNFs were returned")
				Expect(vmVNFs).To(ContainElement(
					cce.VMVNF{
						ID:          vmVNFID,
						Name:        "vm vnf",
						Vendor:      "smart edge",
						Description: "my vm vnf",
						Image:       "http://www.test.com/my_vm_vnf.tar.gz",
						Cores:       4,
						Memory:      1024,
					}))
				Expect(vmVNFs).To(ContainElement(
					cce.VMVNF{

						ID:          vmVNF2ID,
						Name:        "vm vnf",
						Vendor:      "smart edge",
						Description: "my vm vnf",
						Image:       "http://www.test.com/my_vm_vnf.tar.gz",
						Cores:       4,
						Memory:      1024,
					}))
			},
			Entry("GET /vm_vnfs"),
		)
	})

	Describe("GET /vm_vnfs/{id}", func() {
		var (
			vmVNFID string
		)

		BeforeEach(func() {
			vmVNFID = postVMVNFs()
		})

		DescribeTable("200 OK",
			func() {
				vmVNF := getVMVNF(vmVNFID)

				By("Verifying the created VM VNF was returned")
				Expect(vmVNF).To(Equal(
					&cce.VMVNF{
						ID:          vmVNFID,
						Name:        "vm vnf",
						Vendor:      "smart edge",
						Description: "my vm vnf",
						Image:       "http://www.test.com/my_vm_vnf.tar.gz",
						Cores:       4,
						Memory:      1024,
					},
				))
			},
			Entry("GET /vm_vnfs/{id}"),
		)

		DescribeTable("404 Not Found",
			func() {
				By("Sending a GET /vm_vnfs/{id} request")
				resp, err := http.Get(
					fmt.Sprintf("http://127.0.0.1:8080/vm_vnfs/%s",
						uuid.New()))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("GET /vm_vnfs/{id} with nonexistent ID"),
		)
	})

	Describe("PATCH /vm_vnfs", func() {
		var (
			vmVNFID string
		)

		BeforeEach(func() {
			vmVNFID = postVMVNFs()
		})

		DescribeTable("204 No Content",
			func(reqStr string, expectedVNF *cce.VMVNF) {
				By("Sending a PATCH /vm_vnfs request")
				req, err := http.NewRequest(
					http.MethodPatch,
					"http://127.0.0.1:8080/vm_vnfs",
					strings.NewReader(fmt.Sprintf(reqStr, vmVNFID)))
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 204 No Content response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				By("Getting the updated VNF")
				updatedVNF := getVMVNF(vmVNFID)

				By("Verifying the VM VNF was updated")
				expectedVNF.SetID(vmVNFID)
				Expect(updatedVNF).To(Equal(expectedVNF))
			},
			Entry(
				"PATCH /vm_vnfs",
				`
                [
                    {
                        "id": "%s",
                        "name": "vm vnf2",
                        "vendor": "smart edge",
                        "description": "my vm vnf",
                        "image": "http://www.test.com/my_vm_vnf.tar.gz",
                        "cores": 4,
                        "memory": 1024
                    }
                ]`,
				&cce.VMVNF{
					Name:        "vm vnf2",
					Vendor:      "smart edge",
					Description: "my vm vnf",
					Image:       "http://www.test.com/my_vm_vnf.tar.gz",
					Cores:       4,
					Memory:      1024,
				}),
			Entry("PATCH /vm_vnfs with no description",
				`
                [
                    {
                        "id": "%s",
                        "name": "vm vnf2",
                        "vendor": "smart edge",
                        "image": "http://www.test.com/my_vm_vnf.tar.gz",
                        "cores": 4,
                        "memory": 1024
                    }
                ]`,
				&cce.VMVNF{
					Name:        "vm vnf2",
					Vendor:      "smart edge",
					Description: "",
					Image:       "http://www.test.com/my_vm_vnf.tar.gz",
					Cores:       4,
					Memory:      1024,
				}),
		)

		DescribeTable("400 Bad Request",
			func(reqStr string, expectedResp string) {
				By("Sending a PATCH /vm_vnfs request")
				if strings.Contains(reqStr, "%s") {
					reqStr = fmt.Sprintf(reqStr, vmVNFID)
				}
				req, err := http.NewRequest(
					http.MethodPatch,
					"http://127.0.0.1:8080/vm_vnfs",
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
				"PATCH /vm_vnfs without id",
				`
                [
                    {
                        "name": "vm vnf2",
                        "vendor": "smart edge",
                        "description": "my vm vnf",
                        "image": "http://www.test.com/my_vm_app.tar.gz",
                        "cores": 4,
                        "memory": 1024
                    }
                ]`,
				"Validation failed: id not a valid uuid"),
			Entry(
				"PATCH /vm_vnfs without name",
				`
                [
                    {
                        "id": "%s",
                        "vendor": "smart edge",
                        "description": "my vm vnf",
                        "image": "http://www.test.com/my_vm_vnf.tar.gz",
                        "cores": 4,
                        "memory": 1024
                    }
                ]`,
				"Validation failed: name cannot be empty"),
			Entry("PATCH /vm_vnfs without vendor",
				`
                [
                    {
                        "id": "%s",
                        "name": "vm vnf2",
                        "description": "my vm vnf",
                        "image": "http://www.test.com/my_vm_vnf.tar.gz",
                        "cores": 4,
                        "memory": 1024
                    }
                ]`,
				"Validation failed: vendor cannot be empty"),
			Entry("PATCH /vm_vnfs without image",
				`
                [
                    {
                        "id": "%s",
                        "name": "vm vnf2",
                        "vendor": "smart edge",
                        "description": "my vm vnf",
                        "cores": 4,
                        "memory": 1024
                    }
                ]`,
				"Validation failed: image cannot be empty"),
			Entry("PATCH /vm_vnfs with cores not in [1..8]",
				`
                [
                    {
                        "id": "%s",
                        "name": "vm vnf2",
                        "vendor": "smart edge",
                        "description": "my vm vnf",
                        "image": "http://www.test.com/my_vm_vnf.tar.gz",
                        "cores": 9,
                        "memory": 1024
                    }
                ]`,
				"Validation failed: cores must be in [1..8]"),
			Entry("PATCH /vm_vnfs with memory not in [1..16384]",
				`
                [
                    {
                        "id": "%s",
                        "name": "vm vnf2",
                        "vendor": "smart edge",
                        "description": "my vm vnf",
                        "image": "http://www.test.com/my_vm_vnf.tar.gz",
                        "cores": 4,
                        "memory": 16385
                    }
                ]`,
				"Validation failed: memory must be in [1..16384]"),
		)
	})

	Describe("DELETE /vm_vnfs/{id}", func() {
		var (
			vmVNFID string
		)

		BeforeEach(func() {
			vmVNFID = postVMVNFs()
		})

		DescribeTable("200 OK",
			func() {
				By("Sending a DELETE /vm_vnfs/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("http://127.0.0.1:8080/vm_vnfs/%s",
						vmVNFID),
					nil)
				Expect(err).ToNot(HaveOccurred())

				c := http.Client{}
				resp, err := c.Do(req)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying a 200 OK response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				By("Verifying the VM VNF was deleted")

				By("Sending a GET /vm_vnfs/{id} request")
				resp, err = http.Get(
					fmt.Sprintf("http://127.0.0.1:8080/vm_vnfs/%s",
						vmVNFID))

				By("Verifying a 404 Not Found response")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			},
			Entry("DELETE /vm_vnfs/{id}"),
		)

		DescribeTable("404 Not Found",
			func(id string) {
				By("Sending a DELETE /vm_vnfs/{id} request")
				req, err := http.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("http://127.0.0.1:8080/vm_vnfs/%s", id),
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
				"DELETE /vm_vnfs/{id} with nonexistent ID",
				uuid.New()),
		)
	})
})
