// Copyright © 2023 Esko Luontola
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package ip

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"regexp"
	"testing"
)

func TestIP(t *testing.T) {
	Convey("ExternalServiceIPSpec", t, ExternalServiceIPSpec)
	Convey("OutgoingIPSpec", t, OutgoingIPSpec)
	Convey("InterfaceIPSpec", t, InterfaceIPSpec)
}

func ExternalServiceIPSpec() {
	Convey("automatically detects the IP address by calling an external service", func() {
		Convey("response body contains only the IP address", func() {
			ip, err := ExternalServiceIP("https://ifconfig.me/ip")
			So(err, ShouldBeNil)
			So(ip, ShouldMatchPattern, IpAddress)
		})
		Convey("the IP address is wrapped in some HTML and text", func() {
			ip, err := ExternalServiceIP("http://checkip.dyndns.org/")
			So(err, ShouldBeNil)
			So(ip, ShouldMatchPattern, IpAddress)
		})
	})
	Convey("error: response contains no IP", func() {
		ip, err := ExternalServiceIP("https://httpstat.us/200")
		So(err, ShouldBeError, "the response did not contain an IP address: 200 OK")
		So(ip, ShouldEqual, "")
	})
	Convey("error: non-OK status code", func() {
		ip, err := ExternalServiceIP("https://httpstat.us/500")
		So(err, ShouldBeError, `the server returned status 500 Internal Server Error`)
		So(ip, ShouldEqual, "")
	})
}

func OutgoingIPSpec() {
	Convey("automatically detects the outgoing IP address", func() {
		ip, err := OutgoingIP()
		So(err, ShouldBeNil)
		So(ip, ShouldMatchPattern, IpAddress)
	})
}

func InterfaceIPSpec() {
	Convey("returns loopback interface's IP address", func() {
		ip, err := InterfaceIP("lo")
		So(err, ShouldBeNil)
		So(ip, ShouldEqual, "127.0.0.1")
	})
	Convey("returns named interface's IP address", func() {
		ip, err := InterfaceIP("eth0")
		So(err, ShouldBeNil)
		So(ip, ShouldMatchPattern, IpAddress)
	})
}

const IpAddress = `^\d+\.\d+\.\d+.\d+$`

func ShouldMatchPattern(actual interface{}, expected ...interface{}) string {
	pattern := expected[0]
	matched, err := regexp.MatchString(pattern.(string), actual.(string))
	if err != nil {
		return err.Error()
	}
	if !matched {
		return fmt.Sprintf("Expected '%v' match the pattern '%v', but it didn't!", actual, pattern)
	}
	return ""
}
