package source

import (
	"escape.tech/cloudfinder/internal/log"
	"escape.tech/cloudfinder/pkg/provider"
)

type Fastly struct{}

const fastlyFileRangesAPIURL = "https://api.fastly.com/public-ip-list"

type FastlyIPRangeResponse struct {
	Addresses     []string `json:"addresses"`
	IPv6Addresses []string `json:"ipv6_addresses"`
}

func (a Fastly) GetIPRanges() []*IPRange {
	log.Info("Using static Fastly ip ranges")

	var fastlyRanges FastlyIPRangeResponse
	err := LoadFileURLToJSON(fastlyFileRangesAPIURL, &fastlyRanges)
	if err != nil {
		log.Fatal("Failed to load file url to json for Fastly", err)
	}

	fastlyRanges.Addresses = append(fastlyRanges.Addresses, fastlyRanges.IPv6Addresses...)
	ranges := make([]*IPRange, 0)
	for _, cdir := range fastlyRanges.Addresses {
		network, cat := ParseCIDR(cdir)
		ranges = append(ranges, &IPRange{
			Network:  network,
			Cat:      cat,
			Provider: provider.Fastly,
		})
	}
	return ranges
}
