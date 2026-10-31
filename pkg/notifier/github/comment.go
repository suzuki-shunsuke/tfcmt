package github

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/go-github/v93/github"
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
func (g *CommentService) Post(ctx context.Context, body string, opt *PostOptions) error {
	if opt.Number != 0 {
		_, _, err := g.client.API.IssuesCreateComment(
			ctx,
			opt.Number,
			github.IssueCommentRequest{Body: body},
		)
		if err != nil && isBodyTooLongError(err) {
			fallback := g.buildBodyTooLongFallbackComment(body, opt)
			_, _, err2 := g.client.API.IssuesCreateComment(
				ctx,
				opt.Number,
				github.IssueCommentRequest{Body: fallback},
			)
			if err2 != nil {
				return errors.Join(err, fmt.Errorf("post a fallback comment: %w", err2))
			}
			return nil
		}
		return err
	}
	if opt.Revision != "" {
		_, _, err := g.client.API.RepositoriesCreateComment(
			ctx,
			opt.Revision,
			github.CreateCommitCommentRequest{Body: body},
		)
		if err != nil && isBodyTooLongError(err) {
			fallback := g.buildBodyTooLongFallbackComment(body, opt)
			_, _, err2 := g.client.API.RepositoriesCreateComment(
				ctx,
				opt.Revision,
				github.CreateCommitCommentRequest{Body: fallback},
			)
			if err2 != nil {
				return errors.Join(err, fmt.Errorf("post a fallback comment: %w", err2))
			}
			return nil
		}
		return err
	}
	return errors.New("github.comment.post: Number or Revision is required")
}

const bodyTooLongMessage = "Body is too long (maximum is 65536 characters)"

// isBodyTooLongError returns true if the error is a GitHub API 422 error
// caused by the comment body exceeding the maximum length.
// https://github.com/suzuki-shunsuke/tfcmt/issues/1847
func isBodyTooLongError(err error) bool {
	if err == nil {
		return false
	}
	var ghErr *github.ErrorResponse
	if errors.As(err, &ghErr) {
		for _, e := range ghErr.Errors {
			if strings.Contains(e.Message, bodyTooLongMessage) {
				return true
			}
		}
		// Some responses only carry the message at the top level.
		if strings.Contains(ghErr.Message, bodyTooLongMessage) {
			return true
		}
	}
	return strings.Contains(err.Error(), bodyTooLongMessage)
}

// buildBodyTooLongFallbackComment builds a short fallback comment posted when
// the original comment body is too long for GitHub.
// The fallback keeps the plan/apply result summary visible and links to the
// CI log so the omitted detail can still be reached.
func (g *CommentService) buildBodyTooLongFallbackComment(body string, opt *PostOptions) string {
	cfg := g.client.Config
	var link string
	if cfg != nil && cfg.CI != "" {
		link = cfg.CI
	}
	var sb strings.Builder
	if link != "" {
		sb.WriteString("[CI link](")
		sb.WriteString(link)
		sb.WriteString(")\n\n")
	}
	sb.WriteString("> [!WARNING]\n")
	sb.WriteString("> tfcmt couldn't post a comment because the comment body is too long.\n")
	if link != "" {
		sb.WriteString("> So tfcmt posts this comment instead. [Please check the log](")
		sb.WriteString(link)
		sb.WriteString(").\n")
	} else {
		sb.WriteString("> So tfcmt posts this comment instead. Please check the log.\n")
	}
	sb.WriteString("\n")
	sb.WriteString(extractResultSummary(body))
	return sb.String()
}

// extractResultSummary tries to pull the `Plan: N to add, ...` /
// `Apply complete!` line out of a rendered comment body so the fallback still
// shows the result summary. Returns an empty string when nothing is found.
func extractResultSummary(body string) string {
	for _, line := range strings.Split(body, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "Plan: ") ||
			strings.HasPrefix(trimmed, "Apply complete!") ||
			strings.HasPrefix(trimmed, "No changes.") {
			return "<pre><code>" + trimmed + "</code></pre>\n"
		}
	}
	return ""
}

func (g *CommentService) Patch(ctx context.Context, body string, commentID int64) error {
	_, _, err := g.client.API.IssuesUpdateComment(
		ctx,
		commentID,
		github.IssueCommentRequest{Body: body},
	)
	return err
}

type IssueComment struct {
	DatabaseID  int
	Body        string
	IsMinimized bool
}

func (g *CommentService) List(ctx context.Context, owner, repo string, number int) ([]*IssueComment, error) {
	cmts, prErr := g.listPRComment(ctx, owner, repo, number)
	if prErr == nil {
		return cmts, nil
	}
	cmts, err := g.listIssueComment(ctx, owner, repo, number)
	if err == nil {
		return cmts, nil
	}
	return nil, fmt.Errorf("get pull request or issue comments: %w, %v", prErr, err) //nolint:errorlint
}

func (g *CommentService) listIssueComment(ctx context.Context, owner, repo string, number int) ([]*IssueComment, error) { //nolint:dupl
	// https://github.com/shurcooL/githubv4#pagination
	var q struct {
		Repository struct {
			Issue struct {
				Comments struct {
					Nodes    []*IssueComment
					PageInfo struct {
						EndCursor   githubv4.String
						HasNextPage bool
					}
				} `graphql:"comments(first: 100, after: $commentsCursor)"` // 100 per page.
			} `graphql:"issue(number: $issueNumber)"`
		} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
	}
	variables := map[string]any{
		"repositoryOwner": githubv4.String(owner),
		"repositoryName":  githubv4.String(repo),
		"issueNumber":     githubv4.Int(number),    //nolint:gosec // GitHub issue/PR numbers are within int32 range
		"commentsCursor":  (*githubv4.String)(nil), // Null after argument to get first page.
	}

	var allComments []*IssueComment
	for {
		if err := g.client.v4Client.Query(ctx, &q, variables); err != nil {
			return nil, fmt.Errorf("list issue comments by GitHub API: %w", err)
		}
		allComments = append(allComments, q.Repository.Issue.Comments.Nodes...)
		if !q.Repository.Issue.Comments.PageInfo.HasNextPage {
			break
		}
		variables["commentsCursor"] = new(q.Repository.Issue.Comments.PageInfo.EndCursor)
	}
	return allComments, nil
}

func (g *CommentService) listPRComment(ctx context.Context, owner, repo string, number int) ([]*IssueComment, error) { //nolint:dupl
	// https://github.com/shurcooL/githubv4#pagination
	var q struct {
		Repository struct {
			PullRequest struct {
				Comments struct {
					Nodes    []*IssueComment
					PageInfo struct {
						EndCursor   githubv4.String
						HasNextPage bool
					}
				} `graphql:"comments(first: 100, after: $commentsCursor)"` // 100 per page.
			} `graphql:"pullRequest(number: $issueNumber)"`
		} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
	}
	variables := map[string]any{
		"repositoryOwner": githubv4.String(owner),
		"repositoryName":  githubv4.String(repo),
		"issueNumber":     githubv4.Int(number),    //nolint:gosec // GitHub issue/PR numbers are within int32 range
		"commentsCursor":  (*githubv4.String)(nil), // Null after argument to get first page.
	}

	var allComments []*IssueComment
	for {
		if err := g.client.v4Client.Query(ctx, &q, variables); err != nil {
			return nil, fmt.Errorf("list issue comments by GitHub API: %w", err)
		}
		allComments = append(allComments, q.Repository.PullRequest.Comments.Nodes...)
		if !q.Repository.PullRequest.Comments.PageInfo.HasNextPage {
			break
		}
		variables["commentsCursor"] = new(q.Repository.PullRequest.Comments.PageInfo.EndCursor)
	}
	return allComments, nil
}
