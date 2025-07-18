package source

import (
	"fmt"

	"github.com/Escape-Technologies/cloudfinder/internal/log"
	"github.com/Escape-Technologies/cloudfinder/pkg/provider"
)

type Scaleway struct{}

func (a Scaleway) GetProvider() provider.Provider {
	return provider.Scaleway
}

var ScalewayASN = "12876"

func (a Scaleway) GetIPRanges() []*IPRange {
	log.Info("[Scaleway] - Using ranges from ASN list (AS%s)", ScalewayASN)
	ranges, err := getRangesForAsn(ScalewayASN)
	if err != nil {
		msg := fmt.Sprintf("[Scaleway] - Error getting ranges for AS%s:", ScalewayASN)
		log.Error(msg, err)
		return []*IPRange{}
	}
	log.Info("[Scaleway] - Found %d ranges for AS%s", len(ranges), ScalewayASN)
	return ranges
}
