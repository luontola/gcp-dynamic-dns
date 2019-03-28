// Copyright © 2019 Esko Luontola
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
	Convey("OutgoingIPSpec", t, OutgoingIPSpec)
	Convey("InterfaceIPSpec", t, InterfaceIPSpec)
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
