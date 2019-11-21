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
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get FPGA telemetry",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var jobArgs string

		node, _ := cmd.Flags().GetString("node")
		if node == "" {
			fmt.Println(errors.New("target node missing"))
			return
		}

		switch args[0] {
		case "power":
			jobArgs = "./check_if_modules_loaded.sh && fpgainfo power"

		case "temp":
			jobArgs = "./check_if_modules_loaded.sh && fpgainfo temp"

		case "fme":
			jobArgs = "./check_if_modules_loaded.sh && fpgainfo fme"

		case "port":
			fmt.Println(errors.New("Not supported"))
			return

		case "bmc":
			fmt.Println(errors.New("Not supported"))
			return

		case "phy":
			fmt.Println(errors.New("Not supported"))
			return

		case "mac":
			fmt.Println(errors.New("Not supported"))
			return

		default:
			fmt.Println(errors.New("Undefined or missing metric"))
			return
		}

		// retrieve .kube/config file
		kubeconfig := filepath.Join(
			os.Getenv("HOME"), ".kube", "config",
		)

		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// create the clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// edit K8 job with `program` command specifics
		podSpec := &(RSUJob.Spec.Template.Spec)
		containerSpec := &(RSUJob.Spec.Template.Spec.Containers[0])
		containerSpec.Args = []string{jobArgs}
		containerSpec.VolumeMounts = []corev1.VolumeMount{
			{
				Name:      "class",
				MountPath: "/sys/devices",
				ReadOnly:  false,
			},
		}
		podSpec.NodeSelector["kubernetes.io/hostname"] = node
		podSpec.Volumes = []corev1.Volume{
			{
				Name: "class",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: "/sys/devices",
					},
				},
			},
		}

		jobsClient := clientset.BatchV1().Jobs(RSUJobNameSpace)
		k8Job, err := jobsClient.Create(RSUJob)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("Successfully created job:", k8Job.Name)
	},
}

func init() {

	const help = `Get FPGA telemetry

Usage:
  rsu get <metric> -n <target-node>

Metrics:
  power            print power metrics
  temp             print thermal metrics
  fme              print FME information
  port             print accelerator port information
  bmc              print all Board Management Controller sensor values
  phy              print all PHY information
  mac              print MAC information

Flags:
  -h, --help       help
  -n, --node       where the target FPGA card(s) is/are plugged in
`
	// add `get` command
	rsuCmd.AddCommand(getCmd)
	getCmd.Flags().StringP("node", "n", "", "where the target FPGA card is plugged in")
	getCmd.MarkFlagRequired("node")
	getCmd.SetHelpTemplate(help)
}
