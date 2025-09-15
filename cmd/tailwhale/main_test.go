package main

import (
    "bytes"
    "strings"
    "testing"
)

func TestHelp(t *testing.T) {
    var buf bytes.Buffer
    out, errOut = &buf, &buf
    t.Cleanup(func() { out, errOut = nil, nil })

    code := run([]string{"--help"})
    if code != 0 {
        t.Fatalf("expected exit 0, got %d", code)
    }
    if !strings.Contains(buf.String(), "Usage: tailwhale") {
        t.Fatalf("expected usage output, got: %s", buf.String())
    }
}

func TestUnknownCommand(t *testing.T) {
    var buf bytes.Buffer
    out, errOut = &buf, &buf
    t.Cleanup(func() { out, errOut = nil, nil })

    code := run([]string{"frobnicate"})
    if code != 2 {
        t.Fatalf("expected exit 2, got %d", code)
    }
    s := buf.String()
    if !strings.Contains(s, "unknown command") || !strings.Contains(s, "Usage: tailwhale") {
        t.Fatalf("expected error and usage, got: %s", s)
    }
}

func TestList(t *testing.T) {
    var buf bytes.Buffer
    out, errOut = &buf, &buf
    t.Cleanup(func() { out, errOut = nil, nil })

    code := run([]string{"list"})
    if code != 0 {
        t.Fatalf("expected exit 0, got %d", code)
    }
    s := buf.String()
    if !strings.Contains(s, "1 services") || !strings.Contains(s, "demo.host.tn.ts.net") {
        t.Fatalf("unexpected output: %s", s)
    }
}
