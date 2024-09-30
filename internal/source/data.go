package source

import (
	"fmt"
	"net"
	"sync"

	"github.com/Escape-Technologies/cloudfinder/pkg/provider"
)

type IPCat int

const (
	CatIPv4 IPCat = 4
	CatIPv6 IPCat = 6
)

type IPRange struct {
	Network  *net.IPNet        `json:"n"`
	Cat      IPCat             `json:"c"`
	Provider provider.Provider `json:"p"`
}

func (r *IPRange) String() string {
	return r.Network.String() + fmt.Sprint(r.Cat) + r.Provider.String()
}

type IPRangeSource interface {
	GetIPRanges() []*IPRange
}

var AllSources = []IPRangeSource{
	Alibaba{},
	Aws{},
	Azure{},
	Cloudflare{},
	Digitalocean{},
	Fastly{},
	Gcp{},
	Ibm{},
	Linode{},
	Oracle{},
	Ovh{},
	Scaleway{},
	Tencent{},
	Ucloud{},
	Vercel{},
}

func GetAllIPRanges(sources []IPRangeSource) []*IPRange {
	rangeLock := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	ranges := make([]*IPRange, 0)
	for _, source := range sources {
		wg.Add(1)
		go func(s IPRangeSource) {
			sourceRanges := s.GetIPRanges()
			sourceRanges = dedupRanges(sourceRanges)
			rangeLock.Lock()
			ranges = append(ranges, sourceRanges...)
			rangeLock.Unlock()
			wg.Done()
		}(source)
	}
	wg.Wait()
	return ranges
}

func dedupRanges(ranges []*IPRange) []*IPRange {
	existingMap := make(map[string]interface{})
	dedupedRanges := make([]*IPRange, 0)
	for _, r := range ranges {
		// If already seen, skip it
		if _, ok := existingMap[r.String()]; ok {
			continue
		}
		dedupedRanges = append(dedupedRanges, r)
		existingMap[r.String()] = struct{}{}
	}

	return dedupedRanges
}

func (r *IPRange) Overlaps() {

}
