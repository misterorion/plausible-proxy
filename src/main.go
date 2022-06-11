// Credit to https://github.com/mtlynch/plausible-proxy for most of this.
package proxy

import (
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var targetURL = parseURL("https://plausible.io")

func ProxyPlausible(w http.ResponseWriter, r *http.Request) {
	canonicalPath, err := canonicalizePath(r.URL.Path)
	if err != nil {
		http.Error(w, "Unsupported path", http.StatusNotFound)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.Director = func(req *http.Request) {
		req.URL = targetURL
		req.URL.Path = canonicalPath
		req.Host = targetURL.Host

		req.Header.Add("X-Forwarded-Proto", req.Proto)
		req.Header.Add("X-Forwarded-Host", req.Host)
	}

	proxy.ServeHTTP(w, r)
}

func canonicalizePath(path string) (string, error) {
	mappings := map[string]string{
		"/api/event":       "/api/event",
		"/js/plausible.js": "/js/plausible.js",
		"/js/script.js":    "/js/plausible.js",
	}
	for k, v := range mappings {
		if strings.HasSuffix(path, k) {
			return v, nil
		}
	}
	return "", errors.New("unsupported path")
}

func parseURL(u string) *url.URL {
	parsed, err := url.Parse(u)
	if err != nil {
		log.Printf("failed to parse URL %s", u)
		panic(err)
	}
	return parsed
}
