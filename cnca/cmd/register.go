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
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

// registerCmd represents the patch command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register controller to AF services registry",
	Args:  cobra.MaximumNArgs(1),
	Run:   func(cmd *cobra.Command, args []string) {

		var s AfService

		if len(args) < 1 {
			fmt.Println(errors.New("Service name missing"))
			return
		}

		s.AfInstance = args[0]
		s.LocationServices = make([]LocationService, 1)

		if val, _ := cmd.Flags().GetString("dnai"); val != "" {
			s.LocationServices[0].DNAI = val
		}

		if val, _ := cmd.Flags().GetString("dnn"); val != "" {
			s.LocationServices[0].DNN = val
		}

		if val, _ := cmd.Flags().GetString("dns"); val != "" {
			s.LocationServices[0].DNS = val
		}

		srv, err := json.Marshal(s)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(srv))

		// register service
		afID, err := OAM5gRegisterAFService(srv)
		if err != nil {
			klog.Info(err)
			return
		}

		fmt.Printf("Service `%s` registered with AF %s\n", args[0], afID)
	},
}

func init() {

	const help =
`Register controller to NGC AF services registry

Usage:
  cnca register <service-name> [FLAGS]

Example:
  cnca register <service-name>
  cnca register <service-name> --dnai=<DNAI>
  cnca register <service-name> --dnai=<DNAI> --dnn=<DNN> --dns=<DNS>

Flags:
  -h, --help       help
  --dnai           Identifies DNAI
  --dnn            Identifies data network name
  --dns            Identifies DNS
`
	// add `register` command
	cncaCmd.AddCommand(registerCmd)
	registerCmd.Flags().StringP("dnai", "", "", "Identifies DNAI")
	registerCmd.Flags().StringP("dnn", "", "", "Identifies data network name")
	registerCmd.Flags().StringP("dns", "", "", "Identifies DNS")
	registerCmd.SetHelpTemplate(help)
}
