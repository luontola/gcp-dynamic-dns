// Copyright © 2018 Esko Luontola
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package kube

import (
	"fmt"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func ListNodes() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	nodes, err := clientset.CoreV1().Nodes().List(meta.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("There are %d nodes in the cluster\n", len(nodes.Items))
	for _, node := range nodes.Items {
		for _, addr := range node.Status.Addresses {
			if addr.Type == core.NodeExternalIP {
				println("External IP", addr.Address)
			}
		}
	}
}
