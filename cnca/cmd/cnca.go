// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2019 Intel Corporation

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
