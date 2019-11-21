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
		containerSpec.Args = []string{"lspci -vd:0b30"}
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
