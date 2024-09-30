package source

import (
	"fmt"
	"net"
	"sync"

	"github.com/Escape-Technologies/cloudfinder/internal/log"
	"github.com/Escape-Technologies/cloudfinder/pkg/provider"
)

type IPCat int

const (
	CatIPv4 IPCat = 4
	CatIPv6 IPCat = 6
)

type IPRange struct {
	Network  *net.IPNet        `json:"n"`
	Cat      IPCat             `json:"c"`
	Provider provider.Provider `json:"p"`
}

func (r *IPRange) String() string {
	return r.Network.String() + fmt.Sprint(r.Cat) + r.Provider.String()
}

type IPRangeSource interface {
	GetIPRanges() []*IPRange
	GetProvider() provider.Provider
}

var AllSources = []IPRangeSource{
	Alibaba{},
	Aws{},
	Azure{},
	Cloudflare{},
	Digitalocean{},
	Fastly{},
	Gcp{},
	Ibm{},
	Linode{},
	Oracle{},
	Ovh{},
	Scaleway{},
	Tencent{},
	Ucloud{},
	Vercel{},
}

func GetAllIPRanges(sources []IPRangeSource) []*IPRange {
	rangeLock := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	ranges := make([]*IPRange, 0)
	for _, source := range sources {
		wg.Add(1)
		go func(s IPRangeSource) {
			p := s.GetProvider()
			sourceRanges := s.GetIPRanges()
			addProviderToRanges(p, sourceRanges)
			log.Debug("[%s] - Got: %d ranges", p, len(sourceRanges))
			sourceRanges = keepWidestNetworks(sourceRanges)
			log.Debug("[%s] - After deduplication: %d ranges", p, len(sourceRanges))
			rangeLock.Lock()
			ranges = append(ranges, sourceRanges...)
			rangeLock.Unlock()
			wg.Done()
		}(source)
	}
	wg.Wait()
	return ranges
}

// Validate at compile time that we have exactly one range per provider
func init() { //nolint:gochecknoinits
	// Check that len(AllSources) = len(provider.ProviderMap) - 1 (Unknown)
	if len(AllSources) != (len(provider.ProviderMap) - 1) {
		err := fmt.Errorf("len(AllSources) = %d != len(provider.ProviderMap) - 1 = %d. Ensure the enum and sources are up to date", len(AllSources), (len(provider.ProviderMap) - 1))
		log.Fatal("Error", err)
	}

	providers := make(map[provider.Provider]bool)

	// Check that each source reagisters a different provider
	for _, source := range AllSources {
		p := source.GetProvider()
		if _, ok := providers[p]; ok {
			err := fmt.Errorf("provider %s used more than once", p.String())
			log.Fatal("Error", err)
		}
		providers[p] = true
	}

	// Check that each provider has a source
	for p := range provider.ProviderMap {
		if p == provider.Unknown {
			continue
		}

		if _, ok := providers[p]; !ok {
			err := fmt.Errorf("provider %s has no source", p.String())
			log.Fatal("Error", err)
		}
	}
}

func addProviderToRanges(p provider.Provider, ranges []*IPRange) {
	// TODO: Validate that we se the provider only once and exactly once
	for _, r := range ranges {
		r.Provider = p
	}
}

func keepWidestNetworks(ranges []*IPRange) []*IPRange {
	// Map all networks to their address
	m := make(map[string][]*IPRange)
	for _, r := range ranges {
		if _, ok := m[r.Network.IP.String()]; !ok {
			m[r.Network.IP.String()] = make([]*IPRange, 0)
		}
		m[r.Network.IP.String()] = append(m[r.Network.IP.String()], r)
	}

	largestOnlyRanges := make([]*IPRange, 0)

	// Iterate over the map
	for _, v := range m {
		// Should not happen but safeguard
		if len(v) == 0 {
			continue
		}

		largest := v[0]
		for _, r := range v {
			if !largest.Network.Contains(r.Network.IP) {
				largest = r
			}
		}

		largestOnlyRanges = append(largestOnlyRanges, largest)
	}

	return largestOnlyRanges
}
