package github

import (
	"context"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/tfcmt/notifier"
	"github.com/suzuki-shunsuke/tfcmt/terraform"
)

// NotifyService handles communication with the notification related
// methods of GitHub API
type NotifyService service

// Notify posts comment optimized for notifications
func (g *NotifyService) Notify(ctx context.Context, param notifier.ParamExec) (int, error) {
	cfg := g.client.Config
	parser := g.client.Config.Parser
	template := g.client.Config.Template
	var errMsgs []string
	target := cfg.Vars["target"]

	body := param.Stdout
	result := parser.Parse(body)
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
		if result.HasDestroy {
			template = g.client.Config.DestroyWarningTemplate
		}
		if cfg.PR.IsNumber() && cfg.ResultLabels.HasAnyLabelDefined() {
			errMsgs = append(errMsgs, g.updateLabels(ctx, result)...)
		}
	}

	template.SetValue(terraform.CommonTemplate{
		Result:            result.Result,
		Body:              body,
		Link:              cfg.CI,
		UseRawOutput:      cfg.UseRawOutput,
		Vars:              cfg.Vars,
		Templates:         cfg.Templates,
		Stdout:            param.Stdout,
		Stderr:            param.Stderr,
		CombinedOutput:    param.CombinedOutput,
		ExitCode:          param.ExitCode,
		ErrorMessages:     errMsgs,
		CreatedResources:  result.CreatedResources,
		UpdatedResources:  result.UpdatedResources,
		DeletedResources:  result.DeletedResources,
		ReplacedResources: result.ReplacedResources,
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

	postOpt := PostOptions{
		Number:   cfg.PR.Number,
		Revision: cfg.PR.Revision,
	}
	if !g.isHideComment() {
		return result.ExitCode, g.client.Comment.Post(ctx, body, postOpt)
	}
	logE := logrus.WithFields(logrus.Fields{
		"program": "tfcmt",
	})

	embeddedComment := "<!-- tfcmt:plan -->"
	if target != "" {
		embeddedComment = "<!-- tfcmt:plan:" + target + " -->"
	}
	logE.WithFields(logrus.Fields{
		"comment": embeddedComment,
	}).Debug("embedded HTML comment")
	// embed HTML tag to hide old comments
	body += "\n" + embeddedComment

	// get an authenticated user to hide comments only posted by the same user
	login, err := g.client.User.Get(ctx)
	if err != nil {
		logE.WithError(err).Warn("get the authenticated user")
		login = ""
	}
	cmts, err := g.client.Comment.list(ctx, ListOptions{
		PRNumber: cfg.PR.Number,
		Owner:    cfg.Owner,
		Repo:     cfg.Repo,
	})
	if err != nil {
		logE.WithError(err).Error("list pull request comments")
		return result.ExitCode, g.client.Comment.Post(ctx, body, postOpt)
	}
	logE.WithFields(logrus.Fields{
		"count":     len(cmts),
		"owner":     cfg.Owner,
		"repo":      cfg.Repo,
		"pr_number": cfg.PR.Number,
	}).Debug("list of pull request comments")
	extractedComments := []Comment{}
	for _, cmt := range cmts {
		if isExcludedComment(cmt, login) {
			continue
		}
		// hide comments which include `embeddedComment`
		if !strings.Contains(cmt.Body, embeddedComment) {
			continue
		}
		extractedComments = append(extractedComments, cmt)
	}
	logE.WithFields(logrus.Fields{
		"count": len(extractedComments),
	}).Debug("list of commements which will be hidden")
	if err := g.client.Comment.Post(ctx, body, postOpt); err != nil {
		return result.ExitCode, err
	}
	// hide comments
	for _, cmt := range extractedComments {
		if err := g.client.Comment.hide(ctx, cmt.ID); err != nil {
			logE.WithError(err).Error("hide a comment")
		}
		logE.WithFields(logrus.Fields{
			"node_id": cmt.ID,
		}).Debug("hide a comment")
	}
	return result.ExitCode, nil
}

func isExcludedComment(cmt Comment, login string) bool {
	if !cmt.ViewerCanMinimize {
		return true
	}
	if cmt.IsMinimized {
		return true
	}
	// GitHub Actions's GITHUB_TOKEN secret doesn't have a permission to get an authenticated user.
	// So if `login` is empty, we give up filtering comments by login.
	if login != "" && cmt.Author.Login != login {
		return true
	}
	return false
}

func (g *NotifyService) isHideComment() bool {
	if g.client.Config.HideOldComment.Disable {
		return false
	}
	if _, isApply := g.client.Config.Parser.(*terraform.ApplyParser); isApply {
		return false
	}
	if g.client.Config.PR.Number == 0 {
		return false
	}
	return true
}

func (g *NotifyService) updateLabels(ctx context.Context, result terraform.ParseResult) []string {
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
	labels, _, err := g.client.API.IssuesListLabels(ctx, cfg.PR.Number, nil)
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
