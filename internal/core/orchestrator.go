package core

import (
    "context"
    "time"

    "github.com/frnwtr/tailwhale/internal/dockerx"
    tcfg "github.com/frnwtr/tailwhale/internal/traefik"
    ts "github.com/frnwtr/tailwhale/internal/tailscale"
)

// Orchestrator ties discovery to config writers and cert managers.
type Orchestrator struct {
    Provider dockerx.Provider
    Host     string
    Tailnet  string
    Manager  ts.Manager
    // Optional write callback to persist TLS config (e.g., to file)
    WriteTLS func(tcfg.TLSConfig) error
}

// SyncOnce discovers services and returns a TLS config view.
func (o Orchestrator) SyncOnce(ctx context.Context) ([]Service, tcfg.TLSConfig, error) {
    svcs, err := Discover(o.Provider, o.Host, o.Tailnet)
    if err != nil { return nil, nil, err }
    tls := make(tcfg.TLSConfig)
    for _, s := range svcs {
        if o.Manager != nil {
            c, err := o.Manager.Ensure(s.Host)
            if err == nil {
                tls[s.Host] = tcfg.TLSCert{CertFile: c.Path, KeyFile: c.KeyPath}
                continue
            }
        }
        // Placeholder fallback paths
        tls[s.Host] = tcfg.TLSCert{CertFile: "/var/lib/tailwhale/certs/"+s.Name+".crt", KeyFile: "/var/lib/tailwhale/certs/"+s.Name+".key"}
    }
    if o.WriteTLS != nil {
        _ = o.WriteTLS(tls)
    }
    _ = ctx // reserved for future timeouts/cancellations
    return svcs, tls, nil
}

// Watch listens for provider events; falls back to periodic sync if events unavailable.
func (o Orchestrator) Watch(ctx context.Context, interval time.Duration, fn func([]Service, tcfg.TLSConfig)) error {
    // Initial sync
    if svcs, tls, err := o.SyncOnce(ctx); err == nil && fn != nil { fn(svcs, tls) }

    w, err := o.Provider.Watch()
    if err == nil && w != nil {
        defer w.Close()
        // Build cache from initial snapshot
        snapshot, _ := o.Provider.List()
        cache := dockerx.NewCache()
        cache.ApplySnapshot(snapshot)
        // Debounce timer
        debounce := time.NewTimer(0)
        if !debounce.Stop() { <-debounce.C }
        const debounceWindow = 1 * time.Second
        // event loop
        for {
            select {
            case <-ctx.Done():
                return ctx.Err()
            default:
            }
            info, ok, _ := w.Next()
            if !ok { break }
            if info.Event == "destroy" {
                cache.Remove(info.ID)
            } else {
                cache.Upsert(info)
            }
            // reset debounce
            if !debounce.Stop() { select { case <-debounce.C: default: } }
            debounce.Reset(debounceWindow)

            // wait for debounce window to fire before syncing
        debLoop:
            for {
                select {
                case <-ctx.Done():
                    return ctx.Err()
                case <-debounce.C:
                    svcs := DiscoverFromInfos(cache.List(), o.Host, o.Tailnet)
                    tls := make(tcfg.TLSConfig)
                    for _, s := range svcs {
                        if o.Manager != nil {
                            if c, err := o.Manager.Ensure(s.Host); err == nil {
                                tls[s.Host] = tcfg.TLSCert{CertFile: c.Path, KeyFile: c.KeyPath}
                                continue
                            }
                        }
                        tls[s.Host] = tcfg.TLSCert{CertFile: "/var/lib/tailwhale/certs/"+s.Name+".crt", KeyFile: "/var/lib/tailwhale/certs/"+s.Name+".key"}
                    }
                    if o.WriteTLS != nil { _ = o.WriteTLS(tls) }
                    if fn != nil { fn(svcs, tls) }
                    break debLoop
                default:
                    // Accumulate more events until debounce fires, using a worker goroutine to avoid blocking
                    type res struct{ i dockerx.Info; ok bool }
                    ch := make(chan res, 1)
                    go func(){ i, ok, _ := w.Next(); ch <- res{i, ok} }()
                    select {
                    case <-ctx.Done():
                        return ctx.Err()
                    case r := <-ch:
                        if !r.ok { break debLoop }
                        if r.i.Event == "destroy" { cache.Remove(r.i.ID) } else { cache.Upsert(r.i) }
                        if !debounce.Stop() { select { case <-debounce.C: default: } }
                        debounce.Reset(debounceWindow)
                        // continue accumulating
                    case <-time.After(25 * time.Millisecond):
                        // no event very briefly; loop continues until timer fires
                    }
                }
            }
        }
    }

    // Fallback ticker
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-ticker.C:
            svcs, tls, err := o.SyncOnce(ctx)
            if err == nil && fn != nil { fn(svcs, tls) }
        }
    }
}
