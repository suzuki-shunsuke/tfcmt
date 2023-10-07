package controller

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/config"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/notifier"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/notifier/github"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/notifier/localfile"
	tmpl "github.com/suzuki-shunsuke/tfcmt/v4/pkg/template"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/terraform"
)

type Controller struct {
	Config             config.Config
	Parser             terraform.Parser
	Template           *terraform.Template
	ParseErrorTemplate *terraform.Template
}

type Command struct {
	Cmd  string
	Args []string
}

func (ctrl *Controller) renderTemplate(tpl string) (string, error) {
	tmpl, err := template.New("_").Funcs(tmpl.TxtFuncMap()).Parse(tpl)
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

	if !ctrl.Config.Terraform.Plan.WhenAddOrUpdateOnly.DisableLabel {
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
	}

	if !ctrl.Config.Terraform.Plan.WhenDestroy.DisableLabel {
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
	}

	if !ctrl.Config.Terraform.Plan.WhenNoChanges.DisableLabel {
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
	}

	if !ctrl.Config.Terraform.Plan.WhenPlanError.DisableLabel {
		planErrorLabel, err := ctrl.renderTemplate(ctrl.Config.Terraform.Plan.WhenPlanError.Label)
		if err != nil {
			return labels, err
		}
		labels.PlanErrorLabel = planErrorLabel
	}

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
	// Write output to file instead of github comment
	if ctrl.Config.Output != "" {
		client, err := localfile.NewClient(&localfile.Config{
			OutputFile:         ctrl.Config.Output,
			Parser:             ctrl.Parser,
			UseRawOutput:       ctrl.Config.Terraform.UseRawOutput,
			CI:                 ctrl.Config.CI.Link,
			Template:           ctrl.Template,
			ParseErrorTemplate: ctrl.ParseErrorTemplate,
			Vars:               ctrl.Config.Vars,
			EmbeddedVarNames:   ctrl.Config.EmbeddedVarNames,
			Templates:          ctrl.Config.Templates,
		})
		if err != nil {
			return nil, err
		}
		return client.Notify, nil
	}
	client, err := github.NewClient(ctx, &github.Config{
		Token:           ctrl.Config.GitHubToken,
		BaseURL:         ctrl.Config.GHEBaseURL,
		GraphQLEndpoint: ctrl.Config.GHEGraphQLEndpoint,
		Owner:           ctrl.Config.CI.Owner,
		Repo:            ctrl.Config.CI.Repo,
		PR: github.PullRequest{
			Revision: ctrl.Config.CI.SHA,
			Number:   ctrl.Config.CI.PRNumber,
		},
		CI:                 ctrl.Config.CI.Link,
		Parser:             ctrl.Parser,
		UseRawOutput:       ctrl.Config.Terraform.UseRawOutput,
		Template:           ctrl.Template,
		ParseErrorTemplate: ctrl.ParseErrorTemplate,
		ResultLabels:       labels,
		Vars:               ctrl.Config.Vars,
		EmbeddedVarNames:   ctrl.Config.EmbeddedVarNames,
		Templates:          ctrl.Config.Templates,
		Patch:              ctrl.Config.PlanPatch,
		SkipNoChanges:      ctrl.Config.Terraform.Plan.WhenNoChanges.DisableComment,
	})
	if err != nil {
		return nil, err
	}
	return client.Notify, nil
}
