// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package cnca

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an active LTE CUPS userplane or NGC AF subscription",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 2 {
			fmt.Println(errors.New("Missing input(s)"))
			return
		}

		if args[0] == "subscription" {

			// delete subscription
			err := AFDeleteSubscription(args[1])
			if err != nil {
				klog.Info(err)
				return
			}
			fmt.Printf("AF Subscription %s deleted\n", args[1])
			return
		} else if args[0] == "userplane" {

			// delete userplane
			err := LteDeleteUserplane(args[1])
			if err != nil {
				klog.Info(err)
				return
			}
			fmt.Printf("LTE CUPS userplane %s deleted\n", args[1])
			return
		}

		fmt.Println(errors.New("Invalid input(s)"))
	},
}

func init() {

	const help = `Delete an active LTE CUPS userplane or NGC AF subscription
	
Usage:
  cnca delete { userplane <userplane-id> | subscription <subscription-id> }

 Example:
  cnca delete userplane <userplane-id>
  cnca delete subscription <subscription-id>

Flags:
  -h, --help   help
`

	// add `delete` command
	cncaCmd.AddCommand(deleteCmd)
	deleteCmd.SetHelpTemplate(help)
}
