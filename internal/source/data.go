package source

import (
	"fmt"
	"net"
	"sync"

	"escape.tech/cloudfinder/pkg/provider"
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

func (ip *IPRange) String() string {
	return ip.Network.String() + fmt.Sprint(ip.Cat) + ip.Provider.String()
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
			rangeLock.Lock()
			ranges = append(ranges, sourceRanges...)
			rangeLock.Unlock()
			wg.Done()
		}(source)
	}
	wg.Wait()
	return ranges
}
