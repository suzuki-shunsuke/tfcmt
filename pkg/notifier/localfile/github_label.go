package localfile

import (
	"context"

	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/notifier/github"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/terraform"
)

// creates a minimal github.NotifyService and updates labels only.
func UpdateGitHubLabels(ctx context.Context, cfg *GitHubLabelConfig, result terraform.ParseResult) []string {
	if cfg == nil {
		return []string{"GitHubLabelConfig is nil"}
	}

	labels, ok := cfg.Labels.(github.ResultLabels)
	if !ok {
		return []string{"cfg.Labels is not of type github.ResultLabels"}
	}
	ghcfg := &github.Config{
		BaseURL:         cfg.BaseURL,
		GraphQLEndpoint: cfg.GraphQLEndpoint,
		Owner:           cfg.Owner,
		Repo:            cfg.Repo,
		PR: github.PullRequest{
			Revision: cfg.Revision,
			Number:   cfg.PRNumber,
		},
		ResultLabels: labels,
	}
	client, err := github.NewClient(ctx, ghcfg)
	if err != nil {
		return []string{"failed to create GitHub client: " + err.Error()}
	}
	if client.Notify == nil {
		return []string{"GitHub NotifyService is nil"}
	}
	// Call the label update logic directly
	return client.Notify.UpdateLabelsOnly(ctx, result)
}
