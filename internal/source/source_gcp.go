package source

import (
	"errors"

	"github.com/Escape-Technologies/cloudfinder/internal/log"
	"github.com/Escape-Technologies/cloudfinder/pkg/provider"
)

type Gcp struct{}

var gcpFileURLs = []string{
	"https://www.gstatic.com/ipranges/goog.json",
	"https://www.gstatic.com/ipranges/cloud.json",
}

type gcpJSON struct {
	Prefixes []struct {
		IPv4Prefix string `json:"ipv4Prefix"`
		IPv6Prefix string `json:"ipv6Prefix"`
	} `json:"prefixes"`
}

func (a Gcp) GetProvider() provider.Provider {
	return provider.Gcp
}

func (a Gcp) GetIPRanges() []*IPRange {
	ranges := make([]*IPRange, 0)

	for _, gcpFileURL := range gcpFileURLs {
		log.Info("Fetching GCP ip ranges from %s", gcpFileURL)

		var gcpJSON *gcpJSON
		err := LoadFileURLToJSON(gcpFileURL, &gcpJSON)
		if err != nil {
			log.Fatal("Failed to load file url to json for GCP", err)
		}

		for _, prefix := range gcpJSON.Prefixes {
			cidr := prefix.IPv4Prefix
			if cidr == "" {
				cidr = prefix.IPv6Prefix
			}

			if cidr == "" {
				log.Fatal("both ipv4 and ipv6 prefixes are empty", errors.New("must have IP"))
			}

			network, cat := ParseCIDR(cidr)
			ranges = append(ranges, &IPRange{
				Network: network,
				Cat:     cat,
			})
		}
	}

	return ranges
}
