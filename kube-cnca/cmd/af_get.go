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
	Use:  "get",
	Long: "Get active CNCA subscription(s)",

	Run:  func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println(errors.New("Subscription ID missing"))
			return
		}

		if len(args) > 1 {
			fmt.Println("WARNING: Extra args ignored")
		}

		if args[0] == "all" {
			// get all subscriptions
			klog.Info("GET ALL >> TBA <<")
			return
		}

		// get subscription
		sub, err := AFGetSubscription(args[0])
		if err != nil {
			klog.Info(err)
			return
		}

		sub, err = y2j.JSONToYAML(sub)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Subscription `%s`:\n---\n%s\n", args[0], string(sub))
	},
}

func init() {

	const help =
`Get active CNCA subscription(s)

Usage:
  kubectl cnca get [all] [<subscriptionID>]

Example:
  kubectl cnca get <subscriptionID>
  kubectl cnca get all

Flags:
  -h, --help   help
`

	// add `get` command
	cncaCmd.AddCommand(getCmd)
	getCmd.SetHelpTemplate(help)
}
