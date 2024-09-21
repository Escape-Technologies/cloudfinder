package cloud

import (
	"net"
	"time"

	"escape.tech/cloudfinder/internal/static"
	"escape.tech/cloudfinder/internal/tree"
	"escape.tech/cloudfinder/pkg/provider"
)

const httpTimeout = 3 * time.Second

type Resolver interface {
	GetProviderForIP(ip net.IP) provider.Provider
}

type resolver struct {
	ipv4Tree tree.Tree
	ipv6Tree tree.Tree
}

func NewResolver() Resolver {
	return &resolver{
		ipv4Tree: static.LoadIPv4Tree(),
		ipv6Tree: static.LoadIPv6Tree(),
	}
}

func (f *resolver) GetProviderForIP(ip net.IP) provider.Provider {
	if ipv4 := ip.To4(); ipv4 != nil {
		ipRange := f.ipv4Tree.FindIPRange(ip)
		if ipRange == nil {
			return 0 // unspecified
		}
		return ipRange.Provider
	}

	ipRange := f.ipv6Tree.FindIPRange(ip.To16())
	if ipRange == nil {
		return 0 // unspecified
	}
	return ipRange.Provider
}

// func (f *resolver) GetProviderForURL(u *url.URL, httpClient httpclient.Client) inventoryV1.CloudProvider {
// 	// if url.Host is an ip, return the provider for that ip
// 	if u.Hostname() == "" {
// 		log.Error("Failed to parse url", errors.New("no hostname"))
// 		return inventoryV1.CloudProvider_CLOUD_PROVIDER_UNSPECIFIED
// 	}

// 	// log.Info("Looking up ip for %+v", u)
// 	if ip := net.ParseIP(u.Host); ip != nil {
// 		// log.Info("Found ip %v for %v", ip, u)
// 		return f.GetProviderForIP(ip)
// 	}

// 	// dial the url and get the ip from the connection
// 	domain := u.Hostname()
// 	ips, err := net.LookupIP(domain)
// 	if err != nil {
// 		log.Error("Failed to lookup IP", err)
// 		return inventoryV1.CloudProvider_CLOUD_PROVIDER_UNSPECIFIED
// 	}
// 	if len(ips) == 0 {
// 		log.Error("Failed to lookup IP", errors.New("no ips found"))
// 		return inventoryV1.CloudProvider_CLOUD_PROVIDER_UNSPECIFIED
// 	}

// 	// return the provider for the first ip matching a provider
// 	for _, ip := range ips {
// 		// log.Info("Found ip %v for %v", ip, u)
// 		if provider := f.GetProviderForIP(ip); provider != inventoryV1.CloudProvider_CLOUD_PROVIDER_UNSPECIFIED {
// 			return provider
// 		}
// 	}

// 	req := httpclient.ReqFromURL(*u).
// 		WithMethod(http.MethodHead)

// 	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
// 	defer cancel()

// 	resp, err := httpClient.Do(ctx, req)
// 	if err != nil {
// 		err = fmt.Errorf("Failed to do request to '%s': %w", u.String(), err)
// 		log.Warning("Failed to do request", err)
// 		return inventoryV1.CloudProvider_CLOUD_PROVIDER_UNSPECIFIED
// 	}

// 	if resp == nil {
// 		err = fmt.Errorf("Failed to do request to '%s': response is nil", u.String())
// 		log.Warning("Failed to do request", err)
// 		return inventoryV1.CloudProvider_CLOUD_PROVIDER_UNSPECIFIED
// 	}

// 	for key, value := range resp.Req.Header {
// 		if provider := f.GetProviderForHeader(key, value); provider != inventoryV1.CloudProvider_CLOUD_PROVIDER_UNSPECIFIED {
// 			return provider
// 		}
// 	}

// 	return inventoryV1.CloudProvider_CLOUD_PROVIDER_UNSPECIFIED
// }
