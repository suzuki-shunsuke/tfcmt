package controller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/suzuki-shunsuke/go-error-with-exit-code/ecerror"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/mask"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/notifier"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/platform"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/terraform"
)

// Plan sends the notification with notifier
func (c *Controller) Plan(ctx context.Context, logger *slog.Logger, command Command) error {
	if command.Cmd == "" {
		return errors.New("no command specified")
	}
	if err := platform.Complement(&c.Config); err != nil {
		return err
	}

	if err := c.Config.Validate(); err != nil {
		return err
	}

	ntf, err := c.getPlanNotifier(ctx)
	if err != nil {
		return err
	}

	if ntf == nil {
		return errors.New("no notifier specified at all")
	}

	cmd := exec.CommandContext(ctx, command.Cmd, command.Args...) //nolint:gosec
	cmd.Stdin = os.Stdin
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	combinedOutput := &bytes.Buffer{}
	uncolorizedStdout := colorable.NewNonColorable(stdout)
	uncolorizedStderr := colorable.NewNonColorable(stderr)
	uncolorizedCombinedOutput := colorable.NewNonColorable(combinedOutput)
	cmd.Stdout = io.MultiWriter(mask.NewWriter(os.Stdout, c.Config.Masks), uncolorizedStdout, uncolorizedCombinedOutput)
	cmd.Stderr = io.MultiWriter(mask.NewWriter(os.Stderr, c.Config.Masks), uncolorizedStderr, uncolorizedCombinedOutput)
	setCancel(cmd)
	_ = cmd.Run()

	exitCode := 0
	if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}

	if err := c.ExportPlanResult(combinedOutput.String(), exitCode); err != nil {
		return ecerror.Wrap(err, exitCode)
	}

	return ecerror.Wrap(ntf.Plan(ctx, logger, &notifier.ParamExec{
		Stdout:         stdout.String(),
		Stderr:         stderr.String(),
		CombinedOutput: combinedOutput.String(),
		CIName:         c.Config.CI.Name,
		ExitCode:       exitCode,
	}), exitCode)
}

// ExportPlanResult writes machine-readable plan result to file and/or GITHUB_OUTPUT.
func (c *Controller) ExportPlanResult(combinedOutput string, exitCode int) error {
	if c.Config.Terraform.Plan.ResultFile == "" && !c.Config.Terraform.Plan.GitHubOutput {
		return nil
	}

	parser := c.Parser
	if parser == nil {
		parser = terraform.NewPlanParser()
	}

	result := parser.Parse(combinedOutput)
	planResult := terraform.NewPlanResult(exitCode, result, c.Config.Vars)

	if c.Config.Terraform.Plan.ResultFile != "" {
		if err := terraform.WriteResultFile(c.Config.Terraform.Plan.ResultFile, planResult); err != nil {
			return fmt.Errorf("write plan result file: %w", err)
		}
	}

	if c.Config.Terraform.Plan.GitHubOutput {
		outputFile := os.Getenv("GITHUB_OUTPUT")
		if outputFile == "" {
			return errors.New("GITHUB_OUTPUT environment variable is not set")
		}
		if err := terraform.ExportGitHubOutput(outputFile, planResult); err != nil {
			return fmt.Errorf("export plan result to GITHUB_OUTPUT: %w", err)
		}
	}

	return nil
}

const waitDelay = 1000 * time.Hour

func setCancel(cmd *exec.Cmd) {
	cmd.Cancel = func() error {
		return cmd.Process.Signal(os.Interrupt) //nolint:wrapcheck
	}
	cmd.WaitDelay = waitDelay
}
