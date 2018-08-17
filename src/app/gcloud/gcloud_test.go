// Copyright © 2018 Esko Luontola
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package gcloud

import (
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/api/dns/v1"
	"testing"
)

func TestGCloud(t *testing.T) {
	Convey("FilterDnsRecordsByNameSpec", t, FilterDnsRecordsByNameSpec)
}

func FilterDnsRecordsByNameSpec() {
	records := []*DnsRecord{
		{"example", &dns.ResourceRecordSet{Name: "foo.example.com."}},
		{"example", &dns.ResourceRecordSet{Name: "bar.example.com."}},
		{"example", &dns.ResourceRecordSet{Name: "baz.example.com."}},
	}

	Convey("Empty search", func() {
		records = filterDnsRecordsByName(records)

		So(records, ShouldHaveLength, 0)
	})

	Convey("One match", func() {
		records = filterDnsRecordsByName(records, "foo.example.com.")

		So(records, ShouldHaveLength, 1)
		So(records[0].Name, ShouldEqual, "foo.example.com.")
	})

	Convey("Many matches", func() {
		records = filterDnsRecordsByName(records, "bar.example.com.", "baz.example.com.")

		So(records, ShouldHaveLength, 2)
		So(records[0].Name, ShouldEqual, "bar.example.com.")
		So(records[1].Name, ShouldEqual, "baz.example.com.")
	})
}
