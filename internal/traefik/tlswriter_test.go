package traefik

import (
    "strings"
    "testing"
)

func TestMarshalYAMLDeterministic(t *testing.T) {
    cfg := TLSConfig{
        "b.example": {CertFile: "/certs/b.crt", KeyFile: "/certs/b.key"},
        "a.example": {CertFile: "/certs/a.crt", KeyFile: "/certs/a.key"},
    }
    out1 := string(MarshalYAML(cfg))
    out2 := string(MarshalYAML(cfg))
    if out1 != out2 {
        t.Fatalf("non-deterministic output")
    }
    if !strings.Contains(out1, "a.example") || !strings.Contains(out1, "b.example") {
        t.Fatalf("missing hosts: %s", out1)
    }
    if strings.Index(out1, "a.example") > strings.Index(out1, "b.example") {
        t.Fatalf("expected a.example before b.example: %s", out1)
    }
}

