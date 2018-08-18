// Copyright © 2018 Esko Luontola
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package kube

import (
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
)

func configure() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func NodeExternalIPs() ([]string, error) {
	client := configure() // TODO: create a client wrapper (as in gcloud package) if we need more kube operations

	nodes, err := client.CoreV1().Nodes().List(meta.ListOptions{})
	if err != nil {
		return nil, err
	}

	var externalIPs []string
	for _, node := range nodes.Items {
		for _, addr := range node.Status.Addresses {
			if addr.Type == core.NodeExternalIP {
				externalIPs = append(externalIPs, addr.Address)
			}
		}
	}
	return externalIPs, nil
}
