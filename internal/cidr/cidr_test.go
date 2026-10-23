/* SPDX-License-Identifier: GPL-3.0-only */

package cidr

import (
	"net/netip"
	"strings"
	"testing"
)

func TestSubtractHomeLAN(t *testing.T) {
	in := []netip.Prefix{netip.MustParsePrefix("0.0.0.0/0")}
	ex := []netip.Prefix{netip.MustParsePrefix("172.16.1.0/24")}
	got := SubtractPrefixes(in, ex)
	if len(got) == 0 {
		t.Fatal("expected prefixes")
	}
	addr := netip.MustParseAddr("172.16.1.50")
	for _, p := range got {
		if p.Contains(addr) {
			t.Fatalf("excluded still covered by %s", p)
		}
	}
	if !containsAddr(got, "8.8.8.8") {
		t.Fatal("8.8.8.8 should remain")
	}
}

func containsAddr(prefixes []netip.Prefix, ip string) bool {
	a := netip.MustParseAddr(ip)
	for _, p := range prefixes {
		if p.Contains(a) {
			return true
		}
	}
	return false
}

func TestParseGoDebMeta(t *testing.T) {
	text := "# GoDeb-ExcludedIPs: 172.16.1.0/24\n# GoDeb-Trusted-Endpoint: http://172.16.1.20:8080/api/trusted\n[Interface]\n"
	_, excl, trusted := ParseGoDebMeta(text)
	if excl != "172.16.1.0/24" || !strings.HasPrefix(trusted, "http://") {
		t.Fatalf("excl=%q trusted=%q", excl, trusted)
	}
}
