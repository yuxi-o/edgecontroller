// Copyright 2019 Intel Corporation. All rights reserved
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

package main

import (
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new Traffic Influence Subscription at 5G AF",
	Long:  `Create new Traffic Influence Subscription at 5G AF.`,
	Run:   func(cmd *cobra.Command, args []string) {

		// TODO to be provided from .yml
		sub := TrafficInfluSub{
			AFServiceID: "myServiceId",
			AFAppID:     "myAppId",
			AFTransID:   "myTransId",
		}
		// create new Traffic Influence Subscription at AF
		err := AFCreateSubscription(client, sub)
		if err != nil {
			return
		}
	},
}

func init() {
	// add `create` command
	cncaCmd.AddCommand(createCmd)
}
