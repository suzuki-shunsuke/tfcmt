package github

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/go-github/v33/github"
	"github.com/shurcooL/githubv4"
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

type ListOptions struct {
	PRNumber int
	Owner    string
	Repo     string
}

type Comment struct {
	ID     string
	Body   string
	Author struct {
		Login string
	}
	CreatedAt string
	// TODO remove
	IsMinimized       bool
	ViewerCanMinimize bool
}

func (g *CommentService) listIssueComment(ctx context.Context, opt ListOptions) ([]Comment, error) {
	// https://github.com/shurcooL/githubv4#pagination
	var q struct {
		Repository struct {
			Issue struct {
				Comments struct {
					Nodes    []Comment
					PageInfo struct {
						EndCursor   githubv4.String
						HasNextPage bool
					}
				} `graphql:"comments(first: 100, after: $commentsCursor)"` // 100 per page.
			} `graphql:"issue(number: $issueNumber)"`
		} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
	}
	variables := map[string]interface{}{
		"repositoryOwner": githubv4.String(opt.Owner),
		"repositoryName":  githubv4.String(opt.Repo),
		"issueNumber":     githubv4.Int(opt.PRNumber),
		"commentsCursor":  (*githubv4.String)(nil), // Null after argument to get first page.
	}

	var allComments []Comment
	for {
		if err := g.client.v4Client.Query(ctx, &q, variables); err != nil {
			return nil, fmt.Errorf("list issue comments by GitHub API: %w", err)
		}
		allComments = append(allComments, q.Repository.Issue.Comments.Nodes...)
		if !q.Repository.Issue.Comments.PageInfo.HasNextPage {
			break
		}
		variables["commentsCursor"] = githubv4.NewString(q.Repository.Issue.Comments.PageInfo.EndCursor)
	}
	return allComments, nil
}

func (g *CommentService) listPRComment(ctx context.Context, opt ListOptions) ([]Comment, error) {
	// https://github.com/shurcooL/githubv4#pagination
	var q struct {
		Repository struct {
			PullRequest struct {
				Comments struct {
					Nodes    []Comment
					PageInfo struct {
						EndCursor   githubv4.String
						HasNextPage bool
					}
				} `graphql:"comments(first: 100, after: $commentsCursor)"` // 100 per page.
			} `graphql:"pullRequest(number: $issueNumber)"`
		} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
	}
	variables := map[string]interface{}{
		"repositoryOwner": githubv4.String(opt.Owner),
		"repositoryName":  githubv4.String(opt.Repo),
		"issueNumber":     githubv4.Int(opt.PRNumber),
		"commentsCursor":  (*githubv4.String)(nil), // Null after argument to get first page.
	}

	var allComments []Comment
	for {
		if err := g.client.v4Client.Query(ctx, &q, variables); err != nil {
			return nil, fmt.Errorf("list issue comments by GitHub API: %w", err)
		}
		allComments = append(allComments, q.Repository.PullRequest.Comments.Nodes...)
		if !q.Repository.PullRequest.Comments.PageInfo.HasNextPage {
			break
		}
		variables["commentsCursor"] = githubv4.NewString(q.Repository.PullRequest.Comments.PageInfo.EndCursor)
	}
	return allComments, nil
}

func (g *CommentService) list(ctx context.Context, opt ListOptions) ([]Comment, error) {
	cmts, prErr := g.listPRComment(ctx, opt)
	if prErr == nil {
		return cmts, nil
	}
	cmts, err := g.listIssueComment(ctx, opt)
	if err == nil {
		return cmts, nil
	}
	return nil, fmt.Errorf("get pull request or issue comments: %w, %v", prErr, err)
}

func (g *CommentService) hide(ctx context.Context, nodeID string) error {
	var m struct {
		MinimizeComment struct {
			MinimizedComment struct {
				MinimizedReason   githubv4.String
				IsMinimized       githubv4.Boolean
				ViewerCanMinimize githubv4.Boolean
			}
		} `graphql:"minimizeComment(input:$input)"`
	}
	input := githubv4.MinimizeCommentInput{
		Classifier: githubv4.ReportedContentClassifiersOutdated,
		SubjectID:  nodeID,
	}
	if err := g.client.v4Client.Mutate(ctx, &m, input, nil); err != nil {
		return fmt.Errorf("hide an old comment: %w", err)
	}
	return nil
}
