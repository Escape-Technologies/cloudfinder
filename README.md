# cloudfinder  

Detect the cloud / hosting provider of a given IP. Fast, static & offline.  
Cloudfinder offers both a cli and a golang package.

## CLI Usage

Installing the binary:
TODO

From url:
TODO

From cloned repository:  
`go run cmd/cli/cli.go cmd/cli/dial.go <domain, ip, url, ...>`

## PKG Usage

Add dependency:
TODO

Use cloudfinder:

```go
package main

import (
 "net"

 "escape.tech/cloudfinder/pkg/cloud"
 "escape.tech/cloudfinder/pkg/provider"
)

func main() {
 r := cloud.NewResolver()
 ip, err := net.LookupIP("escape.tech")
 if err != nil {
  // do something with your error
 }
 p := r.GetProviderForIP(ip)
 switch p {
 case provider.Aws:
  println("Yay got AWS as expected")
 case provider.Unknown:
  println("Seems like we could not find the provider... Thats weird.")
 default:
  println("Mmmmh seems like the provider is wrong...")
 }
}
```

## pre build

`go run cmd/pre-build/pre-build.go`
