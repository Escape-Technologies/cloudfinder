package source

import (
	"github.com/Escape-Technologies/cloudfinder/internal/log"
	"github.com/Escape-Technologies/cloudfinder/pkg/provider"
)

type Vercel struct{}

func (a Vercel) GetProvider() provider.Provider {
	return provider.Vercel
}

// TODO: Get this range dynamically
// NOTE: this is pretty flaky, as vercel is just a wrapper of AWS. Maybe this provider should be removed
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
