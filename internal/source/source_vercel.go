package source

import (
	"escape.tech/cloudfinder/internal/log"
	"escape.tech/cloudfinder/pkg/provider"
)

type Vercel struct{}

// Source: https://networksdb.io/ip-addresses-of/vercel-inc
const vercelRange = "76.76.21.0/24"

func (a Vercel) GetIPRanges() []*IPRange {
	log.Info("Using static Vercel ip range")

	network, cat := ParseCIDR(vercelRange)
	return []*IPRange{
		{
			Network:  network,
			Cat:      cat,
			Provider: provider.Vercel,
		},
	}
}
