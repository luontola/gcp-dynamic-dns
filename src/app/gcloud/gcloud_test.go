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
	Convey("GroupDnsRecordsByZoneSpec", t, GroupDnsRecordsByZoneSpec)
	Convey("UpdateDnsRecordValuesSpec", t, UpdateDnsRecordValuesSpec)
}

func FilterDnsRecordsByNameSpec() {
	records := []*DnsRecord{
		{"example", &dns.ResourceRecordSet{Name: "foo.example.com."}},
		{"example", &dns.ResourceRecordSet{Name: "bar.example.com."}},
		{"example", &dns.ResourceRecordSet{Name: "baz.example.com."}},
	}

	Convey("Empty search", func() {
		records = filterDnsRecordsByName(records, []string{})

		So(records, ShouldHaveLength, 0)
	})

	Convey("One match", func() {
		records = filterDnsRecordsByName(records, []string{"foo.example.com."})

		So(records, ShouldHaveLength, 1)
		So(records[0].Name, ShouldEqual, "foo.example.com.")
	})

	Convey("Many matches", func() {
		records = filterDnsRecordsByName(records, []string{"bar.example.com.", "baz.example.com."})

		So(records, ShouldHaveLength, 2)
		So(records[0].Name, ShouldEqual, "bar.example.com.")
		So(records[1].Name, ShouldEqual, "baz.example.com.")
	})
}

func GroupDnsRecordsByZoneSpec() {
	records := DnsRecords{
		{"zone1", &dns.ResourceRecordSet{Name: "zone1.com."}},
		{"zone2", &dns.ResourceRecordSet{Name: "zone2.com."}},
		{"zone2", &dns.ResourceRecordSet{Name: "www.zone2.com."}},
	}

	So(records.GroupByZone(), ShouldResemble, map[string]DnsRecords{
		"zone1": {
			{"zone1", &dns.ResourceRecordSet{Name: "zone1.com."}},
		},
		"zone2": {
			{"zone2", &dns.ResourceRecordSet{Name: "zone2.com."}},
			{"zone2", &dns.ResourceRecordSet{Name: "www.zone2.com."}},
		},
	})
}

func UpdateDnsRecordValuesSpec() {
	Convey("one record, one value", func() {
		records := DnsRecords{
			{"zone1", &dns.ResourceRecordSet{Name: "zone1.com.", Rrdatas: []string{"1.1.1.1"}}},
		}

		changes := changesToUpdateDnsRecordValues(records, []string{"2.2.2.2"})

		So(changes, ShouldResemble, &dns.Change{
			Deletions: []*dns.ResourceRecordSet{
				{Name: "zone1.com.", Rrdatas: []string{"1.1.1.1"}},
			},
			Additions: []*dns.ResourceRecordSet{
				{Name: "zone1.com.", Rrdatas: []string{"2.2.2.2"}},
			},
		})
	})

	Convey("multiple records", func() {
		records := DnsRecords{
			{"zone1", &dns.ResourceRecordSet{Name: "zone1.com.", Rrdatas: []string{"1.1.1.1"}}},
			{"zone1", &dns.ResourceRecordSet{Name: "www.zone1.com.", Rrdatas: []string{"1.1.1.1"}}},
		}

		changes := changesToUpdateDnsRecordValues(records, []string{"2.2.2.2"})

		So(changes, ShouldResemble, &dns.Change{
			Deletions: []*dns.ResourceRecordSet{
				{Name: "zone1.com.", Rrdatas: []string{"1.1.1.1"}},
				{Name: "www.zone1.com.", Rrdatas: []string{"1.1.1.1"}},
			},
			Additions: []*dns.ResourceRecordSet{
				{Name: "zone1.com.", Rrdatas: []string{"2.2.2.2"}},
				{Name: "www.zone1.com.", Rrdatas: []string{"2.2.2.2"}},
			},
		})
	})

	Convey("multiple values", func() {
		records := DnsRecords{
			{"zone1", &dns.ResourceRecordSet{Name: "zone1.com.", Rrdatas: []string{"1.1.1.1", "2.2.2.2"}}},
		}

		changes := changesToUpdateDnsRecordValues(records, []string{"3.3.3.3", "4.4.4.4"})

		So(changes, ShouldResemble, &dns.Change{
			Deletions: []*dns.ResourceRecordSet{
				{Name: "zone1.com.", Rrdatas: []string{"1.1.1.1", "2.2.2.2"}},
			},
			Additions: []*dns.ResourceRecordSet{
				{Name: "zone1.com.", Rrdatas: []string{"3.3.3.3", "4.4.4.4"}},
			},
		})
	})

	Convey("some records up to date", func() {
		records := DnsRecords{
			{"zone1", &dns.ResourceRecordSet{Name: "zone1.com.", Rrdatas: []string{"2.2.2.2"}}},
			{"zone1", &dns.ResourceRecordSet{Name: "www.zone1.com.", Rrdatas: []string{"1.1.1.1"}}},
		}

		changes := changesToUpdateDnsRecordValues(records, []string{"2.2.2.2"})

		So(changes, ShouldResemble, &dns.Change{
			Deletions: []*dns.ResourceRecordSet{
				{Name: "www.zone1.com.", Rrdatas: []string{"1.1.1.1"}},
			},
			Additions: []*dns.ResourceRecordSet{
				{Name: "www.zone1.com.", Rrdatas: []string{"2.2.2.2"}},
			},
		})
	})

	Convey("all records up to date", func() {
		records := DnsRecords{
			{"zone1", &dns.ResourceRecordSet{Name: "zone1.com.", Rrdatas: []string{"2.2.2.2"}}},
			{"zone1", &dns.ResourceRecordSet{Name: "www.zone1.com.", Rrdatas: []string{"2.2.2.2"}}},
		}

		changes := changesToUpdateDnsRecordValues(records, []string{"2.2.2.2"})

		So(changes, ShouldEqual, nil)
	})
}
