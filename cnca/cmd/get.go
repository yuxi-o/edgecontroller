// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

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
	Use: "get",
	Short: "Get active LTE CUPS userplane(s) or NGC AF TI subscription(s) " +
		"or NGC AF PFD Transaction(s) or NGC AF PFD Application(s)",
	Args: cobra.MaximumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Println(errors.New("Missing input"))
			return
		}

		if args[0] == "subscription" {

			if pfdCommandCalled == true {
				fmt.Println(errors.New("Invalid input(s)"))
				return
			}

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
		} else if args[0] == "subscriptions" {

			if pfdCommandCalled == true {
				fmt.Println(errors.New("Invalid input(s)"))
				return
			}

			// get subscriptions
			sub, err := AFGetSubscription("all")
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

			if pfdCommandCalled == true {
				fmt.Println(errors.New("Invalid input(s)"))
				return
			}

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
		} else if args[0] == "userplanes" {

			if pfdCommandCalled == true {
				fmt.Println(errors.New("Invalid input(s)"))
				return
			}

			// get userplanes
			up, err := LteGetUserplane("all")
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
		} else if args[0] == "transactions" || args[0] == "transaction" {
			var transId string
			var appId string
			var pfdData []byte
			var err error

			if pfdCommandCalled == false {
				fmt.Println(errors.New("Invalid input(s)"))
				return
			}

			if args[0] == "transaction" && len(args) > 1 {
				transId = args[1]
				if len(args) > 2 && args[2] != "" {
					if args[2] == "application" && len(args) > 3 {
						appId = args[3]
					} else {
						fmt.Println(errors.New("Invalid input(s)"))
						return
					}
				}
			} else if args[0] == "transactions" {
				transId = "all"
			} else {
				fmt.Println(errors.New("Invalid input(s)"))
				return
			}

			if appId != "" {
				// get PFD application
				pfdData, err = AFGetPfdApplication(transId, appId)
				if err != nil {
					klog.Info(err)
					return
				}
			} else {
				// get PFD transaction
				pfdData, err = AFGetPfdTransaction(transId)
				if err != nil {
					klog.Info(err)
					return
				}
			}

			if args[0] == "transactions" && string(pfdData) == "[]" {
				pfdData = []byte("none")
			}

			pfdData, err = y2j.JSONToYAML(pfdData)
			if err != nil {
				fmt.Println(err)
				return
			}

			if appId != "" {
				fmt.Printf("PFD Application: %s\n%s", appId, string(pfdData))
			} else {
				fmt.Printf("PFD Transaction: %s\n%s", transId, string(pfdData))
			}
			return
		}

		fmt.Println(errors.New("Invalid input(s)"))
	},
}

func init() {

	const help = `Get active LTE CUPS userplane(s) or NGC AF TI subscription(s) or NGC AF PFD Transaction(s) or NGC AF PFD Application(s)

Usage:
  cnca {pfd | <none>} get { { userplanes | subscriptions | transactions} | { userplane <userplane-id> | subscription <subscription-id> | transaction <transaction-id> {application <application-id> | <none>}transaction <transaction-id> {application <application-id> | <none>}} }

Example:
  cnca get userplane <userplane-id>
  cnca get subscription <subscription-id>
  cnca get userplanes
  cnca get subscriptions
  cnca pfd get transactions
  cnca pfd get transaction <transaction-id>
  cnca pfd get transaction <transaction-id> application <application-id>

Flags:
  -h, --help   help
`
	// add `get` command
	cncaCmd.AddCommand(getCmd)
	pfdCmd.AddCommand(getCmd)
	getCmd.SetHelpTemplate(help)
}
