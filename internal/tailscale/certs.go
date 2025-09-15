package tailscale

import "time"

// Cert describes a certificate on disk.
type Cert struct {
    Host    string
    Path    string
    KeyPath string
    Expiry  time.Time
}

// Manager issues and renews certificates.
type Manager interface {
    Ensure(host string) (Cert, error)
    Renew(host string) (Cert, error)
}

