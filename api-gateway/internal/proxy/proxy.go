package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Handler forwards to target, mapping stripPrefix to rewritePrefix.
// e.g. stripPrefix "/api/products", rewritePrefix "/products"
func Handler(target *url.URL, stripPrefix, rewritePrefix string) http.Handler {
	proxy := httputil.NewSingleHostReverseProxy(target)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = target.Host
		path := req.URL.Path
		if strings.HasPrefix(path, stripPrefix) {
			suffix := strings.TrimPrefix(path, stripPrefix)
			if suffix == "" || suffix == "/" {
				path = rewritePrefix
			} else {
				path = rewritePrefix + suffix
			}
		}
		req.URL.Path = path
	}
	return proxy
}

func MustParseURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic("invalid proxy URL: " + raw)
	}
	return u
}
