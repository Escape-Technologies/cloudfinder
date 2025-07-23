package source

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Escape-Technologies/cloudfinder/internal/log"
)

const defaultHTTPTimeout = 30 * time.Second

func GetIPCat(ip net.IP) IPCat {
	if ip.To4() != nil {
		return CatIPv4
	}
	return CatIPv6
}

// ParseIpRange parses an ip range string into start, end, and ip type
// eg. ParseIpRange(8.8.4.0/24) -> (ip.Net, IP_TYPE_IPV4)
func ParseCIDR(cidr string) (*net.IPNet, IPCat) {
	ip, network, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}

	return network, GetIPCat(ip)
}

func FileURLToString(url string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultHTTPTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	file, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if file == nil {
		return "", errors.New("response is nil")
	}
	if file.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code is %d", file.StatusCode)
	}
	defer file.Body.Close()

	bodyText, err := io.ReadAll(file.Body)
	if err != nil {
		return "", err
	}

	return string(bodyText), nil
}

func LoadFileURLToJSON(url string, to interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultHTTPTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	file, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if file == nil {
		return errors.New("response is nil")
	}
	if file.StatusCode != http.StatusOK {
		return fmt.Errorf("status code is %d", file.StatusCode)
	}
	defer file.Body.Close()

	return json.NewDecoder(file.Body).Decode(&to)
}

func LoadTextURLToRange(url string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultHTTPTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s", url)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err := fmt.Errorf("failed to read body %s", url)
		panic(err)
	}
	return strings.Split(string(body), "\n"), nil
}

func isPrivateNetwork(n *net.IPNet) bool {
	// Source: https://en.wikipedia.org/wiki/Private_network
	for _, cidr := range []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"100.64.0.0/10",
		"fc00::/7",
	} {
		_, pn, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(err)
		}

		if pn.Contains(n.IP) {
			return true
		}
	}
	return false
}

/// BGP TOOLS

// ensure we fill the map only once
var bgpToolsMutex = sync.Mutex{}

// ASN -> CIDR
var bgpToolsAsnRanges map[string][]*IPRange
var bgpToolsTableURL = "https://bgp.tools/table.txt"

func requestWithRetry(req *http.Request, maxRetries int) (*http.Response, error) {
	TIMEOUT := 45 // to avoid linter cries
	httpClient := &http.Client{
		Timeout: time.Duration(TIMEOUT) * time.Second,
	}

	var globalErr error = nil
	for i := 0; i < maxRetries; i++ { // retry 3 times
		if i > 0 {
			// exponential backoff to avoid overwhelming the server
			two := 2.0 // linter cries a lot
			time.Sleep(time.Duration(math.Pow(two, float64(i))) * time.Second)
		}

		res, err := httpClient.Do(req)
		globalErr = err
		if err != nil {
			if i < maxRetries-1 {
				log.Error("Error getting bgp tools table", err)
				continue
			}
			// last try failed
			log.Fatal("Error getting bgp tools table", err)
		}
		if res.StatusCode != http.StatusOK && i < maxRetries-1 {
			err := fmt.Errorf("status code: %d", res.StatusCode)
			if i < maxRetries-1 {
				log.Error("Error getting bgp tools table", err)
				continue
			}
			// last try failed
			log.Fatal("Error getting bgp tools table", err)
		}
		return res, nil
	}
	return nil, globalErr
}

// Fetches https://bgp.tools/table.txt and parses it into the ASN -> CIDR map
func getRangesForAsn(asn string) []*IPRange {
	bgpToolsMutex.Lock()
	if bgpToolsAsnRanges == nil {
		bgpToolsAsnRanges = make(map[string][]*IPRange)

		// Fill the map
		log.Info("Fetching AS infos from %s", bgpToolsTableURL)
		req, _ := http.NewRequest(http.MethodGet, bgpToolsTableURL, nil) // nolint: noctx
		// bgp tools requires a descriptive user agent in case the program gets out of control
		req.Header.Add("user-agent", "https://github.com/Escape-Technologies/cloudfinder - nohe@escape.tech")
		
		const MaxRetries = 3
		res, err := requestWithRetry(req, MaxRetries)
		if err != nil {
			log.Fatal("Error getting bgp tools table", err)
		}

		scanner := bufio.NewScanner(res.Body)
		// Read lines
		for scanner.Scan() {
			line := scanner.Text()
			// Line is formatted as "<CIDR> <ASN>"
			x := strings.Split(line, " ")
			if len(x) != 2 { // nolint: mnd
				continue
			}
			cidr, asn := x[0], x[1]

			n, cat := ParseCIDR(cidr)
			// Skip private networks
			if isPrivateNetwork(n) {
				continue
			}
			// Fill map
			if _, ok := bgpToolsAsnRanges[asn]; !ok {
				bgpToolsAsnRanges[asn] = make([]*IPRange, 0)
			}
			bgpToolsAsnRanges[asn] = append(bgpToolsAsnRanges[asn], &IPRange{
				Network: n,
				Cat:     cat,
			})
		}
		res.Body.Close()
		log.Info("Got %d AS infos", len(bgpToolsAsnRanges))
	}
	bgpToolsMutex.Unlock()
	if val, ok := bgpToolsAsnRanges[asn]; ok {
		return val
	}
	return []*IPRange{}
}
