package github

import (
	"context"
	"net/http"

	"github.com/mercari/tfnotify/terraform"
)

// NotifyService handles communication with the notification related
// methods of GitHub API
type NotifyService service

// Notify posts comment optimized for notifications
func (g *NotifyService) Notify(ctx context.Context, body string) (exit int, err error) {
	cfg := g.client.Config
	parser := g.client.Config.Parser
	template := g.client.Config.Template

	result := parser.Parse(body)
	if result.Error != nil {
		return result.ExitCode, result.Error
	}
	if result.Result == "" {
		return result.ExitCode, result.Error
	}

	_, isPlan := parser.(*terraform.PlanParser)
	if isPlan {
		if result.HasDestroy && cfg.WarnDestroy {
			// Notify destroy warning as a new comment before normal plan result
			if err = g.notifyDestoryWarning(ctx, body, result); err != nil {
				return result.ExitCode, err
			}
		}
		if cfg.PR.IsNumber() && cfg.ResultLabels.HasAnyLabelDefined() {
			err = g.removeResultLabels(ctx)
			if err != nil {
				return result.ExitCode, err
			}
			var (
				labelToAdd string
				labelColor string
			)

			switch {
			case result.HasAddOrUpdateOnly:
				labelToAdd = cfg.ResultLabels.AddOrUpdateLabel
				labelColor = cfg.ResultLabels.AddOrUpdateLabelColor
			case result.HasDestroy:
				labelToAdd = cfg.ResultLabels.DestroyLabel
				labelColor = cfg.ResultLabels.DestroyLabelColor
			case result.HasNoChanges:
				labelToAdd = cfg.ResultLabels.NoChangesLabel
				labelColor = cfg.ResultLabels.NoChangesLabelColor
			case result.HasPlanError:
				labelToAdd = cfg.ResultLabels.PlanErrorLabel
				labelColor = cfg.ResultLabels.PlanErrorLabelColor
			}

			if labelToAdd != "" {
				labels, _, err := g.client.API.IssuesAddLabels(
					ctx,
					cfg.PR.Number,
					[]string{labelToAdd},
				)
				if err != nil {
					return result.ExitCode, err
				}
				if labelColor != "" {
					// set the color of label
					for _, label := range labels {
						if labelToAdd == label.GetName() {
							if label.GetColor() != labelColor {
								_, _, err := g.client.API.IssuesUpdateLabel(ctx, labelToAdd, labelColor)
								if err != nil {
									return result.ExitCode, err
								}
							}
						}
					}
				}
			}
		}
	}

	template.SetValue(terraform.CommonTemplate{
		Title:        cfg.PR.Title,
		Message:      cfg.PR.Message,
		Result:       result.Result,
		Body:         body,
		Link:         cfg.CI,
		UseRawOutput: cfg.UseRawOutput,
	})
	body, err = template.Execute()
	if err != nil {
		return result.ExitCode, err
	}

	value := template.GetValue()

	if cfg.PR.IsNumber() {
		g.client.Comment.DeleteDuplicates(ctx, value.Title)
	}

	_, isApply := parser.(*terraform.ApplyParser)
	if isApply {
		prNumber, err := g.client.Commits.MergedPRNumber(ctx, cfg.PR.Revision)
		if err == nil {
			cfg.PR.Number = prNumber
		} else if !cfg.PR.IsNumber() {
			commits, err := g.client.Commits.List(ctx, cfg.PR.Revision)
			if err != nil {
				return result.ExitCode, err
			}
			lastRevision, _ := g.client.Commits.lastOne(commits, cfg.PR.Revision)
			cfg.PR.Revision = lastRevision
		}
	}

	return result.ExitCode, g.client.Comment.Post(ctx, body, PostOptions{
		Number:   cfg.PR.Number,
		Revision: cfg.PR.Revision,
	})
}

func (g *NotifyService) notifyDestoryWarning(ctx context.Context, body string, result terraform.ParseResult) error {
	cfg := g.client.Config
	destroyWarningTemplate := g.client.Config.DestroyWarningTemplate
	destroyWarningTemplate.SetValue(terraform.CommonTemplate{
		Title:        cfg.PR.DestroyWarningTitle,
		Message:      cfg.PR.DestroyWarningMessage,
		Result:       result.Result,
		Body:         body,
		Link:         cfg.CI,
		UseRawOutput: cfg.UseRawOutput,
	})
	body, err := destroyWarningTemplate.Execute()
	if err != nil {
		return err
	}

	return g.client.Comment.Post(ctx, body, PostOptions{
		Number:   cfg.PR.Number,
		Revision: cfg.PR.Revision,
	})
}

func (g *NotifyService) removeResultLabels(ctx context.Context) error {
	cfg := g.client.Config
	labels, _, err := g.client.API.IssuesListLabels(ctx, cfg.PR.Number, nil)
	if err != nil {
		return err
	}

	for _, l := range labels {
		labelText := l.GetName()
		if cfg.ResultLabels.IsResultLabel(labelText) {
			resp, err := g.client.API.IssuesRemoveLabel(ctx, cfg.PR.Number, labelText)
			// Ignore 404 errors, which are from the PR not having the label
			if err != nil && resp.StatusCode != http.StatusNotFound {
				return err
			}
		}
	}

	return nil
}
