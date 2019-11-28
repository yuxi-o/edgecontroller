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

package rsu

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func listImages(node string) error {
	var err error
	var cmd *exec.Cmd

	// #nosec
	cmd = exec.Command("ssh", "root@"+node,
		"ls -lh", "/temp/vran_images/", "| awk '{print $6,$7,\"\t\",$5,\"\t\",$9}'")

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	go func() {
		if _, err = io.Copy(os.Stdout, stdout); err != nil {
			fmt.Println(err.Error())
		}
	}()
	go func() {
		if _, err = io.Copy(os.Stderr, stderr); err != nil {
			fmt.Println(err.Error())
		}
	}()

	fmt.Printf("\nAvailable RTL images:\n---------------------")
	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}
	fmt.Printf("\n")
	return nil
}

func listDevices(node string) error {
	var err error
	var cmd *exec.Cmd

	// #nosec
	cmd = exec.Command("ssh", node, "lspci", "-knnd:0b30")

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	go func() {
		if _, err = io.Copy(os.Stdout, stdout); err != nil {
			fmt.Println(err.Error())
		}
	}()
	go func() {
		if _, err = io.Copy(os.Stderr, stderr); err != nil {
			fmt.Println(err.Error())
		}
	}()

	fmt.Printf("FPGA devices installed:\n-----------------------\n")
	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}
	fmt.Printf("\n")
	return nil
}

// discoverCmd represents the discover command
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover FPGA card(s) on a node",
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		node, _ := cmd.Flags().GetString("node")
		if node == "" {
			fmt.Println(errors.New("target node missing"))
			return
		}

		err := listImages(node)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = listDevices(node)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	},
}

func init() {

	const help = `Discover FPGA card(s) on a node

Usage:
  rsu discover -n <target-node>

Flags:
  -h, --help       help
  -n, --node       where the FPGA card(s) to be discovered
`
	// add `discover` command
	rsuCmd.AddCommand(discoverCmd)
	discoverCmd.Flags().StringP("node", "n", "", "where the target FPGA card is plugged in")
	discoverCmd.MarkFlagRequired("node")
	discoverCmd.SetHelpTemplate(help)
}
