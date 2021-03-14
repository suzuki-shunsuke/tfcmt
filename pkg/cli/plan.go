package cli

import (
	"github.com/suzuki-shunsuke/tfcmt/pkg/controller"
	"github.com/suzuki-shunsuke/tfcmt/pkg/terraform"
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

	if err := parseOpts(ctx, &cfg); err != nil {
		return err
	}

	t := &controller.Controller{
		Config:                 cfg,
		Parser:                 terraform.NewPlanParser(),
		Template:               terraform.NewPlanTemplate(cfg.Terraform.Plan.Template),
		DestroyWarningTemplate: terraform.NewDestroyWarningTemplate(cfg.Terraform.Plan.WhenDestroy.Template),
		ParseErrorTemplate:     terraform.NewPlanParseErrorTemplate(cfg.Terraform.Plan.WhenParseError.Template),
	}
	args := ctx.Args()

	return t.Run(ctx.Context, controller.Command{
		Cmd:  args.First(),
		Args: args.Tail(),
	})
}
