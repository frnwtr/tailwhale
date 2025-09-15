package tailscale

import (
    "context"
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "math/big"
    "os"
    "path/filepath"
    "testing"
    "time"
)

type fakeExec struct{ calls [][]string }

func (f *fakeExec) Run(ctx context.Context, name string, args ...string) error {
    _ = ctx
    call := append([]string{name}, args...)
    f.calls = append(f.calls, call)
    // Simulate tailscale cert writing files by touching cert file
    for i := 0; i < len(args)-1; i++ {
        if args[i] == "--cert-file" {
            // write a short-lived cert
            _ = writeTestCert(args[i+1], time.Now().Add(48*time.Hour))
        }
    }
    return nil
}

func writeTestCert(path string, notAfter time.Time) error {
    // minimal self-signed cert
    key, _ := rsa.GenerateKey(rand.Reader, 1024)
    tpl := x509.Certificate{
        SerialNumber: big.NewInt(1),
        Subject: pkix.Name{CommonName: "test"},
        NotBefore: time.Now().Add(-time.Hour),
        NotAfter:  notAfter,
        KeyUsage: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
        BasicConstraintsValid: true,
    }
    der, _ := x509.CreateCertificate(rand.Reader, &tpl, &tpl, &key.PublicKey, key)
    pemBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
    _ = os.MkdirAll(filepath.Dir(path), 0o755)
    return os.WriteFile(path, pemBytes, 0o644)
}

func TestShellManagerEnsureIssuesWhenMissing(t *testing.T){
    dir := t.TempDir()
    m := &ShellManager{Exec: &fakeExec{}, CertDir: dir, MinRemain: 24 * time.Hour}
    c, err := m.Ensure("example.test")
    if err != nil { t.Fatal(err) }
    if time.Until(c.Expiry) <= 0 { t.Fatalf("expected expiry in future, got %v", c.Expiry) }
}

func TestShellManagerEnsureSkipsWhenValid(t *testing.T){
    dir := t.TempDir()
    certPath := filepath.Join(dir, "example.test.crt")
    if err := writeTestCert(certPath, time.Now().Add(72*time.Hour)); err != nil { t.Fatal(err) }
    m := &ShellManager{Exec: &fakeExec{}, CertDir: dir, MinRemain: 24 * time.Hour}
    c, err := m.Ensure("example.test")
    if err != nil { t.Fatal(err) }
    if len(m.Exec.(*fakeExec).calls) != 0 { t.Fatalf("expected no exec calls, got %d", len(m.Exec.(*fakeExec).calls)) }
    if time.Until(c.Expiry) < 48*time.Hour { t.Fatalf("unexpected expiry: %v", c.Expiry) }
}

