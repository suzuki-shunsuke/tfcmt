package cli

import (
	"os"

	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/controller"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/terraform"
	"github.com/urfave/cli/v2"
)

func cmdApply(ctx *cli.Context) error {
	logLevel := ctx.String("log-level")
	setLogLevel(logLevel)

	cfg, err := newConfig(ctx)
	if err != nil {
		return err
	}

	if logLevel == "" {
		logLevel = cfg.Log.Level
		setLogLevel(logLevel)
	}

	if err := parseOpts(ctx, &cfg, os.Environ()); err != nil {
		return err
	}

	t := &controller.Controller{
		Config:             cfg,
		Parser:             terraform.NewApplyParser(),
		Template:           terraform.NewApplyTemplate(cfg.Terraform.Apply.Template),
		ParseErrorTemplate: terraform.NewApplyParseErrorTemplate(cfg.Terraform.Apply.WhenParseError.Template),
	}

	args := ctx.Args()

	return t.Apply(ctx.Context, controller.Command{
		Cmd:  args.First(),
		Args: args.Tail(),
	})
}
