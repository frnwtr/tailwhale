package tailscale

import (
    "strings"
    "testing"
)

func TestParseStatus(t *testing.T){
    json := `{"MagicDNSEnabled": true}`
    s, err := ParseStatus(strings.NewReader(json))
    if err != nil { t.Fatal(err) }
    if !s.MagicDNSEnabled { t.Fatal("expected MagicDNSEnabled true") }
}

