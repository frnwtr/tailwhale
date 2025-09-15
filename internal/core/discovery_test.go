package core

import (
    "testing"

    "github.com/frnwtr/tailwhale/internal/dockerx"
)

func TestDiscoverLabels(t *testing.T){
    p := &dockerx.FakeProvider{Items: []dockerx.Info{
        {ID:"1", Name:"app1", Labels: map[string]string{LabelEnable:"true", LabelMode:"A"}},
        {ID:"2", Name:"app2", Labels: map[string]string{LabelEnable:"false"}},
        {ID:"3", Name:"app3", Labels: map[string]string{LabelEnable:"true", LabelMode:"B", LabelHost:"custom.tn.ts.net"}},
    }}
    svcs, err := Discover(p, "host1", "tn")
    if err != nil { t.Fatal(err) }
    if len(svcs) != 2 { t.Fatalf("expected 2, got %d", len(svcs)) }
    if svcs[0].Name != "app1" { t.Fatalf("sorted by name, got %s", svcs[0].Name) }
    if svcs[0].Host != "app1.host1.tn.ts.net" { t.Fatalf("wrong host: %s", svcs[0].Host) }
    if svcs[1].Name != "app3" || svcs[1].Host != "custom.tn.ts.net" { t.Fatalf("alias not applied: %+v", svcs[1]) }
}

