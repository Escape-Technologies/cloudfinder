package cloud

import (
	"log/slog"
	"net"

	"escape.tech/cloudfinder/internal/log"
	"escape.tech/cloudfinder/internal/static"
	"escape.tech/cloudfinder/internal/tree"
	"escape.tech/cloudfinder/pkg/provider"
)

type Resolver interface {
	GetProviderForIP(ip net.IP) provider.Provider
	WithLogger(logger *slog.Logger)
}

type resolver struct {
	ipv4Tree tree.Tree
	ipv6Tree tree.Tree
}

func NewResolver() Resolver {
	return &resolver{
		ipv4Tree: static.LoadIPv4Tree(),
		ipv6Tree: static.LoadIPv6Tree(),
	}
}

func (f *resolver) WithLogger(logger *slog.Logger) {
	log.Logger = logger
}

func (f *resolver) GetProviderForIP(ip net.IP) provider.Provider {
	if ipv4 := ip.To4(); ipv4 != nil {
		ipRange := f.ipv4Tree.FindIPRange(ip)
		if ipRange == nil {
			return provider.Unknown
		}
		return ipRange.Provider
	}

	ipRange := f.ipv6Tree.FindIPRange(ip.To16())
	if ipRange == nil {
		return provider.Unknown
	}
	return ipRange.Provider
}
