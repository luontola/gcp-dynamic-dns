// Copyright © 2023 Esko Luontola
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
	Convey("FilterDnsRecordsByTypeSpec", t, FilterDnsRecordsByTypeSpec)
	Convey("GroupDnsRecordsByZoneSpec", t, GroupDnsRecordsByZoneSpec)
	Convey("UpdateDnsRecordValuesSpec", t, UpdateDnsRecordValuesSpec)
}

func FilterDnsRecordsByNameSpec() {
	records := []*DnsRecord{
		{ManagedZone: "example", ResourceRecordSet: &dns.ResourceRecordSet{Name: "foo.example.com."}},
		{ManagedZone: "example", ResourceRecordSet: &dns.ResourceRecordSet{Name: "bar.example.com."}},
		{ManagedZone: "example", ResourceRecordSet: &dns.ResourceRecordSet{Name: "baz.example.com."}},
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

func FilterDnsRecordsByTypeSpec() {
	records := []*DnsRecord{
		{ManagedZone: "example", ResourceRecordSet: &dns.ResourceRecordSet{Name: "a1.example.com.", Type: "A"}},
		{ManagedZone: "example", ResourceRecordSet: &dns.ResourceRecordSet{Name: "cname.example.com.", Type: "CNAME"}},
		{ManagedZone: "example", ResourceRecordSet: &dns.ResourceRecordSet{Name: "a2.example.com.", Type: "A"}},
	}

	records = filterDnsRecordsByType(records, "A")

	So(records, ShouldHaveLength, 2)
	So(records[0].Name, ShouldEqual, "a1.example.com.")
	So(records[1].Name, ShouldEqual, "a2.example.com.")
}

func GroupDnsRecordsByZoneSpec() {
	records := DnsRecords{
		{ManagedZone: "zone1", ResourceRecordSet: &dns.ResourceRecordSet{Name: "zone1.com."}},
		{ManagedZone: "zone2", ResourceRecordSet: &dns.ResourceRecordSet{Name: "zone2.com."}},
		{ManagedZone: "zone2", ResourceRecordSet: &dns.ResourceRecordSet{Name: "www.zone2.com."}},
	}

	So(records.GroupByZone(), ShouldResemble, map[string]DnsRecords{
		"zone1": {
			{ManagedZone: "zone1", ResourceRecordSet: &dns.ResourceRecordSet{Name: "zone1.com."}},
		},
		"zone2": {
			{ManagedZone: "zone2", ResourceRecordSet: &dns.ResourceRecordSet{Name: "zone2.com."}},
			{ManagedZone: "zone2", ResourceRecordSet: &dns.ResourceRecordSet{Name: "www.zone2.com."}},
		},
	})
}

func UpdateDnsRecordValuesSpec() {
	Convey("one record, one value", func() {
		records := DnsRecords{
			{ManagedZone: "zone1", ResourceRecordSet: &dns.ResourceRecordSet{Name: "zone1.com.", Rrdatas: []string{"1.1.1.1"}}},
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
		So(ToDnsRecords("zone1", changes), ShouldResemble, DnsRecords{
			{ManagedZone: "zone1", OldRrdatas: []string{"1.1.1.1"},
				ResourceRecordSet: &dns.ResourceRecordSet{Name: "zone1.com.", Rrdatas: []string{"2.2.2.2"}}}})
	})

	Convey("multiple records", func() {
		records := DnsRecords{
			{ManagedZone: "zone1", ResourceRecordSet: &dns.ResourceRecordSet{Name: "zone1.com.", Rrdatas: []string{"1.1.1.1"}}},
			{ManagedZone: "zone1", ResourceRecordSet: &dns.ResourceRecordSet{Name: "www.zone1.com.", Rrdatas: []string{"1.1.1.1"}}},
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
		So(ToDnsRecords("zone1", changes), ShouldResemble, DnsRecords{
			{ManagedZone: "zone1", OldRrdatas: []string{"1.1.1.1"},
				ResourceRecordSet: &dns.ResourceRecordSet{Name: "zone1.com.", Rrdatas: []string{"2.2.2.2"}}},
			{ManagedZone: "zone1", OldRrdatas: []string{"1.1.1.1"},
				ResourceRecordSet: &dns.ResourceRecordSet{Name: "www.zone1.com.", Rrdatas: []string{"2.2.2.2"}}}})
	})

	Convey("multiple values", func() {
		records := DnsRecords{
			{ManagedZone: "zone1", ResourceRecordSet: &dns.ResourceRecordSet{Name: "zone1.com.", Rrdatas: []string{"1.1.1.1", "2.2.2.2"}}},
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
		So(ToDnsRecords("zone1", changes), ShouldResemble, DnsRecords{
			{ManagedZone: "zone1", OldRrdatas: []string{"1.1.1.1", "2.2.2.2"},
				ResourceRecordSet: &dns.ResourceRecordSet{Name: "zone1.com.", Rrdatas: []string{"3.3.3.3", "4.4.4.4"}}}})
	})

	Convey("some records up to date", func() {
		records := DnsRecords{
			{ManagedZone: "zone1", ResourceRecordSet: &dns.ResourceRecordSet{Name: "zone1.com.", Rrdatas: []string{"2.2.2.2"}}},
			{ManagedZone: "zone1", ResourceRecordSet: &dns.ResourceRecordSet{Name: "www.zone1.com.", Rrdatas: []string{"1.1.1.1"}}},
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
		So(ToDnsRecords("zone1", changes), ShouldResemble, DnsRecords{
			{ManagedZone: "zone1", OldRrdatas: []string{"1.1.1.1"},
				ResourceRecordSet: &dns.ResourceRecordSet{Name: "www.zone1.com.", Rrdatas: []string{"2.2.2.2"}}}})
	})

	Convey("all records up to date", func() {
		records := DnsRecords{
			{ManagedZone: "zone1", ResourceRecordSet: &dns.ResourceRecordSet{Name: "zone1.com.", Rrdatas: []string{"2.2.2.2"}}},
			{ManagedZone: "zone1", ResourceRecordSet: &dns.ResourceRecordSet{Name: "www.zone1.com.", Rrdatas: []string{"2.2.2.2"}}},
		}

		changes := changesToUpdateDnsRecordValues(records, []string{"2.2.2.2"})

		So(changes, ShouldBeNil)
	})
}
