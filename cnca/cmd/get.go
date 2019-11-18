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

package cnca

import (
	"errors"
	"fmt"
	y2j "github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get active LTE CUPS userplane(s) or NGC AF subscription(s)",
	Args:  cobra.MaximumNArgs(2),
	Run:   func(cmd *cobra.Command, args []string) {

		if len(args) < 2 {
			fmt.Println(errors.New("Missing input(s)"))
			return
		}

		if args[0] == "subscription" {

			// get subscription
			sub, err := AFGetSubscription(args[1])
			if err != nil {
				klog.Info(err)
				return
			}

			sub, err = y2j.JSONToYAML(sub)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("Active AF Subscription:\n%s", string(sub))
			return
		} else if args[0] == "all" && args[1] == "subscriptions" {

			// get subscriptions
			sub, err := AFGetSubscription(args[0])
			if err != nil {
				klog.Info(err)
				return
			}

			if string(sub) == "[]" {
				sub = []byte("none")
			}

			sub, err = y2j.JSONToYAML(sub)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("Active AF Subscriptions:\n%s", string(sub))
			return
		} else if args[0] == "userplane" {

			// get userplane
			up, err := LteGetUserplane(args[1])
			if err != nil {
				klog.Info(err)
				return
			}

			up, err = y2j.JSONToYAML(up)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("Active LTE CUPS Userplane:\n%s", string(up))
			return
		} else if args[0] == "all" && args[1] == "userplanes" {

			// get userplanes
			up, err := LteGetUserplane(args[0])
			if err != nil {
				klog.Info(err)
				return
			}

			if string(up) == "[]" {
				up = []byte("none")
			}

			up, err = y2j.JSONToYAML(up)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("Active LTE CUPS Userplanes:\n%s", string(up))
			return
		}

		fmt.Println(errors.New("Invalid input(s)"))
	},
}

func init() {

	const help =
`Get active LTE CUPS userplane(s) or NGC AF subscription(s)

Usage:
  cnca get { all { userplanes | subscriptions } | { userplane <userplane-id> | subscription <subscription-id> }

Example:
  cnca get userplane <subscription-id>
  cnca get subscription <subscription-id>
  cnca get all userplanes
  cnca get all subscriptions

Flags:
  -h, --help   help
`
	// add `get` command
	cncaCmd.AddCommand(getCmd)
	getCmd.SetHelpTemplate(help)
}
