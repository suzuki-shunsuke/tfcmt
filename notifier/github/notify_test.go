package github

import (
	"context"
	"testing"

	"github.com/mercari/tfnotify/notifier"
	"github.com/mercari/tfnotify/terraform"
)

func TestNotifyNotify(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		config   Config
		body     string
		ok       bool
		exitCode int
	}{
		{
			name: "case 0",
			// invalid body (cannot parse)
			config: Config{
				Token: "token",
				Owner: "owner",
				Repo:  "repo",
				PR: PullRequest{
					Revision: "abcd",
					Number:   1,
					Message:  "message",
				},
				Parser:   terraform.NewPlanParser(),
				Template: terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
			},
			body:     "body",
			ok:       false,
			exitCode: 0,
		},
		{
			name: "case 1",
			// invalid pr
			config: Config{
				Token: "token",
				Owner: "owner",
				Repo:  "repo",
				PR: PullRequest{
					Revision: "",
					Number:   0,
					Message:  "message",
				},
				Parser:   terraform.NewPlanParser(),
				Template: terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
			},
			body:     "Plan: 1 to add",
			ok:       false,
			exitCode: 0,
		},
		{
			name: "case 2",
			// valid, error
			config: Config{
				Token: "token",
				Owner: "owner",
				Repo:  "repo",
				PR: PullRequest{
					Revision: "",
					Number:   1,
					Message:  "message",
				},
				Parser:   terraform.NewPlanParser(),
				Template: terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
			},
			body:     "Error: hoge",
			ok:       true,
			exitCode: 0,
		},
		{
			name: "case 3",
			// valid, and isPR
			config: Config{
				Token: "token",
				Owner: "owner",
				Repo:  "repo",
				PR: PullRequest{
					Revision: "",
					Number:   1,
					Message:  "message",
				},
				Parser:   terraform.NewPlanParser(),
				Template: terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
			},
			body:     "Plan: 1 to add",
			ok:       true,
			exitCode: 0,
		},
		{
			name: "case 4",
			// valid, and isRevision
			config: Config{
				Token: "token",
				Owner: "owner",
				Repo:  "repo",
				PR: PullRequest{
					Revision: "revision-revision",
					Number:   0,
					Message:  "message",
				},
				Parser:   terraform.NewPlanParser(),
				Template: terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
			},
			body:     "Plan: 1 to add",
			ok:       true,
			exitCode: 0,
		},
		{
			name: "case 5",
			// valid, and contains destroy
			// TODO(dtan4): check two comments were made actually
			config: Config{
				Token: "token",
				Owner: "owner",
				Repo:  "repo",
				PR: PullRequest{
					Revision: "",
					Number:   1,
					Message:  "message",
				},
				Parser:                 terraform.NewPlanParser(),
				Template:               terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				DestroyWarningTemplate: terraform.NewDestroyWarningTemplate(terraform.DefaultDestroyWarningTemplate),
				WarnDestroy:            true,
			},
			body:     "Plan: 1 to add, 1 to destroy",
			ok:       true,
			exitCode: 0,
		},
		{
			name: "case 6",
			// valid with no changes
			// TODO(drlau): check that the label was actually added
			config: Config{
				Token: "token",
				Owner: "owner",
				Repo:  "repo",
				PR: PullRequest{
					Revision: "",
					Number:   1,
					Message:  "message",
				},
				Parser:   terraform.NewPlanParser(),
				Template: terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ResultLabels: ResultLabels{
					AddOrUpdateLabel: "add-or-update",
					DestroyLabel:     "destroy",
					NoChangesLabel:   "no-changes",
					PlanErrorLabel:   "error",
				},
			},
			body:     "No changes. Infrastructure is up-to-date.",
			ok:       true,
			exitCode: 0,
		},
		{
			name: "case 7",
			// valid, contains destroy, but not to notify
			config: Config{
				Token: "token",
				Owner: "owner",
				Repo:  "repo",
				PR: PullRequest{
					Revision: "",
					Number:   1,
					Message:  "message",
				},
				Parser:                 terraform.NewPlanParser(),
				Template:               terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				DestroyWarningTemplate: terraform.NewDestroyWarningTemplate(terraform.DefaultDestroyWarningTemplate),
				WarnDestroy:            false,
			},
			body:     "Plan: 1 to add, 1 to destroy",
			ok:       true,
			exitCode: 0,
		},
		{
			name: "case 8",
			// apply case without merge commit
			config: Config{
				Token: "token",
				Owner: "owner",
				Repo:  "repo",
				PR: PullRequest{
					Revision: "revision",
					Number:   0, // For apply, it is always 0
					Message:  "message",
				},
				Parser:   terraform.NewApplyParser(),
				Template: terraform.NewApplyTemplate(terraform.DefaultApplyTemplate),
			},
			body:     "Apply complete!",
			ok:       true,
			exitCode: 0,
		},
		{
			name: "case 9",
			// apply case as merge commit
			// TODO(drlau): validate cfg.PR.Number = 123
			config: Config{
				Token: "token",
				Owner: "owner",
				Repo:  "repo",
				PR: PullRequest{
					Revision: "Merge pull request #123 from mercari/tfnotify",
					Number:   0, // For apply, it is always 0
					Message:  "message",
				},
				Parser:   terraform.NewApplyParser(),
				Template: terraform.NewApplyTemplate(terraform.DefaultApplyTemplate),
			},
			body:     "Apply complete!",
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
			exitCode, err := client.Notify.Notify(context.Background(), notifier.ParamExec{
				Stdout: testCase.body,
			})
			if (err == nil) != testCase.ok {
				t.Errorf("got error %v", err)
			}
			if exitCode != testCase.exitCode {
				t.Errorf("got %d but want %d", exitCode, testCase.exitCode)
			}
		})
	}
}
