package tailscale

import (
    "os"
    "path/filepath"
    "time"
)

// FileManager is a simple Manager that maps hosts to cert/key paths in a dir.
// It can optionally create placeholder files for local testing.
type FileManager struct {
    Dir           string
    CreateOnEnsure bool
}

func (m *FileManager) Ensure(host string) (Cert, error) {
    cert := Cert{
        Host:    host,
        Path:    filepath.Join(m.Dir, host+".crt"),
        KeyPath: filepath.Join(m.Dir, host+".key"),
        Expiry:  time.Now().Add(24 * time.Hour),
    }
    if m.CreateOnEnsure {
        _ = os.MkdirAll(m.Dir, 0o755)
        _ = os.WriteFile(cert.Path, []byte("dummy cert for "+host+"\n"), 0o644)
        _ = os.WriteFile(cert.KeyPath, []byte("dummy key for "+host+"\n"), 0o600)
    }
    return cert, nil
}

func (m *FileManager) Renew(host string) (Cert, error) {
    // For now, same as Ensure with a bumped expiry.
    c, err := m.Ensure(host)
    if err != nil { return c, err }
    c.Expiry = time.Now().Add(48 * time.Hour)
    return c, nil
}

