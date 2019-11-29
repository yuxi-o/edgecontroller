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
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

// unregisterCmd represents the unregister command
var unregisterCmd = &cobra.Command{
	Use:   "unregister",
	Short: "Un-register controller from AF services registry",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Println(errors.New("AF service ID missing"))
			return
		}
		// unregister AF service
		err := OAM5gUnregisterAFService(args[0])
		if err != nil {
			klog.Info(err)
			return
		}
		fmt.Printf("Service ID `%s` unregistered successfully\n", args[0])
	},
}

func init() {

	const help = `Unregister controller from NGC AF services registry

Usage:
  cnca unregister <service-id>

Flags:
  -h, --help       help
`
	// add `register` command
	cncaCmd.AddCommand(unregisterCmd)
	unregisterCmd.SetHelpTemplate(help)
}
