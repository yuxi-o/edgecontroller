// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package nfd_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/otcshare/common/log"
	"github.com/otcshare/edgecontroller/nfd-master"
	"github.com/otcshare/edgecontroller/pki"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"os"
	pb "sigs.k8s.io/node-feature-discovery/pkg/labeler"
	"testing"
	"time"
)

func TestNfdMaster(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NFD Master test suite")
}

var (
	tempdir string
	ctx     context.Context
	rootCA  *pki.RootCA
)

var _ = BeforeSuite(func() {
	log.SetOutput(GinkgoWriter)

	var err error
	tempdir, err = ioutil.TempDir("", "nfd-master-test")
	Expect(err).NotTo(HaveOccurred())

	rootCA, err = pki.InitRootCA(tempdir)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	os.RemoveAll(tempdir)
})

func createClientTLSCreds(ca *pki.RootCA) credentials.TransportCredentials {

	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	Expect(err).NotTo(HaveOccurred())

	cert, err := ca.NewTLSClientCert(key, "localhost")
	Expect(err).NotTo(HaveOccurred())

	certPool := x509.NewCertPool()
	certPool.AddCert(ca.Cert)

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{cert.Raw},
			PrivateKey:  key,
		}},
		RootCAs: certPool,
	})
}

var _ = Describe("NFD Master Server", func() {
	When("Is initialized with invalid CA certificate file path", func() {
		It("Returns error", func() {
			nfdSrv := &nfd.ServerNFD{
				CaCertPath: "/invalid/path/cert.pem",
				CaKeyPath:  tempdir + "key.pem",
			}
			err := nfdSrv.ServeGRPC(ctx)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(
				"Failed to load CA certificate: unable to read certificate file: open " +
					"/invalid/path/cert.pem: no such file or directory"))
		})
	})
	When("Is initialized with valid CA certificate but invalid CA key path", func() {
		It("Returns error", func() {
			nfdSrv := &nfd.ServerNFD{
				CaCertPath: tempdir + "/cert.pem",
				CaKeyPath:  "/invalid/path/key.pem",
			}
			err := nfdSrv.ServeGRPC(ctx)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(
				"Failed to load CA key: unable to read key file: open " +
					"/invalid/path/key.pem: no such file or directory"))
		})
	})
	When("Is initialized with invalid SNI", func() {
		It("Starts successfully but SetLabels request fails", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			nfdSrv := &nfd.ServerNFD{
				CaCertPath: tempdir + "/cert.pem",
				CaKeyPath:  tempdir + "/key.pem",
				Endpoint:   8082,
				Sni:        "testhost",
			}
			var err error
			go func() {
				err = nfdSrv.ServeGRPC(ctx)
				Expect(err).ShouldNot(HaveOccurred())

			}()

			// wait for server to start
			time.Sleep(1 * time.Second)

			transportCreds := createClientTLSCreds(rootCA)
			conn, err := grpc.Dial("localhost:8082", grpc.WithTransportCredentials(transportCreds))
			Expect(err).ShouldNot(HaveOccurred())
			defer conn.Close()

			client := pb.NewLabelerClient(conn)
			req := pb.SetLabelsRequest{NodeName: "testnode"}

			_, err = client.SetLabels(ctx, &req)
			Expect(err.Error()).To(ContainSubstring(
				"authentication handshake failed: x509: certificate is valid for testhost, not localhost"))
		})
	})
	When("Is initialized with valid CA certificate and key path and valid SNI", func() {
		It("Starts successfully and returns no error on SetLabels request", func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			nfdSrv := &nfd.ServerNFD{
				CaCertPath: tempdir + "/cert.pem",
				CaKeyPath:  tempdir + "/key.pem",
				Endpoint:   8082,
				Sni:        "localhost",
			}
			var err error
			go func() {
				err = nfdSrv.ServeGRPC(ctx)
				Expect(err).ShouldNot(HaveOccurred())

			}()

			// wait for server to start
			time.Sleep(1 * time.Second)

			transportCreds := createClientTLSCreds(rootCA)
			conn, err := grpc.Dial("localhost:8082", grpc.WithTransportCredentials(transportCreds))
			Expect(err).ShouldNot(HaveOccurred())
			defer conn.Close()

			client := pb.NewLabelerClient(conn)
			req := pb.SetLabelsRequest{NodeName: "testnode"}

			resp, err := client.SetLabels(ctx, &req)
			Expect(resp).ShouldNot(Equal(nil))
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
