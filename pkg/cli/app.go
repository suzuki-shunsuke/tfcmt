package cli

import (
	"github.com/suzuki-shunsuke/urfave-cli-help-all/helpall"
	"github.com/urfave/cli/v3"
)

type LDFlags struct {
	Version string
	Commit  string
	Date    string
}

func (f *LDFlags) AppVersion() string {
	return f.Version + " (" + f.Commit + ")"
}

func New(flags *LDFlags) *cli.App {
	app := cli.NewApp()
	app.Name = "tfcmt"
	app.Usage = "Notify the execution result of terraform command"
	app.Version = flags.AppVersion()
	app.ExitErrHandler = func(*cli.Context, error) {}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "owner",
			Usage:   "GitHub Repository owner name",
			EnvVars: []string{"TFCMT_REPO_OWNER"},
		},
		&cli.StringFlag{
			Name:    "repo",
			Usage:   "GitHub Repository name",
			EnvVars: []string{"TFCMT_REPO_NAME"},
		},
		&cli.StringFlag{
			Name:    "sha",
			Usage:   "commit SHA (revision)",
			EnvVars: []string{"TFCMT_SHA"},
		},
		&cli.StringFlag{
			Name:  "build-url",
			Usage: "build url",
		},
		&cli.StringFlag{
			Name:  "log-level",
			Usage: "log level",
		},
		&cli.IntFlag{
			Name:    "pr",
			Usage:   "pull request number",
			EnvVars: []string{"TFCMT_PR_NUMBER"},
		},
		&cli.StringFlag{
			Name:    "config",
			Usage:   "config path",
			EnvVars: []string{"TFCMT_CONFIG"},
		},
		&cli.StringSliceFlag{
			Name:  "var",
			Usage: "template variables. The format of value is '<name>:<value>'. You can refer to the variable in the comment and label template using {{.Vars.<variable name>}}.",
		},
		&cli.StringFlag{
			Name:  "output",
			Usage: "specify file to output result instead of posting a comment",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:      "plan",
			ArgsUsage: " <command> <args>...",
			Usage:     "Run terraform plan and post a comment to GitHub commit, pull request, or issue",
			Description: `Run terraform plan and post a comment to GitHub commit, pull request, or issue.

$ tfcmt [<global options>] plan [-patch] [-skip-no-changes] -- terraform plan [<terraform plan options>]`,
			Action: cmdPlan,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "patch",
					Usage:   "update an existing comment instead of creating a new comment. If there is no existing comment, a new comment is created.",
					EnvVars: []string{"TFCMT_PLAN_PATCH"},
				},
				&cli.BoolFlag{
					Name:    "skip-no-changes",
					Usage:   "If there is no change tfcmt updates a label but doesn't post a comment",
					EnvVars: []string{"TFCMT_SKIP_NO_CHANGES"},
				},
				&cli.BoolFlag{
					Name:    "ignore-warning",
					Usage:   "If skip-no-changes is enabled, comment is posted even if there is a warning. If skip-no-changes is disabled, warning is removed from the comment.",
					EnvVars: []string{"TFCMT_IGNORE_WARNING"},
				},
				&cli.BoolFlag{
					Name:    "disable-label",
					Usage:   "Disable to add or update a label",
					EnvVars: []string{"TFCMT_DISABLE_LABEL"},
				},
			},
		},
		{
			Name:      "apply",
			ArgsUsage: " <command> <args>...",
			Usage:     "Run terraform apply and post a comment to GitHub commit, pull request, or issue",
			Description: `Run terraform apply and post a comment to GitHub commit, pull request, or issue.

$ tfcmt [<global options>] apply -- terraform apply [<terraform apply options>]`,
			Action: cmdApply,
		},
		{
			Name:  "version",
			Usage: "Show version",
			Action: func(ctx *cli.Context) error {
				cli.ShowVersion(ctx)
				return nil
			},
		},
		helpall.New(nil),
	}
	return app
}
