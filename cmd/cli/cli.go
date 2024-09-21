package main

import (
	"context"
	"errors"
	"flag"
	"os"
	"strings"

	"escape.tech/cloudfinder/internal/log"
	"escape.tech/cloudfinder/pkg/cloud"
)

func getArg() string {
	flag.Parse()

	if flag.Arg(0) == "" {
		log.Error("No arg provided", errors.New("URL, IP or domain required as argument"))
		os.Exit(1)
	}
	input := flag.Arg(0)
	return strings.TrimSpace(input)
}

func main() {
	r := cloud.NewResolver()
	input := getArg()
	ips, err := getIPsForURL(context.Background(), input)
	if err != nil {
		log.Error("Failed to get ips, verify input", err)
		os.Exit(1)
	}
	for _, ip := range ips {
		provider := r.GetProviderForIP(ip)
		log.Info("%s (%s): %s", input, ip.String(), provider.String())
	}
}
