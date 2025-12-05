package github

import (
	"log/slog"
	"os"

	"github.com/suzuki-shunsuke/github-comment-metadata/metadata"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

// NotifyService handles communication with the notification related
// methods of GitHub API
type NotifyService service

func (g *NotifyService) getPatchedComment(logger *slog.Logger, comments []*IssueComment, target string) *IssueComment {
	var cmt *IssueComment
	for i, comment := range comments {
		logger := logger.With(
			"comment_database_id", comment.DatabaseID,
			"comment_index", i,
		)
		data := &Metadata{}
		f, err := metadata.Extract(comment.Body, data)
		if err != nil {
			slogerr.WithError(logger, err).Debug("extract metadata from comment")
			continue
		}
		if !f {
			logger.Debug("metadata isn't found")
			continue
		}
		if data.Program != "tfcmt" {
			logger.Debug("Program isn't tfcmt")
			continue
		}
		if data.Command != "plan" {
			logger.Debug("Command isn't plan")
			continue
		}
		if data.Target != target {
			logger.Debug("target is different")
			continue
		}
		if comment.IsMinimized {
			logger.Debug("comment is hidden")
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
	vars := make(map[string]any, len(cfg.EmbeddedVarNames))
	for _, name := range cfg.EmbeddedVarNames {
		vars[name] = cfg.Vars[name]
	}

	data := map[string]any{
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
