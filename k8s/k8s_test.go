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

package k8s_test

import (
	"context"
	"fmt"
	"os/exec"
	"os/user"
	"path"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/smartedgemec/controller-ce/k8s"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	homeDir string
)

// in-order-to run these tests, mini-kube and virtualization
// tools need to be installed to setup mini-kube on travis CI.
var _ = BeforeSuite(func() {
	u, err := user.Current()
	Expect(err).NotTo(HaveOccurred())
	homeDir = u.HomeDir
	// check if kube-ctl installed and in PATH
	cmd := exec.Command("kubectl", "version")
	err = cmd.Run()
	Expect(err).NotTo(HaveOccurred())
	// check if docker installed and in PATH
	cmd = exec.Command("docker", "version")
	err = cmd.Run()
	Expect(err).NotTo(HaveOccurred())
	// check if mini-kube installed and in PATH
	cmd = exec.Command("minikube", "version")
	err = cmd.Run()
	Expect(err).NotTo(HaveOccurred())
	// label node with correct app id
	cmd = exec.Command("kubectl", "label", "nodes", "minikube", "node-uuid=minikube")
	err = cmd.Run()
	Expect(err).NotTo(HaveOccurred())
	// docker pull public image for testing
	cmd = exec.Command("docker", "pull", "nginx:1.12")
	err = cmd.Run()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	// clean up k8s deployments
	cmd := exec.Command("kubectl", "delete", "--all", "deployments", "--namespace=default")
	err := cmd.Run()
	Expect(err).NotTo(HaveOccurred())
	// clean up k8s pods
	cmd = exec.Command("kubectl", "delete", "--all", "pods", "--namespace=default")
	err = cmd.Run()
	Expect(err).NotTo(HaveOccurred())
	// label node with correct app id
	cmd = exec.Command("kubectl", "label", "nodes", "minikube", "node-uuid-")
	err = cmd.Run()
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("K8S", func() {
	Context("API calls to K8S master", func() {
		It("Should deploy, start, stop, restart and undeploy an app from a public docker image", func() {
			kubeConfig := path.Join(homeDir, ".kube", "config")
			config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
			Expect(err).NotTo(HaveOccurred())
			client := k8s.Client{
				Username: config.Username,
				Host:     config.Host,
				APIPath:  config.APIPath,
				CertFile: config.TLSClientConfig.CertFile,
				KeyFile:  config.TLSClientConfig.KeyFile,
				CAFile:   config.TLSClientConfig.CAFile,
			}
			Expect(client.Ping()).To(Succeed())

			app := k8s.App{
				ID:    "1234",
				Name:  "smellascope",
				Image: "nginx:1.12", // commonly used public docker container
			}
			ctx := context.Background()
			ctx, cancel0 := context.WithTimeout(ctx, 5*time.Second)
			defer cancel0()
			err = client.Deploy(ctx, "minikube", &app)
			Expect(err).NotTo(HaveOccurred())

			ctx = context.Background()
			ctx, cancel1 := context.WithTimeout(ctx, 20*time.Second)
			defer cancel1()

			Eventually(func() k8s.LifecycleStatus {
				var status k8s.LifecycleStatus
				status, err = client.Status(ctx, "minikube", app.ID)
				if err != nil {
					fmt.Printf("error checking status: %v", err)
					return k8s.Unknown
				}
				return status
			}, 20*time.Second, 1*time.Second).Should(Equal(k8s.Deployed))

			ctx = context.Background()
			ctx, cancel2 := context.WithTimeout(ctx, 5*time.Second)
			defer cancel2()

			err = client.Start(ctx, "minikube", app.ID)
			Expect(err).NotTo(HaveOccurred())

			ctx = context.Background()
			ctx, cancel3 := context.WithTimeout(ctx, 5*time.Second)
			defer cancel3()

			err = client.Stop(ctx, "minikube", app.ID)
			Expect(err).NotTo(HaveOccurred())

			ctx = context.Background()
			ctx, cancel4 := context.WithTimeout(ctx, 5*time.Second)
			defer cancel4()

			err = client.Restart(ctx, "minikube", app.ID)
			Expect(err).NotTo(HaveOccurred())

			ctx = context.Background()
			ctx, cancel5 := context.WithTimeout(ctx, 5*time.Second)
			defer cancel5()

			err = client.Undeploy(ctx, "minikube", &app)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() error {
				ctx = context.Background()
				var cancel6 context.CancelFunc
				ctx, cancel6 = context.WithTimeout(ctx, 5*time.Second)
				defer cancel6()

				_, err = client.Status(ctx, "minikube", app.ID)
				return err
			}, 40*time.Second, 1*time.Second).Should(HaveOccurred())
		})
	})
})
