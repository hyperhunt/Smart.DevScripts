package wkd

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"golang.org/x/crypto/openpgp"
)

// Discover retrieves keys associated to an email address.
func Discover(addr string) ([]*openpgp.Entity, error) {
	local, domain, err := splitAddress(strings.ToLower(addr))
	if err != nil {
		return nil, err
	}

	_, addrs, err := net.LookupSRV("openpgpkey", "tcp", domain)
	if dnsErr, ok := err.(*net.DNSError); ok {
		if dnsErr.IsTemporary {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	if len(addrs) > 0 {
		addr := addrs[0]
		if addr.Target == domain || strings.HasSuffix(addr.Target, "."+domain) {
			domain = fmt.Sprintf("%v:%v", addr.Target, addr.Port)
		}
	}

	url := "https://" + domain + Base + "/hu/" + hashLocal(local)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return openpgp.ReadKeyRing(resp.Body)
}
