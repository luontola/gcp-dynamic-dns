// Copyright © 2018 Esko Luontola
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package main

// BEFORE RUNNING:
// ---------------
// 1. If not already done, enable the Google Cloud DNS API
//    and check the quota for your project at
//    https://console.developers.google.com/apis/api/dns
// 2. This sample uses Application Default Credentials for authentication.
//    If not already done, install the gcloud CLI from
//    https://cloud.google.com/sdk/ and run
//    `gcloud beta auth application-default login`.
//    For more information, see
//    https://developers.google.com/identity/protocols/application-default-credentials
// 3. Install and update the Go dependencies by running `go get -u` in the
//    project directory.

import (
	"fmt"
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

	req := dnsService.ResourceRecordSets.List(project, managedZone)
	if err := req.Pages(ctx, func(page *dns.ResourceRecordSetsListResponse) error {
		for _, resourceRecordSet := range page.Rrsets {
			// TODO: Change code below to process each `resourceRecordSet` resource:
			fmt.Printf("%#v\n", resourceRecordSet)
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}
