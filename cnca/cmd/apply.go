// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

package cnca

import (
	"errors"
	"fmt"
	y2j "github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/klog"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply LTE CUPS userplane or NGC AF subscription using YAML configuration file",
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		ymlFile, _ := cmd.Flags().GetString("filename")
		if ymlFile == "" {
			fmt.Println(errors.New("YAML file missing"))
			return
		}

		data, err := ioutil.ReadFile(ymlFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		var c Header
		if err = yaml.Unmarshal(data, &c); err != nil {
			fmt.Println(err)
			return
		}

		switch c.Kind {
		case "ngc":
			var s AFTrafficInfluSub
			if err = yaml.Unmarshal(data, &s); err != nil {
				fmt.Println(err)
				return
			}

			sub, err := yaml.Marshal(s.Policy)
			if err != nil {
				fmt.Println(err)
				return
			}

			sub, err = y2j.YAMLToJSON(sub)
			if err != nil {
				fmt.Println(err)
				return
			}

			// create new subscription
			subLoc, err := AFCreateSubscription(sub)
			if err != nil {
				klog.Info(err)
				return
			}
			fmt.Println("Subscription URI:", subLoc)

		case "lte":
			var u LTEUserplane
			if err = yaml.Unmarshal(data, &u); err != nil {
				fmt.Println(err)
				return
			}

			up, err := yaml.Marshal(u.Policy)
			if err != nil {
				fmt.Println(err)
				return
			}

			up, err = y2j.YAMLToJSON(up)
			if err != nil {
				fmt.Println(err)
				return
			}

			// create new LTE userplane
			upID, err := LteCreateUserplane(up)
			if err != nil {
				klog.Info(err)
				return
			}
			fmt.Println("Userplane:", upID)

		default:
			fmt.Println(errors.New("`kind` missing or unknown in YAML file"))
		}
	},
}

func init() {

	const help = `Apply LTE CUPS userplane or NGC AF subscription using YAML configuration file

Usage:
  cnca apply -f <config.yml>

Flags:
  -h, --help       help
  -f, --filename   YAML configuration file
`
	// add `apply` command
	cncaCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringP("filename", "f", "", "YAML configuration file")
	applyCmd.MarkFlagRequired("filename")
	applyCmd.SetHelpTemplate(help)
}
