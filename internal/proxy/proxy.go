package proxy

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sort"
	"strings"
	"titan/internal/core"
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

func StartProxy(errorChannel chan error, container *core.Container) {
	serverConfig := container.ConfigData.Config.Server
	routes, err := buildRoutes(serverConfig.Routes)
	if err != nil {
		errorChannel <- err
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
		httpAddr := fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port)
		log.Printf("Starting HTTP server at %s", httpAddr)
		if err := http.ListenAndServe(httpAddr, httpMux); err != nil {
			errorChannel <- err
		}
	}()

	go func() {
		if serverConfig.SSL.Cert != "" && serverConfig.SSL.Key != "" {
			httpsAddr := fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.SSL.Port)
			log.Printf("Starting HTTPS server at %s", httpsAddr)
			if err := http.ListenAndServeTLS(httpsAddr, serverConfig.SSL.Cert, serverConfig.SSL.Key, httpMux); err != nil {
				errorChannel <- err
			}
		} else {
			errorChannel <- errors.New("TLS configuration missing. Please add valid value and try again")
		}
	}()
}
