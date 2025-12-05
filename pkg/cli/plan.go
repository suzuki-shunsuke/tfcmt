package cli

import (
	"context"
	"log/slog"
	"os"

	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/controller"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/terraform"
	"github.com/urfave/cli/v3"
)

func cmdPlanFunc(logger *slog.Logger, logLevelVar *slog.LevelVar) func(ctx context.Context, cmd *cli.Command) error {
	return func(ctx context.Context, cmd *cli.Command) error {
		logLevel := cmd.String("log-level")
		setLogLevel(logLevelVar, logLevel)

		cfg, err := newConfig(cmd)
		if err != nil {
			return err
		}
		if logLevel == "" {
			logLevel = cfg.Log.Level
			setLogLevel(logLevelVar, logLevel)
		}

		if err := parseOpts(cmd, &cfg, os.Environ()); err != nil {
			return err
		}

		t := &controller.Controller{
			Config:             cfg,
			Parser:             terraform.NewPlanParser(),
			Template:           terraform.NewPlanTemplate(cfg.Terraform.Plan.Template),
			ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(cfg.Terraform.Plan.WhenParseError.Template),
		}
		args := cmd.Args()

		return t.Plan(ctx, logger, controller.Command{
			Cmd:  args.First(),
			Args: args.Tail(),
		})
	}
}
