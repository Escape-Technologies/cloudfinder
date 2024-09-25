package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/Escape-Technologies/cloudfinder/internal/log"
	"github.com/Escape-Technologies/cloudfinder/pkg/cloud"
	"github.com/Escape-Technologies/cloudfinder/pkg/provider"
)

type args struct {
	inputs chan string
	debug  bool
	mode   outputMode
}

func printUsage() {
	// TODO(@nohehf): better usage message
	println("Usage:")
	println("cloudfinder [flags] <ip, host, domain, url> <ip, host, domain, url> ...")
	println("Flags:")
	flag.PrintDefaults()
}

func hasPipe() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return false
	}
	return true
}

// Return inputs from cli args or fallback to pipe (stdin) if possible
func getInputs(in *chan string) {
	// Prioritize args over pipes
	if len(flag.Args()) == 0 && !hasPipe() {
		println("ERROR: No argument provided via arguments or stdin. At least one URL, IP or domain is required")
		printUsage()
		os.Exit(1)
	}

	// Send arguments to channel if available
	if len(flag.Args()) > 0 {
		for _, f := range flag.Args() {
			*in <- f
		}
		close(*in) // Close the channel after sending all arguments
		return
	}

	// Otherwise, read from stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		*in <- scanner.Text()
	}
	close(*in) // Close the channel after reading from stdin

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to read from stdin: %v\n", err)
		os.Exit(1)
	}
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

	switch {
	case json:
		a.mode = outputJson
	case raw:
		a.mode = outputRaw
	}

	a.inputs = make(chan string)
	go getInputs(&a.inputs)
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

	for i := range a.inputs {
		ips, err := getIPsForURL(context.Background(), i)
		if err != nil {
			log.Error("Failed to get ips, verify input", err)
		}

		for _, ip := range ips {
			p := r.GetProviderForIP(ip)
			printOutput(i, ip, p, a.mode)
		}
	}
}

type outputMode int

const (
	outputDefault outputMode = iota
	outputRaw
	outputJson // nolint:revive
)

func marshallOutput(input string, ip net.IP, p provider.Provider) string {
	toMarshall := struct {
		Input string `json:"input"`
		IP    string `json:"ip"`
		P     string `json:"provider"`
	}{
		Input: input,
		IP:    ip.String(),
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
		// Print to stdout
		fmt.Println(marshallOutput(input, ip, p))
	case outputRaw:
		// Print to stdout
		fmt.Println(p.String())
	}
}
