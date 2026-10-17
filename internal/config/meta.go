package config

import "github.com/ZeroDos1988/godeb-wireguard-arch/internal/cidr"

// TunnelMeta is GoDeb metadata embedded as conf comments.
type TunnelMeta struct {
	AllowedIPsBase   string
	ExcludedIPs      string
	TrustedEndpoint  string
}

func ParseMeta(text string) TunnelMeta {
	base, excl, trusted := cidr.ParseGoDebMeta(text)
	return TunnelMeta{
		AllowedIPsBase:  base,
		ExcludedIPs:     excl,
		TrustedEndpoint: trusted,
	}
}
