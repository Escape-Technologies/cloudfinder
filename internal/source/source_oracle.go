package source

import (
	"github.com/Escape-Technologies/cloudfinder/internal/log"
	"github.com/Escape-Technologies/cloudfinder/pkg/provider"
)

type Oracle struct{}

const oracleFileURL = "https://raw.githubusercontent.com/femueller/cloud-ip-ranges/master/oracle-cloud-ip-ranges.json"

type oracleJSON struct {
	Regions []struct {
		Cidrs []struct {
			Cidr string `json:"cidr"`
		} `json:"cidrs"`
	} `json:"regions"`
}

func (a Oracle) GetProvider() provider.Provider {
	return provider.Oracle
}

func (a Oracle) GetIPRanges() []*IPRange {
	log.Info("Using static Oracle ip ranges")
	var oracleJSON oracleJSON
	err := LoadFileURLToJSON(oracleFileURL, &oracleJSON)
	if err != nil {
		log.Fatal("Failed to load file url to json for Oracle", err)
	}

	ranges := make([]*IPRange, 0)
	for _, region := range oracleJSON.Regions {
		for _, cidrs := range region.Cidrs {
			network, cat := ParseCIDR(cidrs.Cidr)
			ranges = append(ranges, &IPRange{
				Network: network,
				Cat:     cat,
			})
		}
	}

	return ranges
}
