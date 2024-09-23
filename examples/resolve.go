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
