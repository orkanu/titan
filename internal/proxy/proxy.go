package proxy

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sort"
	"strings"
	"titan/internal/config"
)

type Route struct {
	Source string
	Target *url.URL
}

// getClientIP extracts the client's IP address from the request
func getClientIP(r *http.Request) string {
	// Get IP address from the remote address (in case there is no X-Forwarded-For)
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func createReverseProxy(target *url.URL, source string) http.HandlerFunc {
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Del("X-Powered-By")
		return nil
	}
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		prefix := strings.TrimRight(source, "/")
		targetBase := strings.TrimRight(target.Path, "/")

		// TODO - the ReplaceAll call works for now but I'd like to improve it as it feels a bit clunky :(
		//   The problem is that after calling "originalDirector(req)" /identity appears twice in the req.Path
		// Prevent double path segment
		suffix := strings.ReplaceAll(req.URL.Path, prefix, "")
		// suffix := strings.TrimPrefix(req.URL.Path, prefix)
		// suffix = strings.TrimPrefix(suffix, prefix)

		if strings.HasPrefix(suffix, "/") {
			suffix = suffix[1:]
		}

		if targetBase != "" {
			req.URL.Path = "/" + strings.TrimRight(targetBase, "/") + "/" + suffix
		} else {
			req.URL.Path = "/" + suffix
		}

		req.URL.Path = strings.ReplaceAll(req.URL.Path, "//", "/")

		// Set the "X-Forwarded-Host" header to the original host
		req.Header.Set("X-Forwarded-Host", req.Host)
		// Keep the original client IP in X-Forwarded-For
		clientIP := getClientIP(req)
		if prior, ok := req.Header["X-Forwarded-For"]; ok {
			clientIP = prior[0] + ", " + clientIP
		}
		req.Header.Set("X-Forwarded-For", clientIP)

		req.Host = target.Host

		// Add X-Forwarded-Proto to indicate if the original request was HTTP or HTTPS
		if req.TLS != nil {
			req.Header.Set("X-Forwarded-Proto", "https")
		} else {
			req.Header.Set("X-Forwarded-Proto", "http")
		}
	}
	return proxy.ServeHTTP
}

func buildRoutes(proxyConfig map[string]struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}) ([]Route, error) {
	routes := make([]Route, 0, len(proxyConfig))
	for _, cfg := range proxyConfig {
		targetURL, err := url.Parse(cfg.Target)
		if err != nil {
			return nil, fmt.Errorf("invalid target URL %q: %w", cfg.Target, err)
		}
		routes = append(routes, Route{
			Source: cfg.Source,
			Target: targetURL,
		})
	}
	// Sort by descending Source length to ensure longest match wins
	sort.Slice(routes, func(i, j int) bool {
		return len(routes[i].Source) > len(routes[j].Source)
	})
	return routes, nil
}

func StartProxy(cfg *config.Config) {
	routes, err := buildRoutes(cfg.Server.Routes)
	if err != nil {
		fmt.Printf("Failed to build routes: %v", err)
		os.Exit(1)
	}

	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		for _, route := range routes {
			if strings.HasPrefix(path, route.Source) {
				createReverseProxy(route.Target, route.Source)(w, r)
				return
			}
		}
		http.NotFound(w, r)
	})

	go func() {
		httpAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		log.Printf("Starting HTTP server at %s", httpAddr)
		if err := http.ListenAndServe(httpAddr, httpMux); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	if cfg.Server.SSL.Cert != "" && cfg.Server.SSL.Key != "" {
		httpsAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.SSL.Port)
		log.Printf("Starting HTTPS server at %s", httpsAddr)
		if err := http.ListenAndServeTLS(httpsAddr, cfg.Server.SSL.Cert, cfg.Server.SSL.Key, httpMux); err != nil {
			log.Fatalf("HTTPS server failed: %v", err)
		}
	} else {
		fmt.Println("TLS configuration missing. Please add valid value and try again")
		os.Exit(1)
	}
}
