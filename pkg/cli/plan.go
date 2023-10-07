package cli

import (
	"os"

	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/controller"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/terraform"
	"github.com/urfave/cli/v2"
)

func cmdPlan(ctx *cli.Context) error {
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
		Parser:             terraform.NewPlanParser(),
		Template:           terraform.NewPlanTemplate(cfg.Terraform.Plan.Template),
		ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(cfg.Terraform.Plan.WhenParseError.Template),
	}
	args := ctx.Args()

	return t.Plan(ctx.Context, controller.Command{
		Cmd:  args.First(),
		Args: args.Tail(),
	})
}
