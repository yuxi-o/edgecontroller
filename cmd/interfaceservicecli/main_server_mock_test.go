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

package main_test

import (
	"context"
	"fmt"
	"net"
	"path/filepath"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/otcshare/edgecontroller/pb/ela"
	"google.golang.org/grpc"
)

type InterfaceServiceServer struct {
	Endpoint string
	server   *grpc.Server

	updateReturnEmpty *empty.Empty
	updateReturnErr   error

	bulkUpdateReturnEmpty *empty.Empty
	bulkUpdateReturnErr   error

	getAllReturnNi  *pb.NetworkInterfaces
	getAllReturnErr error

	getReturnNi  *pb.NetworkInterface
	getReturnErr error
}

func (is *InterfaceServiceServer) StartServer() error {
	fmt.Println("Starting IP API at: ", is.Endpoint)

	tc, err := readTestPKICredentials(filepath.Clean("./certs/s_cert.pem"),
		filepath.Clean("./certs/s_key.pem"),
		filepath.Clean("./certs/cacerts.pem"))
	if err != nil {
		return fmt.Errorf("failed to read pki: %v", err)
	}

	lis, err := net.Listen("tcp", is.Endpoint)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	is.server = grpc.NewServer(grpc.Creds(*tc))
	pb.RegisterInterfaceServiceServer(is.server, is)
	go func() {
		if err := is.server.Serve(lis); err != nil {
			fmt.Printf("API listener exited unexpectedly: %s", err)
		}
	}()
	return nil
}

// GracefulStop shuts down connetions and removes the Unix domain socket
func (is *InterfaceServiceServer) GracefulStop() error {
	is.server.GracefulStop()
	return nil
}

func (is *InterfaceServiceServer) Update(context.Context, *pb.NetworkInterface) (*empty.Empty, error) {
	fmt.Println("@@@ 'Update' from GRPC server @@@")
	return &empty.Empty{}, is.updateReturnErr
}

func (is *InterfaceServiceServer) BulkUpdate(context.Context, *pb.NetworkInterfaces) (*empty.Empty, error) {
	fmt.Println("@@@ 'BulkUpdate' from GRPC server @@@")
	return &empty.Empty{}, is.bulkUpdateReturnErr
}

func (is *InterfaceServiceServer) GetAll(context.Context, *empty.Empty) (*pb.NetworkInterfaces, error) {
	fmt.Println("@@@ 'GetAll' from GRPC server @@@")
	return is.getAllReturnNi, is.getAllReturnErr
}

func (is *InterfaceServiceServer) Get(context.Context, *pb.InterfaceID) (*pb.NetworkInterface, error) {
	fmt.Println("@@@ 'Get' from GRPC server @@@")
	return &pb.NetworkInterface{}, is.getReturnErr
}
