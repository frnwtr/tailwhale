package tailscale

import (
    "context"
    "os/exec"
)

// Executor abstracts command execution for testability.
type Executor interface {
    Run(ctx context.Context, name string, args ...string) error
}

type defaultExec struct{}

func (defaultExec) Run(ctx context.Context, name string, args ...string) error {
    cmd := exec.CommandContext(ctx, name, args...)
    return cmd.Run()
}

