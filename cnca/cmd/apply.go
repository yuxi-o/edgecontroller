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
			subID, err := AFCreateSubscription(sub)
			if err != nil {
				klog.Info(err)
				return
			}
			fmt.Println("Subscription created:", subID)

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
			fmt.Println("Userplane created:", upID)

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
