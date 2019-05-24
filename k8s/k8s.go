// Copyright 2019 Smart-Edge.com, Inc. All rights reserved.
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

package k8s

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/pkg/errors"
	appsV1 "k8s.io/api/apps/v1"
	autoscalingV1 "k8s.io/api/autoscaling/v1"
	apiV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restClient "k8s.io/client-go/rest"
)

// App is kubernetes app
type App struct {
	ID     string
	Name   string
	Vendor string
	Image  string
	Cores  int
	Memory int // in MB
}

const (
	// labelNodeID this label is used to uniquely identify node
	labelNodeID = "node-uuid"
	// labelApplicationPod this label is used to uniquely identify a pod
	labelApplicationPod = "app-pod-id"
)

// Client abstracts out calls to k8s master API
type Client struct {
	// Host must be a host string, a host:port pair, or a URL to the base of the apiserver.
	// If a URL is given then the (optional) Path of that URL represents a prefix that must
	// be appended to all request URIs used to access the apiserver. This allows a frontend
	// proxy to easily relocate all of the apiserver endpoints.
	Host string

	// APIPath is a sub-path that points to an API root.
	APIPath string

	// Authentication info
	Username string
	CertFile string
	KeyFile  string
	CAFile   string

	// NewClientSet creates a new Kubernetes clientset interface. If it is nil,
	// a REST client with TLS will be used. This field is intended for use
	// mocking an external connection.
	NewClientSet func() (kubernetes.Interface, error)

	connectOnce sync.Once
	clientSet   kubernetes.Interface
	err         error
}

// Ping checks the connection to the Kubernetes server.
func (ks *Client) Ping() error {
	ks.connectOnce.Do(ks.init)
	if ks.err != nil {
		return ks.err
	}
	return ks.clientSet.CoreV1().RESTClient().Get().AbsPath("/").Do().Error()
}

func (ks *Client) init() {
	csCreate := ks.NewClientSet
	if csCreate == nil {
		csCreate = func() (kubernetes.Interface, error) {
			return kubernetes.NewForConfig(
				&restClient.Config{
					Host:     ks.Host,
					APIPath:  ks.APIPath,
					Username: ks.Username,
					TLSClientConfig: restClient.TLSClientConfig{
						Insecure: false,
						CertFile: ks.CertFile,
						KeyFile:  ks.KeyFile,
						CAFile:   ks.CAFile,
					},
				},
			)
		}
	}
	ks.clientSet, ks.err = csCreate()
}

// Deploy creates a kubernetes deployment
func (ks *Client) Deploy(ctx context.Context, serial string, app *App) error {
	ks.connectOnce.Do(ks.init)
	if ks.err != nil {
		return ks.err
	}

	// initial checks
	if err := ks.checkNode(serial); err != nil {
		return errors.Wrap(err, "deploy: node available error")
	}
	// make the deployment to the correct node
	if err := ks.deploy(serial, app); err != nil {
		return errors.Wrap(err, "deploy: deployment error")
	}
	return nil
}

// Undeploy cascade deletes a kubernetes deployment
func (ks *Client) Undeploy(ctx context.Context, serial string, app *App) error {
	ks.connectOnce.Do(ks.init)
	if ks.err != nil {
		return ks.err
	}

	// initial checks
	if err := ks.checkNode(serial); err != nil {
		return errors.Wrap(err, "undeploy: node available error")
	}
	// make the deployment to the correct node
	if err := ks.undeploy(serial, app); err != nil {
		return errors.Wrap(err, "undeploy: un-deployment error")
	}
	return nil
}

// check if a node is available for the deployment
func (ks *Client) checkNode(serial string) error {
	nodeList, err := ks.clientSet.CoreV1().Nodes().List(
		metaV1.ListOptions{
			LabelSelector: fmt.Sprintf("%s=%s", labelNodeID, serial),
		},
	)
	if err != nil {
		return errors.Wrap(err, "get kubernetes node list error")
	}
	if len(nodeList.Items) != 1 {
		return errors.Wrap(err, "no nodes or duplicate nodes detected")
	}
	return nil
}

// create a kubernetes deployment
func (ks *Client) deploy(serial string, app *App) error {
	// create naming convention
	deploymentName := deploymentName(app.ID, serial)
	podName := podName(app.ID, serial)
	appName := appName(app.ID, serial)

	// deployment client
	deploymentsClient := ks.clientSet.AppsV1().
		Deployments(apiV1.NamespaceDefault)
	deployment := &appsV1.Deployment{
		ObjectMeta: metaV1.ObjectMeta{
			Name: deploymentName,
		},
		Spec: appsV1.DeploymentSpec{
			// only creates a deployment. no replicas created,
			// to be consistent with docker native deploy
			Replicas: int32Ptr(0),
			Selector: &metaV1.LabelSelector{
				MatchLabels: map[string]string{
					labelApplicationPod: podName,
				},
			},
			Strategy: appsV1.DeploymentStrategy{
				Type: appsV1.RecreateDeploymentStrategyType,
			},
			Template: apiV1.PodTemplateSpec{
				ObjectMeta: metaV1.ObjectMeta{
					Labels: map[string]string{
						labelApplicationPod: podName,
					},
				},
				Spec: apiV1.PodSpec{
					Containers: []apiV1.Container{
						{
							Resources: apiV1.ResourceRequirements{},
							Name:      appName,
							Image:     app.Image,
							// never pull the image, use the local image provided
							ImagePullPolicy: apiV1.PullNever,
						},
					},
					NodeSelector: map[string]string{
						labelNodeID: serial,
					},
				},
			},
		},
	}
	_, err := deploymentsClient.Create(deployment)
	if err != nil {
		return errors.Wrap(err, "create kubernetes deployment error")
	}
	return nil
}

// delete a kubernetes deployment
func (ks *Client) undeploy(serial string, app *App) error {
	// create naming convention
	deploymentName := deploymentName(app.ID, serial)
	// deployment client
	deploymentsClient := ks.clientSet.AppsV1().Deployments(apiV1.NamespaceDefault)
	dpp := metaV1.DeletePropagationForeground
	err := deploymentsClient.Delete(deploymentName, &metaV1.DeleteOptions{
		PropagationPolicy: &dpp,
	})
	if err != nil {
		return errors.Wrap(err, "create kubernetes deployment error")
	}
	return nil
}

func int32Ptr(i int32) *int32 { return &i }

// Start scales up the number of replicas of kubernetes deployment to 1.
func (ks *Client) Start(ctx context.Context, serial, id string) error {
	if _, err := ks.checkDeployment(serial, id); err != nil {
		return errors.Wrap(err, "start: correct amount of deployments not available")
	}
	// create naming convention
	deploymentName := deploymentName(id, serial)
	deploymentsClient := ks.clientSet.AppsV1().
		Deployments(apiV1.NamespaceDefault)

	if _, err := deploymentsClient.UpdateScale(
		deploymentName,
		&autoscalingV1.Scale{
			ObjectMeta: metaV1.ObjectMeta{
				Name:      deploymentName,
				Namespace: apiV1.NamespaceDefault,
			},
			Spec: autoscalingV1.ScaleSpec{Replicas: 1},
		}); err != nil {
		return errors.Wrap(err, "start: update deployment replicas to 1 error")
	}
	return nil
}

// Stop scales down the number of replicas of kubernetes deployment to 0.
func (ks *Client) Stop(ctx context.Context, serial, id string) error {
	if _, err := ks.checkDeployment(serial, id); err != nil {
		return errors.Wrap(err, "stop: checking for kubernetes deployment error")
	}
	deploymentName := deploymentName(id, serial)
	deploymentsClient := ks.clientSet.AppsV1().Deployments(apiV1.NamespaceDefault)
	if _, err := deploymentsClient.UpdateScale(
		deploymentName,
		&autoscalingV1.Scale{
			ObjectMeta: metaV1.ObjectMeta{
				Name:      deploymentName,
				Namespace: apiV1.NamespaceDefault,
			},
			Spec: autoscalingV1.ScaleSpec{Replicas: 0},
		}); err != nil {
		return errors.Wrap(err, "stop: scale down deployments to 0 replicas error")
	}
	return nil
}

// Restart scales down the number of replicas of kubernetes deployment to 0 and then scale up to 1.
func (ks *Client) Restart(ctx context.Context, serial, id string) error {
	if _, err := ks.checkDeployment(serial, id); err != nil {
		return errors.Wrap(err, "restart: correct amount of deployments not available")
	}
	// create naming convention
	deploymentName := deploymentName(id, serial)
	deploymentsClient := ks.clientSet.AppsV1().
		Deployments(apiV1.NamespaceDefault)
	if _, err := deploymentsClient.UpdateScale(
		deploymentName,
		&autoscalingV1.Scale{
			ObjectMeta: metaV1.ObjectMeta{
				Name:      deploymentName,
				Namespace: apiV1.NamespaceDefault,
			},
			Spec: autoscalingV1.ScaleSpec{Replicas: 0},
		},
	); err != nil {
		return errors.Wrap(err,
			"restart: scale down deployments to 0 replicas error")
	}
	if _, err := deploymentsClient.UpdateScale(
		deploymentName,
		&autoscalingV1.Scale{
			ObjectMeta: metaV1.ObjectMeta{
				Name:      deploymentName,
				Namespace: apiV1.NamespaceDefault,
			},
			Spec: autoscalingV1.ScaleSpec{Replicas: 1},
		},
	); err != nil {
		return errors.Wrap(err, "restart: scale up deployments to 1 replicas error")
	}
	return nil
}

// check if a kubernetes deployment is available
func (ks *Client) checkDeployment(serial, id string) (*appsV1.Deployment, error) {
	deploymentName := deploymentName(id, serial)
	deployment, err := ks.clientSet.AppsV1().Deployments(apiV1.NamespaceDefault).
		Get(deploymentName, metaV1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "get kubernetes deployment error")
	}
	if deployment == nil {
		return nil, errors.Wrap(err, "kubernetes deployment not found error")
	}
	return deployment, nil
}

// Status gets the status of kubernetes deployment
func (ks *Client) Status(ctx context.Context, serial, id string) (LifecycleStatus, error) {
	deploymentName := deploymentName(id, serial)
	deploymentsClient := ks.clientSet.AppsV1().Deployments(apiV1.NamespaceDefault)
	deployment, err := deploymentsClient.Get(deploymentName, metaV1.GetOptions{})
	if err != nil {
		return Unknown, errors.Wrap(err, "getting deployment error")
	}

	deployConditions := deployment.Status.Conditions

	filter := func() []appsV1.DeploymentCondition {
		var filtered []appsV1.DeploymentCondition
		for _, condition := range deployConditions {
			if condition.Status == apiV1.ConditionTrue {
				filtered = append(filtered, condition)
			}
		}
		return filtered
	}

	deployConditions = filter()

	if len(deployConditions) == 0 {
		return Stopped, nil
	}

	// incoming condition list not sorted according to last updated timestamp
	sort.Slice(deployConditions, func(i, j int) bool {
		return deployConditions[i].LastUpdateTime.Time.After(
			deployConditions[j].LastUpdateTime.Time,
		)
	})

	switch deployConditions[0].Type {
	case appsV1.DeploymentAvailable:
		return Deployed, nil
	case appsV1.DeploymentProgressing:
		return Deploying, nil
	case appsV1.DeploymentReplicaFailure:
		return Error, nil
	}
	return Unknown, nil
}

// retrieve unique name for k8s deployment
func deploymentName(appID, serial string) string {
	return fmt.Sprintf("deploy-%s-%s", appID, serial)
}

// retrieve unique name for the k8s pod
func podName(appID, serial string) string {
	return fmt.Sprintf("pod-%s-%s", appID, serial)
}

// retrieve unique name for the k8s app
func appName(appID, serial string) string {
	return fmt.Sprintf("app-%s-%s", appID, serial)
}
