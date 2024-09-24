package source

import (
	"escape.tech/cloudfinder/internal/log"
	"escape.tech/cloudfinder/pkg/provider"
)

type Linode struct{}

const LinodeFileURL = "https://raw.githubusercontent.com/devanshbatham/ip2cloud/main/data/linode.txt"

func (a Linode) GetIPRanges() []*IPRange {
	log.Info("Using static Linode ip ranges")

	ranges := make([]*IPRange, 0)
	linodeRanges, err := LoadTextURLToRange(LinodeFileURL)
	if err != nil {
		log.Fatal("Failed to load text url to range for Linode", err)
	}
	for _, cdir := range linodeRanges {
		network, cat := ParseCIDR(cdir)
		ranges = append(ranges, &IPRange{
			Network:  network,
			Cat:      cat,
			Provider: provider.Linode,
		})
	}
	return ranges
}
