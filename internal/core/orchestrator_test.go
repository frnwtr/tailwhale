package core

import (
    "context"
    "path/filepath"
    "testing"

    "github.com/frnwtr/tailwhale/internal/dockerx"
    ts "github.com/frnwtr/tailwhale/internal/tailscale"
)

func TestOrchestratorSyncUsesManagerPaths(t *testing.T){
    p := &dockerx.FakeProvider{Items: []dockerx.Info{{ID:"1", Name:"app1", Labels: map[string]string{LabelEnable:"true", LabelMode:"A"}}}}
    dir := t.TempDir()
    o := Orchestrator{Provider: p, Host: "host1", Tailnet: "tn", Manager: &ts.FileManager{Dir: dir}}
    _, tls, err := o.SyncOnce(context.Background())
    if err != nil { t.Fatal(err) }
    cert := tls["app1.host1.tn.ts.net"]
    if cert.CertFile != filepath.Join(dir, "app1.host1.tn.ts.net.crt") {
        t.Fatalf("unexpected cert path: %s", cert.CertFile)
    }
}

