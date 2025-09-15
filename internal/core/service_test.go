package core

import "testing"

func TestHostnameFor(t *testing.T) {
    got := HostnameFor(ModeA, NameInput{Container: "app", Host: "host1", Tailnet: "tn"})
    if got != "app.host1.tn.ts.net" {
        t.Fatalf("ModeA wrong: %s", got)
    }
    got = HostnameFor(ModeB, NameInput{Container: "app", Tailnet: "tn"})
    if got != "app.tn.ts.net" {
        t.Fatalf("ModeB wrong: %s", got)
    }
    got = HostnameFor(ModeC, NameInput{Host: "host1"})
    if got != "https://host1.ts.net" {
        t.Fatalf("ModeC wrong: %s", got)
    }
}

