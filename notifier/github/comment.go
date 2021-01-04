package github

import (
	"context"
	"errors"

	"github.com/google/go-github/github"
)

// CommentService handles communication with the comment related
// methods of GitHub API
type CommentService service

// PostOptions specifies the optional parameters to post comments to a pull request
type PostOptions struct {
	Number   int
	Revision string
}

// Post posts comment
func (g *CommentService) Post(ctx context.Context, body string, opt PostOptions) error {
	if opt.Number != 0 {
		_, _, err := g.client.API.IssuesCreateComment(
			ctx,
			opt.Number,
			&github.IssueComment{Body: &body},
		)
		return err
	}
	if opt.Revision != "" {
		_, _, err := g.client.API.RepositoriesCreateComment(
			ctx,
			opt.Revision,
			&github.RepositoryComment{Body: &body},
		)
		return err
	}
	return errors.New("github.comment.post: Number or Revision is required")
}

// List lists comments on GitHub issues/pull requests
func (g *CommentService) List(ctx context.Context, number int) ([]*github.IssueComment, error) {
	comments, _, err := g.client.API.IssuesListComments(
		ctx,
		number,
		&github.IssueListCommentsOptions{},
	)
	return comments, err
}
