package source

import (
	"github.com/Escape-Technologies/cloudfinder/internal/log"
	"github.com/Escape-Technologies/cloudfinder/pkg/provider"
)

type Scaleway struct{}

// Source: https://ipinfo.io/AS12876#block-ranges
var scalewayRanges = [...]string{
	"151.115.0.0/18",
	"151.115.64.0/18",
	"163.172.0.0/16",
	"163.172.208.0/20",
	"195.154.0.0/16",
	"212.129.0.0/18",
	"212.47.224.0/19",
	"212.83.128.0/19",
	"212.83.160.0/19",
	"51.15.0.0/16",
	"51.15.0.0/17",
	"51.158.0.0/15",
	"51.158.128.0/17",
	"62.210.0.0/16",
	"62.4.0.0/19",
}

func (a Scaleway) GetProvider() provider.Provider {
	return provider.Scaleway
}

func (a Scaleway) GetIPRanges() []*IPRange {
	log.Info("Using static Scaleway ip ranges")

	ranges := make([]*IPRange, 0)
	for _, cidr := range scalewayRanges {
		network, cat := ParseCIDR(cidr)
		ranges = append(ranges, &IPRange{
			Network: network,
			Cat:     cat,
		})
	}

	return ranges
}
