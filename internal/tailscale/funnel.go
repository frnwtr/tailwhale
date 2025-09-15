package tailscale

import (
    "context"
)

// Funnel provides simple on/off controls via the tailscale CLI.
type Funnel struct{ Exec Executor }

// On enables Tailscale Funnel for the current node or specified service.
// Behavior depends on the host configuration; this is a thin wrapper.
func (f Funnel) On(ctx context.Context, args ...string) error {
    ex := f.Exec
    if ex == nil { ex = defaultExec{} }
    params := append([]string{"funnel", "on"}, args...)
    return ex.Run(ctx, "tailscale", params...)
}

// Off disables Tailscale Funnel.
func (f Funnel) Off(ctx context.Context) error {
    ex := f.Exec
    if ex == nil { ex = defaultExec{} }
    return ex.Run(ctx, "tailscale", "funnel", "off")
}

