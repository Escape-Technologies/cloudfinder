package main

import (
	"fmt"
	"testing"
)

func TestParseIp(t *testing.T) {
	tt := []struct {
		str      string
		hostname string
	}{
		{
			str:      "example.com",
			hostname: "example.com",
		},
		{
			str:      "http://example.com",
			hostname: "example.com",
		},
		{
			str:      "example.com/some/path",
			hostname: "example.com",
		},
		{
			str:      "example.com:8080",
			hostname: "example.com",
		},
		{
			str:      "example.com:8080",
			hostname: "example.com",
		},
		{
			str:      "https://example.com:8080/some/path",
			hostname: "example.com",
		},
		{
			str:      "0.0.0.0",
			hostname: "0.0.0.0",
		},
		{
			str:      "http://0.0.0.0:1234/some/path",
			hostname: "0.0.0.0",
		},
		{
			str:      "0.0.0.0:1234",
			hostname: "0.0.0.0",
		},
	}

	for _, test := range tt {
		t.Run(test.str, func(t *testing.T) {
			t.Parallel()
			hostmame, err := parseHostname(test.str)
			if err != nil {
				t.Error(err)
			}
			if hostmame != test.hostname {
				err := fmt.Errorf("parseHostname(%s) = \"%s\" expected \"%s\"", test.str, hostmame, test.hostname)
				t.Error(err)
			}
		})
	}
}
