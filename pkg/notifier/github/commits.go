package github

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/google/go-github/v43/github"
)

// CommitsService handles communication with the commits related
// methods of GitHub API
type CommitsService service

// List lists commits on a repository
func (g *CommitsService) List(ctx context.Context, revision string) ([]string, error) {
	if revision == "" {
		return []string{}, errors.New("no revision specified")
	}
	commits, _, err := g.client.API.RepositoriesListCommits(
		ctx,
		&github.CommitsListOptions{SHA: revision},
	)
	if err != nil {
		return nil, err
	}
	shas := make([]string, len(commits))
	for i, commit := range commits {
		shas[i] = *commit.SHA
	}
	return shas, nil
}

// Last returns the hash of the previous commit of the given commit
func (g *CommitsService) lastOne(commits []string, revision string) (string, error) {
	if revision == "" {
		return "", errors.New("no revision specified")
	}
	if len(commits) == 0 {
		return "", errors.New("no commits")
	}
	// e.g.
	// a0ce5bf 2018/04/05 20:50:01 (HEAD -> master, origin/master)
	// 5166cfc 2018/04/05 20:40:12
	// 74c4d6e 2018/04/05 20:34:31
	// 9260c54 2018/04/05 20:16:20
	return commits[1], nil
}

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
