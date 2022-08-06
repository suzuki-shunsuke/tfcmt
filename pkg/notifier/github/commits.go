package github

import (
	"context"
	"errors"
)

// CommitsService handles communication with the commits related
// methods of GitHub API
type CommitsService service

type PullRequestState string

const (
	PullRequestStateOpen   = "open"
	PullRequestStateClosed = "closed"
	PullRequestStateAll    = "all"
)

func (g *CommitsService) PRNumber(ctx context.Context, sha string, state PullRequestState) (int, error) {
	prs, _, err := g.client.API.PullRequestsListPullRequestsWithCommit(ctx, sha, nil)
	if err != nil {
		return 0, err
	}
	for _, pr := range prs {
		if pr.GetState() != string(state) {
			continue
		}
		return pr.GetNumber(), nil
	}
	return 0, errors.New("associated pull request isn't found")
}
