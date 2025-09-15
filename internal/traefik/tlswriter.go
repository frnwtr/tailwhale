package traefik

import (
    "bytes"
    "sort"
)

// TLSCert maps hostnames to certificate/key file paths.
type TLSCert struct {
    CertFile string
    KeyFile  string
}

// TLSConfig is a minimal structure for Traefik's file provider.
type TLSConfig map[string]TLSCert

// MarshalYAML deterministically renders a TLSConfig to a Traefik-compatible YAML snippet.
// We avoid importing a YAML lib to keep dependencies minimal at this stage.
func MarshalYAML(cfg TLSConfig) []byte {
    var hosts []string
    for h := range cfg {
        hosts = append(hosts, h)
    }
    sort.Strings(hosts)
    var b bytes.Buffer
    b.WriteString("tls:\n  certificates:\n")
    for _, h := range hosts {
        c := cfg[h]
        b.WriteString("    - certFile: \"" + c.CertFile + "\"\n")
        b.WriteString("      keyFile: \"" + c.KeyFile + "\"\n")
        b.WriteString("      stores:\n        - default\n")
        b.WriteString("      sans:\n        - \"" + h + "\"\n")
    }
    return b.Bytes()
}

