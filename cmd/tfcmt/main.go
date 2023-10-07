package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/apperr"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/cli"
)

var (
	version = ""
	commit  = "" //nolint:gochecknoglobals
	date    = "" //nolint:gochecknoglobals
)

func main() {
	os.Exit(core())
}

func core() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	app := cli.New(&cli.LDFlags{
		Version: version,
		Commit:  commit,
		Date:    date,
	})
	return apperr.HandleExit(app.RunContext(ctx, os.Args))
}
