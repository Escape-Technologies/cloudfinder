package source

import (
	"github.com/Escape-Technologies/cloudfinder/internal/log"
	"github.com/Escape-Technologies/cloudfinder/pkg/provider"
)

type Ucloud struct{}

func (a Ucloud) GetProvider() provider.Provider {
	return provider.Ucloud
}

// UCLOUD INFORMATION TECHNOLOGY (HK) LIMITED
var UcloudASN = "135377"

func (a Ucloud) GetIPRanges() []*IPRange {
	log.Info("[Ucloud] - Using ranges from ASN list (AS%s)", UcloudASN)
	ranges := getRangesForAsn(UcloudASN)
	log.Info("[Ucloud] - Found %d ranges for AS%s", len(ranges), UcloudASN)
	return ranges
}
