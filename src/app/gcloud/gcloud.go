// Copyright © 2018 Esko Luontola
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gcloud

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/dns/v1"
	"log"
	"os"
)

type Client struct {
	project    string
	context    context.Context
	dnsService *dns.Service
}

type DnsRecord struct {
	ManagedZone string
	*dns.ResourceRecordSet
}

func Configure(project string) *Client {
	googleApplicationCredentials := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if googleApplicationCredentials == "" {
		log.Fatal("Environment variable GOOGLE_APPLICATION_CREDENTIALS not set. " +
			"See https://cloud.google.com/docs/authentication/production for instructions.")
	}

	ctx := context.Background()
	client, err := google.DefaultClient(ctx, dns.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	dnsService, err := dns.New(client)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		project:    project,
		context:    ctx,
		dnsService: dnsService,
	}
}

func (this *Client) DnsRecordsByName(names []string) ([]*DnsRecord, error) {
	records, err := this.DnsRecords()
	if err != nil {
		return nil, err
	}
	return filterDnsRecordsByName(records, names), nil
}

func filterDnsRecordsByName(records []*DnsRecord, names []string) []*DnsRecord {
	if len(names) == 0 {
		return []*DnsRecord{}
	}
	var results []*DnsRecord
	for _, record := range records {
		if equalsAny(record.Name, names) {
			results = append(results, record)
		}
	}
	return results
}

func equalsAny(haystack string, needles []string) bool {
	for _, needle := range needles {
		if haystack == needle {
			return true
		}
	}
	return false
}

func (this *Client) DnsRecords() ([]*DnsRecord, error) {
	zones, err := this.ManagedZones()
	if err != nil {
		return nil, err
	}
	var records []*DnsRecord
	for _, zone := range zones {
		rrsets, err := this.ResourceRecordSets(zone.Name)
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

func (this *Client) ManagedZones() ([]*dns.ManagedZone, error) {
	var results []*dns.ManagedZone
	err := this.dnsService.ManagedZones.List(this.project).Pages(this.context, func(page *dns.ManagedZonesListResponse) error {
		for _, zone := range page.ManagedZones {
			results = append(results, zone)
		}
		return nil
	})
	return results, err
}

func (this *Client) ResourceRecordSets(managedZone string) ([]*dns.ResourceRecordSet, error) {
	var results []*dns.ResourceRecordSet
	req := this.dnsService.ResourceRecordSets.List(this.project, managedZone)
	err := req.Pages(this.context, func(page *dns.ResourceRecordSetsListResponse) error {
		for _, rrset := range page.Rrsets {
			results = append(results, rrset)
		}
		return nil
	})
	return results, err
}
