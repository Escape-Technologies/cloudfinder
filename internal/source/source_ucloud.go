package source

import (
	"github.com/Escape-Technologies/cloudfinder/internal/log"
	"github.com/Escape-Technologies/cloudfinder/pkg/provider"
)

type Ucloud struct{}

const ucloudFileURL = "https://raw.githubusercontent.com/devanshbatham/ip2cloud/main/data/ucloud.txt"

func (a Ucloud) GetProvider() provider.Provider {
	return provider.Ucloud
}

func (a Ucloud) GetIPRanges() []*IPRange {
	log.Info("Using static Ucloud ip ranges")

	ranges := make([]*IPRange, 0)
	ucloudRanges, err := LoadTextURLToRange(ucloudFileURL)
	if err != nil {
		log.Fatal("Failed to load text url to range for ucloud", err)
	}
	for _, cidr := range ucloudRanges {
		network, cat := ParseCIDR(cidr)
		ranges = append(ranges, &IPRange{
			Network: network,
			Cat:     cat,
		})
	}
	return ranges
}
