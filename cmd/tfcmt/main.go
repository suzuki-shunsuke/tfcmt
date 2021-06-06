package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/suzuki-shunsuke/tfcmt/pkg/apperr"
	"github.com/suzuki-shunsuke/tfcmt/pkg/cli"
)

func main() {
	os.Exit(core())
}

func core() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	app := cli.New()
	return apperr.HandleExit(app.RunContext(ctx, os.Args))
}
