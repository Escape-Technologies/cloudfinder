//go:generate go-enum --marshal --noprefix --nocomments
package provider

/*
ENUM(
Unknown
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
Akamai
)
*/
type Provider int

// re-export ProviderMap
var ProviderMap = _ProviderMap
