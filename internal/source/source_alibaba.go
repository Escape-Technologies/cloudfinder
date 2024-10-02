package source

import (
	"github.com/Escape-Technologies/cloudfinder/internal/log"
	"github.com/Escape-Technologies/cloudfinder/pkg/provider"
)

type Alibaba struct{}

var AlibabaASNs = []string{
	// Alibaba AS45102
	"45102",
	// Alibaba (china) AS37963
	"37963",
}

func (a Alibaba) GetProvider() provider.Provider {
	return provider.Alibaba
}

func (a Alibaba) GetIPRanges() []*IPRange {
	ranges := make([]*IPRange, 0)
	for _, asn := range AlibabaASNs {
		log.Info("[Alibaba] - Using ranges from ASN list (AS%s)", asn)
		_ranges := getRangesForAsn(asn)
		ranges = append(ranges, _ranges...)
		log.Info("[Alibaba] - Found %d ranges for AS%s", len(ranges), asn)
	}
	return ranges
}
