package source

import (
	"fmt"

	"github.com/Escape-Technologies/cloudfinder/internal/log"
	"github.com/Escape-Technologies/cloudfinder/pkg/provider"
)

type Ovh struct{}

func (a Ovh) GetProvider() provider.Provider {
	return provider.Ovh
}

// Note: there is also AS22598 for OVHTelecom. But this probably doesn't expose hosting services
var OvhCloudASN = "16276"

func (a Ovh) GetIPRanges() []*IPRange {
	log.Info("[Ovh] - Using ranges from ASN list (AS%s)", OvhCloudASN)
	ranges, err := getRangesForAsn(OvhCloudASN)
	if err != nil {
		msg := fmt.Sprintf("[Ovh] - Error getting ranges for AS%s:", OvhCloudASN)
		log.Error(msg, err)
		return []*IPRange{}
	}
	log.Info("[Ovh] - Found %d ranges for AS%s", len(ranges), OvhCloudASN)
	return ranges
}
