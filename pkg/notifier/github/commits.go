package github

import (
	"context"
	"errors"
)

// CommitsService handles communication with the commits related
// methods of GitHub API
type CommitsService service

func (g *CommitsService) MergedPRNumber(ctx context.Context, sha string) (int, error) {
	prs, _, err := g.client.API.PullRequestsListPullRequestsWithCommit(ctx, sha, nil)
	if err != nil {
		return 0, err
	}
	for _, pr := range prs {
		if pr.GetState() != "closed" {
			continue
		}
		return pr.GetNumber(), nil
	}
	return 0, errors.New("associated pull request isn't found")
}
