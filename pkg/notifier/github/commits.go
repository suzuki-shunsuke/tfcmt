package github

import (
	"context"
	"errors"
	"strconv"
	"strings"
)

// CommitsService handles communication with the commits related
// methods of GitHub API
type CommitsService service

func (g *CommitsService) MergedPRNumber(ctx context.Context, revision string) (int, error) {
	commit, _, err := g.client.API.RepositoriesGetCommit(ctx, revision)
	if err != nil {
		return 0, err
	}

	message := commit.Commit.GetMessage()
	if !strings.HasPrefix(message, "Merge pull request #") {
		return 0, errors.New("not a merge commit")
	}

	message = strings.TrimPrefix(message, "Merge pull request #")
	i := strings.Index(message, " from")
	if i >= 0 {
		return strconv.Atoi(message[0:i])
	}

	return 0, errors.New("not a merge commit")
}
