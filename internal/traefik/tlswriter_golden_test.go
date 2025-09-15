package traefik

import (
    "testing"
)

func TestMarshalYAMLGolden(t *testing.T){
    cfg := TLSConfig{
        "b.example": {CertFile: "/certs/b.crt", KeyFile: "/certs/b.key"},
        "a.example": {CertFile: "/certs/a.crt", KeyFile: "/certs/a.key"},
    }
    got := string(MarshalYAML(cfg))
    want := "tls:\n  certificates:\n" +
        "    - certFile: \"/certs/a.crt\"\n" +
        "      keyFile: \"/certs/a.key\"\n" +
        "      stores:\n        - default\n" +
        "      sans:\n        - \"a.example\"\n" +
        "    - certFile: \"/certs/b.crt\"\n" +
        "      keyFile: \"/certs/b.key\"\n" +
        "      stores:\n        - default\n" +
        "      sans:\n        - \"b.example\"\n"
    if got != want {
        t.Fatalf("unexpected YAML\n--- got ---\n%q\n--- want ---\n%q", got, want)
    }
}

