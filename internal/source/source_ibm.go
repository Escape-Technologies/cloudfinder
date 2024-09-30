package source

import (
	"github.com/Escape-Technologies/cloudfinder/internal/log"
	"github.com/Escape-Technologies/cloudfinder/pkg/provider"
)

type Ibm struct{}

const ibmFileURL = "https://raw.githubusercontent.com/devanshbatham/ip2cloud/main/data/ibm.txt"

func (a Ibm) GetProvider() provider.Provider {
	return provider.Ibm
}

func (a Ibm) GetIPRanges() []*IPRange {
	log.Info("Using static Ibm ip ranges")

	ranges := make([]*IPRange, 0)
	ibmRanges, err := LoadTextURLToRange(ibmFileURL)
	if err != nil {
		log.Fatal("Failed to load text url to range for IBM", err)
	}
	for _, cdir := range ibmRanges {
		network, cat := ParseCIDR(cdir)
		ranges = append(ranges, &IPRange{
			Network: network,
			Cat:     cat,
		})
	}
	return ranges
}
