package github

import (
	"context"
	"testing"

	"github.com/suzuki-shunsuke/tfcmt/pkg/notifier"
	"github.com/suzuki-shunsuke/tfcmt/pkg/terraform"
)

func TestNotifyNotify(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name      string
		config    *Config
		ok        bool
		exitCode  int
		paramExec *notifier.ParamExec
	}{
		{
			name: "case 0",
			// invalid body (cannot parse)
			config: &Config{
				Token: "token",
				Owner: "owner",
				Repo:  "repo",
				PR: &PullRequest{
					Revision: "abcd",
					Number:   1,
				},
				Parser:             terraform.NewPlanParser(),
				Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: &notifier.ParamExec{
				Stdout:   "body",
				ExitCode: 1,
			},
			ok:       true,
			exitCode: 1,
		},
		{
			name: "case 1",
			// invalid pr
			config: &Config{
				Token: "token",
				Owner: "owner",
				Repo:  "repo",
				PR: &PullRequest{
					Revision: "",
					Number:   0,
				},
				Parser:             terraform.NewPlanParser(),
				Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: &notifier.ParamExec{
				Stdout:   "Plan: 1 to add",
				ExitCode: 0,
			},
			ok:       false,
			exitCode: 0,
		},
		{
			name: "case 2",
			// valid, error
			config: &Config{
				Token: "token",
				Owner: "owner",
				Repo:  "repo",
				PR: &PullRequest{
					Revision: "",
					Number:   1,
				},
				Parser:             terraform.NewPlanParser(),
				Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: &notifier.ParamExec{
				Stdout:   "Error: hoge",
				ExitCode: 1,
			},
			ok:       true,
			exitCode: 1,
		},
		{
			name: "case 3",
			// valid, and isPR
			config: &Config{
				Token: "token",
				Owner: "owner",
				Repo:  "repo",
				PR: &PullRequest{
					Revision: "",
					Number:   1,
				},
				Parser:             terraform.NewPlanParser(),
				Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: &notifier.ParamExec{
				Stdout:   "Plan: 1 to add",
				ExitCode: 2,
			},
			ok:       true,
			exitCode: 2,
		},
		{
			name: "case 4",
			// valid, and isRevision
			config: &Config{
				Token: "token",
				Owner: "owner",
				Repo:  "repo",
				PR: &PullRequest{
					Revision: "revision-revision",
					Number:   0,
				},
				Parser:             terraform.NewPlanParser(),
				Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: &notifier.ParamExec{
				Stdout:   "Plan: 1 to add",
				ExitCode: 2,
			},
			ok:       true,
			exitCode: 2,
		},
		{
			name: "case 5",
			// valid, and contains destroy
			// TODO(dtan4): check two comments were made actually
			config: &Config{
				Token: "token",
				Owner: "owner",
				Repo:  "repo",
				PR: &PullRequest{
					Revision: "",
					Number:   1,
				},
				Parser:             terraform.NewPlanParser(),
				Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: &notifier.ParamExec{
				Stdout:   "Plan: 1 to add, 1 to destroy",
				ExitCode: 2,
			},
			ok:       true,
			exitCode: 2,
		},
		{
			name: "case 6",
			// valid with no changes
			// TODO(drlau): check that the label was actually added
			config: &Config{
				Token: "token",
				Owner: "owner",
				Repo:  "repo",
				PR: &PullRequest{
					Revision: "",
					Number:   1,
				},
				Parser:             terraform.NewPlanParser(),
				Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
				ResultLabels: &ResultLabels{
					AddOrUpdateLabel: "add-or-update",
					DestroyLabel:     "destroy",
					NoChangesLabel:   "no-changes",
					PlanErrorLabel:   "error",
				},
			},
			paramExec: &notifier.ParamExec{
				Stdout:   "No changes. Infrastructure is up-to-date.",
				ExitCode: 0,
			},
			ok:       true,
			exitCode: 0,
		},
		{
			name: "case 7",
			// valid, contains destroy, but not to notify
			config: &Config{
				Token: "token",
				Owner: "owner",
				Repo:  "repo",
				PR: &PullRequest{
					Revision: "",
					Number:   1,
				},
				Parser:             terraform.NewPlanParser(),
				Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: &notifier.ParamExec{
				Stdout:   "Plan: 1 to add, 1 to destroy",
				ExitCode: 2,
			},
			ok:       true,
			exitCode: 2,
		},
		{
			name: "case 8",
			// apply case without merge commit
			config: &Config{
				Token: "token",
				Owner: "owner",
				Repo:  "repo",
				PR: &PullRequest{
					Revision: "revision",
					Number:   0, // For apply, it is always 0
				},
				Parser:             terraform.NewApplyParser(),
				Template:           terraform.NewApplyTemplate(terraform.DefaultApplyTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: &notifier.ParamExec{
				Stdout:   "Apply complete!",
				ExitCode: 0,
			},
			ok:       true,
			exitCode: 0,
		},
		{
			name: "case 9",
			// apply case as merge commit
			// TODO(drlau): validate cfg.PR.Number = 123
			config: &Config{
				Token: "token",
				Owner: "owner",
				Repo:  "repo",
				PR: &PullRequest{
					Revision: "Merge pull request #123 from suzuki-shunsuke/tfcmt",
					Number:   0, // For apply, it is always 0
				},
				Parser:             terraform.NewApplyParser(),
				Template:           terraform.NewApplyTemplate(terraform.DefaultApplyTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
			},
			paramExec: &notifier.ParamExec{
				Stdout:   "Apply complete!",
				ExitCode: 0,
			},
			ok:       true,
			exitCode: 0,
		},
	}

	for i, testCase := range testCases {
		testCase := testCase
		if testCase.name == "" {
			t.Fatalf("testCase.name is required: index: %d", i)
		}
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			client, err := NewClient(context.Background(), testCase.config)
			if err != nil {
				t.Fatal(err)
			}
			api := newFakeAPI()
			client.API = &api
			exitCode, err := client.Notify.Notify(context.Background(), testCase.paramExec)
			if (err == nil) != testCase.ok {
				t.Errorf("got error %v", err)
			}
			if exitCode != testCase.exitCode {
				t.Errorf("got %d but want %d", exitCode, testCase.exitCode)
			}
		})
	}
}
