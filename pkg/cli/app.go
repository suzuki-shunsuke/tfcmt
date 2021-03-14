package cli

import (
	"github.com/suzuki-shunsuke/tfcmt/pkg/constant"
	"github.com/urfave/cli/v2"
)

func New() *cli.App {
	app := cli.NewApp()
	app.Name = "tfcmt"
	app.Usage = "Notify the execution result of terraform command"
	app.Version = constant.Version
	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "owner", Usage: "GitHub Repository owner name"},
		&cli.StringFlag{Name: "repo", Usage: "GitHub Repository name"},
		&cli.StringFlag{Name: "sha", Usage: "commit SHA (revision)"},
		&cli.StringFlag{Name: "build-url", Usage: "build url"},
		&cli.StringFlag{Name: "log-level", Usage: "log level"},
		&cli.IntFlag{Name: "pr", Usage: "pull request number"},
		&cli.StringFlag{Name: "config", Usage: "config path"},
		&cli.StringSliceFlag{Name: "var", Usage: "template variables. The format of value is '<name>:<value>'"},
	}
	app.Commands = []*cli.Command{
		{
			Name:   "plan",
			Usage:  "Run terraform plan and post a comment to GitHub commit or pull request",
			Action: cmdPlan,
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
