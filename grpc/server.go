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
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	cce "github.com/smartedgemec/controller-ce"
	"github.com/smartedgemec/controller-ce/pb"
)

const (
	// SNI is the server name for TLS when connecting to the Controller post-enrollment
	SNI = "v1.community.controller.mec"
	// EnrollmentSNI is the server name for TLS when connecting to the Controller for enrollment
	EnrollmentSNI = "v1.enroll.community.controller.mec"

	// This is the gRPC full RPC path (format: /${package}.${service}/${rpc})
	// for the authentication endpoint. The proto is defined here:
	// https://github.com/smartedgemec/schema/blob/master/pb/auth.proto
	//
	// This path is used to allow Appliances connected (to the EnrollmentSNI)
	// without a valid client certificate to hit the enrollment endpoint to
	// receive a certificate. It is similar to how a REST app may require a
	// session token for API paths other than /login.
	enrollmentMethod = "/openness.auth.AuthService/RequestCredentials"

	// TODO confirm this with Intel - see https://github.com/smartedgemec/controller-ce/pull/61/files#r285201296
	// The appliance's port that it listens on for gRPC connections from the
	// Controller. This will be removed in the future when the Appliance is
	// assumed to not be routable from the Controller in all cases. Instead the
	// Appliance will be required to open two TCP streams, one for outgoing
	// RPCs to the Controller and one for inbound.
	nodeGRPCPort = "8081"
)

// Server wraps grpc.Server
type Server struct {
	controller *cce.Controller
	grpc       *grpc.Server
}

// NewServer creates a new Server.
func NewServer(controller *cce.Controller, conf *tls.Config) *Server {
	s := &Server{
		controller: controller,
		grpc: grpc.NewServer(
			grpc.Creds(credentials.NewTLS(conf)),
			grpc.UnaryInterceptor(
				func(
					ctx context.Context,
					req interface{},
					info *grpc.UnaryServerInfo,
					handler grpc.UnaryHandler,
				) (resp interface{}, err error) {
					// apply checkAuth middleware
					if err := checkAuth(ctx,
						info.FullMethod); err != nil {
						return nil, err
					}
					return handler(ctx, req)
				},
			),
			grpc.StreamInterceptor(
				func(
					srv interface{},
					ss grpc.ServerStream,
					info *grpc.StreamServerInfo,
					handler grpc.StreamHandler,
				) error {
					// apply checkAuth middleware
					if err := checkAuth(ss.Context(),
						info.FullMethod); err != nil {
						return err
					}
					return handler(srv, ss)
				},
			),
		),
	}

	pb.RegisterAuthServiceServer(s.grpc, s)
	// TODO: register more services

	return s
}

// checkAuth is a middleware, applied inside the unary and stream interceptors,
// to ensure that if the enrollment server config was used (i.e. no client cert
// was provided) that only the enrollment endpoint is authorized.
func checkAuth(ctx context.Context, method string) error {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return fmt.Errorf("expected peer info in gRPC context")
	}
	tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "expected peer auth to be TLS, got %T", p.AuthInfo)
	}
	switch tlsInfo.State.ServerName {
	case EnrollmentSNI:
		if method != enrollmentMethod {
			return status.Errorf(codes.PermissionDenied, "unauthorized RPC: %s", method)
		}
		return nil
	case SNI:
		return nil
	default:
		return fmt.Errorf("unexpected server name: %s", tlsInfo.State.ServerName)
	}
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
func (s *Server) RequestCredentials(ctx context.Context, id *pb.Identity) (*pb.Credentials, error) { // nolint: gocyclo
	// Parse and validate CSR
	csr := id.GetCsr()
	if csr == "" {
		return nil, status.Error(codes.InvalidArgument, "CSR cannot be empty")
	}
	csrPEM, _ := pem.Decode([]byte(csr))
	if csrPEM == nil {
		return nil, status.Error(codes.InvalidArgument, "unable to decode CSR")
	}
	certReq, err := x509.ParseCertificateRequest(csrPEM.Bytes)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "error parsing CSR: %v", err)
	}

	// Node's identity is base64-encoded (w/o padding) MD5 hash of the public key data
	hash := md5.Sum(certReq.RawSubjectPublicKeyInfo)
	serial := base64.RawURLEncoding.EncodeToString(hash[:])

	// Verify the Node's pre-approval by public key data
	entities, err := s.controller.PersistenceService.Filter(ctx, &cce.Node{}, []cce.Filter{{
		Field: "entity->>'$.serial'",
		Value: serial,
	}})
	if err != nil || len(entities) == 0 {
		if err != nil {
			log.Errf("error getting node approval: %v", err)
		}
		return nil, status.Errorf(codes.Unauthenticated, "node %s not approved", serial)
	}
	node := entities[0].(*cce.Node)

	// Sign cert request
	cert, err := s.controller.AuthorityService.SignCSR(
		certReq.Raw,
		&x509.Certificate{
			Subject: pkix.Name{CommonName: node.ID},
		})
	if err != nil {
		log.Errf("Failed to sign CSR: %v", err)
		return nil, status.Error(codes.Internal, "unable to sign CSR")
	}
	certPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		})

	// Get signer chain for response
	caChain, err := s.controller.AuthorityService.CAChain()
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to get CA chain")
	}
	if len(caChain) == 0 {
		log.Errf("Failed to get CA chain: CA chain is empty")
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

	// Store Node credentials
	creds := &cce.Credentials{
		ID:          node.ID,
		Certificate: string(certPEM),
	}
	if err = s.controller.PersistenceService.Create(ctx, creds); err != nil {
		log.Errf("Failed to store Node credentials: %v", err)
		return nil, status.Error(codes.Internal, "unable to store credentials")
	}

	// Get the Node's IP address
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "missing peer data from context")
	}
	nodeIP, _, err := net.SplitHostPort(p.Addr.String())
	if nodeIP == "" || err != nil {
		return nil, status.Errorf(codes.Internal, "bad remote address in peer data: %s: %v",
			p.Addr.String(), err)
	}

	// Store the Node's address
	nodeWithTarget := &cce.Node{
		ID:         node.ID,
		Name:       node.Name,
		Location:   node.Location,
		Serial:     node.Serial,
		GRPCTarget: net.JoinHostPort(nodeIP, nodeGRPCPort),
	}
	if err := s.controller.PersistenceService.BulkUpdate(ctx, []cce.Persistable{
		nodeWithTarget,
	}); err != nil {
		log.Errf("Failed to store Node address: %v", err)
		return nil, status.Error(codes.Internal, "unable to store node address")
	}

	return &pb.Credentials{
		Certificate: creds.Certificate,
		CaChain:     chainPEM,
		CaPool:      caPoolPEM,
	}, nil
}
