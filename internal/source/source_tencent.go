package source

import (
	"escape.tech/cloudfinder/internal/log"
	"escape.tech/cloudfinder/pkg/provider"
)

type Tencent struct{}

const tencentFileURL = "https://raw.githubusercontent.com/devanshbatham/ip2cloud/main/data/tencent.txt"

func (a Tencent) GetIPRanges() []*IPRange {
	log.Info("Using static Tencent ip ranges")

	ranges := make([]*IPRange, 0)
	tencentRanges, err := LoadTextURLToRange(tencentFileURL)
	if err != nil {
		log.Fatal("Failed to load text url to range for Tencent", err)
	}
	for _, cdir := range tencentRanges {
		network, cat := ParseCIDR(cdir)
		ranges = append(ranges, &IPRange{
			Network:  network,
			Cat:      cat,
			Provider: provider.Tencent,
		})
	}
	return ranges
}
