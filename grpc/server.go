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

package grpc

import (
	"context"
	"encoding/pem"
	"log"
	"net"

	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cce "github.com/smartedgemec/controller-ce"
	"github.com/smartedgemec/controller-ce/pb"
)

// Server wraps grpc.Server
type Server struct {
	controller *cce.Controller
	grpc       *grpc.Server
}

// NewServer creates a new Server.
func NewServer(controller *cce.Controller) *Server {
	s := &Server{
		controller: controller,
		grpc:       grpc.NewServer(),
	}

	pb.RegisterAuthServiceServer(s.grpc, s)

	return s
}

// Serve wraps grpc.Server.Serve.
func (s *Server) Serve(lis net.Listener) error {
	return s.grpc.Serve(lis)
}

// GracefulStop wraps grpc.Server.GracefulStop.
func (s *Server) GracefulStop() {
	s.grpc.GracefulStop()
}

// Stop wraps grpc.Server.Stop.
func (s *Server) Stop() {
	s.grpc.Stop()
}

// RequestCredentials requests authentication endpoint credentials.
func (s *Server) RequestCredentials(
	ctx context.Context,
	id *pb.Identity,
) (*pb.Credentials, error) {
	csr := id.GetCsr()
	if csr == "" {
		return nil, status.Error(codes.InvalidArgument, "CSR cannot be empty")
	}

	csrPEM, _ := pem.Decode([]byte(csr))
	if csrPEM == nil {
		return nil, status.Error(codes.InvalidArgument, "unable to decode CSR")
	}

	// TODO: INTC-432: Verify the Node's public key

	cert, err := s.controller.AuthorityService.SignCSR(csrPEM.Bytes)
	if err != nil {
		log.Printf("Failed to sign CSR: %s\n", err.Error())
		return nil, status.Error(codes.Internal, "unable to sign CSR")
	}

	certPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		})

	caChain, err := s.controller.AuthorityService.CAChain()
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to get CA chain")
	}
	if len(caChain) == 0 {
		log.Println("Failed to get CA chain: CA chain is empty")
		return nil, status.Error(codes.Internal, "CA chain is empty")
	}

	// Encode each certificate in CA chain in PEM
	var chainPEM []string
	for _, caCert := range caChain {
		caPEM := pem.EncodeToMemory(
			&pem.Block{
				Type:  "CERTIFICATE",
				Bytes: caCert.Raw,
			},
		)
		chainPEM = append(chainPEM, string(caPEM))
	}

	// Add the root CA to the Node's CA pool
	caPoolPEM := chainPEM[len(chainPEM)-1:]

	creds := &cce.Credentials{
		ID:          cert.Subject.CommonName,
		Certificate: string(certPEM),
	}

	if err = s.controller.PersistenceService.Create(ctx, creds); err != nil {
		log.Printf("Failed to store credentials: %s\n", err.Error())
		return nil, status.Error(codes.Internal, "unable to store credentials")
	}

	// TODO: INTC-431: Store the Node's IP address

	return &pb.Credentials{
		Certificate: creds.Certificate,
		CaChain:     chainPEM,
		CaPool:      caPoolPEM,
	}, nil
}
