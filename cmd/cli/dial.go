package main

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"
)

func parseHostname(URL string) (string, error) {
	// If there is no scheme, net/url parsing will fail, parsing the host as path.
	// In that case you add a leading //
	if !(strings.Contains(URL, "//")) {
		URL = "//" + URL
	}

	u, err := url.Parse(URL)
	if err != nil {
		return "", fmt.Errorf("failed to parse url \"%s\": %w", URL, err)
	}

	if u.Host == "" {
		return "", fmt.Errorf("failed to get host from url \"%s\"", URL)
	}

	return u.Hostname(), nil
}

func getIPsForURL(ctx context.Context, URL string) ([]net.IP, error) {
	hostname, err := parseHostname(URL)
	if err != nil {
		return nil, fmt.Errorf("could not get ips for url \"%s\": %w", URL, err)
	}
	// If we already have an ip return it (no need to check the DNS)
	ip := net.ParseIP(hostname)
	if ip != nil {
		return []net.IP{ip}, nil
	}
	ips, err := net.DefaultResolver.LookupIP(ctx, "ip", hostname)
	if err != nil {
		return nil, fmt.Errorf("could not get ips for url \"%s\": %w", URL, err)
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("DNS lookup did not find any ips for host \"%s\"", hostname)
	}
	return ips, nil
}
