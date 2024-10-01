package source

import (
	"encoding/json"

	"github.com/Escape-Technologies/cloudfinder/internal/log"
	"github.com/Escape-Technologies/cloudfinder/pkg/provider"
)

type Ibm struct{}

// This url is sourced from: https://ibm.biz/cidr-calculator, itself referenced in the official docs: https://cloud.ibm.com/docs/cloud-infrastructure?topic=cloud-infrastructure-ibm-cloud-ip-ranges
const ibmFileURL = "https://raw.githubusercontent.com/dprosper/cidr-calculator/main/data/datacenters.json"

type cat []*struct {
	CidrBlocks []string `json:"cidr_blocks"`
}

type ibmJSON struct {
	DataCenters []struct {
		Public        cat `json:"front_end_public_network"`
		LoadBalancers cat `json:"load_balancers_ips"`
		Service       cat `json:"service_network"`
		FileBlock     cat `json:"file_block"`
		Icons         cat `json:"icos"`
		AdvMon        cat `json:"advmon"`
		RheLs         cat `json:"rhe_ls"`
		Ims           cat `json:"ims"`
	} `json:"data_centers"`
}

func (a Ibm) GetProvider() provider.Provider {
	return provider.Ibm
}

func (a Ibm) GetIPRanges() []*IPRange {
	log.Info("Fetching ibm ranges")

	ranges := make([]*IPRange, 0)
	j := &ibmJSON{}
	err := LoadFileURLToJSON(ibmFileURL, j)
	if err != nil {
		log.Fatal("Failed to load json from url for IBM", err)
	}

	// Much nesting lol
	for _, d := range j.DataCenters {
		// convert to map to iterate over fields easily
		b, _ := json.Marshal(d)
		m := map[string]cat{}
		err = json.Unmarshal(b, &m)
		if err != nil {
			log.Fatal("Failed to load json from url for IBM", err)
		}

		for _, c := range m {
			for _, cidrs := range c {
				for _, cidr := range cidrs.CidrBlocks {
					network, cat := ParseCIDR(cidr)
					if isPrivateNetwork(network) {
						continue
					}
					ranges = append(ranges, &IPRange{
						Network: network,
						Cat:     cat,
					})
				}
			}
		}
	}

	return ranges
}
