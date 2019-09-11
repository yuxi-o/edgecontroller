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
	"log"
	"os/exec"
	"os/user"
	"path"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cce "github.com/otcshare/edgecontroller"
	"github.com/otcshare/edgecontroller/k8s"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	homeDir string
	nodeID  = "d9912f2c-bd5a-411b-8a76-195fa9401f63"
	appID   = "99459845-422d-4b32-8395-e8f50fd34792"
)

// In order to run these tests, mini-kube and virtualization
// tools need to be installed to setup mini-kube on travis CI.
var _ = BeforeSuite(func() {
	u, err := user.Current()
	Expect(err).NotTo(HaveOccurred())
	homeDir = u.HomeDir
	// check if kube-ctl installed and in PATH
	cmd := exec.Command("kubectl", "version")
	Expect(cmd.Run()).To(Succeed())
	// check if docker installed and in PATH
	cmd = exec.Command("docker", "version")
	Expect(cmd.Run()).To(Succeed())
	// check if mini-kube installed and in PATH
	cmd = exec.Command("minikube", "version")
	Expect(cmd.Run()).To(Succeed())
	// label node with correct app id
	cmd = exec.Command("kubectl", "label", "nodes", "minikube", fmt.Sprintf("node-id=%s", nodeID))
	Expect(cmd.Run()).To(Succeed())
	// docker pull public image for testing
	cmd = exec.Command("docker", "pull", "nginx:1.12")
	Expect(cmd.Run()).To(Succeed())
})

var _ = AfterSuite(func() {
	// clean up k8s deployments
	cmd := exec.Command("kubectl", "delete", "--all", "deployments", "--namespace=default")
	Expect(cmd.Run()).To(Succeed())
	// clean up k8s pods
	cmd = exec.Command("kubectl", "delete", "--all", "pods", "--namespace=default")
	Expect(cmd.Run()).To(Succeed())
	// label node with correct app id
	cmd = exec.Command("kubectl", "label", "nodes", "minikube", "node-id-")
	Expect(cmd.Run()).To(Succeed())
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

			trafficPolicy := &cce.TrafficPolicyKubeOVN{
				ID:   "374bd735-8be6-42c3-a7d1-41fbb61542e0",
				Name: "traffic policy for app",
				Ingress: []*cce.IngressRule{
					{
						From: []*cce.IPBlock{
							{
								CIDR:   "192.168.1.0/24",
								Except: []string{"192.168.1.0/30"},
							},
						},
						Ports: []*cce.Port{
							{
								Port:     80,
								Protocol: "tcp",
							},
						},
					},
				},
				Egress: []*cce.EgressRule{
					{
						To: []*cce.IPBlock{
							{
								CIDR:   "10.16.0.0/16",
								Except: []string{"10.16.0.0/24"},
							},
						},
						Ports: []*cce.Port{
							{
								Port:     443,
								Protocol: "tcp",
							},
						},
					},
				},
			}

			app := k8s.App{
				ID:     appID,
				Image:  "nginx:1.12", // commonly used public docker container
				Cores:  1,
				Memory: 100,
				Ports: []*k8s.PortProto{
					{
						Port:     8080,
						Protocol: "tcp",
					},
				},
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			Expect(client.Deploy(ctx, nodeID, app)).To(Succeed())

			ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			Eventually(func() k8s.LifecycleStatus {
				var status k8s.LifecycleStatus
				status, err = client.Status(ctx, nodeID, appID)
				if err != nil {
					log.Printf("error checking status: %v", err)
					return k8s.Unknown
				}
				return status
			}, 10*time.Second, 1*time.Second).Should(Equal(k8s.Deployed))

			ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			Expect(client.Start(ctx, nodeID, appID)).To(Succeed())

			ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			Expect(client.ApplyNetworkPolicy(ctx, nodeID, appID, trafficPolicy.ToK8s())).To(Succeed())

			ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			Expect(client.DeleteNetworkPolicy(ctx, nodeID, appID)).To(Succeed())

			ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			Expect(client.Stop(ctx, nodeID, appID)).To(Succeed())

			ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			Expect(client.Restart(ctx, nodeID, appID)).To(Succeed())

			ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			Expect(client.Undeploy(ctx, nodeID, appID)).To(Succeed())

			Eventually(func() error {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				_, err = client.Status(ctx, nodeID, appID)
				return err
			}, 40*time.Second, 1*time.Second).Should(HaveOccurred())
		})
	})
})
