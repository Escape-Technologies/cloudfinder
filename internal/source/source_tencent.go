package source

import (
	"github.com/Escape-Technologies/cloudfinder/internal/log"
	"github.com/Escape-Technologies/cloudfinder/pkg/provider"
)

type Tencent struct{}

const tencentFileURL = "https://raw.githubusercontent.com/devanshbatham/ip2cloud/main/data/tencent.txt"

func (a Tencent) GetProvider() provider.Provider {
	return provider.Tencent
}

func (a Tencent) GetIPRanges() []*IPRange {
	log.Info("Using static Tencent ip ranges")

	ranges := make([]*IPRange, 0)
	tencentRanges, err := LoadTextURLToRange(tencentFileURL)
	if err != nil {
		log.Fatal("Failed to load text url to range for Tencent", err)
	}
	for _, cidr := range tencentRanges {
		network, cat := ParseCIDR(cidr)
		ranges = append(ranges, &IPRange{
			Network: network,
			Cat:     cat,
		})
	}
	return ranges
}
