// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package main_test

import (
	"context"
	"fmt"
	"net"
	"path/filepath"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/open-ness/edgecontroller/pb/ela"
	"google.golang.org/grpc"
)

type InterfaceServiceServer struct {
	Endpoint string
	server   *grpc.Server

	updateReturnErr     error
	bulkUpdateReturnErr error

	getAllReturnNi  *pb.NetworkInterfaces
	getAllReturnErr error

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
