package cli

import (
	"github.com/suzuki-shunsuke/tfcmt/pkg/controller"
	"github.com/suzuki-shunsuke/tfcmt/pkg/terraform"
	"github.com/urfave/cli/v2"
)

func cmdApply(ctx *cli.Context) error {
	cfg, err := newConfig(ctx)
	if err != nil {
		return err
	}
	t := &controller.Controller{
		Config:             cfg,
		Context:            ctx,
		Parser:             terraform.NewApplyParser(),
		Template:           terraform.NewApplyTemplate(cfg.Terraform.Apply.Template),
		ParseErrorTemplate: terraform.NewApplyParseErrorTemplate(cfg.Terraform.Apply.WhenParseError.Template),
	}
	return t.Run(ctx.Context)
}
