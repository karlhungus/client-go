/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"fmt"
	"time"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		deployment, err := clientset.ExtensionsV1beta1().Deployments("default").Get("demo", metav1.GetOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There this deployment is, %v\n", deployment)
		fmt.Printf("there are %d replicas", *deployment.Spec.Replicas)

		fmt.Printf("create a secret!")

		secret, err := clientset.CoreV1().Secrets("default").Get("demo-secret", metav1.GetOptions{}) //	("default").List(metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Didn't find secret, creating")
			sec := &apiv1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "demo-secret",
				},
			}

			newSecret, err := clientset.CoreV1().Secrets("default").Create(sec)
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("Created Secret: %s", newSecret)
		} else {
			fmt.Printf("found secrete: %s, destroying", secret)
			err2 := clientset.CoreV1().Secrets("default").Delete("demo-secret", &metav1.DeleteOptions{})
			if err2 != nil {
				fmt.Printf("failed to delete: %s", err2)
			}
		}

		pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		// Examples for error handling:
		// - Use helper functions like e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		_, err = clientset.CoreV1().Pods("default").Get("example-xxxxx", metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod not found\n")
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found pod\n")
		}

		time.Sleep(10 * time.Second)
	}
}
