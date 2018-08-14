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
	managedZone := "luontola-fi" // TODO: parameterize

	println(project)
	println(managedZone)

	zones, err := getManagedZones(project, ctx, dnsService)
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(zones)

	rrsets, err := getResourceRecordSets(project, managedZone, ctx, dnsService)
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(rrsets)
}

func getManagedZones(project string, ctx context.Context, dnsService *dns.Service) ([]*dns.ManagedZone, error) {
	var results []*dns.ManagedZone
	if err := dnsService.ManagedZones.List(project).Pages(ctx, func(page *dns.ManagedZonesListResponse) error {
		for _, zone := range page.ManagedZones {
			results = append(results, zone)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return results, nil
}

func getResourceRecordSets(project string, managedZone string, ctx context.Context, dnsService *dns.Service) ([]*dns.ResourceRecordSet, error) {
	var results []*dns.ResourceRecordSet
	req := dnsService.ResourceRecordSets.List(project, managedZone)
	if err := req.Pages(ctx, func(page *dns.ResourceRecordSetsListResponse) error {
		for _, rrset := range page.Rrsets {
			results = append(results, rrset)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return results, nil
}
