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

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply new CNCA subscription using YAML configuration file",
	Args:   cobra.MaximumNArgs(0),
	Run:   func(cmd *cobra.Command, args []string) {

		ymlFile, _ := cmd.Flags().GetString("filename")
		if ymlFile == "" {
			fmt.Println(errors.New("CNCA yaml file missing"))
			return
		}

		data, err := ioutil.ReadFile(ymlFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		var p TrafficPolicy
		if err = yaml.Unmarshal(data, &p); err != nil {
			fmt.Println(err)
			return
		}

		sub, err := yaml.Marshal(p.Policy)
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
	},
}

func init() {

	const help =
`Apply new CNCA subscription using YAML configuration file

Usage:
  cnca apply -f <CNCAConfig.yml>

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
