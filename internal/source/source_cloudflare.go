package source

import (
	"strings"

	"github.com/Escape-Technologies/cloudfinder/internal/log"
	"github.com/Escape-Technologies/cloudfinder/pkg/provider"
)

type Cloudflare struct{}

var cloudflareFileUrls = []string{
	"https://www.cloudflare.com/ips-v4/#",
	"https://www.cloudflare.com/ips-v6/#",
}

func (a Cloudflare) GetProvider() provider.Provider {
	return provider.Cloudflare
}

func (a Cloudflare) GetIPRanges() []*IPRange {
	ranges := make([]*IPRange, 0)
	for _, cloudflareFileURL := range cloudflareFileUrls {
		log.Info("Fetching cloudflare ip ranges from %s", cloudflareFileURL)

		content, err := FileURLToString(cloudflareFileURL)
		if err != nil {
			log.Fatal("Failed to read cloudflare ip ranges", err)
		}
		content = strings.ReplaceAll(content, " ", "")
		ips := strings.Split(content, "\n")
		for _, ip := range ips {
			if ip == "" {
				continue
			}
			network, cat := ParseCIDR(ip)
			ranges = append(ranges, &IPRange{
				Network: network,
				Cat:     cat,
			})
		}
	}
	return ranges
}
