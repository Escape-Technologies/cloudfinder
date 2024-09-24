package source

import (
	"encoding/csv"
	"strings"

	"escape.tech/cloudfinder/internal/log"
	"escape.tech/cloudfinder/pkg/provider"
)

type Digitalocean struct{}

const doFileURL = "https://digitalocean.com/geo/google.csv"

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
		cdir := line[0]
		network, cat := ParseCIDR(cdir)
		ranges = append(ranges, &IPRange{
			Network:  network,
			Cat:      cat,
			Provider: provider.Digitalocean,
		})
	}

	return ranges
}
