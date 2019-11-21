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
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// default values
var (
	RSUJobNameSpace = "default"
	privileged      = true
)

// RSUJob struct to hold RSU job specification for K8
var RSUJob = &batchv1.Job{
	TypeMeta: metav1.TypeMeta{
		Kind:       "Job",
		APIVersion: "v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:   "rsu-pi-pod",
		Labels: make(map[string]string),
	},
	Spec: batchv1.JobSpec{
		// Optional: Parallelism:,
		// Optional: Completions:,
		// Optional: ActiveDeadlineSeconds:,
		// Optional: Selector:,
		// Optional: ManualSelector:,
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Name:   "rsu-pi",
				Labels: make(map[string]string),
			},
			Spec: corev1.PodSpec{
				InitContainers: []corev1.Container{}, // Doesn't seem obligatory(?)...
				Containers: []corev1.Container{
					{
						Name:    "pi",
						Image:   "perl",
						Command: []string{"perl", "-Mbignum=bpi", "-wle", "print bpi(2000)"},
						SecurityContext: &corev1.SecurityContext{
							Privileged: &privileged,
						},
						ImagePullPolicy: corev1.PullPolicy(corev1.PullIfNotPresent),
						Env:             []corev1.EnvVar{},
						VolumeMounts:    []corev1.VolumeMount{},
					},
				},
				RestartPolicy:    corev1.RestartPolicyOnFailure,
				Volumes:          []corev1.Volume{},
				ImagePullSecrets: []corev1.LocalObjectReference{},
			},
		},
	},
}

// programCmd represents the program command
var programCmd = &cobra.Command{
	Use:   "program",
	Short: "Program an FPGA device on a target node with an RTL image",
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		RTLFile, _ := cmd.Flags().GetString("filename")
		if RTLFile == "" {
			fmt.Println(errors.New("RTL image file missing"))
			return
		}

		// get absolute path
		RTLFile, _ = filepath.Abs(RTLFile)
		fmt.Println(RTLFile)

		node, _ := cmd.Flags().GetString("node")
		if node == "" {
			fmt.Println(errors.New("target node missing"))
			return
		}

		dev, _ := cmd.Flags().GetString("device")
		if dev == "" {
			fmt.Println(errors.New("target PCI device missing"))
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

	const help = `Program an FPGA device on a target node with an RTL image

Usage:
  rsu program -f <RTL-img-file> -n <target-node> -d <target-device>

Flags:
  -h, --help       help
  -f, --filename   RTL image file
  -n, --node       where the target FPGA card is plugged in
  -d, --device     PCI ID of the target FPGA card
`
	// add `program` command
	rsuCmd.AddCommand(programCmd)
	programCmd.Flags().StringP("filename", "f", "", "RTL image file")
	programCmd.MarkFlagRequired("filename")
	programCmd.Flags().StringP("node", "n", "", "where the target FPGA card is plugged in")
	programCmd.MarkFlagRequired("node")
	programCmd.Flags().StringP("device", "d", "", "PCI ID of the target FPGA card")
	programCmd.MarkFlagRequired("device")
	programCmd.SetHelpTemplate(help)
}
