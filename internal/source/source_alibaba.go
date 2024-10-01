package source

import (
	"github.com/Escape-Technologies/cloudfinder/internal/log"
	"github.com/Escape-Technologies/cloudfinder/pkg/provider"
)

type Alibaba struct{}

const alibabaFileURL = "https://raw.githubusercontent.com/devanshbatham/ip2cloud/main/data/aliyun.txt"

func (a Alibaba) GetProvider() provider.Provider {
	return provider.Alibaba
}

func (a Alibaba) GetIPRanges() []*IPRange {
	log.Info("Using static Alibaba ip ranges")

	ranges := make([]*IPRange, 0)
	alibabaRanges, err := LoadTextURLToRange(alibabaFileURL)
	if err != nil {
		log.Fatal("Failed to load text url to range for Alibaba", err)
	}
	for _, cdir := range alibabaRanges {
		network, cat := ParseCIDR(cdir)
		ranges = append(ranges, &IPRange{
			Network: network,
			Cat:     cat,
		})
	}
	return ranges
}
