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
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5" //nolint:gosec
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	cce "github.com/smartedgemec/controller-ce"
	cceGRPC "github.com/smartedgemec/controller-ce/grpc"
	authpb "github.com/smartedgemec/controller-ce/pb/auth"
	"github.com/smartedgemec/controller-ce/pki"
	"github.com/smartedgemec/controller-ce/swagger"
)

const adminPass = "word"

var (
	cmd    *exec.Cmd
	ctrl   *gexec.Session
	node   *gexec.Session
	nodeIn io.WriteCloser

	authSvcCli authpb.AuthServiceClient
	apiCli     *apiClient

	conf     *tls.Config
	telemDir string

	controllerRootPEM []byte
)

var _ = BeforeSuite(func() {
	logger := grpclog.NewLoggerV2(GinkgoWriter, GinkgoWriter, GinkgoWriter)
	grpclog.SetLoggerV2(logger)
	startup()
	initAuthSvcCli()
})

var _ = AfterSuite(func() {
	// time.Sleep(10 * time.Minute)
	shutdown()
})

func TestApplicationClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller CE API Suite")
}

func initAuthSvcCli() {
	timeoutCtx, cancel := context.WithTimeout(
		context.Background(), 2*time.Second)
	defer cancel()

	caPool := x509.NewCertPool()
	Expect(caPool.AppendCertsFromPEM(controllerRootPEM)).To(BeTrue(),
		"should load Controller self-signed root into trust pool")
	tlsCreds := credentials.NewClientTLSFromCert(caPool, cceGRPC.EnrollmentSNI)

	conn, err := grpc.DialContext(
		timeoutCtx,
		net.JoinHostPort("127.0.0.1", "8081"),
		grpc.WithTransportCredentials(tlsCreds),
		grpc.WithBlock())
	Expect(err).ToNot(HaveOccurred(), "Dial failed: %v", err)

	authSvcCli = authpb.NewAuthServiceClient(conn)
}

func startup() {
	By("Building the controller")
	exe, err := gexec.Build("github.com/smartedgemec/controller-ce/cmd/cce")
	Expect(err).ToNot(HaveOccurred(), "Problem building service")

	By("Creating a temp dir for telemetry output files")
	tmpdir, err := ioutil.TempDir(".", "telemetry")
	Expect(err).NotTo(HaveOccurred())
	telemDir, err = filepath.Abs(tmpdir)
	Expect(err).NotTo(HaveOccurred())

	By("Starting the controller")
	cmd = exec.Command(exe,
		"-log-level", "debug",
		"-dsn", "root:beer@tcp(:8083)/controller_ce",
		"-httpPort", "8080",
		"-grpcPort", "8081",
		"-elaPort", "42101",
		"-evaPort", "42102",
		"-syslogPort", "6514",
		"-statsdPort", "8125",
		"-syslog-path", filepath.Join(telemDir, "syslog.log"),
		"-statsd-path", filepath.Join(telemDir, "statsd.log"),
		"-adminPass", adminPass)
	ctrl, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred(), "Problem starting service")

	By("Reading the Controller self-signed CA from output")
	Eventually(ctrl.Err, 3).Should(gbytes.Say(
		`-----END CERTIFICATE-----`),
		"Service did not print CA cert in time")
	certMatches := regexp.MustCompile(
		`(?s)-----BEGIN CERTIFICATE-----.*?-----END CERTIFICATE-----`,
	).FindAll(ctrl.Err.Contents(), -1)
	Expect(certMatches).To(HaveLen(1),
		"Service did not print a single CA cert")
	controllerRootPEM = certMatches[0]
	conf = loadTLSConfig(filepath.Join(".", "certificates", "ca"))

	By("Verifying that the controller started successfully")
	Eventually(ctrl.Err, 3).Should(gbytes.Say(
		"Controller CE ready"),
		"Service did not start in time")

	By("Requesting an authentication token from the controller")
	apiCli = &apiClient{
		Token: authToken(),
	}

	By("Building the node")
	exe, err = gexec.Build(
		"github.com/smartedgemec/controller-ce/test/node/grpc")
	Expect(err).ToNot(HaveOccurred(), "Problem building node")

	cmd = exec.Command(exe,
		"-ela-port", "42101",
		"-eva-port", "42102",
	)
	nodeIn, err = cmd.StdinPipe()
	Expect(err).ToNot(HaveOccurred(), "Problem creating node stdin pipe")

	By("Starting the node")
	node, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred(), "Problem starting node")
}

func shutdown() {
	if ctrl != nil {
		By("Stopping the controller service")
		ctrl.Kill()
	}
	if node != nil {
		By("Stopping the test node")
		node.Kill()
	}
	if nodeIn != nil {
		nodeIn.Close()
	}
	if telemDir != "" {
		By("Cleaning up telemetry output")
		Expect(os.RemoveAll(telemDir)).To(Succeed())
	}
}

func clearGRPCTargetsTable() {
	By("Connecting to the database")
	db, err := sql.Open("mysql", "root:beer@tcp(:8083)/controller_ce?multiStatements=true")
	Expect(err).ToNot(HaveOccurred())

	defer func() {
		Expect(db.Close()).To(Succeed())
	}()

	By("Pinging the database")
	err = db.Ping()
	Expect(err).ToNot(HaveOccurred())

	timeoutCtx, cancel := context.WithTimeout(
		context.Background(), 2*time.Second)
	defer cancel()

	By("Executing the delete query")
	_, err = db.ExecContext(
		timeoutCtx,
		"DELETE FROM node_grpc_targets")
	Expect(err).ToNot(HaveOccurred())
}

func authToken() string {
	payload, err := json.Marshal(
		struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{"admin", adminPass})
	Expect(err).ToNot(HaveOccurred())

	req, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/auth",
		bytes.NewReader(payload),
	)
	Expect(err).ToNot(HaveOccurred())

	resp, err := new(http.Client).Do(req)
	Expect(err).ToNot(HaveOccurred())
	defer resp.Body.Close()

	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	var auth struct {
		Token string `json:"token"`
	}
	Expect(json.NewDecoder(resp.Body).Decode(&auth)).To(Succeed())
	Expect(auth.Token).ToNot(BeEmpty())

	return auth.Token
}

type respBody struct {
	ID string
}

func postApps(appType string) (id string) {
	By("Sending a POST /apps request")
	resp, err := apiCli.Post(
		"http://127.0.0.1:8080/apps",
		"application/json",
		strings.NewReader(fmt.Sprintf(`
			{
				"type": "%s",
				"name": "%s app",
				"version": "latest",
				"vendor": "smart edge",
				"description": "my %s app",
				"cores": 4,
				"memory": 1024,
				"ports": [{"port": 80, "protocol": "tcp"}],
				"source": "http://www.test.com/my_%s_app.tar.gz"
			}`, appType, appType, appType, appType)))
	Expect(err).ToNot(HaveOccurred())
	defer resp.Body.Close()

	By("Verifying a 201 Created response")
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var rb respBody

	By("Unmarshaling the response")
	Expect(json.Unmarshal(body, &rb)).To(Succeed())

	return rb.ID
}

func getApp(id string) *swagger.AppDetail {
	By("Sending a GET /apps/{app_id} request")
	resp, err := apiCli.Get(
		fmt.Sprintf("http://127.0.0.1:8080/apps/%s", id))
	Expect(err).ToNot(HaveOccurred())
	defer resp.Body.Close()

	By("Verifying a 200 OK response")
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var app swagger.AppDetail

	By("Unmarshaling the response")
	Expect(json.Unmarshal(body, &app)).To(Succeed())

	return &app
}

func postDNSConfigs() (id string) {
	By("Sending a POST /dns_configs request")
	resp, err := apiCli.Post(
		"http://127.0.0.1:8080/dns_configs",
		"application/json",
		strings.NewReader(`
			{
				"name": "dns config 123",
				"a_records": [{
					"name": "a record 1",
					"description": "description 1",
					"ips": [
						"172.16.55.43",
						"172.16.55.44"
					]
				}],
				"forwarders": [{
					"name": "forwarder 1",
					"description": "description 1",
					"ip": "8.8.8.8"
				}, {
					"name": "forwarder 2",
					"description": "description 2",
					"ip": "1.1.1.1"
				}]
			}`))
	Expect(err).ToNot(HaveOccurred())
	defer resp.Body.Close()

	By("Verifying a 201 Created response")
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var rb respBody

	By("Unmarshaling the response")
	Expect(json.Unmarshal(body, &rb)).To(Succeed())

	return rb.ID
}

func getDNSConfig(id string) *cce.DNSConfig {
	By("Sending a GET /dns_configs/{id} request")
	resp, err := apiCli.Get(
		fmt.Sprintf("http://127.0.0.1:8080/dns_configs/%s", id))
	Expect(err).ToNot(HaveOccurred())
	defer resp.Body.Close()

	By("Verifying a 200 OK response")
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var dnsConfig cce.DNSConfig

	By("Unmarshaling the response")
	Expect(json.Unmarshal(body, &dnsConfig)).To(Succeed())

	return &dnsConfig
}

func postDNSConfigsAppAliases(
	dnsConfigID string,
	appID string,
) (id string) {
	By("Sending a POST /dns_configs_app_aliases request")
	resp, err := apiCli.Post(
		"http://127.0.0.1:8080/dns_configs_app_aliases",
		"application/json",
		strings.NewReader(fmt.Sprintf(`
			{
				"dns_config_id": "%s",
				"name": "dns config app alias",
				"description": "my dns config app alias",
				"app_id": "%s"
			}`, dnsConfigID, appID)))
	Expect(err).ToNot(HaveOccurred())
	defer resp.Body.Close()

	By("Verifying a 201 Created response")
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var rb respBody

	By("Unmarshaling the response")
	Expect(json.Unmarshal(body, &rb)).To(Succeed())

	return rb.ID
}

func getDNSConfigsAppAlias(id string) *cce.DNSConfigAppAlias {
	By("Sending a GET /dns_configs_app_aliases/{id} request")
	resp, err := apiCli.Get(
		fmt.Sprintf("http://127.0.0.1:8080/dns_configs_app_aliases/%s", id))
	Expect(err).ToNot(HaveOccurred())
	defer resp.Body.Close()

	By("Verifying a 200 OK response")
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var dnsConfigAppAlias cce.DNSConfigAppAlias

	By("Unmarshaling the response")
	Expect(json.Unmarshal(body, &dnsConfigAppAlias)).To(Succeed())

	return &dnsConfigAppAlias
}

type nodeConfig struct {
	nodeID string
	serial string
	key    *ecdsa.PrivateKey
	creds  *authpb.Credentials
}

func createAndRegisterNode() *nodeConfig {
	By("Generating node private key")
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	Expect(err).ToNot(HaveOccurred())

	By("Creating a CSR with private key")
	csrDER, err := x509.CreateCertificateRequest(
		rand.Reader,
		&x509.CertificateRequest{},
		key,
	)
	Expect(err).ToNot(HaveOccurred())

	By("Parsing the CSR")
	certReq, err := x509.ParseCertificateRequest(csrDER)
	Expect(err).ToNot(HaveOccurred())

	By("Encoding the CSR in PEM")
	csrPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE REQUEST",
			Bytes: csrDER,
		})

	By("Pre-approving Node by serial")
	hash := md5.Sum(certReq.RawSubjectPublicKeyInfo) //nolint:gosec
	serial := base64.RawURLEncoding.EncodeToString(hash[:])
	nodeID := postNodesSerial(serial)

	By("Resetting the node")
	Expect(cmd.Process.Signal(syscall.SIGABRT)).To(Succeed(), "Problem resetting node")
	Expect(fmt.Fprintln(nodeIn, nodeID)).To(Equal(len(nodeID) + 1))

	By("Verifying that the node started successfully")
	Eventually(node.Err, 3).Should(gbytes.Say(
		"test-node: listening on port: 4210[12]"),
		"Node did not start in time")
	Eventually(node.Err, 3).Should(gbytes.Say(
		"test-node: listening on port: 4210[12]"),
		"Node did not start in time")

	By("Requesting credentials from auth service")
	creds, err := authSvcCli.RequestCredentials(
		context.TODO(),
		&authpb.Identity{
			Csr: string(csrPEM),
		},
	)
	Expect(err).ToNot(HaveOccurred())

	return &nodeConfig{
		nodeID: nodeID,
		serial: serial,
		key:    key,
		creds:  creds,
	}
}

func postNodesSerial(serial string) (id string) {
	By("Sending a POST /nodes request")
	resp, err := apiCli.Post(
		"http://127.0.0.1:8080/nodes",
		"application/json",
		strings.NewReader(fmt.Sprintf(`
			{
				"name": "Test Node 1",
				"location": "Localhost port 42101",
				"serial": "%s"
			}`, serial)))
	Expect(err).ToNot(HaveOccurred())
	defer resp.Body.Close()

	By("Verifying a 201 Created response")
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var rb respBody

	By("Unmarshaling the response")
	Expect(json.Unmarshal(body, &rb)).To(Succeed())

	return rb.ID
}

func getNode(id string) *cce.NodeResp {
	By("Sending a GET /nodes/{id} request")
	resp, err := apiCli.Get(
		fmt.Sprintf("http://127.0.0.1:8080/nodes/%s", id))
	Expect(err).ToNot(HaveOccurred())
	defer resp.Body.Close()

	By("Verifying a 200 OK response")
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var nodeResp cce.NodeResp

	By("Unmarshaling the response")
	Expect(json.Unmarshal(body, &nodeResp)).To(Succeed())

	return &nodeResp
}

func postNodesApps(nodeID, appID string) (id string) {
	By("Sending a POST /nodes_apps request")
	resp, err := apiCli.Post(
		"http://127.0.0.1:8080/nodes_apps",
		"application/json",
		strings.NewReader(fmt.Sprintf(`
			{
				"node_id": "%s",
				"app_id": "%s"
			}`, nodeID, appID)))
	Expect(err).ToNot(HaveOccurred())
	defer resp.Body.Close()

	By("Verifying a 201 Created response")
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var rb respBody

	By("Unmarshaling the response")
	Expect(json.Unmarshal(body, &rb)).To(Succeed())

	return rb.ID
}

func getNodeApp(id string) *cce.NodeAppResp {
	By("Sending a GET /nodes_apps/{id} request")
	resp, err := apiCli.Get(
		fmt.Sprintf("http://127.0.0.1:8080/nodes_apps/%s", id))

	By("Verifying a 200 OK response")
	Expect(err).ToNot(HaveOccurred())
	defer resp.Body.Close()
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var nodeAppResp cce.NodeAppResp

	By("Unmarshaling the response")
	Expect(json.Unmarshal(body, &nodeAppResp)).To(Succeed())

	return &nodeAppResp
}

func getNodeApps(nodeID, appID string) []*cce.NodeApp {
	By("Sending a GET /nodes_apps request")
	url := fmt.Sprintf("http://127.0.0.1:8080/nodes_apps?node_id=%s", nodeID)
	if appID != "" {
		url += "&app_id=" + appID
	}
	resp, err := apiCli.Get(url)

	By("Verifying a 200 OK response")
	Expect(err).ToNot(HaveOccurred())
	defer resp.Body.Close()
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var nodeApps []*cce.NodeApp

	By("Unmarshaling the response")
	Expect(json.Unmarshal(body, &nodeApps)).To(Succeed())

	return nodeApps
}

func postNodesDNSConfigs(nodeID, dnsConfigID string) (id string) {
	By("Sending a POST /nodes_dns_configs request")
	resp, err := apiCli.Post(
		"http://127.0.0.1:8080/nodes_dns_configs",
		"application/json",
		strings.NewReader(fmt.Sprintf(`
			{
				"node_id": "%s",
				"dns_config_id": "%s"
			}`, nodeID, dnsConfigID)))
	Expect(err).ToNot(HaveOccurred())
	defer resp.Body.Close()

	By("Verifying a 201 Created response")
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var rb respBody

	By("Unmarshaling the response")
	Expect(json.Unmarshal(body, &rb)).To(Succeed())

	return rb.ID
}

func getNodeDNSConfig(id string) *cce.NodeDNSConfig {
	By("Sending a GET /nodes_dns_configs/{id} request")
	resp, err := apiCli.Get(
		fmt.Sprintf("http://127.0.0.1:8080/nodes_dns_configs/%s", id))
	Expect(err).ToNot(HaveOccurred())
	defer resp.Body.Close()

	By("Verifying a 200 OK response")
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var nodeDNSConfig cce.NodeDNSConfig

	By("Unmarshaling the response")
	Expect(json.Unmarshal(body, &nodeDNSConfig)).To(Succeed())

	return &nodeDNSConfig
}

func postTrafficPolicies() (id string) {
	By("Sending a POST /traffic_policies request")
	resp, err := apiCli.Post(
		"http://127.0.0.1:8080/traffic_policies",
		"application/json",
		strings.NewReader(`
		{
			"rules": [{
				"description": "test-rule-1",
				"priority": 1,
				"source": {
					"description": "test-source-1",
					"macs": {
						"mac_addresses": [
							"F0-59-8E-7B-36-8A",
							"23-20-8E-15-89-D1",
							"35-A4-38-73-35-45"
						]
					},
					"ip": {
						"address": "223.1.1.0",
						"mask": 16,
						"begin_port": 2000,
						"end_port": 2012,
						"protocol": "tcp"
					},
					"gtp": {
						"address": "10.6.7.2",
						"mask": 12,
						"imsis": [
							"310150123456789",
							"310150123456790",
							"310150123456791"
						]
					}
				},
				"destination": {
					"description": "test-destination-1",
					"macs": {
						"mac_addresses": [
							"7D-C2-3A-1C-63-D9",
							"E9-6B-D1-D2-1A-6B",
							"C8-32-A9-43-85-55"
						]
					},
					"ip": {
						"address": "64.1.1.0",
						"mask": 16,
						"begin_port": 1000,
						"end_port": 1012,
						"protocol": "tcp"
					},
					"gtp": {
						"address": "108.6.7.2",
						"mask": 4,
						"imsis": [
							"310150123456792",
							"310150123456793",
							"310150123456794"
						]
					}
				},
				"target": {
					"description": "test-target-1",
					"action": "accept",
					"mac": {
						"mac_address": "C7-5A-E7-98-1B-A3"
					},
					"ip": {
						"address": "123.2.3.4",
						"port": 1600
					}
				}
			}]
		}`))
	Expect(err).ToNot(HaveOccurred())
	defer resp.Body.Close()

	By("Verifying a 201 Created response")
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var rb respBody

	By("Unmarshaling the response")
	Expect(json.Unmarshal(body, &rb)).To(Succeed())

	return rb.ID
}

func getTrafficPolicy(id string) *cce.TrafficPolicy {
	By("Sending a GET /traffic_policies/{id} request")
	resp, err := apiCli.Get(
		fmt.Sprintf("http://127.0.0.1:8080/traffic_policies/%s", id))
	Expect(err).ToNot(HaveOccurred())
	defer resp.Body.Close()

	By("Verifying a 200 OK response")
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var trafficPolicy cce.TrafficPolicy

	By("Unmarshaling the response")
	Expect(json.Unmarshal(body, &trafficPolicy)).To(Succeed())

	return &trafficPolicy
}

func postNodesAppsTrafficPolicies(
	nodeAppID string,
	trafficPolicyID string,
) (id string) {
	By("Sending a POST /nodes_apps_traffic_policies request")
	resp, err := apiCli.Post(
		"http://127.0.0.1:8080/nodes_apps_traffic_policies",
		"application/json",
		strings.NewReader(fmt.Sprintf(`
			{
				"nodes_apps_id": "%s",
				"traffic_policy_id": "%s"
			}`, nodeAppID, trafficPolicyID)))
	Expect(err).ToNot(HaveOccurred())
	defer resp.Body.Close()

	By("Verifying a 201 Created response")
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var rb respBody

	By("Unmarshaling the response")
	Expect(json.Unmarshal(body, &rb)).To(Succeed())

	return rb.ID
}

func getNodeAppTrafficPolicy(id string) *cce.NodeAppTrafficPolicy {
	By("Sending a GET /nodes_apps_traffic_policies/{id} request")
	resp, err := apiCli.Get(
		fmt.Sprintf("http://127.0.0.1:8080/nodes_apps_traffic_policies/%s", id))
	Expect(err).ToNot(HaveOccurred())
	defer resp.Body.Close()

	By("Verifying a 200 OK response")
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var nodeAppTrafficPolicy cce.NodeAppTrafficPolicy

	By("Unmarshaling the response")
	Expect(json.Unmarshal(body, &nodeAppTrafficPolicy)).To(Succeed())

	return &nodeAppTrafficPolicy
}

func loadTLSConfig(dir string) *tls.Config {
	key, err := pki.LoadKey(filepath.Join(dir, "key.pem"))
	Expect(err).NotTo(HaveOccurred())
	cert, err := pki.LoadCertificate(filepath.Join(dir, "cert.pem"))
	Expect(err).NotTo(HaveOccurred())
	certPool := x509.NewCertPool()
	certPool.AddCert(cert)
	return &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{cert.Raw},
			PrivateKey:  key,
			Leaf:        cert,
		}},
		RootCAs:      certPool,
		MinVersion:   tls.VersionTLS12,
		CipherSuites: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256},
	}
}
