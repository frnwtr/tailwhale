package main

import (
    "context"
    "encoding/json"
    "flag"
    "fmt"
    "io"
    "os"
    "time"

    "github.com/frnwtr/tailwhale/internal/core"
    "github.com/frnwtr/tailwhale/internal/dockerx"
    "github.com/frnwtr/tailwhale/internal/fsx"
    "github.com/frnwtr/tailwhale/internal/appconfig"
    traefik "github.com/frnwtr/tailwhale/internal/traefik"
    ts "github.com/frnwtr/tailwhale/internal/tailscale"
)

// Version is set at build time via -ldflags if desired.
var Version = "0.0.0-dev"

var out io.Writer = os.Stdout
var errOut io.Writer = os.Stderr

func usage() {
    fmt.Fprintln(out, "TailWhale CLI")
    fmt.Fprintln(out, "Usage: tailwhale <command> [flags]")
    fmt.Fprintln(out)
    fmt.Fprintln(out, "Commands:")
    fmt.Fprintln(out, "  list        Show exposed services")
    fmt.Fprintln(out, "  sync        Perform a full sync")
    fmt.Fprintln(out, "  watch       Run in daemon/watch mode")
    fmt.Fprintln(out)
    fmt.Fprintln(out, "Flags:")
    fmt.Fprintln(out, "  -h, --help  Show help")
    fmt.Fprintf(out, "  --version   Show version (%s)\n", Version)
}

func run(args []string) int {
    if len(args) == 0 {
        usage()
        return 0
    }

    // Global flags
    switch args[0] {
    case "-h", "--help", "help":
        usage()
        return 0
    case "--version", "version":
        fmt.Fprintln(out, Version)
        return 0
    }

    // Subcommands
    switch args[0] {
    case "list":
        fs := flag.NewFlagSet("list", flag.ContinueOnError)
        fs.SetOutput(errOut)
        jsonOut := fs.Bool("json", false, "output JSON")
        fromFile := fs.String("from-file", "", "load containers from JSON file (for testing)")
        if err := fs.Parse(args[1:]); err != nil {
            return 2
        }
        var provider dockerx.Provider
        if *fromFile != "" {
            provider = &dockerx.FileProvider{Path: *fromFile}
        } else {
            provider = dockerx.NewProvider()
        }
        svcs, err := core.Discover(provider, "host", "tn")
        if err != nil { fmt.Fprintln(errOut, err); return 1 }
        if *jsonOut {
            enc := json.NewEncoder(out)
            enc.SetIndent("", "  ")
            _ = enc.Encode(svcs)
        } else {
            fmt.Fprintf(out, "%d services\n", len(svcs))
            for _, s := range svcs {
                fmt.Fprintf(out, "- %s (%s) %s\n", s.Name, s.ID, s.Host)
            }
        }
        return 0
    case "sync":
        fs := flag.NewFlagSet("sync", flag.ContinueOnError)
        fs.SetOutput(errOut)
        cfgPath := fs.String("config", "", "path to JSON config file")
        host := fs.String("host", "host", "host name for mode A/C")
        tailnet := fs.String("tailnet", "tn", "tailnet name")
        tlsPath := fs.String("tls-path", "traefik/tls.yml", "path to write Traefik TLS yaml")
        certDir := fs.String("cert-dir", "/var/lib/tailwhale/certs", "directory for issued certs (stub)")
        if err := fs.Parse(args[1:]); err != nil {
            return 2
        }
        // Merge config file (if provided) with flags (flags override)
        if *cfgPath != "" {
            if c, err := appconfig.Load(*cfgPath); err == nil {
                if fs.Lookup("host").Value.String() == "host" && c.Host != "" { *host = c.Host }
                if fs.Lookup("tailnet").Value.String() == "tn" && c.Tailnet != "" { *tailnet = c.Tailnet }
                if fs.Lookup("tls-path").Value.String() == "traefik/tls.yml" && c.TLSPath != "" { *tlsPath = c.TLSPath }
                if fs.Lookup("cert-dir").Value.String() == "/var/lib/tailwhale/certs" && c.CertDir != "" { *certDir = c.CertDir }
            }
        }
        orch := core.Orchestrator{Provider: &dockerx.FakeProvider{}, Host: *host, Tailnet: *tailnet, Manager: &ts.FileManager{Dir: *certDir}}
        svcs, tls, err := orch.SyncOnce(context.Background())
        if err != nil { fmt.Fprintln(errOut, err); return 1 }
        fmt.Fprintf(out, "Synced %d services\n", len(svcs))
        data := traefik.MarshalYAML(tls)
        if err := fsx.WriteFileAtomic(*tlsPath, data, 0o644); err != nil {
            fmt.Fprintf(errOut, "failed to write %s: %v\n", *tlsPath, err)
            return 1
        }
        fmt.Fprintf(out, "Wrote %s (%d bytes)\n", *tlsPath, len(data))
        return 0
    case "watch":
        fs := flag.NewFlagSet("watch", flag.ContinueOnError)
        fs.SetOutput(errOut)
        cfgPath := fs.String("config", "", "path to JSON config file")
        host := fs.String("host", "host", "host name for mode A/C")
        tailnet := fs.String("tailnet", "tn", "tailnet name")
        tlsPath := fs.String("tls-path", "traefik/tls.yml", "path to write Traefik TLS yaml")
        certDir := fs.String("cert-dir", "/var/lib/tailwhale/certs", "directory for issued certs (stub)")
        interval := fs.Duration("interval", 10*time.Second, "sync interval (fallback)")
        if err := fs.Parse(args[1:]); err != nil {
            return 2
        }
        if *cfgPath != "" {
            if c, err := appconfig.Load(*cfgPath); err == nil {
                if fs.Lookup("host").Value.String() == "host" && c.Host != "" { *host = c.Host }
                if fs.Lookup("tailnet").Value.String() == "tn" && c.Tailnet != "" { *tailnet = c.Tailnet }
                if fs.Lookup("tls-path").Value.String() == "traefik/tls.yml" && c.TLSPath != "" { *tlsPath = c.TLSPath }
                if fs.Lookup("cert-dir").Value.String() == "/var/lib/tailwhale/certs" && c.CertDir != "" { *certDir = c.CertDir }
            }
        }
        provider := dockerx.NewProvider()
        orch := core.Orchestrator{Provider: provider, Host: *host, Tailnet: *tailnet}
        // Configure tailscale manager and TLS writer
        orch.Manager = &ts.FileManager{Dir: *certDir}
        orch.WriteTLS = func(cfg traefik.TLSConfig) error {
            data := traefik.MarshalYAML(cfg)
            return fsx.WriteFileAtomic(*tlsPath, data, 0o644)
        }
        ctx, cancel := context.WithCancel(context.Background())
        defer cancel()
        fmt.Fprintln(out, "watching for container changes...")
        _ = orch.Watch(ctx, *interval, func(svcs []core.Service, tlsCfg traefik.TLSConfig){
            _ = tlsCfg // already written via WriteTLS; optionally print summary
            fmt.Fprintf(out, "synced %d services\n", len(svcs))
        })
        return 0
    default:
        fmt.Fprintf(errOut, "unknown command: %s\n\n", args[0])
        usage()
        return 2
    }
}

func main() {
    os.Exit(run(os.Args[1:]))
}

type coreTLS = traefik.TLSConfig

func coreYaml(t traefik.TLSConfig) []byte {
    // tiny helper to print YAML-like output without external deps
    type tlsCert struct{ CertFile, KeyFile string }
    type tlsBlock struct{ Certificates []tlsCert }
    // naive: order unspecified here; kept for preview only
    var buf []byte
    buf = append(buf, []byte("tls:\n  certificates:\n")...)
    for _, c := range t {
        buf = append(buf, []byte("    - certFile: \""+c.CertFile+"\"\n      keyFile: \""+c.KeyFile+"\"\n")...)
    }
    _ = tlsBlock{}
    return buf
}
