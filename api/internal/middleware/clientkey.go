package middleware

import (
	"net"
	"net/http"
	"net/netip"
)

// flyClientIPHeader is set by Fly Proxy to the client IP it saw.
// https://fly.io/docs/networking/request-headers/
const flyClientIPHeader = "Fly-Client-IP"

// clientKey returns the rate-limit bucket key for a request.
//
// With trustFly it uses Fly-Client-IP when that header holds exactly one valid
// IP address. In every other case it uses the host of r.RemoteAddr.
// X-Forwarded-For and similar headers are never read: Fly documents them as
// spoofable. IPv6 addresses map to their /64 prefix so one host cannot rotate
// through its own subnet to get a fresh bucket.
func clientKey(r *http.Request, trustFly bool) string {
	if trustFly {
		if vals := r.Header.Values(flyClientIPHeader); len(vals) == 1 {
			if addr, err := netip.ParseAddr(vals[0]); err == nil {
				return normalizeAddr(addr)
			}
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	addr, err := netip.ParseAddr(host)
	if err != nil {
		return host
	}
	return normalizeAddr(addr)
}

func normalizeAddr(addr netip.Addr) string {
	addr = addr.Unmap().WithZone("")
	if addr.Is6() {
		return netip.PrefixFrom(addr, 64).Masked().String()
	}
	return addr.String()
}
