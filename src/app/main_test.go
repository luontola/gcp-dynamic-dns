// Copyright © 2018 Esko Luontola
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestGCloud(t *testing.T) {
	Convey("ParamDnsNamesSpec", t, ParamDnsNamesSpec)
}

func ParamDnsNamesSpec() {
	DnsNames := "DNS_NAMES"
	defer os.Unsetenv(DnsNames)

	Convey("one name", func() {
		os.Setenv(DnsNames, "domain1.example.com.")
		So(paramDnsNames(), ShouldResemble, []string{"domain1.example.com."})
	})

	Convey("many names", func() {
		os.Setenv(DnsNames, "domain1.example.com. domain2.example.com. domain3.example.com.")
		So(paramDnsNames(), ShouldResemble, []string{"domain1.example.com.", "domain2.example.com.", "domain3.example.com."})
	})
}
