# cloudfinder  

Detect the cloud / hosting provider of a given host. Fast, static & offline.  
Cloudfinder offers both a cli and a golang package.
Cloud provider ranges are also tracked and can be found in `./ranges`.

## CLI Usage

### Installation

Run install script: `curl -sSL https://raw.githubusercontent.com/Escape-Technologies/cloudfinder/main/install.sh | sh`
Or directly download a [release](https://github.com/Escape-Technologies/cloudfinder/releases/latest).

### Usage

```bash
cloudfinder [flags] <ip, host, domain, url> <ip, host, domain, url> ...
Flags:
  -debug
        enable debug mode
  -h    print help
  -help
        print help
  -json
        output json
  -raw
        output raw provider string
  -v    print version number
  -version
        print version number
```

## Examples

```bash
cloudfinder escape.tech
[15:06:39.755] INFO: escape.tech (13.39.28.216): Aws
[15:06:39.756] INFO: escape.tech (13.37.196.127): Aws
[15:06:39.756] INFO: escape.tech (13.36.180.15): Aws
```

You can provide multiple inputs:

```bash
cloudfinder escape.tech jobs.escape.tech
[15:31:34.602] INFO: escape.tech (13.39.28.216): Aws
[15:31:34.603] INFO: escape.tech (13.37.196.127): Aws
[15:31:34.603] INFO: escape.tech (13.36.180.15): Aws
[15:31:34.623] INFO: jobs.escape.tech (52.6.1.219): Aws
[15:31:34.623] INFO: jobs.escape.tech (44.212.166.106): Aws
[15:31:34.623] INFO: jobs.escape.tech (52.55.10.55): Aws
```

Or take the input from stdin:

```bash
echo "escape.tech" | cloudfinder 
[15:07:43.573] INFO: escape.tech (13.39.28.216): Aws
[15:07:43.573] INFO: escape.tech (13.37.196.127): Aws
[15:07:43.573] INFO: escape.tech (13.36.180.15): Aws
```

Output can also be raw text:

```bash
cloudfinder --raw escape.tech
escape.tech,13.37.196.127,Aws
escape.tech,13.39.28.216,Aws
escape.tech,13.36.180.15,Aws
```

Or JSON:

```bash
cloudfinder --json escape.tech
{"input":"escape.tech","ip":"13.37.196.127","provider":"Aws"}
{"input":"escape.tech","ip":"13.36.180.15","provider":"Aws"}
{"input":"escape.tech","ip":"13.39.28.216","provider":"Aws"}
```

### Example: using with subfinder

You can pipe the output of external tools into cloudfinder. Here is an example using [subfinder](https://github.com/projectdiscovery/subfinder) to enumerate all subdomains of a given domain, and then finding their cloud providers.

```bash
# run subfinder pipe to cloudfinder and use jq to collect into a single json
subfinder -d "escape.tech" | cloudfinder --json | jq -s '.'

# You'll get errors on stderr for domains that are not exposed and you'll get on stdout:
[
  {
    "input": "www.jobs.escape.tech",
    "ip": "52.6.1.219",
    "provider": "Aws"
  },
  # ...
  {
    "input": "www.docs.escape.tech",
    "ip": "76.76.21.123",
    "provider": "Vercel"
  }
]
```

## Go Package Usage

Add dependency: `go get github.com/Escape-Technologies/cloudfinder@latest`

Use cloudfinder:

```go
package main

import (
 "net"

 "github.com/Escape-Technologies/cloudfinder/pkg/cloud"
 "github.com/Escape-Technologies/cloudfinder/pkg/provider"
)

func main() {
 r := cloud.NewResolver()
 ip, err := net.LookupIP("escape.tech")
 if err != nil {
  // do something with your error
 }
 if len(ip) == 0 {
  // do something with your error
 }
 for _, i := range ip {
  p := r.GetProviderForIP(i)
  switch p {
  case provider.Aws:
   println("Yay got AWS as expected")
  case provider.Unknown:
   println("Seems like we could not find the provider... Thats weird.")
  default:
   println("Mmmmh seems like the provider is wrong...")
  }
 }
}
```
