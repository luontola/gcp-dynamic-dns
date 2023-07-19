// Copyright © 2023 Esko Luontola
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package ip

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"time"
)

func ExternalServiceIP(url string) (string, error) {
	client := &http.Client{
		Timeout: time.Minute,
	}
	response, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return "", fmt.Errorf("the server returned status %v", response.Status)
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	body := string(bodyBytes)

	ipAddressPattern := regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)
	found := ipAddressPattern.FindString(body)
	if found != "" && net.ParseIP(found) != nil {
		return found, nil
	}
	return "", fmt.Errorf("the response did not contain an IP address: %v", body)
}

func OutgoingIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := localAddr.IP.String()
	return ip, nil
}

func InterfaceIP(name string) (string, error) {
	ifi, err := net.InterfaceByName(name)
	if err != nil {
		return "", err
	}
	addrs, err := ifi.Addrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		// an interface may have multiple addresses, but we're interested only in the IPv4 address, not IPv6 addresses
		ip4 := addr.(*net.IPNet).IP.To4()
		if ip4 != nil {
			return ip4.String(), nil
		}
	}
	return "", errors.New("interface had no IPv4 addresses")
}
