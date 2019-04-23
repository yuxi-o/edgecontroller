// Copyright 2019 Smart-Edge.com, Inc. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	nodegmock "github.com/smartedgemec/controller-ce/mock/node/grpc"
	"github.com/smartedgemec/controller-ce/pb"
	"google.golang.org/grpc"
)

const name = "test-node"

func main() {
	var (
		err  error
		port uint
	)
	log.Print(name, ": starting")

	// CLI flags
	flag.UintVar(&port, "port", 8080, "Port for service to listen on")
	flag.Parse()

	// Set up channels to capture SIGINT and SIGTERM
	sigChan := make(chan os.Signal, 2)
	defer close(sigChan)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create the mock node
	mockNode := nodegmock.NewMockNode()

	// Register the services with the grpc server
	server := grpc.NewServer()
	pb.RegisterApplicationDeploymentServiceServer(server, mockNode.AppDeploySvc)
	pb.RegisterApplicationLifecycleServiceServer(server, mockNode.AppLifeSvc)
	pb.RegisterApplicationPolicyServiceServer(server, mockNode.AppPolicySvc)
	pb.RegisterVNFDeploymentServiceServer(server, mockNode.VNFDeploySvc)
	pb.RegisterVNFLifecycleServiceServer(server, mockNode.VNFLifeSvc)
	pb.RegisterInterfaceServiceServer(server, mockNode.InterfaceSvc)
	pb.RegisterInterfacePolicyServiceServer(server, mockNode.IfPolicySvc)
	pb.RegisterZoneServiceServer(server, mockNode.ZoneSvc)

	// Shut down the server gracefully
	go func() {
		<-sigChan
		log.Printf("%s: shutting down", name)

		server.GracefulStop()
	}()

	// Start the listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("Error listening on port:", err)
	}
	defer listener.Close()

	// Start the server
	log.Print(name, ": listening on port: ",
		listener.Addr().(*net.TCPAddr).Port)
	err = server.Serve(listener)
	if err != nil && context.Canceled == nil {
		log.Fatal("Error starting GRPC server:", err)
	}
}
