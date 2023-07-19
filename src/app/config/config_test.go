// Copyright © 2023 Esko Luontola
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package config

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	Convey("ConfigSpec", t, ConfigSpec)
}

func ConfigSpec() {
	GoogleProject := "GOOGLE_PROJECT"
	os.Setenv(GoogleProject, "dummy")
	defer os.Unsetenv(GoogleProject)

	DnsNames := "DNS_NAMES"
	os.Setenv(DnsNames, "dummy")
	defer os.Unsetenv(DnsNames)

	ServiceUrls := "SERVICE_URLS"
	defer os.Unsetenv(ServiceUrls)

	Convey(DnsNames, func() {
		os.Setenv(DnsNames, "domain1.example.com. domain2.example.com. domain3.example.com.")
		conf := FromEnv()
		So(conf.DnsNames, ShouldResemble, []string{"domain1.example.com.", "domain2.example.com.", "domain3.example.com."})
	})

	Convey("NextServiceUrl rotates through all URLs", func() {
		os.Setenv(ServiceUrls, "http://url1 http://url2")
		conf := FromEnv()
		So(conf.NextServiceUrl(), ShouldEqual, "http://url1")
		So(conf.NextServiceUrl(), ShouldEqual, "http://url2")
		So(conf.NextServiceUrl(), ShouldEqual, "http://url1")
	})
}
