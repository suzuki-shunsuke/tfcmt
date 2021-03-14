package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/go-ci-env/cienv"
	"github.com/suzuki-shunsuke/tfcmt/pkg/apperr"
	"github.com/suzuki-shunsuke/tfcmt/pkg/config"
	"github.com/suzuki-shunsuke/tfcmt/pkg/notifier"
	"github.com/suzuki-shunsuke/tfcmt/pkg/notifier/github"
	"github.com/suzuki-shunsuke/tfcmt/pkg/platform"
	"github.com/suzuki-shunsuke/tfcmt/pkg/terraform"
	"github.com/urfave/cli/v2"
)

type tfcmt struct {
	config                 config.Config
	context                *cli.Context
	parser                 terraform.Parser
	template               *terraform.Template
	destroyWarningTemplate *terraform.Template
	parseErrorTemplate     *terraform.Template
}

func (t *tfcmt) renderTemplate(tpl string) (string, error) {
	tmpl, err := template.New("_").Funcs(sprig.TxtFuncMap()).Parse(tpl)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, map[string]interface{}{
		"Vars": t.config.Vars,
	}); err != nil {
		return "", fmt.Errorf("render a label template: %w", err)
	}
	return buf.String(), nil
}

func (t *tfcmt) renderGitHubLabels() (github.ResultLabels, error) { //nolint:cyclop
	labels := github.ResultLabels{
		AddOrUpdateLabelColor: t.config.Terraform.Plan.WhenAddOrUpdateOnly.Color,
		DestroyLabelColor:     t.config.Terraform.Plan.WhenDestroy.Color,
		NoChangesLabelColor:   t.config.Terraform.Plan.WhenNoChanges.Color,
		PlanErrorLabelColor:   t.config.Terraform.Plan.WhenPlanError.Color,
	}

	target, ok := t.config.Vars["target"]
	if !ok {
		target = ""
	}

	if labels.AddOrUpdateLabelColor == "" {
		labels.AddOrUpdateLabelColor = "1d76db" // blue
	}
	if labels.DestroyLabelColor == "" {
		labels.DestroyLabelColor = "d93f0b" // red
	}
	if labels.NoChangesLabelColor == "" {
		labels.NoChangesLabelColor = "0e8a16" // green
	}

	if t.config.Terraform.Plan.WhenAddOrUpdateOnly.Label == "" {
		if target == "" {
			labels.AddOrUpdateLabel = "add-or-update"
		} else {
			labels.AddOrUpdateLabel = target + "/add-or-update"
		}
	} else {
		addOrUpdateLabel, err := t.renderTemplate(t.config.Terraform.Plan.WhenAddOrUpdateOnly.Label)
		if err != nil {
			return labels, err
		}
		labels.AddOrUpdateLabel = addOrUpdateLabel
	}

	if t.config.Terraform.Plan.WhenDestroy.Label == "" {
		if target == "" {
			labels.DestroyLabel = "destroy"
		} else {
			labels.DestroyLabel = target + "/destroy"
		}
	} else {
		destroyLabel, err := t.renderTemplate(t.config.Terraform.Plan.WhenDestroy.Label)
		if err != nil {
			return labels, err
		}
		labels.DestroyLabel = destroyLabel
	}

	if t.config.Terraform.Plan.WhenNoChanges.Label == "" {
		if target == "" {
			labels.NoChangesLabel = "no-changes"
		} else {
			labels.NoChangesLabel = target + "/no-changes"
		}
	} else {
		nochangesLabel, err := t.renderTemplate(t.config.Terraform.Plan.WhenNoChanges.Label)
		if err != nil {
			return labels, err
		}
		labels.NoChangesLabel = nochangesLabel
	}

	planErrorLabel, err := t.renderTemplate(t.config.Terraform.Plan.WhenPlanError.Label)
	if err != nil {
		return labels, err
	}
	labels.PlanErrorLabel = planErrorLabel

	return labels, nil
}

func (t *tfcmt) getNotifier(ctx context.Context, ci platform.CI) (notifier.Notifier, error) {
	labels := github.ResultLabels{}
	if !t.config.Terraform.Plan.DisableLabel {
		a, err := t.renderGitHubLabels()
		if err != nil {
			return nil, err
		}
		labels = a
	}
	client, err := github.NewClient(ctx, github.Config{
		Token:   t.config.Notifier.Github.Token,
		BaseURL: t.config.Notifier.Github.BaseURL,
		Owner:   t.config.Notifier.Github.Repository.Owner,
		Repo:    t.config.Notifier.Github.Repository.Name,
		PR: github.PullRequest{
			Revision: ci.PR.Revision,
			Number:   ci.PR.Number,
		},
		CI:                     ci.URL,
		Parser:                 t.parser,
		UseRawOutput:           t.config.Terraform.UseRawOutput,
		Template:               t.template,
		DestroyWarningTemplate: t.destroyWarningTemplate,
		ParseErrorTemplate:     t.parseErrorTemplate,
		ResultLabels:           labels,
		Vars:                   t.config.Vars,
		Templates:              t.config.Templates,
	})
	if err != nil {
		return nil, err
	}
	return client.Notify, nil
}

// Run sends the notification with notifier
func (t *tfcmt) Run(ctx context.Context) error { //nolint:cyclop
	ciname := t.config.CI
	if t.context.String("ci") != "" {
		ciname = t.context.String("ci")
	}
	ciname = strings.ToLower(ciname)
	ci, err := platform.Get(ciname)
	if err != nil {
		return err
	}
	if sha := t.context.String("sha"); sha != "" {
		ci.PR.Revision = sha
	}
	if pr := t.context.Int("pr"); pr != 0 {
		ci.PR.Number = pr
	}
	if ci.PR.Number == 0 {
		// support suzuki-shunsuke/ci-info
		if prS := os.Getenv("CI_INFO_PR_NUMBER"); prS != "" {
			a, err := strconv.Atoi(prS)
			if err != nil {
				return fmt.Errorf("parse CI_INFO_PR_NUMBER %s: %w", prS, err)
			}
			ci.PR.Number = a
		}
	}
	if buildURL := t.context.String("build-url"); buildURL != "" {
		ci.URL = buildURL
	}

	if ci.PR.Revision == "" && ci.PR.Number == 0 {
		return errors.New("pull request number or SHA (revision) is needed")
	}

	ntf, err := t.getNotifier(ctx, ci)
	if err != nil {
		return err
	}

	if ntf == nil {
		return errors.New("no notifier specified at all")
	}

	args := t.context.Args()
	cmd := exec.CommandContext(ctx, args.First(), args.Tail()...) //nolint:gosec
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	combinedOutput := &bytes.Buffer{}
	uncolorizedStdout := colorable.NewNonColorable(stdout)
	uncolorizedStderr := colorable.NewNonColorable(stderr)
	uncolorizedCombinedOutput := colorable.NewNonColorable(combinedOutput)
	cmd.Stdout = io.MultiWriter(os.Stdout, uncolorizedStdout, uncolorizedCombinedOutput)
	cmd.Stderr = io.MultiWriter(os.Stderr, uncolorizedStderr, uncolorizedCombinedOutput)
	_ = cmd.Run()

	return apperr.NewExitError(ntf.Notify(ctx, notifier.ParamExec{
		Stdout:         stdout.String(),
		Stderr:         stderr.String(),
		CombinedOutput: combinedOutput.String(),
		Cmd:            cmd,
		Args:           args,
		CIName:         ciname,
		ExitCode:       cmd.ProcessState.ExitCode(),
	}))
}

func main() {
	app := cli.NewApp()
	app.Name = "tfcmt"
	app.Usage = "Notify the execution result of terraform command"
	app.Version = version
	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "ci", Usage: "name of CI to run tfcmt"},
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

	ctx, cancel := context.WithCancel(context.Background())
	go handleSignal(cancel)

	os.Exit(apperr.HandleExit(app.RunContext(ctx, os.Args)))
}

func setLogLevel(logLevel string) {
	if logLevel == "" {
		return
	}
	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"log_level": logLevel,
		}).WithError(err).Error("the log level is invalid")
	}
	logrus.SetLevel(lvl)
}

func parseVarOpts(vars []string, varsM map[string]string) error {
	for _, v := range vars {
		a := strings.Index(v, ":")
		if a == -1 {
			return errors.New("the value of var option is invalid. the format should be '<name>:<value>': " + v)
		}
		varsM[v[:a]] = v[a+1:]
	}
	return nil
}

func newConfig(ctx *cli.Context) (config.Config, error) { //nolint:cyclop
	cfg := config.Config{}
	confPath, err := cfg.Find(ctx.String("config"))
	if err != nil {
		return cfg, err
	}
	if confPath != "" {
		if err := cfg.LoadFile(confPath); err != nil {
			return cfg, err
		}
	}
	vars := ctx.StringSlice("var")
	vm := make(map[string]string, len(vars))
	if err := parseVarOpts(vars, vm); err != nil {
		return cfg, err
	}
	cfg.Vars = vm

	if owner := ctx.String("owner"); owner != "" {
		cfg.Notifier.Github.Repository.Owner = owner
	}
	if repo := ctx.String("repo"); repo != "" {
		cfg.Notifier.Github.Repository.Name = repo
	}

	var platform cienv.Platform
	if cfg.CI == "" {
		platform = cienv.Get()
		if platform != nil {
			cfg.CI = platform.CI()
		}
	} else {
		platform = cienv.GetByName(cfg.CI)
	}
	if platform != nil {
		if cfg.Notifier.Github.Repository.Owner == "" {
			cfg.Notifier.Github.Repository.Owner = platform.RepoOwner()
		}
		if cfg.Notifier.Github.Repository.Name == "" {
			cfg.Notifier.Github.Repository.Name = platform.RepoName()
		}
	}

	if err := cfg.Validation(); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func cmdPlan(ctx *cli.Context) error {
	logLevel := ctx.String("log-level")
	setLogLevel(logLevel)

	cfg, err := newConfig(ctx)
	if err != nil {
		return err
	}
	if logLevel == "" {
		logLevel = cfg.Log.Level
		setLogLevel(logLevel)
	}

	t := &tfcmt{
		config:                 cfg,
		context:                ctx,
		parser:                 terraform.NewPlanParser(),
		template:               terraform.NewPlanTemplate(cfg.Terraform.Plan.Template),
		destroyWarningTemplate: terraform.NewDestroyWarningTemplate(cfg.Terraform.Plan.WhenDestroy.Template),
		parseErrorTemplate:     terraform.NewPlanParseErrorTemplate(cfg.Terraform.Plan.WhenParseError.Template),
	}
	return t.Run(ctx.Context)
}

func cmdApply(ctx *cli.Context) error {
	cfg, err := newConfig(ctx)
	if err != nil {
		return err
	}
	t := &tfcmt{
		config:             cfg,
		context:            ctx,
		parser:             terraform.NewApplyParser(),
		template:           terraform.NewApplyTemplate(cfg.Terraform.Apply.Template),
		parseErrorTemplate: terraform.NewApplyParseErrorTemplate(cfg.Terraform.Apply.WhenParseError.Template),
	}
	return t.Run(ctx.Context)
}
