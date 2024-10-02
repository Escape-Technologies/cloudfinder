package source

import (
	"github.com/Escape-Technologies/cloudfinder/internal/log"
	"github.com/Escape-Technologies/cloudfinder/pkg/provider"
)

type Tencent struct{}

func (a Tencent) GetProvider() provider.Provider {
	return provider.Tencent
}

var TencentASNs = []string{
	// Tencent Cloud
	"132591",
	// Tencent Global
	"132203",
	// Tencent-CN
	"45090",
}

func (a Tencent) GetIPRanges() []*IPRange {
	ranges := make([]*IPRange, 0)
	for _, asn := range TencentASNs {
		log.Info("[Tencent] - Using ranges from ASN list (AS%s)", asn)
		_ranges := getRangesForAsn(asn)
		ranges = append(ranges, _ranges...)
		log.Info("[Tencent] - Found %d ranges for AS%s", len(ranges), asn)
	}
	return ranges
}
