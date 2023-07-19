// Copyright © 2023 Esko Luontola
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"app/config"
	"app/gcloud"
	"app/ip"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	command := "help"
	if len(os.Args) == 2 {
		command = os.Args[1]
	}
	conf := config.FromEnv()
	switch command {
	case "sync":
		sync(conf)
	case "sync-once":
		syncOnce(conf)
	case "list-ip":
		listIP(conf)
	case "list-dns":
		listDns(conf)
	default:
		printHelp()
		os.Exit(1)
	}
	os.Exit(0)
}

func printHelp() {
	fmt.Printf("%v <command>\n", os.Args[0])
	println("Available commands:")
	println("  sync        Update DNS records continuously")
	println("  sync-once   Update DNS records once")
	println("  list-ip     Print current IP address")
	println("  list-dns    Print current DNS records")
	println("  help        Print this help")
}

// commands

func sync(conf *config.Config) {
	client := gcloud.Configure(conf.GoogleProject)

	var previousIP string
	for {
		currentIP, err := readCurrentIP(conf)

		if err != nil {
			log.Println("WARN: Failed to read the current IP:", err)
		} else if currentIP != previousIP {
			handleChangedIP(currentIP, conf.DnsNames, client)
			previousIP = currentIP
		} else {
			log.Printf("Current IP is %v\n", currentIP)
		}
		if conf.Mode == "service" {
			time.Sleep(time.Minute * 5)
		} else {
			time.Sleep(time.Minute)
		}
	}
}

func syncOnce(conf *config.Config) {
	client := gcloud.Configure(conf.GoogleProject)

	currentIP, err := readCurrentIP(conf)
	if err != nil {
		log.Fatal("Failed to read the current IP: ", err)
	}
	handleChangedIP(currentIP, conf.DnsNames, client)
}

func handleChangedIP(currentIP string, dnsNames []string, client *gcloud.Client) {
	log.Printf("Updating IP %v to DNS records %v\n", currentIP, dnsNames)
	records := readDnsRecords(client, dnsNames)
	updated := updateDnsRecords(client, records, []string{currentIP})

	if len(updated) == 0 {
		log.Println("Nothing to update")
	} else {
		log.Printf("Updated %d DNS records:\n", len(updated))
		for _, record := range updated {
			log.Printf("    %s -> %v\n", record.Name, record.Rrdatas)
		}
	}
}

func listIP(conf *config.Config) {
	currentIP, err := readCurrentIP(conf)
	if err != nil {
		log.Fatal("Failed to read the current IP: ", err)
	}
	log.Printf("Current IP is %v\n", currentIP)
}

func listDns(conf *config.Config) {
	client := gcloud.Configure(conf.GoogleProject)
	records := readDnsRecords(client, conf.DnsNames)
	for _, record := range records {
		println(record.Name, record.Type, record.Ttl, " ", strings.Join(record.Rrdatas, " "))
	}
}

// operations

func readCurrentIP(conf *config.Config) (string, error) {
	var currentIP string
	var err error
	mode := conf.Mode
	switch mode {
	case "service":
		url := conf.NextServiceUrl()
		currentIP, err = ip.ExternalServiceIP(url)
		if err != nil {
			err = fmt.Errorf("failure using external service %v: %w", url, err)
		}
	case "interface":
		if name := conf.InterfaceName; name != "" {
			currentIP, err = ip.InterfaceIP(name)
		} else {
			currentIP, err = ip.OutgoingIP()
		}
	default:
		log.Fatal("Invalid MODE:", mode)
	}
	return currentIP, err
}

func readDnsRecords(client *gcloud.Client, dnsNames []string) gcloud.DnsRecords {
	records, err := client.DnsRecordsByNameAndType(dnsNames, "A")
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
