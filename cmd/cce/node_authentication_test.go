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
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	cceGRPC "github.com/smartedgemec/controller-ce/grpc"
	"github.com/smartedgemec/controller-ce/pb"
)

var _ = Describe("Node Auth Service", func() {
	var authSvcCli pb.AuthServiceClient

	BeforeEach(func() {
		timeoutCtx, cancel := context.WithTimeout(
			ctx, 2*time.Second)
		defer cancel()

		caPool := x509.NewCertPool()
		Expect(caPool.AppendCertsFromPEM(controllerRootPEM)).To(BeTrue(),
			"should load Controller self-signed root into trust pool")
		tlsCreds := credentials.NewClientTLSFromCert(caPool, cceGRPC.EnrollmentSNI)

		conn, err := grpc.DialContext(
			timeoutCtx,
			fmt.Sprintf("%s:%d", "127.0.0.1", 8081),
			grpc.WithTransportCredentials(tlsCreds),
			grpc.WithBlock())
		Expect(err).ToNot(HaveOccurred(), "Dial failed: %v", err)

		authSvcCli = pb.NewAuthServiceClient(conn)
	})

	Describe("RequestCredentials", func() {
		Describe("Success", func() {
			It("Should return auth credentials", func() {
				By("Generating node private key")
				key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
				Expect(err).ToNot(HaveOccurred())

				By("Creating a certificate signing request with private key")
				csrDER, err := x509.CreateCertificateRequest(
					rand.Reader,
					&x509.CertificateRequest{},
					key,
				)
				Expect(err).ToNot(HaveOccurred())

				By("Encoding certificate signing request in PEM")
				csrPEM := pem.EncodeToMemory(
					&pem.Block{
						Type:  "CERTIFICATE REQUEST",
						Bytes: csrDER,
					})

				By("Requesting credentials from auth service")
				credentials, err := authSvcCli.RequestCredentials(
					ctx,
					&pb.Identity{
						Csr: string(csrPEM),
					},
				)
				Expect(err).ToNot(HaveOccurred())

				By("Validating the returned credentials")
				Expect(credentials).ToNot(BeNil())
				Expect(credentials.Certificate).ToNot(BeNil())
				Expect(credentials.CaChain).ToNot(BeEmpty())

				By("Decoding PEM-encoded client certificate")
				certBlock, rest := pem.Decode([]byte(credentials.Certificate))
				Expect(certBlock).ToNot(BeNil())
				Expect(rest).To(BeEmpty())

				By("Parsing the client certificate")
				cert, err := x509.ParseCertificate(certBlock.Bytes)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying certificate was signed with the node's key")
				pubKeyDER, err := x509.MarshalPKIXPublicKey(key.Public())
				Expect(err).ToNot(HaveOccurred())
				Expect(cert.RawSubjectPublicKeyInfo).To(Equal(pubKeyDER))

				By("Verifying the CN is derived from the public key")
				hash := md5.Sum(pubKeyDER)
				cn := base64.RawURLEncoding.EncodeToString(hash[:])
				Expect(cert.Subject.CommonName).To(Equal(cn))

				By("Decoding CA certificates chain to DER")
				var chainDER []byte
				for _, ca := range credentials.CaChain {
					block, _ := pem.Decode([]byte(ca))
					Expect(block).ToNot(BeNil())
					chainDER = append(chainDER, block.Bytes...)
				}

				By("Parsing the CA certificates chain")
				chainCerts, err := x509.ParseCertificates(chainDER)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the certificate was signed by the Controller CA")
				Expect(cert.CheckSignatureFrom(chainCerts[0])).To(Succeed())

				By("Decoding CA certificates pool to DER")
				var poolDER []byte
				for _, ca := range credentials.CaPool {
					block, _ := pem.Decode([]byte(ca))
					Expect(block).ToNot(BeNil())
					poolDER = append(poolDER, block.Bytes...)
				}

				By("Parsing the CA certificates pool")
				poolCerts, err := x509.ParseCertificates(poolDER)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the CA pool contains the Controller CA")
				Expect(poolCerts).To(ContainElement(chainCerts[0]))
			})
		})
	})

	Describe("Errors", func() {
		It("Should return an error if payload is empty", func() {
			By("Requesting credentials from auth service")
			credentials, err := authSvcCli.RequestCredentials(
				ctx,
				&pb.Identity{
					Csr: "",
				},
			)
			Expect(err).To(HaveOccurred())
			Expect(credentials).To(BeNil())
		})

		It("Should return an error if payload is invalid", func() {
			By("Requesting credentials from auth service")
			credentials, err := authSvcCli.RequestCredentials(
				ctx,
				&pb.Identity{
					Csr: "123",
				},
			)
			Expect(err).To(HaveOccurred())
			Expect(credentials).To(BeNil())
		})
	})
})
