package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/controller"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/terraform"
)

func actionApply(ctx context.Context, logger *slogutil.Logger, args *ApplyArgs) error {
	if err := logger.SetLevel(args.LogLevel); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}

	cfg, err := newConfig(args.Config)
	if err != nil {
		return err
	}

	if args.LogLevel == "" {
		if err := logger.SetLevel(cfg.Log.Level); err != nil {
			return fmt.Errorf("set log level: %w", err)
		}
	}

	if err := parseOptsApply(args, &cfg, os.Environ()); err != nil {
		return err
	}

	t := &controller.Controller{
		Config:             cfg,
		Parser:             terraform.NewApplyParser(),
		Template:           terraform.NewApplyTemplate(cfg.Terraform.Apply.Template),
		ParseErrorTemplate: terraform.NewApplyParseErrorTemplate(cfg.Terraform.Apply.WhenParseError.Template),
	}

	return t.Apply(ctx, logger.Logger, controller.Command{
		Cmd:  args.Command,
		Args: args.CommandArgs,
	})
}
