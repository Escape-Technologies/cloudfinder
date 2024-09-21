package source

import (
	"escape.tech/cloudfinder/internal/log"
	"escape.tech/cloudfinder/pkg/provider"
)

type Ucloud struct{}

const ucloudFileURL = "https://raw.githubusercontent.com/devanshbatham/ip2cloud/main/data/ucloud.txt"

func (a Ucloud) GetIPRanges() []*IPRange {
	log.Info("Using static Ucloud ip ranges")

	ranges := make([]*IPRange, 0)
	ucloudRanges, err := LoadTextURLToRange(ucloudFileURL)
	if err != nil {
		log.Fatal("Failed to load text url to range for ucloud", err)
	}
	for _, cdir := range ucloudRanges {
		network, cat := ParseCIDR(cdir)
		ranges = append(ranges, &IPRange{
			Network:  network,
			Cat:      cat,
			Provider: provider.Ucloud,
		})
	}
	return ranges
}
