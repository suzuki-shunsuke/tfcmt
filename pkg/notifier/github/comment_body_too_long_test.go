//nolint:revive
package github

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/go-github/v93/github"
)

func TestIsBodyTooLongError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil",
			err:  nil,
			want: false,
		},
		{
			name: "unrelated error",
			err:  errors.New("network timeout"),
			want: false,
		},
		{
			name: "422 with body too long message in Errors",
			err: &github.ErrorResponse{
				Message: "Validation Failed",
				Errors: []github.Error{
					{Message: "Body is too long (maximum is 65536 characters)"},
				},
			},
			want: true,
		},
		{
			name: "422 with body too long message at top level",
			err: &github.ErrorResponse{
				Message: "Validation Failed [{Resource:IssueComment Field:data Code:unprocessable Message:Body is too long (maximum is 65536 characters)}]",
			},
			want: true,
		},
		{
			name: "422 with another validation message",
			err: &github.ErrorResponse{
				Message: "Validation Failed",
				Errors: []github.Error{
					{Message: "some field is invalid"},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := isBodyTooLongError(tt.err); got != tt.want {
				t.Errorf("isBodyTooLongError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

func TestExtractResultSummary(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "plan line",
			body: "## Plan Result\n\nsome text\nPlan: 20000 to add, 0 to change, 0 to destroy.\nmore text",
			want: "<pre><code>Plan: 20000 to add, 0 to change, 0 to destroy.</code></pre>\n",
		},
		{
			name: "apply line",
			body: "Apply complete! Resources: 0 added, 0 changed, 0 destroyed.",
			want: "<pre><code>Apply complete! Resources: 0 added, 0 changed, 0 destroyed.</code></pre>\n",
		},
		{
			name: "no changes",
			body: "No changes. Your infrastructure matches the configuration.",
			want: "<pre><code>No changes. Your infrastructure matches the configuration.</code></pre>\n",
		},
		{
			name: "nothing found",
			body: "some random output without a result line",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := extractResultSummary(tt.body); got != tt.want {
				t.Errorf("extractResultSummary(%q) = %q, want %q", tt.body, got, tt.want)
			}
		})
	}
}

func TestPostBodyTooLongFallback(t *testing.T) { //nolint:tparallel
	t.Setenv("GITHUB_TOKEN", "xxx")
	cfg := newFakeConfig()
	cfg.CI = "https://github.com/foo/bar/actions/runs/123"
	client, err := NewClient(t.Context(), &cfg)
	if err != nil {
		t.Fatal(err)
	}

	var posted []string
	tooLong := true // first call fails with body-too-long, second succeeds
	api := fakeAPI{
		FakeIssuesCreateComment: func(ctx context.Context, number int, comment github.IssueCommentRequest) (*github.IssueComment, *github.Response, error) {
			posted = append(posted, comment.Body)
			if tooLong {
				tooLong = false
				return nil, nil, &github.ErrorResponse{
					Message: "Validation Failed",
					Errors: []github.Error{
						{Message: "Body is too long (maximum is 65536 characters)"},
					},
				}
			}
			return &github.IssueComment{ID: new(int64(1))}, nil, nil
		},
	}
	client.API = &api

	body := "## Plan Result\n\n" + strings.Repeat("x", 70000) + "\nPlan: 20000 to add, 0 to change, 0 to destroy."
	opt := &PostOptions{Number: 1}
	if err := client.Comment.Post(t.Context(), body, opt); err != nil {
		t.Fatalf("Post should succeed with a fallback comment, got %v", err)
	}
	if len(posted) != 2 {
		t.Fatalf("two comments should be posted (original + fallback), got %d", len(posted))
	}
	fallback := posted[1]
	if !strings.Contains(fallback, "tfcmt couldn't post a comment because the comment body is too long") {
		t.Errorf("fallback comment should explain the reason, got %q", fallback)
	}
	if !strings.Contains(fallback, "https://github.com/foo/bar/actions/runs/123") {
		t.Errorf("fallback comment should link the CI log, got %q", fallback)
	}
	if !strings.Contains(fallback, "Plan: 20000 to add, 0 to change, 0 to destroy.") {
		t.Errorf("fallback comment should keep the result summary, got %q", fallback)
	}
	if len(fallback) >= 65536 {
		t.Errorf("fallback comment must be under the limit, got %d chars", len(fallback))
	}
}
