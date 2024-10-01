package source

import (
	"encoding/csv"
	"strings"

	"github.com/Escape-Technologies/cloudfinder/internal/log"
	"github.com/Escape-Technologies/cloudfinder/pkg/provider"
)

type Digitalocean struct{}

const doFileURL = "https://digitalocean.com/geo/google.csv"

func (a Digitalocean) GetProvider() provider.Provider {
	return provider.Digitalocean
}

func (a Digitalocean) GetIPRanges() []*IPRange {
	log.Info("Fetching do ip ranges from %s", doFileURL)

	content, err := FileURLToString(doFileURL)
	if err != nil {
		log.Fatal("Failed to read digitalocean ip ranges", err)
	}
	data, err := csv.NewReader(strings.NewReader(content)).ReadAll()
	if err != nil {
		log.Fatal("Failed to read csv", err)
	}

	ranges := make([]*IPRange, 0)
	for _, line := range data {
		cidr := line[0]
		network, cat := ParseCIDR(cidr)
		ranges = append(ranges, &IPRange{
			Network: network,
			Cat:     cat,
		})
	}

	return ranges
}
