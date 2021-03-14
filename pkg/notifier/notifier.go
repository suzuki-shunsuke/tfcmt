package notifier

import (
	"context"
	"os/exec"

	"github.com/urfave/cli/v2"
)

// Notifier is a notification interface
type Notifier interface {
	Notify(ctx context.Context, param ParamExec) (int, error)
}

type ParamExec struct {
	Stdout         string
	Stderr         string
	CombinedOutput string
	CIName         string
	Args           cli.Args
	Cmd            *exec.Cmd
	ExitCode       int
}
