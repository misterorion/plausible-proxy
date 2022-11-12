package proxy

import (
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

var targetURL = parseURL("https://plausible.io")

var mappings = map[string]string{
	"/api/event":       "/api/event",
	"/js/plausible.js": "/js/plausible.js",
	"/js/script.js":    "/js/plausible.js",
}

// init registers the function target and handler
func init() {
	functions.HTTP("plausibleProxy", plausibleProxy)
}

// plausibleProxy sends request data to plausible.io
func plausibleProxy(rw http.ResponseWriter, req *http.Request) {
	canonicalPath, err := canonicalizePath(req.URL.Path)
	if err != nil {
		http.Error(rw, "Unsupported path", http.StatusNotFound)
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

	proxy.ServeHTTP(rw, req)
}

// canonicalizePath returns a path string mapped to an input path string
func canonicalizePath(path string) (string, error) {
	for k, v := range mappings {
		if strings.HasSuffix(path, k) {
			return v, nil
		}
	}
	return "", errors.New("unsupported path")
}

// parseURL returns a URL structure and panics if there is an error
func parseURL(u string) *url.URL {
	parsed, err := url.Parse(u)
	if err != nil {
		log.Printf("failed to parse URL %s", u)
		panic(err)
	}
	return parsed
}
