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
	Use:  "patch",
	Long: "Patch an active CNCA subscription using YAML configuration file",
	Run:  func(cmd *cobra.Command, args []string) {

		var data []byte
		var err error

		if len(args) < 1 {
			fmt.Println(errors.New("Subscription ID missing"))
			return
		}

		if len(args) > 1 {
			fmt.Println("WARNING: Extra args ignored")
		}

		if ymlFile, _ := cmd.Flags().GetString("filename"); ymlFile != "" {
			data, err = ioutil.ReadFile(ymlFile)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		var p TrafficPolicy
		if err = yaml.Unmarshal(data, &p); err != nil {
			fmt.Println(err)
			return
		}

		if val, _ := cmd.Flags().GetString("serviceId"); val != "" {
			p.Policy.AfServiceID = val
		}

		if val, _ := cmd.Flags().GetString("appId"); val != "" {
			p.Policy.AfAppID = val
		}

		if val, _ := cmd.Flags().GetString("transId"); val != "" {
			p.Policy.AfTransID = val
		}

		if val, _ := cmd.Flags().GetString("dnn"); val != "" {
			p.Policy.Dnn = val
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

		// patch subscription
		err = AFPatchSubscription(args[0], sub)
		if err != nil {
			klog.Info(err)
			return
		}

		fmt.Printf("Subscription %s patched\n", args[0])
	},
}

func init() {

	const help =
`Patch an active CNCA subscription using YAML configuration file

Usage:
  kubectl cnca patch <subscriptionID> [FLAGS]

Example:
  kubectl cnca patch <subscriptionID> -f <CNCAConfig.yml>
  kubectl cnca patch <subscriptionID> --appId=<AppID>

Flags:
  -h, --help       help
  -f, --filename   YAML configuration file
  --serviceId      Identifies a service on behalf of which the AF is issuing
                     the request
  --appId          Identifies an application
  --transId        Identifies an NEF Northbound interface transaction, generated
                     by the AF
  --dnn            Identifies data network name
`
	// add `patch` command
	cncaCmd.AddCommand(patchCmd)
	patchCmd.Flags().StringP("filename", "f", "", "YAML configuration file")
	patchCmd.Flags().StringP("serviceId", "", "", "ServiceID")
	patchCmd.Flags().StringP("appId", "", "", "AppID")
	patchCmd.Flags().StringP("transId", "", "", "TransID")
	patchCmd.Flags().StringP("dnn", "", "", "DNN")
	patchCmd.SetHelpTemplate(help)
}
