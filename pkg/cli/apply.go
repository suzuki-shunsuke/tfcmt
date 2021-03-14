package cli

import (
	"github.com/suzuki-shunsuke/tfcmt/pkg/controller"
	"github.com/suzuki-shunsuke/tfcmt/pkg/platform"
	"github.com/suzuki-shunsuke/tfcmt/pkg/terraform"
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

	var ci platform.CI
	if err := parseOpts(ctx, &ci); err != nil {
		return err
	}

	t := &controller.Controller{
		Config:             cfg,
		Parser:             terraform.NewApplyParser(),
		Template:           terraform.NewApplyTemplate(cfg.Terraform.Apply.Template),
		ParseErrorTemplate: terraform.NewApplyParseErrorTemplate(cfg.Terraform.Apply.WhenParseError.Template),
		CI:                 ci,
	}

	args := ctx.Args()

	return t.Run(ctx.Context, controller.Command{
		Cmd:  args.First(),
		Args: args.Tail(),
	})
}
