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
	Use:          "cnca",
	Long:         "Core Newtwork Configuration Agent (CNCA) command line",
	SilenceUsage: true,
}

// Execute CNCA agent
func Execute() error {
	return cncaCmd.Execute()
}
