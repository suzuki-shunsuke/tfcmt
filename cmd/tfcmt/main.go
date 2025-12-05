package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"

	"github.com/suzuki-shunsuke/slog-error/slogerr"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/apperr"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/cli"
)

var (
	version = ""
	commit  = "" //nolint:gochecknoglobals
	date    = "" //nolint:gochecknoglobals
)

func main() {
	if code := core(); code != 0 {
		os.Exit(code)
	}
}

func core() int {
	logLevelVar := &slog.LevelVar{}
	logger := slogutil.New(&slogutil.InputNew{
		Name:    "tfcmt",
		Version: version,
		Out:     os.Stderr,
		Level:   logLevelVar,
	})
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	app := cli.New(&cli.LDFlags{
		Version: version,
		Commit:  commit,
		Date:    date,
	}, logger, logLevelVar)
	if err := app.Run(ctx, os.Args); err != nil {
		var exitErr *apperr.ExitError
		if errors.As(err, &exitErr) {
			if exitErr.Error() != "" {
				slogerr.WithError(logger, err).Error("tfcmt failed")
			}
			if code := exitErr.ExitCode(); code != 0 {
				return code
			}
			if exitErr.Error() == "" {
				return apperr.ExitCodeOK
			}
			return apperr.ExitCodeError
		}
		slogerr.WithError(logger, err).Error("tfcmt failed")
		return 1
	}
	return 0
}
