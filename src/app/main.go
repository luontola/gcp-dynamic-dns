// Copyright © 2018 Esko Luontola
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"app/gcloud"
	"app/kube"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"log"
	"os"
	"strings"
)

func main() {
	command := "help"
	if len(os.Args) == 2 {
		command = os.Args[1]
	}
	switch command {
	case "list-nodes":
		listNodes()
		os.Exit(0)
	case "list-dns":
		listDns()
		os.Exit(0)
	default:
		printHelp()
		os.Exit(1)
	}

	// TODO: check which DNS records have a different IP
	// TODO: update DNS records with new IPs
}

func listNodes() {
	ips := kube.NodeExternalIPs()
	println("External IPs", strings.Join(ips, ", "))
}

func printHelp() {
	fmt.Printf("%v <command>\n", os.Args[0])
	println("Available commands:")
	println("  list-nodes  Print list of Kubernetes cluster nodes")
	println("  list-dns    Print list of DNS records")
	println("  help        Print this help")
}

func listDns() {
	// TODO: parameterize
	names := []string{"k8s-test1.luontola.fi.", "k8s-test2.luontola.fi.", "k8s-test3.luontola.fi."}
	project := os.Getenv("GOOGLE_PROJECT")
	if project == "" {
		log.Fatal("Environment variable GOOGLE_PROJECT not set.")
	}

	client := gcloud.Configure(project)
	records, err := client.DnsRecordsByName(names)
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(records)
}
