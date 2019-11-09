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
	"github.com/spf13/cobra"
)

// cncaCmd represents the base command when called without any subcommands
var cncaCmd = &cobra.Command{
	Use:  "cnca",
	Long: "Kubernetes CNCA configuration command line",
	SilenceUsage: true,
}

func init() {

	const usage =
`Usage:
  kubectl cnca [command] [flags]

Available Commands:
  apply       Apply new CNCA traffic policy
  get         Get an active CNCA traffic policy
  get all     Get all active CNCA traffic policies
  patch       Patch an active CNCA traffic policy
  delete      Delete an active CNCA traffic policy
  help        Help about any command

Flags:
  -h, --help   help

Use "kubectl [command] --help" for more information about a command.
`

	cncaCmd.SetUsageTemplate(usage)
}  

// Execute CNCA agent
func Execute() error {
	return cncaCmd.Execute()
}
