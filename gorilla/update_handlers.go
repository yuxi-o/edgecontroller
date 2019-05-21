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

func handleUpdateNodesApps(
	ctx context.Context,
	ps cce.PersistenceService,
	e cce.Validatable,
) error {
	nodeCC, err := connectNode(ctx, ps, &e.(*cce.NodeAppReq).NodeApp)
	if err != nil {
		return err
	}

	log.Println(nodeCC.Node)

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

	return nil
}

func handleUpdateNodesVNFs(
	ctx context.Context,
	ps cce.PersistenceService,
	e cce.Validatable,
) error {
	nodeCC, err := connectNode(ctx, ps, &e.(*cce.NodeVNFReq).NodeVNF)
	if err != nil {
		return err
	}

	log.Println(nodeCC.Node)

	switch e.(*cce.NodeVNFReq).Cmd {
	case "start":
		err = nodeCC.VNFLifeSvcCli.Start(ctx, e.(*cce.NodeVNFReq).VNFID)
	case "stop":
		err = nodeCC.VNFLifeSvcCli.Stop(ctx, e.(*cce.NodeVNFReq).VNFID)
	case "restart":
		err = nodeCC.VNFLifeSvcCli.Restart(ctx, e.(*cce.NodeVNFReq).VNFID)
	}
	if err != nil {
		return err
	}

	return nil
}
