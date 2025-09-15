package tailscale

import (
    "encoding/json"
    "io"
)

// Status contains a minimal subset of tailscale status --json we care about.
type Status struct {
    MagicDNSEnabled bool `json:"MagicDNSEnabled"`
}

// ParseStatus parses `tailscale status --json` output (or a subset) into Status.
func ParseStatus(r io.Reader) (Status, error) {
    var s Status
    dec := json.NewDecoder(r)
    err := dec.Decode(&s)
    return s, err
}

