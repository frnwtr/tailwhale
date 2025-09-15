package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io"
    "os"

    "github.com/frnwtr/tailwhale/internal/core"
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
        if err := fs.Parse(args[1:]); err != nil {
            return 2
        }
        // Skeleton demo data using core naming
        demo := []core.Service{{
            ID:      "demo-1",
            Name:    "demo",
            Exposed: false,
            Mode:    core.ModeA,
            Host:    core.HostnameFor(core.ModeA, core.NameInput{Container: "demo", Host: "host", Tailnet: "tn"}),
        }}
        if *jsonOut {
            enc := json.NewEncoder(out)
            enc.SetIndent("", "  ")
            _ = enc.Encode(demo)
        } else {
            fmt.Fprintf(out, "%d services\n", len(demo))
            for _, s := range demo {
                fmt.Fprintf(out, "- %s (%s) %s\n", s.Name, s.ID, s.Host)
            }
        }
        return 0
    case "sync":
        fs := flag.NewFlagSet("sync", flag.ContinueOnError)
        fs.SetOutput(errOut)
        if err := fs.Parse(args[1:]); err != nil {
            return 2
        }
        fmt.Fprintln(out, "Sync not implemented yet")
        return 0
    case "watch":
        fs := flag.NewFlagSet("watch", flag.ContinueOnError)
        fs.SetOutput(errOut)
        if err := fs.Parse(args[1:]); err != nil {
            return 2
        }
        fmt.Fprintln(out, "Watch not implemented yet")
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
