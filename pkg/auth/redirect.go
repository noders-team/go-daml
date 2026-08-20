package auth

import (
	"fmt"
	"net"
	"net/http"
)

func DenyPrivateRedirects(req *http.Request, via []*http.Request) error {
	if len(via) >= 10 {
		return fmt.Errorf("stopped after 10 redirects")
	}

	host := req.URL.Hostname()
	ips, err := net.DefaultResolver.LookupIPAddr(req.Context(), host)
	if err != nil {
		return fmt.Errorf("redirect target %q could not be resolved: %w", host, err)
	}
	if len(ips) == 0 {
		return fmt.Errorf("redirect target %q resolved to no addresses", host)
	}

	for _, ip := range ips {
		if isDisallowedRedirectIP(ip.IP) {
			return fmt.Errorf("redirect to internal/private address %q is not allowed", ip.IP)
		}
	}

	return nil
}

func isDisallowedRedirectIP(ip net.IP) bool {
	return ip.IsLoopback() ||
		ip.IsPrivate() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsUnspecified()
}
