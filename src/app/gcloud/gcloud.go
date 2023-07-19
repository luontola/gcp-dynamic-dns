// Copyright © 2023 Esko Luontola
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gcloud

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/dns/v1"
	"log"
	"os"
	"reflect"
	"strings"
)

type Client struct {
	project    string
	context    context.Context
	dnsService *dns.Service
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

func (this *Client) DnsRecordsByNameAndType(names []string, recordType string) (DnsRecords, error) {
	records, err := this.DnsRecords()
	if err != nil {
		return nil, err
	}
	found := filterDnsRecordsByName(records, names)
	found = filterDnsRecordsByType(found, recordType)
	if len(found) != len(names) {
		return nil, errors.New(fmt.Sprintf("Expected DNS records <%v> of type <%v>, but only found <%v> of them from the available <%v>",
			strings.Join(names, ", "),
			recordType,
			strings.Join(found.NamesAndTypes(), ", "),
			strings.Join(records.NamesAndTypes(), ", ")))
	}
	return found, nil
}

func filterDnsRecordsByName(records DnsRecords, names []string) DnsRecords {
	if len(names) == 0 {
		return DnsRecords{}
	}
	var results DnsRecords
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

func filterDnsRecordsByType(records DnsRecords, recordType string) DnsRecords {
	var results DnsRecords
	for _, record := range records {
		if record.Type == recordType {
			results = append(results, record)
		}
	}
	return results
}

func (this *Client) DnsRecords() (DnsRecords, error) {
	zones, err := this.ManagedZones()
	if err != nil {
		return nil, err
	}
	var records DnsRecords
	for _, zone := range zones {
		rrsets, err := this.ResourceRecordSets(zone.Name)
		if err != nil {
			return nil, err
		}
		for _, rrset := range rrsets {
			record := &DnsRecord{ManagedZone: zone.Name, ResourceRecordSet: rrset}
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

func (this *Client) UpdateDnsRecords(records DnsRecords, newValues []string) (DnsRecords, error) {
	var updated DnsRecords
	for managedZone, recordsInZone := range records.GroupByZone() {
		changes := changesToUpdateDnsRecordValues(recordsInZone, newValues)
		if changes == nil {
			continue
		}
		resp, err := this.dnsService.Changes.Create(this.project, managedZone, changes).Context(this.context).Do()
		if err != nil {
			return nil, err
		}

		updated = append(updated, ToDnsRecords(managedZone, resp.Additions)...)
	}
	return updated, nil
}

func changesToUpdateDnsRecordValues(records DnsRecords, newValues []string) *dns.Change {
	changes := &dns.Change{}
	for _, record := range records {
		if reflect.DeepEqual(record.Rrdatas, newValues) {
			continue
		}
		changes.Deletions = append(changes.Deletions, record.ResourceRecordSet)
		addition := *record.ResourceRecordSet
		addition.Rrdatas = newValues
		changes.Additions = append(changes.Additions, &addition)
	}
	if len(changes.Additions) == 0 {
		return nil
	}
	return changes
}

type DnsRecord struct {
	ManagedZone string
	OldRrdatas  []string
	*dns.ResourceRecordSet
}

func (record DnsRecord) NameAndType() string {
	return fmt.Sprintf("%v %v", record.Name, record.Type)
}

type DnsRecords []*DnsRecord

func ToDnsRecords(managedZone string, records []*dns.ResourceRecordSet) DnsRecords {
	var results DnsRecords
	for _, record := range records {
		results = append(results, &DnsRecord{ManagedZone: managedZone, ResourceRecordSet: record})
	}
	return results
}

func (records DnsRecords) GroupByZone() map[string]DnsRecords {
	byZone := make(map[string]DnsRecords)
	for _, record := range records {
		byZone[record.ManagedZone] = append(byZone[record.ManagedZone], record)
	}
	return byZone
}

func (records DnsRecords) NamesAndTypes() []string {
	names := make([]string, len(records))
	for i, record := range records {
		names[i] = record.NameAndType()
	}
	return names
}
