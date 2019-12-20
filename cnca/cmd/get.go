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
	Use:   "get",
	Short: "Get active LTE CUPS userplane(s) or NGC AF subscription(s)",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Println(errors.New("Missing input"))
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
		} else if args[0] == "subscriptions" {

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
		}

		fmt.Println(errors.New("Invalid input(s)"))
	},
}

func init() {

	const help = `Get active LTE CUPS userplane(s) or NGC AF subscription(s)

Usage:
  cnca get { { userplanes | subscriptions } | { userplane <userplane-id> | subscription <subscription-id> } }

Example:
  cnca get userplane <subscription-id>
  cnca get subscription <subscription-id>
  cnca get userplanes
  cnca get subscriptions

Flags:
  -h, --help   help
`
	// add `get` command
	cncaCmd.AddCommand(getCmd)
	getCmd.SetHelpTemplate(help)
}
