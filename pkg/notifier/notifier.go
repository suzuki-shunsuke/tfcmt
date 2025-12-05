package notifier

import (
	"context"
	"log/slog"
)

// Notifier is a notification interface
type Notifier interface {
	Apply(ctx context.Context, logger *slog.Logger, param *ParamExec) error
	Plan(ctx context.Context, logger *slog.Logger, param *ParamExec) error
}

type ParamExec struct {
	Stdout         string
	Stderr         string
	CombinedOutput string
	CIName         string
	ExitCode       int
}
