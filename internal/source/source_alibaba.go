package source

import (
	"fmt"

	"github.com/Escape-Technologies/cloudfinder/internal/log"
	"github.com/Escape-Technologies/cloudfinder/pkg/provider"
)

type Alibaba struct{}

var AlibabaASNs = []string{
	// Alibaba cloud
	"24429",
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
		_ranges, err := getRangesForAsn(asn)
		if err != nil {
			msg := fmt.Sprintf("[Alibaba] - Error getting ranges for AS%s:", asn)
			log.Error(msg, err)
			continue
		}
		ranges = append(ranges, _ranges...)
		log.Info("[Alibaba] - Found %d ranges for AS%s", len(ranges), asn)
	}
	return ranges
}
