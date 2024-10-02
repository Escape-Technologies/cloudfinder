package source

import (
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
	ranges := getRangesForAsn(OvhCloudASN)
	log.Info("[Ovh] - Found %d ranges for AS%s", len(ranges), OvhCloudASN)
	return ranges
}
