package core

import (
    "context"
    "time"

    "github.com/frnwtr/tailwhale/internal/dockerx"
    tcfg "github.com/frnwtr/tailwhale/internal/traefik"
)

// Orchestrator ties discovery to config writers and cert managers.
type Orchestrator struct {
    Provider dockerx.Provider
    Host     string
    Tailnet  string
}

// SyncOnce discovers services and returns a TLS config view.
func (o Orchestrator) SyncOnce(ctx context.Context) ([]Service, tcfg.TLSConfig, error) {
    svcs, err := Discover(o.Provider, o.Host, o.Tailnet)
    if err != nil { return nil, nil, err }
    tls := make(tcfg.TLSConfig)
    for _, s := range svcs {
        // Placeholder: paths are deterministic placeholders
        tls[s.Host] = tcfg.TLSCert{CertFile: "/var/lib/tailwhale/certs/"+s.Name+".crt", KeyFile: "/var/lib/tailwhale/certs/"+s.Name+".key"}
    }
    _ = ctx // reserved for future timeouts/cancellations
    return svcs, tls, nil
}

// Watch is a simple loop that periodically syncs. Real impl will use events.
func (o Orchestrator) Watch(ctx context.Context, interval time.Duration, fn func([]Service, tcfg.TLSConfig)) error {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-ticker.C:
            svcs, tls, err := o.SyncOnce(ctx)
            if err == nil && fn != nil {
                fn(svcs, tls)
            }
        }
    }
}

