package source

import (
	"github.com/Escape-Technologies/cloudfinder/internal/log"
	"github.com/Escape-Technologies/cloudfinder/pkg/provider"
)

type Aws struct{}

const awsFileURL = "https://ip-ranges.amazonaws.com/ip-ranges.json"

type awsJSON struct {
	Prefixes []struct {
		IPPrefix string `json:"ip_prefix"`
		Region   string `json:"region"`
		Service  string `json:"service"`
	} `json:"prefixes"`
}

func (a Aws) GetProvider() provider.Provider {
	return provider.Aws
}

func (a Aws) GetIPRanges() []*IPRange {
	log.Info("Fetching AWS ip ranges from %s", awsFileURL)

	var awsJSON *awsJSON
	err := LoadFileURLToJSON(awsFileURL, &awsJSON)
	if err != nil {
		log.Fatal("Failed to load file url to json for AWS", err)
	}

	ranges := make([]*IPRange, 0)
	for _, prefix := range awsJSON.Prefixes {
		network, cat := ParseCIDR(prefix.IPPrefix)
		ranges = append(ranges, &IPRange{
			Network: network,
			Cat:     cat,
		})
	}

	return ranges
}
