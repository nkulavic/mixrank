package mixrank

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Download streams an HTTPS result URL without forwarding the MixRank API key.
// DNS answers are checked at connection time to prevent private-network requests.
func Download(ctx context.Context, resultURL string, dst io.Writer, maxBytes int64) error {
	u, e := url.Parse(strings.TrimSpace(resultURL))
	if e != nil || u.Scheme != "https" || u.Host == "" || u.User != nil {
		return errors.New("download requires an HTTPS result URL")
	}
	if maxBytes <= 0 {
		return errors.New("positive download byte limit required")
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = nil
	transport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		host, port, e := net.SplitHostPort(address)
		if e != nil {
			return nil, errors.New("invalid download host")
		}
		ips, e := net.DefaultResolver.LookupIPAddr(ctx, host)
		if e != nil {
			return nil, errors.New("download DNS lookup failed")
		}
		if len(ips) == 0 {
			return nil, errors.New("download host unresolved")
		}
		for _, ip := range ips {
			if !publicIP(ip.IP) {
				return nil, errors.New("download host resolves to a non-public address")
			}
		}
		return (&net.Dialer{Timeout: 20 * time.Second}).DialContext(ctx, network, net.JoinHostPort(ips[0].IP.String(), port))
	}
	defer transport.CloseIdleConnections()
	cl := &http.Client{Transport: transport, Timeout: 10 * time.Minute, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if req.URL.Scheme != "https" || req.URL.User != nil || len(via) > 5 {
			return errors.New("download redirect refused")
		}
		return nil
	}}
	req, e := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if e != nil {
		return errors.New("invalid download URL")
	}
	res, e := cl.Do(req)
	if e != nil {
		return errors.New("download connection failed")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New("download did not return HTTP 200")
	}
	n, e := io.Copy(dst, io.LimitReader(res.Body, maxBytes+1))
	if e != nil {
		return errors.New("download interrupted; retain partial output")
	}
	if n > maxBytes {
		return errors.New("download exceeded byte limit; retain partial output")
	}
	return nil
}
func publicIP(ip net.IP) bool {
	if !ip.IsGlobalUnicast() || ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsUnspecified() {
		return false
	}
	for _, cidr := range []string{"100.64.0.0/10", "192.0.0.0/24", "192.0.2.0/24", "198.18.0.0/15", "198.51.100.0/24", "203.0.113.0/24", "2001:db8::/32"} {
		_, n, _ := net.ParseCIDR(cidr)
		if n.Contains(ip) {
			return false
		}
	}
	return true
}
