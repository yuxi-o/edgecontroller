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

package cce

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/otcshare/common/proxy/progutil"
	"github.com/otcshare/edgecontroller/jose"
	"github.com/otcshare/edgecontroller/k8s"
)

// Our network callback helper
var PrefaceLis *progutil.PrefaceListener

// OrchestrationMode global level orchestration mode for application deployment
type OrchestrationMode int

const (
	// OrchestrationModeNative uses Docker on the node to control application
	// container instances
	OrchestrationModeNative OrchestrationMode = iota
	// OrchestrationModeKubernetes uses an external Kubernetes master to
	// control application container instances on nodes
	OrchestrationModeKubernetes
	// OrchestrationModeKubernetesOVN uses an external Kubernetes master to
	// control application container instances on nodes. NTS functionality is
	// replaced with OVN/OVS. Traffic Policies are replaced with Kubernetes
	// Network Policies.
	OrchestrationModeKubernetesOVN
)

// Controller aggregates controller services.
type Controller struct {
	OrchestrationMode OrchestrationMode
	// must not be nil if
	// OrchestrationModeKubernetes or OrchestrationModeKubernetesOVN
	KubernetesClient   *k8s.Client
	PersistenceService PersistenceService
	AuthorityService   AuthorityService
	TokenService       *jose.JWSTokenIssuer
	AdminCreds         *AuthCreds

	// The edge node's port that it listens on for gRPC connections from the
	// Controller and serves Mm5-related endpoints for application and network
	// policy configuration.
	//
	// If ELAPort is empty the default of 42101 is used.
	ELAPort string

	// The edge node's port that it listens on for gRPC connections from the
	// Controller and serves Mm6-related endpoints for app deployment and
	// lifecycle commands.
	//
	// If EVAPort is empty the default of 42102 is used.
	EVAPort string

	// EdgeNodeCreds are the transport credentials for connecting to an edge
	// node. The server name will be overridden.
	EdgeNodeCreds *tls.Config
}

// PersistenceService manages entity persistence. The methods with zv parameters take a zero-value Persistable for
// reflectively creating new instances of the concrete type. In the case of Delete it is used to get the table name.
type PersistenceService interface {
	Create(ctx context.Context, e Persistable) error
	Read(ctx context.Context, id string, zv Persistable) (e Persistable, err error)
	ReadAll(ctx context.Context, zv Persistable) (ps []Persistable, err error)
	Filter(ctx context.Context, zv Filterable, fs []Filter) (ps []Persistable, err error)
	BulkUpdate(ctx context.Context, ps []Persistable) error
	Delete(ctx context.Context, id string, zv Persistable) (ok bool, err error)
}

// Validatable can be validated.
type Validatable interface {
	Validate() error
}

// Persistable can be persisted.
type Persistable interface {
	GetTableName() string
	GetID() string
	SetID(id string)
}

// Filterable is a Persistable that can be filtered.
type Filterable interface {
	Persistable
	FilterFields() []string
}

// ReqEntity is a request entity.
type ReqEntity interface {
	Validate() error
	GetTableName() string
}

// RespEntity is a response entity.
type RespEntity interface {
}

// NodeEntity has a node ID.
type NodeEntity interface {
	GetNodeID() string
}

// Filter filters queries in PersistenceService.Filter.
type Filter struct {
	Field string
	Value string
}

func getIp(ctx context.Context, ps PersistenceService, nodeId string) (string, error) {
	targets, err := ps.Filter(ctx, &NodeGRPCTarget{},
		[]Filter{
			{
				Field: "node_id",
				Value: nodeId,
			},
		})
	if err != nil {
		return "", err
	}

	if len(targets) > 0 {
		target := targets[0].(*NodeGRPCTarget).GRPCTarget
		return target, nil
	}
	return "", fmt.Errorf("IP for %v not found", nodeId)
}

// Inform the proxy we're serving this host
func RegisterToProxy(ctx context.Context, ps PersistenceService, nodeId string) {
	ip, err := getIp(ctx, ps, nodeId)
	if err != nil {
		return
	}
	PrefaceLis.RegisterHost(ip)
}
