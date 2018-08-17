// Copyright © 2018 Esko Luontola
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"app/gcloud"
	"app/kube"
	"fmt"
	"os"
)

func main() {
	command := "help"
	if len(os.Args) == 2 {
		command = os.Args[1]
	}
	switch command {
	case "list-nodes":
		kube.ListNodes()
		os.Exit(0)
	case "list-dns":
		gcloud.ListDns()
		os.Exit(0)
	default:
		fmt.Printf("%v <command>\n", os.Args[0])
		println("Available commands:")
		println("  list-nodes  Print list of Kubernetes cluster nodes")
		println("  list-dns    Print list of DNS records")
		println("  help        Print this help")
		os.Exit(1)
	}

	// TODO: check which DNS records have a different IP
	// TODO: update DNS records with new IPs
}
