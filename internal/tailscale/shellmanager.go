package tailscale

import (
    "context"
    "crypto/x509"
    "encoding/pem"
    "errors"
    "io/fs"
    "os"
    "path/filepath"
    "time"
)

// ShellManager uses the `tailscale` CLI to ensure and renew certificates.
// It assumes the tailscale CLI is installed and authenticated on the host.
type ShellManager struct {
    Exec   Executor
    CertDir string
    // Minimum remaining validity; if a cert expires sooner than this, Renew() should be called.
    MinRemain time.Duration
}

func (m *ShellManager) ensureExec() Executor {
    if m.Exec != nil { return m.Exec }
    return defaultExec{}
}

func (m *ShellManager) Ensure(host string) (Cert, error) {
    c := Cert{
        Host:    host,
        Path:    filepath.Join(m.CertDir, host+".crt"),
        KeyPath: filepath.Join(m.CertDir, host+".key"),
    }
    // If existing and valid, return it.
    if exp, ok := readCertExpiry(c.Path); ok {
        c.Expiry = exp
        if m.MinRemain == 0 || time.Until(exp) > m.MinRemain {
            return c, nil
        }
    }
    // Missing or near-expiry: run `tailscale cert` to (re)issue.
    if err := os.MkdirAll(m.CertDir, 0o755); err != nil { return c, err }
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    args := []string{"cert", "--cert-file", c.Path, "--key-file", c.KeyPath, host}
    if err := m.ensureExec().Run(ctx, "tailscale", args...); err != nil {
        return c, err
    }
    if exp, ok := readCertExpiry(c.Path); ok { c.Expiry = exp }
    return c, nil
}

func (m *ShellManager) Renew(host string) (Cert, error) {
    // Delegate to Ensure which reissues if close to expiry or missing.
    if m.MinRemain == 0 {
        m.MinRemain = 24 * time.Hour
    }
    return m.Ensure(host)
}

func readCertExpiry(path string) (time.Time, bool) {
    b, err := os.ReadFile(path)
    if err != nil { return time.Time{}, false }
    // Try PEM first
    var block *pem.Block
    var rest = b
    for {
        block, rest = pem.Decode(rest)
        if block == nil { break }
        if block.Type == "CERTIFICATE" {
            if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
                return cert.NotAfter, true
            }
        }
    }
    // Try raw DER
    if cert, err := x509.ParseCertificate(b); err == nil {
        return cert.NotAfter, true
    }
    return time.Time{}, false
}

var ErrNotExist = errors.New("file does not exist")

func fileExists(path string) error {
    _, err := os.Stat(path)
    if err == nil { return nil }
    if errors.Is(err, fs.ErrNotExist) { return ErrNotExist }
    return err
}

