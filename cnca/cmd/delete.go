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
	Use: "delete",
	Short: "Delete an active LTE CUPS userplane or NGC AF TI subscription or " +
		"NGC AF PFD Transaction or NGC AF PFD Application",
	Args: cobra.MaximumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 2 {
			fmt.Println(errors.New("Missing input(s)"))
			return
		}

		if args[0] == "subscription" {

			if pfdCommandCalled == true {
				fmt.Println(errors.New("Invalid input(s)"))
				return
			}

			// delete subscription
			err := AFDeleteSubscription(args[1])
			if err != nil {
				klog.Info(err)
				return
			}
			fmt.Printf("AF Subscription %s deleted\n", args[1])
			return
		} else if args[0] == "userplane" {

			if pfdCommandCalled == true {
				fmt.Println(errors.New("Invalid input(s)"))
				return
			}

			// delete userplane
			err := LteDeleteUserplane(args[1])
			if err != nil {
				klog.Info(err)
				return
			}
			fmt.Printf("LTE CUPS userplane %s deleted\n", args[1])
			return
		} else if args[0] == "transaction" && args[1] != "" {

			if pfdCommandCalled == false {
				fmt.Println(errors.New("Invalid input(s)"))
				return
			}

			if len(args) > 2 {
				if args[2] == "application" && len(args) > 3 {
					// delete PFD application
					err := AFDeletePfdApplication(args[1], args[3])
					if err != nil {
						klog.Info(err)
						return
					}
					fmt.Printf("AF PFD Application %s deleted\n", args[3])
					return
				}
			} else {
				// delete PFD transaction
				err := AFDeletePfdTransaction(args[1])
				if err != nil {
					klog.Info(err)
					return
				}
				fmt.Printf("AF PFD Transaction %s deleted\n", args[1])
				return
			}
		}

		fmt.Println(errors.New("Invalid input(s)"))
	},
}

func init() {

	const help = `Delete an active LTE CUPS userplane or NGC AF TI subscription or NGC AF PFD Transaction or NGC AF PFD Application
	
Usage:
  cnca {pfd | <none>} delete { userplane <userplane-id> | subscription <subscription-id> | transaction <transaction-id> {application <application-id> | <none>}  }

 Example:
  cnca delete userplane <userplane-id>
  cnca delete subscription <subscription-id>
  cnca pfd delete transaction <transaction-id>
  cnca pfd delete transaction <transaction-id> application <application-id> 

Flags:
  -h, --help   help
`

	// add `delete` command
	cncaCmd.AddCommand(deleteCmd)
	pfdCmd.AddCommand(deleteCmd)
	deleteCmd.SetHelpTemplate(help)
}
