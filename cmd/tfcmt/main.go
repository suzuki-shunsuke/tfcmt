package main

import (
	"context"
	"os"

	"github.com/suzuki-shunsuke/tfcmt/pkg/apperr"
	"github.com/suzuki-shunsuke/tfcmt/pkg/cli"
	"github.com/suzuki-shunsuke/tfcmt/pkg/signal"
)

func main() {
	app := cli.New()
	ctx, cancel := context.WithCancel(context.Background())
	go signal.Handle(cancel)

	os.Exit(apperr.HandleExit(app.RunContext(ctx, os.Args)))
}
