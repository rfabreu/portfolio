package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func keyReq(remoteAddr string, headers map[string][]string) *http.Request {
	r := httptest.NewRequest(http.MethodPost, "/api/chat", nil)
	r.RemoteAddr = remoteAddr
	for k, vs := range headers {
		for _, v := range vs {
			r.Header.Add(k, v)
		}
	}
	return r
}

func TestClientKey_TrustOnUsesFlyClientIP(t *testing.T) {
	a := clientKey(keyReq("172.16.0.1:1000", map[string][]string{"Fly-Client-IP": {"203.0.113.1"}}), true)
	b := clientKey(keyReq("172.16.0.1:1000", map[string][]string{"Fly-Client-IP": {"203.0.113.2"}}), true)
	if a == b {
		t.Errorf("same RemoteAddr with different Fly-Client-IP must differ, both %q", a)
	}
	if a != "203.0.113.1" {
		t.Errorf("key = %q, want 203.0.113.1", a)
	}
}

func TestClientKey_TrustOffIgnoresFlyClientIP(t *testing.T) {
	a := clientKey(keyReq("172.16.0.1:1000", map[string][]string{"Fly-Client-IP": {"203.0.113.1"}}), false)
	b := clientKey(keyReq("172.16.0.1:1000", map[string][]string{"Fly-Client-IP": {"203.0.113.2"}}), false)
	if a != b || a != "172.16.0.1" {
		t.Errorf("trust off must key on RemoteAddr: got %q and %q", a, b)
	}
}

func TestClientKey_InvalidFlyClientIPFallsBack(t *testing.T) {
	cases := map[string][]string{
		"list":        {"1.2.3.4, 5.6.7.8"},
		"garbage":     {"garbage"},
		"empty":       {""},
		"host:port":   {"1.2.3.4:80"},
		"two headers": {"1.2.3.4", "5.6.7.8"},
		"cidr":        {"1.2.3.4/24"},
	}
	for name, vals := range cases {
		got := clientKey(keyReq("172.16.0.1:1000", map[string][]string{"Fly-Client-IP": vals}), true)
		if got != "172.16.0.1" {
			t.Errorf("%s: key = %q, want fallback 172.16.0.1", name, got)
		}
	}
}

func TestClientKey_MissingFlyClientIPFallsBack(t *testing.T) {
	if got := clientKey(keyReq("172.16.0.1:1000", nil), true); got != "172.16.0.1" {
		t.Errorf("key = %q, want 172.16.0.1", got)
	}
}

func TestClientKey_NeverReadsForwardedHeaders(t *testing.T) {
	h := map[string][]string{
		"X-Forwarded-For": {"203.0.113.9"},
		"X-Real-Ip":       {"203.0.113.9"},
		"Forwarded":       {"for=203.0.113.9"},
	}
	for _, trust := range []bool{true, false} {
		if got := clientKey(keyReq("172.16.0.1:1000", h), trust); got != "172.16.0.1" {
			t.Errorf("trust=%v: key = %q, want 172.16.0.1", trust, got)
		}
	}
}

func TestClientKey_IPv6SharesSlash64(t *testing.T) {
	k1 := clientKey(keyReq("[2001:db8::1]:1000", nil), false)
	k2 := clientKey(keyReq("[2001:db8::2]:1000", nil), false)
	k3 := clientKey(keyReq("[2001:db8:0:1::1]:1000", nil), false)
	if k1 != k2 {
		t.Errorf("same /64 must share a key: %q vs %q", k1, k2)
	}
	if k1 == k3 {
		t.Errorf("different /64 must not share a key: %q", k1)
	}
	if k1 != "2001:db8::/64" {
		t.Errorf("key = %q, want 2001:db8::/64", k1)
	}
}

func TestClientKey_IPv6ViaFlyClientIPAndZone(t *testing.T) {
	a := clientKey(keyReq("172.16.0.1:1", map[string][]string{"Fly-Client-IP": {"2001:db8::1"}}), true)
	b := clientKey(keyReq("172.16.0.1:1", map[string][]string{"Fly-Client-IP": {"2001:db8::ffff"}}), true)
	if a != b {
		t.Errorf("same /64 via header must share: %q vs %q", a, b)
	}
	z1 := clientKey(keyReq("172.16.0.1:1", map[string][]string{"Fly-Client-IP": {"fe80::1%eth0"}}), true)
	z2 := clientKey(keyReq("172.16.0.1:1", map[string][]string{"Fly-Client-IP": {"fe80::1%eth1"}}), true)
	if z1 != z2 {
		t.Errorf("zone must be dropped: %q vs %q", z1, z2)
	}
}

func TestClientKey_IPv4MappedEqualsIPv4(t *testing.T) {
	a := clientKey(keyReq("[::ffff:1.2.3.4]:1000", nil), false)
	b := clientKey(keyReq("1.2.3.4:1000", nil), false)
	if a != b || a != "1.2.3.4" {
		t.Errorf("mapped %q vs plain %q, want both 1.2.3.4", a, b)
	}
	c := clientKey(keyReq("172.16.0.1:1", map[string][]string{"Fly-Client-IP": {"::ffff:1.2.3.4"}}), true)
	if c != "1.2.3.4" {
		t.Errorf("mapped via header = %q, want 1.2.3.4", c)
	}
}

func TestClientKey_RemoteAddrWithoutPort(t *testing.T) {
	if got := clientKey(keyReq("1.2.3.4", nil), false); got != "1.2.3.4" {
		t.Errorf("key = %q, want 1.2.3.4", got)
	}
}
