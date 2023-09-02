package github

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/tfcmt/pkg/notifier"
	"github.com/suzuki-shunsuke/tfcmt/pkg/terraform"
)

// Plan posts comment optimized for notifications
func (g *NotifyService) Plan(ctx context.Context, param *notifier.ParamExec) (int, error) { //nolint:cyclop
	cfg := g.client.Config
	parser := g.client.Config.Parser
	template := g.client.Config.Template
	var errMsgs []string

	if cfg.PR.Number == 0 && cfg.PR.Revision != "" {
		if prNumber, err := g.client.Commits.PRNumber(ctx, cfg.PR.Revision, PullRequestStateOpen); err == nil {
			cfg.PR.Number = prNumber
		}
	}

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

	if cfg.PR.IsNumber() && cfg.ResultLabels.HasAnyLabelDefined() {
		errMsgs = append(errMsgs, g.updateLabels(ctx, result)...)
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
		MovedResources:         result.MovedResources,
		ImportedResources:      result.ImportedResources,
	})
	body, err := template.Execute()
	if err != nil {
		return result.ExitCode, err
	}

	logE := logrus.WithFields(logrus.Fields{
		"program": "tfcmt",
	})

	embeddedComment, err := getEmbeddedComment(cfg, param.CIName, true)
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
			if err := g.client.Comment.Post(ctx, body, &PostOptions{
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

	if result.HasNoChanges && result.Warning == "" && len(errMsgs) == 0 && cfg.SkipNoChanges {
		logE.Debug("skip posting a comment because there is no change")
		return result.ExitCode, nil
	}

	logE.Debug("create a comment")
	if err := g.client.Comment.Post(ctx, body, &PostOptions{
		Number:   cfg.PR.Number,
		Revision: cfg.PR.Revision,
	}); err != nil {
		return result.ExitCode, err
	}
	return result.ExitCode, nil
}
