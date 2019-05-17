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

// PersistenceService manages entity persistence. The methods with EntityModel
// parameters take a zero-value Entity for reflectively creating new instances
// of the concrete type. In the case of Delete it is used to get the table name.
type PersistenceService interface {
	Create(ctx context.Context, e Entity) error
	Read(ctx context.Context, id string, zv EntityModel) (e Entity, err error)
	ReadAll(ctx context.Context, zv EntityModel) (es []Entity, err error)
	Filter(ctx context.Context,
		zv EntityModel, fs []Filter) (es []Entity, err error)
	BulkUpdate(ctx context.Context, es []Entity) error
	Delete(ctx context.Context, id string, zv EntityModel) (ok bool, err error)
}

// Entity is a persistable resource that has a table name and an ID and that can
// be validated.
type Entity interface {
	GetTableName() string
	GetID() string
	SetID(id string)
	Validate() error
}

// EntityModel is a placeholder for zero-value Entity objects. See the
// PersistenceService for details on its usage.
type EntityModel interface {
	GetTableName() string
}

// JoinEntity is a resource that joins a Node to another Entity.
type JoinEntity interface {
	GetNodeID() string
}

// Filter filters queries in PersistenceService.Filter.
type Filter struct {
	Field string
	Value string
}
