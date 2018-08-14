// Copyright © 2018 Esko Luontola
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"github.com/davecgh/go-spew/spew"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/dns/v1"
	"log"
	"os"
)

func main() {
	googleApplicationCredentials := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if googleApplicationCredentials == "" {
		log.Fatal("Environment variable GOOGLE_APPLICATION_CREDENTIALS not set. " +
			"See https://cloud.google.com/docs/authentication/production for instructions.")
	}

	ctx := context.Background()

	c, err := google.DefaultClient(ctx, dns.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	dnsService, err := dns.New(c)
	if err != nil {
		log.Fatal(err)
	}

	project := "www-prod-204113" // TODO: parameterize

	records, err := getDnsRecords(project, ctx, dnsService)
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(records)
}

type DnsRecord struct {
	ManagedZone string
	*dns.ResourceRecordSet
}

func getDnsRecords(project string, ctx context.Context, dnsService *dns.Service) ([]*DnsRecord, error) {
	zones, err := getManagedZones(project, ctx, dnsService)
	if err != nil {
		return nil, err
	}
	var records []*DnsRecord
	for _, zone := range zones {
		rrsets, err := getResourceRecordSets(project, zone.Name, ctx, dnsService)
		if err != nil {
			return nil, err
		}
		for _, rrset := range rrsets {
			record := &DnsRecord{zone.Name, rrset}
			records = append(records, record)
		}

	}
	return records, nil
}

func getManagedZones(project string, ctx context.Context, dnsService *dns.Service) ([]*dns.ManagedZone, error) {
	var results []*dns.ManagedZone
	err := dnsService.ManagedZones.List(project).Pages(ctx, func(page *dns.ManagedZonesListResponse) error {
		for _, zone := range page.ManagedZones {
			results = append(results, zone)
		}
		return nil
	})
	return results, err
}

func getResourceRecordSets(project string, managedZone string, ctx context.Context, dnsService *dns.Service) ([]*dns.ResourceRecordSet, error) {
	var results []*dns.ResourceRecordSet
	req := dnsService.ResourceRecordSets.List(project, managedZone)
	err := req.Pages(ctx, func(page *dns.ResourceRecordSetsListResponse) error {
		for _, rrset := range page.Rrsets {
			results = append(results, rrset)
		}
		return nil
	})
	return results, err
}
