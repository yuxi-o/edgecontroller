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
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// default values
var (
	privileged   = true
	backoffLimit = int32(0)
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
			ObjectMeta: metav1.ObjectMeta{
				Name: "rsu-pi",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:    "fpga-opae",
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
