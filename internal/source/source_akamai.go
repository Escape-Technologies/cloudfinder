package source

import (
	"github.com/Escape-Technologies/cloudfinder/internal/log"
	"github.com/Escape-Technologies/cloudfinder/pkg/provider"
)

type Akamai struct{}

func (a Akamai) GetProvider() provider.Provider {
	return provider.Akamai
}

// source: https://github.com/SecOps-Institute/Akamai-ASN-and-IPs-List/blob/master/akamai_asn_list.lst
var AkamaiASNs = []string{
	"12222",
	"16625",
	"16702",
	"17204",
	"18680",
	"18717",
	"20189",
	"20940",
	"21342",
	"21357",
	"21399",
	"22207",
	"22452",
	"23454",
	"23455",
	"23903",
	"24319",
	"26008",
	"30675",
	"31107",
	"31108",
	"31109",
	"31110",
	"31377",
	"33047",
	"33905",
	"34164",
	"34850",
	"35204",
	"35993",
	"35994",
	"36183",
	"39836",
	"43639",
	"55409",
	"55770",
	// "63949", Linode
	"133103",
	"393560",
}

func (a Akamai) GetIPRanges() []*IPRange {
	ranges := make([]*IPRange, 0)
	for _, asn := range AkamaiASNs {
		log.Info("[Akamai] - Using ranges from ASN list (AS%s)", asn)
		_ranges := getRangesForAsn(asn)
		ranges = append(ranges, _ranges...)
		log.Info("[Akamai] - Found %d ranges for AS%s", len(ranges), asn)
	}
	return ranges
}
