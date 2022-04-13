package notifier

import (
	"context"
	"os/exec"
)

// Notifier is a notification interface
type Notifier interface {
	Notify(ctx context.Context, param *ParamExec) (int, error)
}

type ParamExec struct {
	Stdout         string
	Stderr         string
	CombinedOutput string
	CIName         string
	Cmd            *exec.Cmd
	ExitCode       int
}
