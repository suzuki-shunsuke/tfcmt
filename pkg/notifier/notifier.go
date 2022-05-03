package notifier

import (
	"context"
	"os/exec"
)

// Notifier is a notification interface
type Notifier interface {
	Apply(ctx context.Context, param ParamExec) (int, error)
	Plan(ctx context.Context, param ParamExec) (int, error)
}

type ParamExec struct {
	Stdout         string
	Stderr         string
	CombinedOutput string
	CIName         string
	Cmd            *exec.Cmd
	ExitCode       int
}
