package github

import (
	"context"

	"github.com/google/go-github/v55/github"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/terraform"
)

type fakeAPI struct {
	API
	FakeIssuesCreateComment                    func(ctx context.Context, number int, comment *github.IssueComment) (*github.IssueComment, *github.Response, error)
	FakeIssuesListLabels                       func(ctx context.Context, number int, opts *github.ListOptions) ([]*github.Label, *github.Response, error)
	FakeIssuesAddLabels                        func(ctx context.Context, number int, labels []string) ([]*github.Label, *github.Response, error)
	FakeIssuesRemoveLabel                      func(ctx context.Context, number int, label string) (*github.Response, error)
	FakeRepositoriesCreateComment              func(ctx context.Context, sha string, comment *github.RepositoryComment) (*github.RepositoryComment, *github.Response, error)
	FakeRepositoriesListCommits                func(ctx context.Context, opt *github.CommitsListOptions) ([]*github.RepositoryCommit, *github.Response, error)
	FakeRepositoriesGetCommit                  func(ctx context.Context, sha string) (*github.RepositoryCommit, *github.Response, error)
	FakePullRequestsListPullRequestsWithCommit func(ctx context.Context, sha string, opt *github.ListOptions) ([]*github.PullRequest, *github.Response, error)
}

func (g *fakeAPI) IssuesCreateComment(ctx context.Context, number int, comment *github.IssueComment) (*github.IssueComment, *github.Response, error) {
	return g.FakeIssuesCreateComment(ctx, number, comment)
}

func (g *fakeAPI) IssuesListLabels(ctx context.Context, number int, opt *github.ListOptions) ([]*github.Label, *github.Response, error) {
	return g.FakeIssuesListLabels(ctx, number, opt)
}

func (g *fakeAPI) IssuesAddLabels(ctx context.Context, number int, labels []string) ([]*github.Label, *github.Response, error) {
	return g.FakeIssuesAddLabels(ctx, number, labels)
}

func (g *fakeAPI) IssuesRemoveLabel(ctx context.Context, number int, label string) (*github.Response, error) {
	return g.FakeIssuesRemoveLabel(ctx, number, label)
}

func (g *fakeAPI) RepositoriesCreateComment(ctx context.Context, sha string, comment *github.RepositoryComment) (*github.RepositoryComment, *github.Response, error) {
	return g.FakeRepositoriesCreateComment(ctx, sha, comment)
}

func (g *fakeAPI) RepositoriesListCommits(ctx context.Context, opt *github.CommitsListOptions) ([]*github.RepositoryCommit, *github.Response, error) {
	return g.FakeRepositoriesListCommits(ctx, opt)
}

func (g *fakeAPI) RepositoriesGetCommit(ctx context.Context, sha string) (*github.RepositoryCommit, *github.Response, error) {
	return g.FakeRepositoriesGetCommit(ctx, sha)
}

func (g *fakeAPI) PullRequestsListPullRequestsWithCommit(ctx context.Context, sha string, opt *github.ListOptions) ([]*github.PullRequest, *github.Response, error) {
	return g.FakePullRequestsListPullRequestsWithCommit(ctx, sha, opt)
}

func newFakeAPI() fakeAPI {
	return fakeAPI{
		FakeIssuesCreateComment: func(ctx context.Context, number int, comment *github.IssueComment) (*github.IssueComment, *github.Response, error) {
			return &github.IssueComment{
				ID:   github.Int64(371748792),
				Body: github.String("comment 1"),
			}, nil, nil
		},
		FakeIssuesListLabels: func(ctx context.Context, number int, opts *github.ListOptions) ([]*github.Label, *github.Response, error) {
			labels := []*github.Label{
				{
					ID:   github.Int64(371748792),
					Name: github.String("label 1"),
				},
				{
					ID:   github.Int64(371765743),
					Name: github.String("label 2"),
				},
			}
			return labels, nil, nil
		},
		FakeIssuesAddLabels: func(ctx context.Context, number int, labels []string) ([]*github.Label, *github.Response, error) {
			return nil, nil, nil
		},
		FakeIssuesRemoveLabel: func(ctx context.Context, number int, label string) (*github.Response, error) {
			return nil, nil
		},
		FakeRepositoriesCreateComment: func(ctx context.Context, sha string, comment *github.RepositoryComment) (*github.RepositoryComment, *github.Response, error) {
			return &github.RepositoryComment{
				ID:       github.Int64(28427394),
				CommitID: github.String("04e0917e448b662c2b16330fad50e97af16ff27a"),
				Body:     github.String("comment 1"),
			}, nil, nil
		},
		FakeRepositoriesListCommits: func(ctx context.Context, opt *github.CommitsListOptions) ([]*github.RepositoryCommit, *github.Response, error) {
			commits := []*github.RepositoryCommit{
				{
					SHA: github.String("04e0917e448b662c2b16330fad50e97af16ff27a"),
				},
				{
					SHA: github.String("04e0917e448b662c2b16330fad50e97af16ff27b"),
				},
				{
					SHA: github.String("04e0917e448b662c2b16330fad50e97af16ff27c"),
				},
			}
			return commits, nil, nil
		},
		FakeRepositoriesGetCommit: func(ctx context.Context, sha string) (*github.RepositoryCommit, *github.Response, error) {
			return &github.RepositoryCommit{
				SHA: github.String(sha),
				Commit: &github.Commit{
					Message: github.String(sha),
				},
			}, nil, nil
		},
		FakePullRequestsListPullRequestsWithCommit: func(ctx context.Context, sha string, opt *github.ListOptions) ([]*github.PullRequest, *github.Response, error) {
			return []*github.PullRequest{
				{
					State:  github.String("open"),
					Number: github.Int(1),
				},
				{
					State:  github.String("closed"),
					Number: github.Int(2),
				},
			}, nil, nil
		},
	}
}

func newFakeConfig() Config {
	return Config{
		Token: "token",
		Owner: "owner",
		Repo:  "repo",
		PR: PullRequest{
			Revision: "abcd",
			Number:   1,
		},
		Parser:   terraform.NewPlanParser(),
		Template: terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
	}
}
