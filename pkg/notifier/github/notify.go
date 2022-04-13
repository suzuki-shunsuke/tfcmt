package github

import (
	"context"
	"net/http"
	"os"

	"github.com/google/go-github/v39/github"
	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/github-comment-metadata/metadata"
	"github.com/suzuki-shunsuke/tfcmt/pkg/notifier"
	"github.com/suzuki-shunsuke/tfcmt/pkg/terraform"
)

// NotifyService handles communication with the notification related
// methods of GitHub API
type NotifyService service

// Notify posts comment optimized for notifications
func (g *NotifyService) Notify(ctx context.Context, param *notifier.ParamExec) (int, error) { //nolint:cyclop
	cfg := g.client.Config
	parser := g.client.Config.Parser
	template := g.client.Config.Template
	var errMsgs []string

	result := parser.Parse(param.CombinedOutput)
	result.ExitCode = param.ExitCode
	if result.HasParseError {
		template = g.client.Config.ParseErrorTemplate
	} else {
		if result.Error != nil {
			return result.ExitCode, result.Error
		}
		if result.Result == "" {
			return result.ExitCode, result.Error
		}
	}

	_, isPlan := parser.(*terraform.PlanParser)
	if isPlan {
		if cfg.PR.IsNumber() && cfg.ResultLabels.HasAnyLabelDefined() {
			errMsgs = append(errMsgs, g.updateLabels(ctx, result)...)
		}
	}

	template.SetValue(terraform.CommonTemplate{
		Result:                 result.Result,
		ChangedResult:          result.ChangedResult,
		ChangeOutsideTerraform: result.OutsideTerraform,
		Warning:                result.Warning,
		HasDestroy:             result.HasDestroy,
		Link:                   cfg.CI,
		UseRawOutput:           cfg.UseRawOutput,
		Vars:                   cfg.Vars,
		Templates:              cfg.Templates,
		Stdout:                 param.Stdout,
		Stderr:                 param.Stderr,
		CombinedOutput:         param.CombinedOutput,
		ExitCode:               param.ExitCode,
		ErrorMessages:          errMsgs,
		CreatedResources:       result.CreatedResources,
		UpdatedResources:       result.UpdatedResources,
		DeletedResources:       result.DeletedResources,
		ReplacedResources:      result.ReplacedResources,
	})
	body, err := template.Execute()
	if err != nil {
		return result.ExitCode, err
	}
	if _, isApply := parser.(*terraform.ApplyParser); isApply {
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

	logE := logrus.WithFields(logrus.Fields{
		"program": "tfcmt",
	})

	embeddedComment, err := getEmbeddedComment(&cfg, param.CIName, isPlan)
	if err != nil {
		return result.ExitCode, err
	}
	logE.WithFields(logrus.Fields{
		"comment": embeddedComment,
	}).Debug("embedded HTML comment")
	// embed HTML tag to hide old comments
	body += embeddedComment

	if cfg.Patch && cfg.PR.Number != 0 {
		logE.Debug("try patching")
		comments, err := g.client.Comment.List(ctx, cfg.Owner, cfg.Repo, cfg.PR.Number)
		if err != nil {
			logE.WithError(err).Debug("list comments")
			if err := g.client.Comment.Post(ctx, body, PostOptions{
				Number:   cfg.PR.Number,
				Revision: cfg.PR.Revision,
			}); err != nil {
				return result.ExitCode, err
			}
			return result.ExitCode, nil
		}
		logE.WithField("size", len(comments)).Debug("list comments")
		comment := g.getPatchedComment(logE, comments, cfg.Vars["target"])
		if comment != nil {
			if comment.Body == body {
				logE.Debug("comment isn't changed")
				return result.ExitCode, nil
			}
			logE.WithField("comment_id", comment.DatabaseID).Debug("patch a comment")
			if err := g.client.Comment.Patch(ctx, body, int64(comment.DatabaseID)); err != nil {
				return result.ExitCode, err
			}
			return result.ExitCode, nil
		}
	}

	logE.Debug("create a comment")
	if err := g.client.Comment.Post(ctx, body, PostOptions{
		Number:   cfg.PR.Number,
		Revision: cfg.PR.Revision,
	}); err != nil {
		return result.ExitCode, err
	}
	return result.ExitCode, nil
}

func (g *NotifyService) getPatchedComment(logE *logrus.Entry, comments []*IssueComment, target string) *IssueComment {
	var cmt *IssueComment
	for i, comment := range comments {
		logE := logE.WithFields(logrus.Fields{
			"comment_database_id": comment.DatabaseID,
			"comment_index":       i,
		})
		data := &Metadata{}
		f, err := metadata.Extract(comment.Body, data)
		if err != nil {
			logE.WithError(err).Debug("extract metadata from comment")
			continue
		}
		if !f {
			logE.Debug("metadata isn't found")
			continue
		}
		if data.Program != "tfcmt" {
			logE.Debug("Program isn't tfcmt")
			continue
		}
		if data.Target != target {
			logE.Debug("target is different")
			continue
		}
		if comment.IsMinimized {
			logE.Debug("comment is hidden")
			continue
		}
		cmt = comment
	}
	return cmt
}

type Metadata struct {
	Target  string
	Program string
}

func getEmbeddedComment(cfg *Config, ciName string, isPlan bool) (string, error) {
	vars := make(map[string]interface{}, len(cfg.EmbeddedVarNames))
	for _, name := range cfg.EmbeddedVarNames {
		vars[name] = cfg.Vars[name]
	}

	data := map[string]interface{}{
		"Program":  "tfcmt",
		"Vars":     vars,
		"SHA1":     cfg.PR.Revision,
		"PRNumber": cfg.PR.Number,
	}
	if target := cfg.Vars["target"]; target != "" {
		data["Target"] = target
	}
	if isPlan {
		data["Command"] = "plan"
	} else {
		data["Command"] = "apply"
	}
	if err := metadata.SetCIEnv(ciName, os.Getenv, data); err != nil {
		return "", err
	}
	embeddedComment, err := metadata.Convert(data)
	if err != nil {
		return "", err
	}
	return embeddedComment, nil
}

func (g *NotifyService) updateLabels(ctx context.Context, result terraform.ParseResult) []string { //nolint:cyclop
	cfg := g.client.Config
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

	errMsgs := []string{}

	logE := logrus.WithFields(logrus.Fields{
		"program": "tfcmt",
	})

	currentLabelColor, err := g.removeResultLabels(ctx, labelToAdd)
	if err != nil {
		msg := "remove labels: " + err.Error()
		logE.WithError(err).Error("remove labels")
		errMsgs = append(errMsgs, msg)
	}

	if labelToAdd == "" {
		return errMsgs
	}

	if currentLabelColor == "" {
		labels, _, err := g.client.API.IssuesAddLabels(ctx, cfg.PR.Number, []string{labelToAdd})
		if err != nil {
			msg := "add a label " + labelToAdd + ": " + err.Error()
			logE.WithError(err).WithFields(logrus.Fields{
				"label": labelToAdd,
			}).Error("add a label")
			errMsgs = append(errMsgs, msg)
		}
		if labelColor != "" {
			// set the color of label
			for _, label := range labels {
				if labelToAdd == label.GetName() {
					if label.GetColor() != labelColor {
						if _, _, err := g.client.API.IssuesUpdateLabel(ctx, labelToAdd, labelColor); err != nil {
							msg := "update a label color (name: " + labelToAdd + ", color: " + labelColor + "): " + err.Error()
							logE.WithError(err).WithFields(logrus.Fields{
								"label": labelToAdd,
								"color": labelColor,
							}).Error("update a label color")
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
			logE.WithError(err).WithFields(logrus.Fields{
				"label": labelToAdd,
				"color": labelColor,
			}).Error("update a label color")
			errMsgs = append(errMsgs, msg)
		}
	}
	return errMsgs
}

func (g *NotifyService) removeResultLabels(ctx context.Context, label string) (string, error) {
	cfg := g.client.Config
	// A Pull Request can have 100 labels the maximum
	labels, _, err := g.client.API.IssuesListLabels(ctx, cfg.PR.Number, &github.ListOptions{
		PerPage: 100, //nolint:gomnd
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
