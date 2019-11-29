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
	"path/filepath"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// default values
var (
	privileged   = true
	backoffLimit = int32(0)
	namespace    = "default"
	jobTimeout   = 300 //seconds
)

// RSUJob struct to hold RSU job specification for K8
var RSUJob = &batchv1.Job{
	TypeMeta: metav1.TypeMeta{
		Kind:       "Job",
		APIVersion: "v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name: "fpga-opae-job",
	},
	Spec: batchv1.JobSpec{
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:    "fpga-opae", // to be edited by command
						Image:   "fpga-opae-pacn3000:1.0",
						Command: []string{"/bin/bash", "-c", "--"},
						Args:    []string{""}, // to be added by command
						SecurityContext: &corev1.SecurityContext{
							Privileged: &privileged,
						},
						ImagePullPolicy: corev1.PullPolicy(corev1.PullNever),
						Env:             []corev1.EnvVar{},
						VolumeMounts:    []corev1.VolumeMount{}, // to be added by command
					},
				},
				RestartPolicy:    corev1.RestartPolicyNever,
				Volumes:          []corev1.Volume{}, // to be added by command
				ImagePullSecrets: []corev1.LocalObjectReference{},
				NodeSelector:     make(map[string]string), // to be added by command
			},
		},
		BackoffLimit: &backoffLimit,
	},
}

func k8LogCmd(pod string) (*exec.Cmd, error) {
	var err error
	var cmd *exec.Cmd

	// #nosec
	cmd = exec.Command("kubectl", "logs", "-f", pod)

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

	err = cmd.Start()
	if err != nil {
		return cmd, err
	}
	return cmd, nil
}

// PrintJobLogs prints logs from k8 pod belonging to the given job
func PrintJobLogs(clientset *kubernetes.Clientset, job *batchv1.Job) (*exec.Cmd, error) {
	var cmd *exec.Cmd
	// get pod of job based on labels
	set := labels.Set(job.Spec.Selector.MatchLabels)
	listOptions := metav1.ListOptions{LabelSelector: set.AsSelector().String()}
	podsClient := clientset.CoreV1().Pods(namespace)
	pods, err := podsClient.List(listOptions)
	if len(pods.Items) < 1 {
		return cmd, errors.New("Failed to retrieve pod")
	}
	// we should have 1 pod for the job
	pod := pods.Items[0]
	// wait for pod to create container
	for {
		k8Pod, _ := podsClient.Get(pod.Name, metav1.GetOptions{})
		if k8Pod.Status.Phase != corev1.PodPending {
			break
		}
		time.Sleep(1 * time.Second)
	}
	// print logs
	cmd, err = k8LogCmd(pod.Name)
	if err != nil {
		return cmd, err
	}
	return cmd, nil
}

// GetK8Clientset returns the clientset for kubernetes
func GetK8Clientset() (*kubernetes.Clientset, error) {
	var clientset *kubernetes.Clientset
	// retrieve .kube/config file
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return clientset, err
	}
	// create the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return clientset, err
	}
	return clientset, nil
}
