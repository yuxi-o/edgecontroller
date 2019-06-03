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

	cce "github.com/smartedgemec/controller-ce"
)

func handleGetNodesApps(ctx context.Context, ps cce.PersistenceService, e cce.Persistable) (cce.RespEntity, error) {
	nodeCC, err := connectNode(ctx, ps, e.(*cce.NodeApp))
	if err != nil {
		return nil, err
	}

	log.Debug(nodeCC.Node)

	ctrl := getController(ctx)

	var status string
	switch ctrl.OrchestrationMode {
	case cce.OrchestrationModeNative:
		s, err := nodeCC.AppLifeSvcCli.GetStatus(ctx, e.(*cce.NodeApp).AppID)
		if err != nil {
			return nil, err
		}
		status = s.String()
	case cce.OrchestrationModeKubernetes:
		s, err := ctrl.KubernetesClient.Status(ctx, e.(*cce.NodeApp).NodeID, e.(*cce.NodeApp).AppID)
		if err != nil {
			return nil, err
		}
		status = s.String()
	}

	return &cce.NodeAppResp{
		NodeApp: *e.(*cce.NodeApp),
		Status:  status,
	}, nil
}
