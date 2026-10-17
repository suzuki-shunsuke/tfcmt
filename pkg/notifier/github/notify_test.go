package github

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/go-github/v92/github"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/notifier"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/terraform"
)

// commentCounter counts comment posts without changing the fake responses.
type commentCounter struct {
	fakeAPI

	issues int
	repos  int
}

func (g *commentCounter) IssuesCreateComment(ctx context.Context, number int, comment github.IssueCommentRequest) (*github.IssueComment, *github.Response, error) {
	g.issues++
	return g.fakeAPI.IssuesCreateComment(ctx, number, comment)
}

func (g *commentCounter) RepositoriesCreateComment(ctx context.Context, sha string, comment github.CreateCommitCommentRequest) (*github.RepositoryComment, *github.Response, error) {
	g.repos++
	return g.fakeAPI.RepositoriesCreateComment(ctx, sha, comment)
}

func TestNotifyApply(t *testing.T) { //nolint:tparallel
	t.Setenv("GITHUB_TOKEN", "xxx")
	logger := slog.New(slog.DiscardHandler)
	testCases := []struct {
		name          string
		config        Config
		ok            bool
		comments      int
		checkComments bool
		paramExec     notifier.ParamExec
	}{
		{
			name: "case 8",
			// apply case without merge commit
			config: Config{
				Owner: "owner",
				Repo:  "repo",
				PR: PullRequest{
					Revision: "revision",
					Number:   0, // For apply, it is always 0
				},
				Parser:             terraform.NewApplyParser(),
				Template:           terraform.NewApplyTemplate(terraform.DefaultApplyTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: notifier.ParamExec{
				Stdout:   "Apply complete!",
				ExitCode: 0,
			},
			ok: true,
		},
		{
			name: "case 9",
			// apply case as merge commit
			// TODO(drlau): validate cfg.PR.Number = 123
			config: Config{
				Owner: "owner",
				Repo:  "repo",
				PR: PullRequest{
					Revision: "Merge pull request #123 from suzuki-shunsuke/tfcmt",
					Number:   0, // For apply, it is always 0
				},
				Parser:             terraform.NewApplyParser(),
				Template:           terraform.NewApplyTemplate(terraform.DefaultApplyTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: notifier.ParamExec{
				Stdout:   "Apply complete!",
				ExitCode: 0,
			},
			ok: true,
		},
		{
			name: "skip comment when no resources change",
			config: Config{
				Owner: "owner",
				Repo:  "repo",
				PR: PullRequest{
					Revision: "revision",
					Number:   1,
				},
				Parser:             terraform.NewApplyParser(),
				Template:           terraform.NewApplyTemplate(terraform.DefaultApplyTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
				SkipNoChanges:      true,
			},
			paramExec: notifier.ParamExec{
				CombinedOutput: "Apply complete! Resources: 0 added, 0 changed, 0 destroyed.\n",
				ExitCode:       0,
			},
			ok:            true,
			checkComments: true,
			comments:      0,
		},
		{
			name: "post comment when resources change and skip is enabled",
			config: Config{
				Owner: "owner",
				Repo:  "repo",
				PR: PullRequest{
					Revision: "revision",
					Number:   1,
				},
				Parser:             terraform.NewApplyParser(),
				Template:           terraform.NewApplyTemplate(terraform.DefaultApplyTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
				SkipNoChanges:      true,
			},
			paramExec: notifier.ParamExec{
				CombinedOutput: "Apply complete! Resources: 1 added, 0 changed, 0 destroyed.\n",
				ExitCode:       0,
			},
			ok:            true,
			checkComments: true,
			comments:      1,
		},
		{
			name: "post comment when skip is disabled and no resources change",
			config: Config{
				Owner: "owner",
				Repo:  "repo",
				PR: PullRequest{
					Revision: "revision",
					Number:   1,
				},
				Parser:             terraform.NewApplyParser(),
				Template:           terraform.NewApplyTemplate(terraform.DefaultApplyTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: notifier.ParamExec{
				CombinedOutput: "Apply complete! Resources: 0 added, 0 changed, 0 destroyed.\n",
				ExitCode:       0,
			},
			ok:            true,
			checkComments: true,
			comments:      1,
		},
		{
			name: "post comment when apply fails and skip is enabled",
			config: Config{
				Owner: "owner",
				Repo:  "repo",
				PR: PullRequest{
					Revision: "revision",
					Number:   1,
				},
				Parser:             terraform.NewApplyParser(),
				Template:           terraform.NewApplyTemplate(terraform.DefaultApplyTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
				SkipNoChanges:      true,
			},
			paramExec: notifier.ParamExec{
				CombinedOutput: "Error: failed to send enable services request\n",
				ExitCode:       1,
			},
			ok:            true,
			checkComments: true,
			comments:      1,
		},
		{
			name: "post comment when apply output cannot be parsed and skip is enabled",
			config: Config{
				Owner: "owner",
				Repo:  "repo",
				PR: PullRequest{
					Revision: "revision",
					Number:   1,
				},
				Parser:             terraform.NewApplyParser(),
				Template:           terraform.NewApplyTemplate(terraform.DefaultApplyTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
				SkipNoChanges:      true,
			},
			paramExec: notifier.ParamExec{
				CombinedOutput: "not terraform output\n",
				ExitCode:       1,
			},
			ok:            true,
			checkComments: true,
			comments:      1,
		},
	}

	for i, testCase := range testCases {
		if testCase.name == "" {
			t.Fatalf("testCase.name is required: index: %d", i)
		}
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			cfg := testCase.config
			client, err := NewClient(t.Context(), &cfg)
			if err != nil {
				t.Fatal(err)
			}
			api := &commentCounter{fakeAPI: newFakeAPI()}
			client.API = api
			paramExec := testCase.paramExec
			if err := client.Notify.Apply(t.Context(), logger, &paramExec); (err == nil) != testCase.ok {
				t.Errorf("got error %v", err)
			}
			if testCase.checkComments && api.issues+api.repos != testCase.comments {
				t.Errorf("IssuesCreateComment=%d RepositoriesCreateComment=%d, want %d comments", api.issues, api.repos, testCase.comments)
			}
		})
	}
}

func TestNotifyPlan(t *testing.T) { //nolint:tparallel
	t.Setenv("GITHUB_TOKEN", "xxx")
	logger := slog.New(slog.DiscardHandler)
	testCases := []struct {
		name      string
		config    Config
		ok        bool
		paramExec notifier.ParamExec
	}{
		{
			name: "case 0",
			// invalid body (cannot parse)
			config: Config{
				Owner: "owner",
				Repo:  "repo",
				PR: PullRequest{
					Revision: "abcd",
					Number:   1,
				},
				Parser:             terraform.NewPlanParser(),
				Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: notifier.ParamExec{
				Stdout:   "body",
				ExitCode: 1,
			},
			ok: true,
		},
		{
			name: "case 1",
			// invalid pr
			config: Config{
				Owner: "owner",
				Repo:  "repo",
				PR: PullRequest{
					Revision: "",
					Number:   0,
				},
				Parser:             terraform.NewPlanParser(),
				Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: notifier.ParamExec{
				Stdout:   "Plan: 1 to add",
				ExitCode: 0,
			},
			ok: false,
		},
		{
			name: "case 2",
			// valid, error
			config: Config{
				Owner: "owner",
				Repo:  "repo",
				PR: PullRequest{
					Revision: "",
					Number:   1,
				},
				Parser:             terraform.NewPlanParser(),
				Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: notifier.ParamExec{
				Stdout:   "Error: hoge",
				ExitCode: 1,
			},
			ok: true,
		},
		{
			name: "case 3",
			// valid, and isPR
			config: Config{
				Owner: "owner",
				Repo:  "repo",
				PR: PullRequest{
					Revision: "",
					Number:   1,
				},
				Parser:             terraform.NewPlanParser(),
				Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: notifier.ParamExec{
				Stdout:   "Plan: 1 to add",
				ExitCode: 2,
			},
			ok: true,
		},
		{
			name: "case 4",
			// valid, and isRevision
			config: Config{
				Owner: "owner",
				Repo:  "repo",
				PR: PullRequest{
					Revision: "revision-revision",
					Number:   0,
				},
				Parser:             terraform.NewPlanParser(),
				Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: notifier.ParamExec{
				Stdout:   "Plan: 1 to add",
				ExitCode: 2,
			},
			ok: true,
		},
		{
			name: "case 5",
			// valid, and contains destroy
			// TODO(dtan4): check two comments were made actually
			config: Config{
				Owner: "owner",
				Repo:  "repo",
				PR: PullRequest{
					Revision: "",
					Number:   1,
				},
				Parser:             terraform.NewPlanParser(),
				Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: notifier.ParamExec{
				Stdout:   "Plan: 1 to add, 1 to destroy",
				ExitCode: 2,
			},
			ok: true,
		},
		{
			name: "case 6",
			// valid with no changes
			// TODO(drlau): check that the label was actually added
			config: Config{
				Owner: "owner",
				Repo:  "repo",
				PR: PullRequest{
					Revision: "",
					Number:   1,
				},
				Parser:             terraform.NewPlanParser(),
				Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
				ResultLabels: ResultLabels{
					AddOrUpdateLabel: "add-or-update",
					DestroyLabel:     "destroy",
					NoChangesLabel:   "no-changes",
					PlanErrorLabel:   "error",
				},
			},
			paramExec: notifier.ParamExec{
				Stdout:   "No changes. Infrastructure is up-to-date.",
				ExitCode: 0,
			},
			ok: true,
		},
		{
			name: "case 7",
			// valid, contains destroy, but not to notify
			config: Config{
				Owner: "owner",
				Repo:  "repo",
				PR: PullRequest{
					Revision: "",
					Number:   1,
				},
				Parser:             terraform.NewPlanParser(),
				Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: notifier.ParamExec{
				Stdout:   "Plan: 1 to add, 1 to destroy",
				ExitCode: 2,
			},
			ok: true,
		},
	}

	for i, testCase := range testCases {
		if testCase.name == "" {
			t.Fatalf("testCase.name is required: index: %d", i)
		}
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			cfg := testCase.config
			client, err := NewClient(t.Context(), &cfg)
			if err != nil {
				t.Fatal(err)
			}
			api := newFakeAPI()
			client.API = &api
			paramExec := testCase.paramExec
			if err := client.Notify.Plan(t.Context(), logger, &paramExec); (err == nil) != testCase.ok {
				t.Errorf("got error %v", err)
			}
		})
	}
}
