// Copyright © 2023 Esko Luontola
// This software is released under the Apache License 2.0.
// The license text is at http://www.apache.org/licenses/LICENSE-2.0

package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	Mode                string
	ServiceUrls         []string
	nextServiceUrlIndex int
	InterfaceName       string
	DnsNames            []string
	GoogleProject       string
}

func FromEnv() *Config {
	config := &Config{
		Mode:                envOrDefault("MODE", "service"),
		ServiceUrls:         strings.Fields(envOrDefault("SERVICE_URLS", "https://ifconfig.me/ip http://checkip.dyndns.org/ http://ip1.dynupdate.no-ip.com/")),
		nextServiceUrlIndex: 0,
		InterfaceName:       envOrDefault("INTERFACE_NAME", ""),
		DnsNames:            strings.Fields(envOrFail("DNS_NAMES")),
		GoogleProject:       envOrFail("GOOGLE_PROJECT"),
	}
	return config
}

func (config *Config) NextServiceUrl() string {
	urls := config.ServiceUrls
	index := config.nextServiceUrlIndex
	url := urls[index]
	config.nextServiceUrlIndex = (index + 1) % len(urls)
	return url
}

func envOrDefault(key string, defaultValue string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return defaultValue
	}
	return v
}

func envOrFail(key string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		log.Fatal("Environment variable ", key, " was not set")
	}
	return v
}
