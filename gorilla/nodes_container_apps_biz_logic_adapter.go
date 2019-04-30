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
	"log"

	"github.com/smartedgemec/controller-ce/grpc/node"
)

type nodesContainerAppsBLA struct {
}

func (bla *nodesContainerAppsBLA) create(node node.ClientConn) error {
	log.Println("BLA: deploying container app")
	log.Println(node)
	return nil
}

func (bla *nodesContainerAppsBLA) getByFilter() error {
	log.Println("BLA: getting container apps status")
	return nil
}

func (bla *nodesContainerAppsBLA) getByID() error {
	log.Println("BLA: getting container app status")
	return nil
}

func (bla *nodesContainerAppsBLA) bulkUpdate(nodes []node.ClientConn) error {
	log.Println("BLA: updating container apps")
	log.Println(nodes)
	return nil
}

func (bla *nodesContainerAppsBLA) delete() error {
	log.Println("BLA: undeploying container app")
	return nil
}
