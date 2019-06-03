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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cceGRPC "github.com/smartedgemec/controller-ce/grpc"
	"github.com/smartedgemec/controller-ce/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var _ = Describe("Kubernetes Read Metadata Service", func() {
	var (
		cvaSvcCli pb.ControllerVirtualizationAgentClient
		nodeID    string
	)

	BeforeEach(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		caPool := x509.NewCertPool()
		Expect(caPool.AppendCertsFromPEM(controllerRootPEM)).To(BeTrue(),
			"should load Controller self-signed root into trust pool")
		tlsCreds := credentials.NewClientTLSFromCert(caPool, cceGRPC.EnrollmentSNI)

		authConn, err := grpc.DialContext(
			ctx,
			fmt.Sprintf("%s:%d", "127.0.0.1", 8081),
			grpc.WithTransportCredentials(tlsCreds),
			grpc.WithBlock())
		Expect(err).ToNot(HaveOccurred(), "Dial failed: %v", err)
		authSvcCli := pb.NewAuthServiceClient(authConn)

		key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
		Expect(err).ToNot(HaveOccurred())

		By("Creating a certificate signing request with private key")
		csrDER, err := x509.CreateCertificateRequest(
			rand.Reader,
			&x509.CertificateRequest{},
			key,
		)
		Expect(err).ToNot(HaveOccurred())
		certReq, err := x509.ParseCertificateRequest(csrDER)
		Expect(err).ToNot(HaveOccurred())

		By("Encoding certificate signing request in PEM")
		csrPEM := pem.EncodeToMemory(
			&pem.Block{
				Type:  "CERTIFICATE REQUEST",
				Bytes: csrDER,
			})

		By("Pre-approving Node by serial")
		hash := md5.Sum(certReq.RawSubjectPublicKeyInfo)
		serial := base64.RawURLEncoding.EncodeToString(hash[:])

		approveNodeEnrollment(serial)

		ctx, cancel = context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		By("Requesting credentials from auth service")
		creds, err := authSvcCli.RequestCredentials(
			ctx,
			&pb.Identity{
				Csr: string(csrPEM),
			},
		)
		Expect(err).ToNot(HaveOccurred())

		block, _ := pem.Decode([]byte(creds.Certificate))
		if block == nil {
			Fail("failed to parse certificate PEM")
		}

		certBlock, _ := pem.Decode([]byte(creds.Certificate))
		Expect(certBlock).NotTo(BeNil(), "error decoding certificate in enrollment response")

		x509Cert, err := x509.ParseCertificate(certBlock.Bytes)
		Expect(err).ToNot(HaveOccurred())

		nodeID = x509Cert.Subject.CommonName

		cert := tls.Certificate{Certificate: [][]byte{certBlock.Bytes}, PrivateKey: key}

		t := &tls.Config{
			RootCAs:      caPool,
			Certificates: []tls.Certificate{cert},
			ServerName:   cceGRPC.SNI,
		}

		ctx, cancel = context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		cvaConn, err := grpc.DialContext(
			ctx,
			fmt.Sprintf("%s:%d", "127.0.0.1", 8081),
			grpc.WithTransportCredentials(credentials.NewTLS(t)),
			grpc.WithBlock())
		Expect(err).ToNot(HaveOccurred(), "Dial failed: %v", err)

		cvaSvcCli = pb.NewControllerVirtualizationAgentClient(cvaConn)

		// label node with correct id
		Expect(exec.Command("kubectl",
			"label", "nodes", "minikube", fmt.Sprintf("node-id=%s", nodeID)).Run()).To(Succeed())
	})

	AfterEach(func() {
		// un-label node with id
		Expect(exec.Command("kubectl", "label", "nodes", "minikube", "node-id-").Run()).To(Succeed())
		// clean up all k8s deployments
		cmd := exec.Command("kubectl", "delete", "--all", "deployments,pods", "--namespace=default")
		Expect(cmd.Run()).To(Succeed())
	})

	Describe("Get Pod Information By IP", func() {
		Context("Success", func() {
			It("Should return pod information", func() {
				By("Deploying an application to Kubernetes")
				deployApp(nodeID)

				By("Generating IP address of the pod deployed")
				var ip string
				count := 0
				Eventually(func() net.IP {
					count++
					By(fmt.Sprintf("Attempt #%d: Verifying if ip assigned to pod is valid", count))

					out, err := exec.Command("kubectl",
						"get", "pods", "-o=jsonpath='{.items[0].status.podIP}'").Output()
					Expect(err).ToNot(HaveOccurred())

					ip = strings.Trim(string(out), "'")
					return net.ParseIP(ip)
				}, 15*time.Second, 1*time.Second).ShouldNot(BeNil())

				ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
				defer cancel()

				By(fmt.Sprintf("Requesting container info for pod with ip: %s", ip))
				containerInfo, err := cvaSvcCli.GetContainerByIP(
					ctx,
					&pb.ContainerIP{
						Ip: ip,
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(containerInfo.Id).To(Equal(appID))
			})
		})
		Context("Error", func() {
			It("Should return an error if request contains no correct IP address", func() {
				ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
				defer cancel()

				By("Requesting container info for pod with ip: ''")
				_, err := cvaSvcCli.GetContainerByIP(
					ctx,
					&pb.ContainerIP{
						Ip: "",
					},
				)
				Expect(err).To(HaveOccurred())
			})

			It("Should return an error if request contains an IP that is not assigned to a pod", func() {
				ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
				defer cancel()

				impossibleIP := "0.0.0.0"
				By("Requesting container info for pod with ip: 0.0.0.0")
				_, err := cvaSvcCli.GetContainerByIP(
					ctx,
					&pb.ContainerIP{
						Ip: impossibleIP,
					},
				)
				Expect(err).To(MatchError("rpc error: code = Internal desc = unable to get pod name by ip"))
			})
		})
	})
})
