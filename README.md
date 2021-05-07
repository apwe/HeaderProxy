# Header Proxy plugin for traefik

This plugin proxies requests to desired destination (outside of traefik) based on headers provided in the configuration.

## Dev `traefik.yml` configuration file for traefik

```yml
pilot:
  token: [REDACTED]

experimental:
  devPlugin:
    goPath: /home/user/go
    moduleName: github.com/apwe/headerproxy

entryPoints:
  http:
    address: ":8000"
    forwardedHeaders:
      insecure: true

providers:
  file:
    filename: rules-headerproxy.yaml 
```

## How to dev
```bash
$ docker run -d --network host containous/whoami -port 5000
# traefik --config-file traefik.yml
```
## How to use

### Headers

Provide header names:
- 'TenantHeader'  : a header controlling the tenant that the requests will be proxied to 
- 'ServiceHeader' : a header controlling a service from the tenant that the request will be proxied to

```yml
          Headers:
            TenantHeader: 'Tenant'
            ServiceHeader: 'Service'
```

So, in order to match the configuration, you need to hit traefik with a request containing 'Tenant: tenant-name' and 'Service: service-name' headers.


### Tenants

Provide tenant names as the keys - these are the values of your 'TenantHeader: $$$%%#^$' header passed to traefik:

```yml
          Tenants:
            tenant-name-a:
...
            tenant-name-b:
...
```

### Services

Provide service names as the keys - these are the values of your 'ServiceHeader: &$^#$&#$' header passed to traefik

```yml
            tenant-name-a:
              Services:
                ServiceA:
...
                ServiceB: 
...
                ServiceC: 
...
            tenant-name-b:
              Services:
                ServiceA: 
...
                ServiceB:
...
                ServiceC: 
...
```

### Destination

Provide value to the 'Destination' key with the url of the service you want to proxy the requests.
urls MUST contain protocol scheme, as the plugin is parsing the url and using the scheme in the communication.

```yml
                ServiceA: 
                  Destination: 'https://example-a.com'
                ServiceB: 
                  Destination: 'https://example-b.com'
                ServiceC: 
                  Destination: 'https://example-c.com'
```

### Example configuration for traefik file provider

```yml
http:
  routers:
    my-router:
      entryPoints:
      - http
      middlewares:
      - headerproxy
      service: service-whoami
      rule: Path(`/whoami`)

  services:
    service-whoami:
      loadBalancer:
        servers:
        - url: http://localhost:5000/
      passHostHeader: false
  middlewares:
    headerproxy:
      plugin:
        dev:
          Headers:
            TenantHeader: 'Tenant'
            ServiceHeader: 'Service'
          Tenants:
            tenant-name-a:
              Services:
                ServiceA: 
                  Destination: 'https://example-a.com'
                ServiceB: 
                  Destination: 'https://example-b.com'
                ServiceC: 
                  Destination: 'https://example-c.com'
            tenant-name-b:
              Services:
                ServiceA: 
                  Destination: 'https://example-a.com'
                ServiceB: 
                  Destination: 'https://example-b.com'
```