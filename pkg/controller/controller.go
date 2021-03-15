package controller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/mattn/go-colorable"
	"github.com/suzuki-shunsuke/tfcmt/pkg/apperr"
	"github.com/suzuki-shunsuke/tfcmt/pkg/config"
	"github.com/suzuki-shunsuke/tfcmt/pkg/notifier"
	"github.com/suzuki-shunsuke/tfcmt/pkg/notifier/github"
	"github.com/suzuki-shunsuke/tfcmt/pkg/platform"
	"github.com/suzuki-shunsuke/tfcmt/pkg/terraform"
)

type Controller struct {
	Config                 config.Config
	Parser                 terraform.Parser
	Template               *terraform.Template
	DestroyWarningTemplate *terraform.Template
	ParseErrorTemplate     *terraform.Template
}

type Command struct {
	Cmd  string
	Args []string
}

// Run sends the notification with notifier
func (ctrl *Controller) Run(ctx context.Context, command Command) error {
	if err := platform.Complement(&ctrl.Config); err != nil {
		return err
	}

	if err := ctrl.Config.Validate(); err != nil {
		return err
	}

	ntf, err := ctrl.getNotifier(ctx)
	if err != nil {
		return err
	}

	if ntf == nil {
		return errors.New("no notifier specified at all")
	}

	cmd := exec.CommandContext(ctx, command.Cmd, command.Args...) //nolint:gosec
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	combinedOutput := &bytes.Buffer{}
	uncolorizedStdout := colorable.NewNonColorable(stdout)
	uncolorizedStderr := colorable.NewNonColorable(stderr)
	uncolorizedCombinedOutput := colorable.NewNonColorable(combinedOutput)
	cmd.Stdout = io.MultiWriter(os.Stdout, uncolorizedStdout, uncolorizedCombinedOutput)
	cmd.Stderr = io.MultiWriter(os.Stderr, uncolorizedStderr, uncolorizedCombinedOutput)
	_ = cmd.Run()

	return apperr.NewExitError(ntf.Notify(ctx, notifier.ParamExec{
		Stdout:         stdout.String(),
		Stderr:         stderr.String(),
		CombinedOutput: combinedOutput.String(),
		Cmd:            cmd,
		CIName:         ctrl.Config.CI.Name,
		ExitCode:       cmd.ProcessState.ExitCode(),
	}))
}

func (ctrl *Controller) renderTemplate(tpl string) (string, error) {
	tmpl, err := template.New("_").Funcs(sprig.TxtFuncMap()).Parse(tpl)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, map[string]interface{}{
		"Vars": ctrl.Config.Vars,
	}); err != nil {
		return "", fmt.Errorf("render a label template: %w", err)
	}
	return buf.String(), nil
}

func (ctrl *Controller) renderGitHubLabels() (github.ResultLabels, error) { //nolint:cyclop
	labels := github.ResultLabels{
		AddOrUpdateLabelColor: ctrl.Config.Terraform.Plan.WhenAddOrUpdateOnly.Color,
		DestroyLabelColor:     ctrl.Config.Terraform.Plan.WhenDestroy.Color,
		NoChangesLabelColor:   ctrl.Config.Terraform.Plan.WhenNoChanges.Color,
		PlanErrorLabelColor:   ctrl.Config.Terraform.Plan.WhenPlanError.Color,
	}

	target, ok := ctrl.Config.Vars["target"]
	if !ok {
		target = ""
	}

	if labels.AddOrUpdateLabelColor == "" {
		labels.AddOrUpdateLabelColor = "1d76db" // blue
	}
	if labels.DestroyLabelColor == "" {
		labels.DestroyLabelColor = "d93f0b" // red
	}
	if labels.NoChangesLabelColor == "" {
		labels.NoChangesLabelColor = "0e8a16" // green
	}

	if ctrl.Config.Terraform.Plan.WhenAddOrUpdateOnly.Label == "" {
		if target == "" {
			labels.AddOrUpdateLabel = "add-or-update"
		} else {
			labels.AddOrUpdateLabel = target + "/add-or-update"
		}
	} else {
		addOrUpdateLabel, err := ctrl.renderTemplate(ctrl.Config.Terraform.Plan.WhenAddOrUpdateOnly.Label)
		if err != nil {
			return labels, err
		}
		labels.AddOrUpdateLabel = addOrUpdateLabel
	}

	if ctrl.Config.Terraform.Plan.WhenDestroy.Label == "" {
		if target == "" {
			labels.DestroyLabel = "destroy"
		} else {
			labels.DestroyLabel = target + "/destroy"
		}
	} else {
		destroyLabel, err := ctrl.renderTemplate(ctrl.Config.Terraform.Plan.WhenDestroy.Label)
		if err != nil {
			return labels, err
		}
		labels.DestroyLabel = destroyLabel
	}

	if ctrl.Config.Terraform.Plan.WhenNoChanges.Label == "" {
		if target == "" {
			labels.NoChangesLabel = "no-changes"
		} else {
			labels.NoChangesLabel = target + "/no-changes"
		}
	} else {
		nochangesLabel, err := ctrl.renderTemplate(ctrl.Config.Terraform.Plan.WhenNoChanges.Label)
		if err != nil {
			return labels, err
		}
		labels.NoChangesLabel = nochangesLabel
	}

	planErrorLabel, err := ctrl.renderTemplate(ctrl.Config.Terraform.Plan.WhenPlanError.Label)
	if err != nil {
		return labels, err
	}
	labels.PlanErrorLabel = planErrorLabel

	return labels, nil
}

func (ctrl *Controller) getNotifier(ctx context.Context) (notifier.Notifier, error) {
	labels := github.ResultLabels{}
	if !ctrl.Config.Terraform.Plan.DisableLabel {
		a, err := ctrl.renderGitHubLabels()
		if err != nil {
			return nil, err
		}
		labels = a
	}
	client, err := github.NewClient(ctx, github.Config{
		Token:   ctrl.Config.GitHubToken,
		BaseURL: ctrl.Config.GHEBaseURL,
		Owner:   ctrl.Config.CI.Owner,
		Repo:    ctrl.Config.CI.Repo,
		PR: github.PullRequest{
			Revision: ctrl.Config.CI.SHA,
			Number:   ctrl.Config.CI.PRNumber,
		},
		CI:                     ctrl.Config.CI.Link,
		Parser:                 ctrl.Parser,
		UseRawOutput:           ctrl.Config.Terraform.UseRawOutput,
		Template:               ctrl.Template,
		DestroyWarningTemplate: ctrl.DestroyWarningTemplate,
		ParseErrorTemplate:     ctrl.ParseErrorTemplate,
		ResultLabels:           labels,
		Vars:                   ctrl.Config.Vars,
		Templates:              ctrl.Config.Templates,
	})
	if err != nil {
		return nil, err
	}
	return client.Notify, nil
}
