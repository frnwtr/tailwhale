package tailscale

import (
    "context"
    "strings"
    "testing"
)

type recExec struct{ calls [][]string }

func (r *recExec) Run(ctx context.Context, name string, args ...string) error {
    _ = ctx
    r.calls = append(r.calls, append([]string{name}, args...))
    return nil
}

func TestFunnelOnOff(t *testing.T){
    exec := &recExec{}
    f := Funnel{Exec: exec}
    if err := f.On(context.Background(), "80"); err != nil { t.Fatal(err) }
    if err := f.Off(context.Background()); err != nil { t.Fatal(err) }
    if len(exec.calls) != 2 { t.Fatalf("expected 2 calls, got %d", len(exec.calls)) }
    got := strings.Join(exec.calls[0], " ")
    if !strings.Contains(got, "tailscale funnel on 80") { t.Fatalf("unexpected call: %s", got) }
}

