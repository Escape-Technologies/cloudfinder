package source

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
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
	for _, cdir := range []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"100.64.0.0/10",
		"fc00::/7",
	} {
		_, pn, err := net.ParseCIDR(cdir)
		if err != nil {
			panic(err)
		}

		if pn.Contains(n.IP) {
			return true
		}
	}
	return false
}
