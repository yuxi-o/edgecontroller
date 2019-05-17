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

	"github.com/smartedgemec/controller-ce/jose"
)

// Controller aggregates controller services.
type Controller struct {
	PersistenceService PersistenceService
	AuthorityService   AuthorityService
	TokenService       *jose.JWSTokenIssuer
	AdminCreds         *AuthCreds
}

// PersistenceService manages entity persistence. The methods with zv parameters take a zero-value Persistable for
// reflectively creating new instances of the concrete type. In the case of Delete it is used to get the table name.
type PersistenceService interface {
	Create(ctx context.Context, e Persistable) error
	Read(ctx context.Context, id string, zv Persistable) (e Persistable, err error)
	ReadAll(ctx context.Context, zv Persistable) (ps []Persistable, err error)
	Filter(ctx context.Context, zv Persistable, fs []Filter) (ps []Persistable, err error)
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
