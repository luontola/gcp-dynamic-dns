// Copyright © 2019 Esko Luontola
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"app/gcloud"
	"app/ip"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"github.com/davecgh/go-spew/spew"
	"log"
	"os"
	"strings"
	"time"
)


// Credit to Brad Peabody and Dougnukem
// https://stackoverflow.com/a/50825191
var privateIPBlocks []*net.IPNet

func init() {
    for _, cidr := range []string{
	        "127.0.0.0/8",    // IPv4 loopback
	        "10.0.0.0/8",     // RFC1918
	        "172.16.0.0/12",  // RFC1918
	        "192.168.0.0/16", // RFC1918
	        "169.254.0.0/16", // RFC3927 link-local
	        "::1/128",        // IPv6 loopback
	        "fe80::/10",      // IPv6 link-local
	        "fc00::/7",       // IPv6 unique local addr
    } {
       _, block, err := net.ParseCIDR(cidr)
       if err != nil {
            panic(fmt.Errorf("parse error on %q: %v", cidr, err))
       }
       privateIPBlocks = append(privateIPBlocks, block)
    }
}

func isPrivateIP(ip net.IP) bool {
    if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
	        return true
    }

    for _, block := range privateIPBlocks {
        if block.Contains(ip) {
            return true
        }
    }
    return false
}

func main() {
	command := "help"
	if len(os.Args) == 2 {
		command = os.Args[1]
	}
	switch command {
	case "sync":
		sync()
		os.Exit(0)
	case "sync-once":
		syncOnce()
		os.Exit(0)
	case "list-ip":
		listIP()
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
	println("  sync        Update DNS records continuously")
	println("  sync-once   Update DNS records once")
	println("  list-ip     Print current IP address")
	println("  list-dns    Print current DNS records")
	println("  help        Print this help")
}

// commands

func sync() {
	names := paramDnsNames()
	project := paramGoogleProject()
	client := gcloud.Configure(project)

	var previousIP string
	for {
		currentIP := readCurrentIP()
		if currentIP != previousIP {
			handleChangedIP(currentIP, names, client)
			previousIP = currentIP
		}
		time.Sleep(time.Minute)
	}
}

func syncOnce() {
	names := paramDnsNames()
	project := paramGoogleProject()
	client := gcloud.Configure(project)

	currentIP := readCurrentIP()
	handleChangedIP(currentIP, names, client)
}

func handleChangedIP(currentIP string, names []string, client *gcloud.Client) {
	log.Printf("Current IP is %v\n", currentIP)

	log.Printf("Updating DNS records %v\n", names)
	records := readDnsRecords(client, names)
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

func listIP() {
	currentIP := readCurrentIP()

	log.Printf("Current IP is %v\n", currentIP)
}

func listDns() {
	names := paramDnsNames()
	project := paramGoogleProject()

	client := gcloud.Configure(project)
	records := readDnsRecords(client, names)

	spew.Dump(records)
}

// parameters

func paramInterfaceName() string {
	return os.Getenv("INTERFACE_NAME")
}

func paramDnsNames() []string {
	names := os.Getenv("DNS_NAMES")
	if names == "" {
		log.Fatal("Environment variable DNS_NAMES not set.")
	}
	return strings.Split(names, " ")
}

func paramGoogleProject() string {
	project := os.Getenv("GOOGLE_PROJECT")
	if project == "" {
		log.Fatal("Environment variable GOOGLE_PROJECT not set.")
	}
	return project
}

// operations

func readCurrentIP() string {
	var currentIP string
	var err error
	if name := paramInterfaceName(); name != "" {
		currentIP, err = ip.InterfaceIP(name)
	} else {
		currentIP, err = ip.OutgoingIP()
	}
	if err != nil {
		log.Fatal("Failed to read current IP: ", err)
	}

	if isPrivateIP(net.ParseIP(currentIP)) {
		//get IP from external source
		resp, err := http.Get("https://ifconfig.me/ip")
		if err != nil {
			log.Fatal("Failed to issue http.Get request to https://ifconfig.me/ip: ", err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Error reading response from https://ifconfig.me/ip: ", err)
		}
		currentIP = string(body)
		defer resp.Body.Close()
	}

	return currentIP
}

func readDnsRecords(client *gcloud.Client, names []string) gcloud.DnsRecords {
	records, err := client.DnsRecordsByNameAndType(names, "A")
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
