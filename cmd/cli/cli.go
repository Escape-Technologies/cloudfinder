package main

import (
	"context"
	"encoding/json"
	"flag"
	"log/slog"
	"net"
	"os"
	"strings"

	"escape.tech/cloudfinder/internal/log"
	"escape.tech/cloudfinder/pkg/cloud"
	"escape.tech/cloudfinder/pkg/provider"
)

type args struct {
	input string
	debug bool
	mode  outputMode
}

func printUsage() {
	println("Usage:")
	println("cloudfinder [flags] <ip, host, domain, url>")
	println("Flags:")
	flag.PrintDefaults()
}

func parseArgs() args {
	a := args{}
	a.mode = outputDefault
	flag.BoolVar(&a.debug, "debug", false, "enable debug mode")

	var json, raw, help bool
	flag.BoolVar(&json, "json", false, "output json")
	flag.BoolVar(&raw, "raw", false, "output raw provider string")

	flag.BoolVar(&help, "help", false, "print help")
	flag.BoolVar(&help, "h", false, "print help")

	flag.Parse()

	if help {
		printUsage()
		os.Exit(0)
	}

	if flag.Arg(0) == "" {
		println("ERROR: No argument provided. URL, IP or domain is required")
		printUsage()
		os.Exit(1)
	}

	switch {
	case json:
		a.mode = outputJson
	case raw:
		a.mode = outputRaw
	}

	a.input = strings.TrimSpace(flag.Arg(0))
	return a
}

func main() {
	r := cloud.NewResolver()
	a := parseArgs()

	// Set the right logger
	if a.debug {
		r.WithLogger(log.NewLogger(slog.LevelDebug))
	} else {
		r.WithLogger(log.NewPrettyLogger(slog.LevelInfo))
	}

	ips, err := getIPsForURL(context.Background(), a.input)
	if err != nil {
		log.Error("Failed to get ips, verify input", err)
		os.Exit(1)
	}

	for _, ip := range ips {
		p := r.GetProviderForIP(ip)
		printOutput(a.input, ip, p, a.mode)
	}
}

type outputMode int

const (
	outputDefault outputMode = iota
	outputRaw
	outputJson
)

func marshallOutput(input string, ip net.IP, p provider.Provider) string {
	toMarshall := struct {
		Input string `json:"input"`
		Ip    string `json:"ip"`
		P     string `json:"provider"`
	}{
		Input: input,
		Ip:    ip.String(),
		P:     p.String(),
	}
	bytes, err := json.Marshal(toMarshall)
	if err != nil {
		log.Fatal("could not output JSON", err)
	}
	return string(bytes)
}

func printOutput(input string, ip net.IP, p provider.Provider, mode outputMode) {
	switch mode {
	case outputDefault:
		log.Info("%s (%s): %s", input, ip.String(), p.String())
	case outputJson:
		println(marshallOutput(input, ip, p))
	case outputRaw:
		println(p.String())
	}
}
