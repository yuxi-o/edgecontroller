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
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:  "delete",
	Long: "Delete an active CNCA subscription",

	Run:  func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Println(errors.New("Subscription ID missing"))
			return
		}

		if len(args) > 1 {
			fmt.Println("WARNING: Extra args ignored")
		}

		// delete subscription
		err := AFDeleteSubscription(args[0])
		if err != nil {
			klog.Info(err)
			return
		}
		fmt.Printf("Subscription %s deleted\n", args[0])
	},
}

func init() {

	const help =
`Delete an active CNCA subscription
	
Usage:
  kubectl cnca delete <subscriptionID>

Flags:
  -h, --help   help
`
	
	// add `delete` command
	cncaCmd.AddCommand(deleteCmd)
	deleteCmd.SetHelpTemplate(help)
}
