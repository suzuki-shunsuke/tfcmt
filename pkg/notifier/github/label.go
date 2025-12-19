package github

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/go-github/v80/github"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/terraform"
)

func (g *NotifyService) UpdateLabels(ctx context.Context, logger *slog.Logger, result terraform.ParseResult) []string { //nolint:cyclop
	cfg := g.client.Config
	var (
		labelToAdd string
		labelColor string
	)
	if cfg.PR.Number == 0 {
		return nil
	}

	// If an override label is specified, use it instead of the calculated label
	if cfg.ResultLabels.OverrideLabel != "" {
		labelToAdd = cfg.ResultLabels.OverrideLabel
		// Don't remove other labels when using override label
		return g.addLabelOnly(ctx, logger, labelToAdd, labelColor)
	}

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
	case result.HasError:
		labelToAdd = cfg.ResultLabels.PlanErrorLabel
		labelColor = cfg.ResultLabels.PlanErrorLabelColor
	}

	errMsgs := []string{}

	currentLabelColor, err := g.removeResultLabels(ctx, labelToAdd)
	if err != nil {
		msg := "remove labels: " + err.Error()
		slogerr.WithError(logger, err).Error("remove labels")
		errMsgs = append(errMsgs, msg)
	}

	if labelToAdd == "" {
		return errMsgs
	}

	if len(labelToAdd) > 50 { //nolint:mnd
		return append(errMsgs, fmt.Sprintf("failed to add a label %s: label name is too long (max: 50)", labelToAdd))
	}

	if currentLabelColor == "" {
		labels, _, err := g.client.API.IssuesAddLabels(ctx, cfg.PR.Number, []string{labelToAdd})
		if err != nil {
			msg := "add a label " + labelToAdd + ": " + err.Error()
			slogerr.WithError(logger, err).Error("add a label", "label", labelToAdd)
			errMsgs = append(errMsgs, msg)
		}
		if labelColor != "" {
			// set the color of label
			for _, label := range labels {
				if labelToAdd == label.GetName() {
					if label.GetColor() != labelColor {
						if _, _, err := g.client.API.IssuesUpdateLabel(ctx, labelToAdd, labelColor); err != nil {
							msg := "update a label color (name: " + labelToAdd + ", color: " + labelColor + "): " + err.Error()
							slogerr.WithError(logger, err).Error("update a label color", "label", labelToAdd, "color", labelColor)
							errMsgs = append(errMsgs, msg)
						}
					}
				}
			}
		}
	} else if labelColor != "" && labelColor != currentLabelColor {
		// set the color of label
		if _, _, err := g.client.API.IssuesUpdateLabel(ctx, labelToAdd, labelColor); err != nil {
			msg := "update a label color (name: " + labelToAdd + ", color: " + labelColor + "): " + err.Error()
			slogerr.WithError(logger, err).Error("update a label color", "label", labelToAdd, "color", labelColor)
			errMsgs = append(errMsgs, msg)
		}
	}
	return errMsgs
}

func (g *NotifyService) removeResultLabels(ctx context.Context, label string) (string, error) {
	cfg := g.client.Config
	// A Pull Request can have 100 labels the maximum
	labels, _, err := g.client.API.IssuesListLabels(ctx, cfg.PR.Number, &github.ListOptions{
		PerPage: 100, //nolint:mnd
	})
	if err != nil {
		return "", err
	}

	labelColor := ""
	for _, l := range labels {
		labelText := l.GetName()
		if labelText == label {
			labelColor = l.GetColor()
			continue
		}
		if cfg.ResultLabels.IsResultLabel(labelText) {
			resp, err := g.client.API.IssuesRemoveLabel(ctx, cfg.PR.Number, labelText)
			// Ignore 404 errors, which are from the PR not having the label
			if err != nil && resp.StatusCode != http.StatusNotFound {
				return labelColor, err
			}
		}
	}

	return labelColor, nil
}

// addLabelOnly adds a label without removing other labels.
// This is used when an override label is specified to allow multiple labels to accumulate.
func (g *NotifyService) addLabelOnly(ctx context.Context, logger *slog.Logger, labelToAdd, labelColor string) []string {
	cfg := g.client.Config
	errMsgs := []string{}

	if labelToAdd == "" {
		return errMsgs
	}

	if len(labelToAdd) > 50 { //nolint:mnd
		return append(errMsgs, fmt.Sprintf("failed to add a label %s: label name is too long (max: 50)", labelToAdd))
	}

	// Check if label already exists
	labels, _, err := g.client.API.IssuesListLabels(ctx, cfg.PR.Number, &github.ListOptions{
		PerPage: 100, //nolint:mnd
	})
	if err != nil {
		msg := "list labels: " + err.Error()
		slogerr.WithError(logger, err).Error("list labels")
		return append(errMsgs, msg)
	}

	labelExists := false
	currentLabelColor := ""
	for _, label := range labels {
		if label.GetName() == labelToAdd {
			labelExists = true
			currentLabelColor = label.GetColor()
			break
		}
	}

	if !labelExists {
		_, _, err := g.client.API.IssuesAddLabels(ctx, cfg.PR.Number, []string{labelToAdd})
		if err != nil {
			msg := "add a label " + labelToAdd + ": " + err.Error()
			slogerr.WithError(logger, err).Error("add a label", "label", labelToAdd)
			errMsgs = append(errMsgs, msg)
			return errMsgs
		}
	}

	if labelColor != "" && labelColor != currentLabelColor {
		if _, _, err := g.client.API.IssuesUpdateLabel(ctx, labelToAdd, labelColor); err != nil {
			msg := "update a label color (name: " + labelToAdd + ", color: " + labelColor + "): " + err.Error()
			slogerr.WithError(logger, err).Error("update a label color", "label", labelToAdd, "color", labelColor)
			errMsgs = append(errMsgs, msg)
		}
	}

	return errMsgs
}
