package appconfig

import (
    "encoding/json"
    "os"
)

// Config holds runtime settings. Flags override file values.
type Config struct {
    Host    string `json:"host"`
    Tailnet string `json:"tailnet"`
    TLSPath string `json:"tlsPath"`
    CertDir string `json:"certDir"`
}

// Load reads a JSON config file. If path is empty, returns zero Config.
func Load(path string) (Config, error) {
    var c Config
    if path == "" { return c, nil }
    b, err := os.ReadFile(path)
    if err != nil { return c, err }
    if err := json.Unmarshal(b, &c); err != nil { return c, err }
    return c, nil
}

