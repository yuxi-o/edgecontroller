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

package gorilla

import (
	"context"
	"log"

	cce "github.com/smartedgemec/controller-ce"
)

func handleDeleteNodesApps(
	ctx context.Context,
	ps cce.PersistenceService,
	e cce.Persistable,
) error {
	nodeCC, err := connectNode(ctx, ps, e.(*cce.NodeApp))
	if err != nil {
		return err
	}

	log.Println(nodeCC.Node)

	if err := nodeCC.AppDeploySvcCli.Undeploy(ctx, e.(*cce.NodeApp).AppID); err != nil {
		return err
	}

	return nil
}
