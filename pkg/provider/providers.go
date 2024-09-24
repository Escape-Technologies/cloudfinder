//go:generate stringer -type=Provider
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
