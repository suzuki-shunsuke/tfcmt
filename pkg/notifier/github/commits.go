package github

import (
	"context"
	"errors"

	"github.com/google/go-github/v49/github"
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
	prs, _, err := g.client.API.PullRequestsListPullRequestsWithCommit(ctx, sha, &github.PullRequestListOptions{
		State: string(state),
		Sort:  "updated",
	})
	if err != nil {
		return 0, err
	}
	if len(prs) == 0 {
		return 0, errors.New("associated pull request isn't found")
	}
	return prs[0].GetNumber(), nil
}
