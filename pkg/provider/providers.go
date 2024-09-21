//go:generate stringer -type=Provider

// idea: https://chatgpt.com/share/e/66ed8fef-442c-8008-810d-c029aaf42bec
package provider

type Provider int

const (
	Unknown Provider = iota
	Aws
	Alibaba
	Azure
	Cloudflare
	Digitalocean
	Fastly
	Gcp
	Ibm
	Linode
	Oracle
	Ovh
	Scaleway
	Tencent
	Ucloud
	Vercel
)

// func GetAllIPRanges() []*IPRange {
// 	rangeLock := &sync.Mutex{}
// 	wg := &sync.WaitGroup{}
// 	ranges := make([]*IPRange, 0)
// 	for _, source := range sourceRegistry {
// 		wg.Add(1)
// 		go func(s IPRangeSource) {
// 			sourceRanges := s.GetIPRanges()
// 			rangeLock.Lock()
// 			ranges = append(ranges, sourceRanges...)
// 			rangeLock.Unlock()
// 			wg.Done()
// 		}(source)
// 	}
// 	wg.Wait()
// 	return ranges
// }
