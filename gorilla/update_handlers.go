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

func handleUpdateNodesApps(ctx context.Context, ps cce.PersistenceService, e cce.Validatable) error { // nolint: gocyclo
	nodeCC, err := connectNode(ctx, ps, &e.(*cce.NodeAppReq).NodeApp)
	if err != nil {
		return err
	}

	log.Debug(nodeCC.Node)

	ctrl := getController(ctx)

	switch ctrl.OrchestrationMode {
	case cce.OrchestrationModeNative:
		switch e.(*cce.NodeAppReq).Cmd {
		case "start":
			err = nodeCC.AppLifeSvcCli.Start(ctx, e.(*cce.NodeAppReq).AppID)
		case "stop":
			err = nodeCC.AppLifeSvcCli.Stop(ctx, e.(*cce.NodeAppReq).AppID)
		case "restart":
			err = nodeCC.AppLifeSvcCli.Restart(ctx, e.(*cce.NodeAppReq).AppID)
		}
		if err != nil {
			return err
		}
	case cce.OrchestrationModeKubernetes:
		switch e.(*cce.NodeAppReq).Cmd {
		case "start":
			err = ctrl.KubernetesClient.Start(ctx,
				e.(*cce.NodeAppReq).NodeApp.NodeID, e.(*cce.NodeAppReq).NodeApp.AppID)
		case "stop":
			err = ctrl.KubernetesClient.Stop(ctx,
				e.(*cce.NodeAppReq).NodeApp.NodeID, e.(*cce.NodeAppReq).NodeApp.AppID)
		case "restart":
			err = ctrl.KubernetesClient.Restart(ctx,
				e.(*cce.NodeAppReq).NodeApp.NodeID, e.(*cce.NodeAppReq).NodeApp.AppID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
