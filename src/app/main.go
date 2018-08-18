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
)

func main() {
	command := "help"
	if len(os.Args) == 2 {
		command = os.Args[1]
	}
	switch command {
	case "sync":
		sync()
		os.Exit(0)
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
}

func printHelp() {
	fmt.Printf("%v <command>\n", os.Args[0])
	println("Available commands:")
	println("  sync        Update DNS records with Kubernetes cluster nodes")
	println("  list-nodes  Print list of Kubernetes cluster nodes")
	println("  list-dns    Print list of DNS records")
	println("  help        Print this help")
}

// commands

func sync() {
	names := paramDnsNames()
	project := paramGoogleProject()

	nodeIPs := readNodeIPs()
	log.Printf("Kubernetes node IPs are %v\n", nodeIPs)

	log.Printf("Updating DNS records %v\n", names)
	client := gcloud.Configure(project)
	records := readDnsRecords(client, names)
	updated := updateDnsRecords(client, records, nodeIPs)

	if len(updated) == 0 {
		log.Println("Nothing to update")
	} else {
		log.Printf("Updated %d DNS records:\n", len(updated))
		for _, record := range updated {
			log.Printf("    %s -> %v\n", record.Name, record.Rrdatas)
		}
	}
}

func listNodes() {
	nodeIPs := readNodeIPs()

	log.Printf("Kubernetes node IPs are %v\n", nodeIPs)
}

func listDns() {
	names := paramDnsNames()
	project := paramGoogleProject()

	client := gcloud.Configure(project)
	records := readDnsRecords(client, names)

	spew.Dump(records)
}

// parameters

func paramDnsNames() []string {
	// TODO: parameterize
	names := []string{"k8s-test1.luontola.fi.", "k8s-test2.luontola.fi.", "k8s-test3.luontola.fi."}
	return names
}

func paramGoogleProject() string {
	project := os.Getenv("GOOGLE_PROJECT")
	if project == "" {
		log.Fatal("Environment variable GOOGLE_PROJECT not set.")
	}
	return project
}

// operations

func readNodeIPs() []string {
	nodeIPs, err := kube.NodeExternalIPs()
	if err != nil {
		log.Fatal("Failed to read Kubernetes node IPs: ", err)
	}
	return nodeIPs
}

func readDnsRecords(client *gcloud.Client, names []string) gcloud.DnsRecords {
	records, err := client.DnsRecordsByName(names)
	if err != nil {
		log.Fatal("Failed to read DNS records: ", err)
	}
	return records
}

func updateDnsRecords(client *gcloud.Client, records gcloud.DnsRecords, newValues []string) gcloud.DnsRecords {
	updated, err := client.UpdateDnsRecords(records, newValues)
	if err != nil {
		log.Fatal("Failed to update DNS records: ", err)
	}
	return updated
}
