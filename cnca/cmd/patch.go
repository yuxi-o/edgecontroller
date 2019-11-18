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
	"fmt"
	"errors"
	y2j "github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/klog"
)

// patchCmd represents the patch command
var patchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Patch an active LTE CUPS userplane or NGC AF subscription using YAML configuration file",
	Args:  cobra.MaximumNArgs(1),
	Run:   func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Println(errors.New("LTE CUPS userplane or NGC AF subscription ID missing"))
			return
		}

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
	
			// patch subscription
			err = AFPatchSubscription(args[0], sub)
			if err != nil {
				klog.Info(err)
				return
			}
			fmt.Printf("Subscription %s patched\n", args[0])

		case "lte":
			var u LTEUserplane
			if err = yaml.Unmarshal(data, &u); err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(u.Policy)

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
	
			// patch userplane
			err = LtePatchUserplane(args[0], up)
			if err != nil {
				klog.Info(err)
				return
			}
			fmt.Printf("Subscription %s patched\n", args[0])

		default:
			fmt.Println(errors.New("`kind` missing or unknown in YAML file"))
		}
	},
}

func init() {

	const help =
`Patch an active LTE CUPS userplane or NGC AF subscription using YAML configuration file

Usage:
  cnca patch { <userplane-id> | <subscription-id> } -f <config.yml>

Flags:
  -h, --help       help
  -f, --filename   YAML configuration file
`
	// add `patch` command
	cncaCmd.AddCommand(patchCmd)
	patchCmd.Flags().StringP("filename", "f", "", "YAML configuration file")
	applyCmd.MarkFlagRequired("filename")
	patchCmd.SetHelpTemplate(help)
}
