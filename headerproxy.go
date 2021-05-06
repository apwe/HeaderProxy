package headerproxy

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Service struct
type Service struct {
	Destination string `yaml:"Destination"`
}

// Tenant struct
type Tenant struct {
	Services map[string]Service `yaml:"Services"`
}

type Header struct {
	TenantHeader string `yaml:"TenantHeader"`
	ServiceHeader string `yaml:"ServiceHeader"`
}

// Config struct holds configuration to be passed to the plugin
type Config struct {
	Tenants map[string]Tenant // `yaml:"Tenants"`
	Headers []Header
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		// Tenants: []Tenant{},
		Tenants: map[string]Tenant{},
		Headers: []Header{},
	}
}

// HostReplacement holds the necessary components of a Traefik plugin
type HostReplacement struct {
	tenants map[string]Tenant
	headers []Header
	next    http.Handler
	name    string
}

// New instantiates and returns the required components used to handle a HTTP request
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &HostReplacement{
		tenants: config.Tenants,
		headers: config.Headers,
		next:    next,
		name:    name,
	}, nil
}

// Iterate over headers (if provided in the configuration) to match the ones that come in the request &
// match the headers & send a request to a specified tenant and service configuration
func (u *HostReplacement) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, header := range u.headers {
		tenant, TenantHeaderExist := u.tenants[req.Header.Get(header.TenantHeader)]
		if TenantHeaderExist {
			service, ServiceHeaderExist := tenant.Services[req.Header.Get(header.ServiceHeader)]
			if ServiceHeaderExist {
				ProxyReq(rw, req, service)
			}
		} else {
			u.next.ServeHTTP(rw, req)
		}
	}
}

// Send the request
func ProxyReq(rw http.ResponseWriter, req *http.Request, service Service) {
	director := func(req *http.Request) {
		address := service.Destination
		dest, _ := url.Parse(address)
		req.Header.Add("X-Origin-Host", dest.Host)
		req.URL.Scheme = dest.Scheme
		req.URL.Host = dest.Host
	}

	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(rw, req)
}
