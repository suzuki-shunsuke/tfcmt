package github

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/github-comment-metadata/metadata"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/terraform"
)

// NotifyService handles communication with the notification related
// methods of GitHub API
type NotifyService service

// UpdateLabelsOnly allows updating PR labels without posting comments
func (g *NotifyService) UpdateLabelsOnly(ctx context.Context, result terraform.ParseResult) []string {
	return g.updateLabels(ctx, result)
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
		if data.Command != "plan" {
			logE.Debug("Command isn't plan")
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
	Command string
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
