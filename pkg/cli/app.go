package cli

import (
	"github.com/urfave/cli/v2"
)

type LDFlags struct {
	Version string
	Commit  string
	Date    string
}

func (flags *LDFlags) AppVersion() string {
	return flags.Version + " (" + flags.Commit + ")"
}

func New(flags *LDFlags) *cli.App {
	app := cli.NewApp()
	app.Name = "tfcmt"
	app.Usage = "Notify the execution result of terraform command"
	app.Version = flags.AppVersion()
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
			Usage: "template variables. The format of value is '<name>:<value>'",
		},
		&cli.StringFlag{
			Name:  "output",
			Usage: "specify file to output result instead of post comment",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:   "plan",
			Usage:  "Run terraform plan and post a comment to GitHub commit or pull request",
			Action: cmdPlan,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "patch",
					Usage: "update an existing comment instead of creating a new comment. If there is no existing comment, a new comment is created.",
				},
				&cli.BoolFlag{
					Name:  "skip-no-changes",
					Usage: "If there is no change tfcmt updates a label but doesn't post a comment",
				},
			},
		},
		{
			Name:   "apply",
			Usage:  "Run terraform apply and post a comment to GitHub commit or pull request",
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
	}
	return app
}
