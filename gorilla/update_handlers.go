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
	"net/http"

	"github.com/pkg/errors"
	cce "github.com/smartedgemec/controller-ce"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func handleUpdateNodes(
	ctx context.Context,
	ps cce.PersistenceService,
	e cce.Validatable,
) (statusCode int, err error) {
	ctrl := getController(ctx)
	nodePort := ctrl.ELAPort
	if nodePort == "" {
		nodePort = defaultELAPort
	}
	nodeCC, err := connectNode(ctx, ps, &e.(*cce.NodeReq).Node, nodePort, ctrl.EdgeNodeCreds)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if e.(*cce.NodeReq).NetworkInterfaces != nil {
		if err := nodeCC.IfaceSvcCli.BulkUpdate(ctx, e.(*cce.NodeReq).NetworkInterfaces); err != nil {
			if s, ok := status.FromError(errors.Cause(err)); ok {
				if s.Code() == codes.NotFound {
					return http.StatusNotFound, errors.New(s.Message())
				}
			}
			return http.StatusInternalServerError, err
		}
	}

	for _, nitp := range e.(*cce.NodeReq).TrafficPolicies {
		tp, err := ps.Read(ctx, nitp.TrafficPolicyID, &cce.TrafficPolicy{})
		if err != nil {
			return http.StatusInternalServerError, err
		}
		if tp == nil {
			// If nil, set an empty policy
			tp = &cce.TrafficPolicy{}
		}
		if err := nodeCC.IfacePolicySvcCli.Set(ctx, nitp.NetworkInterfaceID, tp.(*cce.TrafficPolicy)); err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return 0, nil
}

func handleUpdateNodesApps( //nolint: gocyclo
	ctx context.Context,
	ps cce.PersistenceService,
	e cce.Validatable,
) (statusCode int, err error) {
	ctrl := getController(ctx)
	nodePort := ctrl.EVAPort
	if nodePort == "" {
		nodePort = defaultEVAPort
	}
	nodeCC, err := connectNode(ctx, ps, &e.(*cce.NodeAppReq).NodeApp, nodePort, ctrl.EdgeNodeCreds)
	if err != nil {
		return http.StatusInternalServerError, err
	}

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
			return http.StatusInternalServerError, err
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
			return http.StatusInternalServerError, err
		}
	}

	return 0, nil
}
