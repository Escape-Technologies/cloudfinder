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

// Digital ocean csv has been broken, with lines containing less than the usual 5 fields
// This will remove incorrect lines
func fixCsv(in string) string {
	var result []string
	lines := strings.Split(in, "\n")

	for _, line := range lines {
		if strings.Count(line, ",") == 4 { //nolint:mnd
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

func (a Digitalocean) GetIPRanges() []*IPRange {
	log.Info("Fetching do ip ranges from %s", doFileURL)

	content, err := FileURLToString(doFileURL)
	if err != nil {
		log.Fatal("Failed to read digitalocean ip ranges", err)
	}
	content = fixCsv(content)
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
