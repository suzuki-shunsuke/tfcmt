package cli

import (
	"context"

	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
	"github.com/urfave/cli/v3"
)

type GlobalArgs struct {
	Owner    string
	Repo     string
	SHA      string
	BuildURL string
	LogLevel string
	PR       int
	Config   string
	Var      []string
	Output   string
}

type PlanArgs struct {
	*GlobalArgs

	Patch              bool
	PatchCount         int
	SkipNoChanges      bool
	SkipNoChangesCount int
	IgnoreWarning      bool
	IgnoreWarningCount int
	DisableLabel       bool
	DisableLabelCount  int
	Label              string
	Command            string
	CommandArgs        []string
}

type ApplyArgs struct {
	*GlobalArgs

	Label       string
	Command     string
	CommandArgs []string
}

func Run(ctx context.Context, logger *slogutil.Logger, env *urfave.Env) error {
	globalArgs := &GlobalArgs{}
	planArgs := &PlanArgs{GlobalArgs: globalArgs}
	applyArgs := &ApplyArgs{GlobalArgs: globalArgs}

	return urfave.Command(env, &cli.Command{ //nolint:wrapcheck
		Name:           "tfcmt",
		Usage:          "Notify the execution result of terraform command",
		ExitErrHandler: func(context.Context, *cli.Command, error) {},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "owner",
				Usage:       "GitHub Repository owner name",
				Sources:     cli.EnvVars("TFCMT_REPO_OWNER"),
				Destination: &globalArgs.Owner,
			},
			&cli.StringFlag{
				Name:        "repo",
				Usage:       "GitHub Repository name",
				Sources:     cli.EnvVars("TFCMT_REPO_NAME"),
				Destination: &globalArgs.Repo,
			},
			&cli.StringFlag{
				Name:        "sha",
				Usage:       "commit SHA (revision)",
				Sources:     cli.EnvVars("TFCMT_SHA"),
				Destination: &globalArgs.SHA,
			},
			&cli.StringFlag{
				Name:        "build-url",
				Usage:       "build url",
				Destination: &globalArgs.BuildURL,
			},
			&cli.StringFlag{
				Name:        "log-level",
				Usage:       "log level",
				Destination: &globalArgs.LogLevel,
			},
			&cli.IntFlag{
				Name:        "pr",
				Usage:       "pull request number",
				Sources:     cli.EnvVars("TFCMT_PR_NUMBER"),
				Destination: &globalArgs.PR,
			},
			&cli.StringFlag{
				Name:        "config",
				Usage:       "config path",
				Sources:     cli.EnvVars("TFCMT_CONFIG"),
				Destination: &globalArgs.Config,
			},
			&cli.StringSliceFlag{
				Name:        "var",
				Usage:       "template variables. The format of value is '<name>:<value>'. You can refer to the variable in the comment and label template using {{.Vars.<variable name>}}.",
				Destination: &globalArgs.Var,
			},
			&cli.StringFlag{
				Name:        "output",
				Usage:       "specify file to output result instead of posting a comment",
				Destination: &globalArgs.Output,
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "plan",
				ArgsUsage: " <command> <args>...",
				Usage:     "Run terraform plan and post a comment to GitHub commit, pull request, or issue",
				Description: `Run terraform plan and post a comment to GitHub commit, pull request, or issue.

$ tfcmt [<global options>] plan [-patch] [-skip-no-changes] -- terraform plan [<terraform plan options>]`,
				Action: func(ctx context.Context, _ *cli.Command) error {
					return actionPlan(ctx, logger, planArgs)
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "patch",
						Usage:       "update an existing comment instead of creating a new comment. If there is no existing comment, a new comment is created.",
						Sources:     cli.EnvVars("TFCMT_PLAN_PATCH"),
						Destination: &planArgs.Patch,
						Config: cli.BoolConfig{
							Count: &planArgs.PatchCount,
						},
					},
					&cli.BoolFlag{
						Name:        "skip-no-changes",
						Usage:       "If there is no change tfcmt updates a label but doesn't post a comment",
						Sources:     cli.EnvVars("TFCMT_SKIP_NO_CHANGES"),
						Destination: &planArgs.SkipNoChanges,
						Config: cli.BoolConfig{
							Count: &planArgs.SkipNoChangesCount,
						},
					},
					&cli.BoolFlag{
						Name:        "ignore-warning",
						Usage:       "If skip-no-changes is enabled, comment is posted even if there is a warning. If skip-no-changes is disabled, warning is removed from the comment.",
						Sources:     cli.EnvVars("TFCMT_IGNORE_WARNING"),
						Destination: &planArgs.IgnoreWarning,
						Config: cli.BoolConfig{
							Count: &planArgs.IgnoreWarningCount,
						},
					},
					&cli.BoolFlag{
						Name:        "disable-label",
						Usage:       "Disable to add or update a label",
						Sources:     cli.EnvVars("TFCMT_DISABLE_LABEL"),
						Destination: &planArgs.DisableLabel,
						Config: cli.BoolConfig{
							Count: &planArgs.DisableLabelCount,
						},
					},
					&cli.StringFlag{
						Name:        "label",
						Usage:       "Override the default calculated label with a custom label. This allows assigning multiple labels across runs.",
						Sources:     cli.EnvVars("TFCMT_LABEL"),
						Destination: &planArgs.Label,
					},
				},
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:        "command",
						Destination: &planArgs.Command,
					},
					&cli.StringArgs{
						Name:        "args",
						Destination: &planArgs.CommandArgs,
						Max:         -1,
					},
				},
			},
			{
				Name:      "apply",
				ArgsUsage: " <command> <args>...",
				Usage:     "Run terraform apply and post a comment to GitHub commit, pull request, or issue",
				Description: `Run terraform apply and post a comment to GitHub commit, pull request, or issue.

$ tfcmt [<global options>] apply -- terraform apply [<terraform apply options>]`,
				Action: func(ctx context.Context, _ *cli.Command) error {
					return actionApply(ctx, logger, applyArgs)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "label",
						Usage:       "Override the default calculated label with a custom label. This allows assigning multiple labels across runs.",
						Sources:     cli.EnvVars("TFCMT_LABEL"),
						Destination: &applyArgs.Label,
					},
				},
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:        "command",
						Destination: &applyArgs.Command,
					},
					&cli.StringArgs{
						Name:        "args",
						Destination: &applyArgs.CommandArgs,
						Max:         -1,
					},
				},
			},
		},
	}).Run(ctx, env.Args)
}
