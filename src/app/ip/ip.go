// Copyright © 2023 Esko Luontola
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package ip

import (
	"context"
	"errors"
	"fmt"
	"github.com/huin/goupnp/dcps/internetgateway2"
	"golang.org/x/sync/errgroup"
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

type RouterClient interface {
	GetExternalIPAddress() (
		NewExternalIPAddress string,
		err error,
	)
}

func detectRouterClients(ctx context.Context) ([]RouterClient, error) {
	tasks, _ := errgroup.WithContext(ctx)
	// request each type of client in parallel
	var ip1Clients []*internetgateway2.WANIPConnection1
	tasks.Go(func() error {
		var err error
		ip1Clients, _, err = internetgateway2.NewWANIPConnection1Clients()
		return err
	})
	var ip2Clients []*internetgateway2.WANIPConnection2
	tasks.Go(func() error {
		var err error
		ip2Clients, _, err = internetgateway2.NewWANIPConnection2Clients()
		return err
	})
	var ppp1Clients []*internetgateway2.WANPPPConnection1
	tasks.Go(func() error {
		var err error
		ppp1Clients, _, err = internetgateway2.NewWANPPPConnection1Clients()
		return err
	})

	if err := tasks.Wait(); err != nil {
		return nil, err
	}

	var clients []RouterClient
	for _, client := range ip2Clients {
		clients = append(clients, client)
	}
	for _, client := range ip1Clients {
		clients = append(clients, client)
	}
	for _, client := range ppp1Clients {
		clients = append(clients, client)
	}
	if len(clients) == 0 {
		return nil, errors.New("no UPnP services found")
	}
	return clients, nil
}

var routerClients []RouterClient

func UpnpRouterIP() (string, error) {
	var err error
	if routerClients == nil { // detect the internet gateway only once
		routerClients, err = detectRouterClients(context.Background())
		if err != nil {
			return "", err
		}
	}
	// if multiple clients were found, try them all until one which works was found
	for _, routerClient := range routerClients {
		var ip string
		ip, err = routerClient.GetExternalIPAddress()
		if err == nil && ip != "" { // sometimes the IP is an empty string even though there is no error
			routerClients = []RouterClient{routerClient} // remember the client which worked and keep using only it
			return ip, nil
		}
	}
	// none of the clients worked
	routerClients = nil // re-detect the internet gateway to improve resiliency
	return "", fmt.Errorf("could not get the external IP address through UPnP: %w", err)
}
