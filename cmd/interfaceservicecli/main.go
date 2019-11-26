// Copyright 2019 Intel Corporation. All rights reserved.
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

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/otcshare/edgecontroller/pb/ela"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type cliFlags struct {
	CertsDir    string
	Endpoint    string
	ServiceName string
	Timeout     int
	Cmd         string
	Val         string
}

// Cfg stores flags passed to CLI
var Cfg cliFlags

func init() {
	flag.StringVar(&Cfg.Endpoint, "endpoint", "", "Interface service endpoint")
	flag.StringVar(&Cfg.ServiceName, "servicename", "interfaceservice.openness", "Name of server in certificate")
	flag.StringVar(&Cfg.Cmd, "cmd", "help", "Interface service command")
	flag.StringVar(&Cfg.Val, "val", "", "Interface service command parameters")
	flag.StringVar(&Cfg.CertsDir, "certsdir", "./certs/client/interfaceservice", "Directory of key and certificate")
	flag.IntVar(&Cfg.Timeout, "timeout", 3, "Timeout value for grpc call (in seconds)")
}

func getTransportCredentials() (*credentials.TransportCredentials, error) {
	crtPath := filepath.Clean(filepath.Join(Cfg.CertsDir, "cert.pem"))
	keyPath := filepath.Clean(filepath.Join(Cfg.CertsDir, "key.pem"))
	caPath := filepath.Clean(filepath.Join(Cfg.CertsDir, "root.pem"))

	cert, err := tls.LoadX509KeyPair(crtPath, keyPath)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, err
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, errors.Errorf("Failed append CA certs from %s", caPath)
	}

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
		ServerName:   Cfg.ServiceName,
	})

	return &creds, nil
}

func createConnection(ctx context.Context) *grpc.ClientConn {
	tc, err := getTransportCredentials()
	if err != nil {
		fmt.Println("Error when creating transport credentials: " + err.Error())
		os.Exit(1)
	}

	conn, err := grpc.DialContext(ctx, Cfg.Endpoint,
		grpc.WithTransportCredentials(*tc), grpc.WithBlock())

	if err != nil {
		fmt.Println("Error when dialing: " + Cfg.Endpoint + " err:" + err.Error())
		os.Exit(1)
	}

	return conn
}

func printHelp() {
	fmt.Print(`
    Get or attach/detach network interfaces to OVS on remote edge node

    -endpoint      Endpoint to be requested
    -servicename   Name to be used as server name for TLS handshake
    -cmd           Supported commands: get, attach, detach
    -val           PCI address for attach and detach commands. Multiple addresses can be passed
                   and must be separated by commas: -val=0000:00:00.0,0000:00:00.1
    -certsdir      Directory where cert.pem and key.pem for client and root.pem for CA resides   
    -timeout       Timeout value [s] for grpc requests

	`)
}

func splitAndValidatePCIFormat(val string) []string {
	devs := strings.Split(val, ",")
	var validPCIs []string

	// 0000:00:00.0
	for _, dev := range devs {
		s := strings.Split(dev, ":")
		if len(s) == 3 && len(s[0]) == 4 && len(s[1]) == 2 && len(s[2]) == 4 {
			validPCIs = append(validPCIs, dev)
		} else {
			fmt.Println("Invalid PCI address: " + dev + ". Skipping...")
		}
	}
	return validPCIs
}

func updateInterfaces(driver pb.NetworkInterface_InterfaceDriver, pcis string) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Cfg.Timeout)*time.Second)
	defer cancel()

	conn := createConnection(ctx)
	defer conn.Close()

	client := pb.NewInterfaceServiceClient(conn)

	var ifsReq []*pb.NetworkInterface

	ifsActual, err := client.GetAll(ctx, &empty.Empty{})
	if err != nil {
		return err
	}

	addr := splitAndValidatePCIFormat(pcis)
	for _, a := range addr {
		found := false
		for _, i := range ifsActual.GetNetworkInterfaces() {
			if i.GetId() == a {
				i.Driver = driver
				ifsReq = append(ifsReq, i)
				found = true
			}
		}

		if !found {
			fmt.Println("Interface: " + a + " not found. Skipping...")
		}
	}

	if len(ifsReq) != 0 {
		_, err = client.BulkUpdate(ctx, &pb.NetworkInterfaces{
			NetworkInterfaces: ifsReq,
		})
	}

	if err != nil {
		return err
	}

	op := "attached"
	if driver == pb.NetworkInterface_KERNEL {
		op = "detached"
	}

	for _, i := range ifsReq {
		fmt.Println("Interface: " + i.GetId() + " successfully " + op)
	}

	return nil
}

func printInterfaces() error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Cfg.Timeout)*time.Second)
	defer cancel()

	conn := createConnection(ctx)
	defer conn.Close()

	client := pb.NewInterfaceServiceClient(conn)

	var err error
	ifs, err := client.GetAll(ctx, &empty.Empty{})
	if err != nil {
		return err
	}

	if len(ifs.GetNetworkInterfaces()) == 0 {
		return errors.Errorf("No interfaces on node found")
	}

	for _, dev := range ifs.GetNetworkInterfaces() {
		drv := "detached"
		if dev.GetDriver() == 1 {
			drv = "attached"
		}
		fmt.Printf("%s  |  %s  |  %s\n", dev.GetId(), dev.GetMacAddress(), drv)
	}

	return nil
}

func main() {
	flag.Parse()

	if err := StartCli(); err != nil {
		fmt.Println("Error when executing command: [" + Cfg.Cmd + "] err: " + err.Error())
		os.Exit(1)
	}
}

// StartCli handles command and arguments to call corresponding CLI function
func StartCli() error {
	var err error

	switch Cfg.Cmd {
	case "attach":
		err = updateInterfaces(pb.NetworkInterface_USERSPACE, Cfg.Val)
	case "detach":
		err = updateInterfaces(pb.NetworkInterface_KERNEL, Cfg.Val)
	case "get":
		err = printInterfaces()
	case "help", "h", "":
		printHelp()
	default:
		fmt.Println("Unrecognized action: " + Cfg.Cmd)
		printHelp()
	}

	return err
}
