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

	t := &controller.Controller{
		Config:                 cfg,
		Context:                ctx,
		Parser:                 terraform.NewPlanParser(),
		Template:               terraform.NewPlanTemplate(cfg.Terraform.Plan.Template),
		DestroyWarningTemplate: terraform.NewDestroyWarningTemplate(cfg.Terraform.Plan.WhenDestroy.Template),
		ParseErrorTemplate:     terraform.NewPlanParseErrorTemplate(cfg.Terraform.Plan.WhenParseError.Template),
	}
	return t.Run(ctx.Context)
}
