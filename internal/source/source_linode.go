package source

import (
	"github.com/Escape-Technologies/cloudfinder/internal/log"
	"github.com/Escape-Technologies/cloudfinder/pkg/provider"
)

type Linode struct{}

func (a Linode) GetProvider() provider.Provider {
	return provider.Linode
}

var LinodeASNs = []string{
	// Linode AS63949
	"63949",
	// Linode CorpNet
	"48337",
}

func (a Linode) GetIPRanges() []*IPRange {
	ranges := make([]*IPRange, 0)
	for _, asn := range LinodeASNs {
		log.Info("[Linode] - Using ranges from ASN list (AS%s)", asn)
		_ranges := getRangesForAsn(asn)
		ranges = append(ranges, _ranges...)
		log.Info("[Linode] - Found %d ranges for AS%s", len(ranges), asn)
	}
	return ranges
}
