package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/mattn/go-colorable"
	"github.com/mercari/tfnotify/config"
	"github.com/mercari/tfnotify/notifier"
	"github.com/mercari/tfnotify/notifier/github"
	"github.com/mercari/tfnotify/terraform"
	"github.com/urfave/cli/v2"
)

type tfnotify struct {
	config                 config.Config
	context                *cli.Context
	parser                 terraform.Parser
	template               terraform.Template
	destroyWarningTemplate terraform.Template
	warnDestroy            bool
}

func getCI(ciname string) (CI, error) {
	var ci CI
	switch ciname {
	case "circleci", "circle-ci":
		return circleci()
	case "codebuild":
		return codebuild()
	case "github-actions":
		return githubActions(), nil
	case "cloud-build", "cloudbuild":
		return cloudbuild()
	case "":
		return ci, errors.New("CI service: required (e.g. circleci)")
	default:
		return ci, fmt.Errorf("CI service %s: not supported yet", ciname)
	}
}

func (t *tfnotify) renderTemplate(tpl string) (string, error) {
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

func (t *tfnotify) renderGitHubLabels() (github.ResultLabels, error) {
	labels := github.ResultLabels{
		AddOrUpdateLabelColor: t.config.Terraform.Plan.WhenAddOrUpdateOnly.Color,
		DestroyLabelColor:     t.config.Terraform.Plan.WhenDestroy.Color,
		NoChangesLabelColor:   t.config.Terraform.Plan.WhenNoChanges.Color,
		PlanErrorLabelColor:   t.config.Terraform.Plan.WhenPlanError.Color,
	}

	addOrUpdateLabel, err := t.renderTemplate(t.config.Terraform.Plan.WhenAddOrUpdateOnly.Label)
	if err != nil {
		return labels, err
	}
	labels.AddOrUpdateLabel = addOrUpdateLabel

	destroyLabel, err := t.renderTemplate(t.config.Terraform.Plan.WhenDestroy.Label)
	if err != nil {
		return labels, err
	}
	labels.DestroyLabel = destroyLabel

	nochangesLabel, err := t.renderTemplate(t.config.Terraform.Plan.WhenNoChanges.Label)
	if err != nil {
		return labels, err
	}
	labels.NoChangesLabel = nochangesLabel

	planErrorLabel, err := t.renderTemplate(t.config.Terraform.Plan.WhenPlanError.Label)
	if err != nil {
		return labels, err
	}
	labels.PlanErrorLabel = planErrorLabel

	return labels, nil
}

func (t *tfnotify) getNotifier(ctx context.Context, ci CI, selectedNotifier string) (notifier.Notifier, error) {
	switch selectedNotifier {
	case "github":
		labels, err := t.renderGitHubLabels()
		if err != nil {
			return nil, err
		}
		client, err := github.NewClient(ctx, github.Config{
			Token:   t.config.Notifier.Github.Token,
			BaseURL: t.config.Notifier.Github.BaseURL,
			Owner:   t.config.Notifier.Github.Repository.Owner,
			Repo:    t.config.Notifier.Github.Repository.Name,
			PR: github.PullRequest{
				Revision:              ci.PR.Revision,
				Number:                ci.PR.Number,
				Title:                 t.context.String("title"),
				Message:               t.context.String("message"),
				DestroyWarningTitle:   t.context.String("destroy-warning-title"),
				DestroyWarningMessage: t.context.String("destroy-warning-message"),
			},
			CI:                     ci.URL,
			Parser:                 t.parser,
			UseRawOutput:           t.config.Terraform.UseRawOutput,
			Template:               t.template,
			DestroyWarningTemplate: t.destroyWarningTemplate,
			WarnDestroy:            t.warnDestroy,
			ResultLabels:           labels,
			Vars:                   t.config.Vars,
		})
		if err != nil {
			return nil, err
		}
		return client.Notify, nil
	case "":
		return nil, errors.New("notifier is missing")
	default:
		return nil, fmt.Errorf("%s: not supported notifier yet", selectedNotifier)
	}
}

// Run sends the notification with notifier
func (t *tfnotify) Run(ctx context.Context) error {
	ciname := t.config.CI
	if t.context.String("ci") != "" {
		ciname = t.context.String("ci")
	}
	ciname = strings.ToLower(ciname)
	ci, err := getCI(ciname)
	if err != nil {
		return err
	}

	selectedNotifier := t.config.GetNotifierType()
	if t.context.String("notifier") != "" {
		selectedNotifier = t.context.String("notifier")
	}

	ntf, err := t.getNotifier(ctx, ci, selectedNotifier)
	if err != nil {
		return err
	}

	if ntf == nil {
		return errors.New("no notifier specified at all")
	}

	args := t.context.Args()
	cmd := exec.CommandContext(ctx, args.First(), args.Tail()...) //nolint:gosec
	stdout := &bytes.Buffer{}
	uncolorize := colorable.NewNonColorable(stdout)
	stderr := &bytes.Buffer{}
	combinedOutput := &bytes.Buffer{}
	cmd.Stdout = io.MultiWriter(os.Stdout, uncolorize, combinedOutput)
	cmd.Stderr = io.MultiWriter(os.Stderr, stderr, combinedOutput)
	_ = cmd.Run()

	return NewExitError(ntf.Notify(ctx, notifier.ParamExec{
		Stdout:         stdout.String(),
		Stderr:         stderr.String(),
		CombinedOutput: combinedOutput.String(),
		Cmd:            cmd,
		Args:           args,
		ExitCode:       cmd.ProcessState.ExitCode(),
	}))
}

func main() {
	app := cli.NewApp()
	app.Name = "tfcmt"
	app.Usage = "Notify the execution result of terraform command"
	app.Version = version
	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "ci", Usage: "name of CI to run tfnotify"},
		&cli.StringFlag{Name: "config", Usage: "config path"},
		&cli.StringFlag{Name: "notifier", Usage: "notification destination"},
		&cli.StringSliceFlag{Name: "var", Usage: "template variables. The format of value is '<name>:<value>'"},
	}
	app.Commands = []*cli.Command{
		{
			Name:   "plan",
			Usage:  "Parse stdin as a plan result",
			Action: cmdPlan,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "title, t",
					Usage: "Specify the title to use for notification",
				},
				&cli.StringFlag{
					Name:  "message, m",
					Usage: "Specify the message to use for notification",
				},
				&cli.StringFlag{
					Name:  "destroy-warning-title",
					Usage: "Specify the title to use for destroy warning notification",
				},
				&cli.StringFlag{
					Name:  "destroy-warning-message",
					Usage: "Specify the message to use for destroy warning notification",
				},
			},
		},
		{
			Name:   "apply",
			Usage:  "Parse stdin as a apply result",
			Action: cmdApply,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "title, t",
					Usage: "Specify the title to use for notification",
				},
				&cli.StringFlag{
					Name:  "message, m",
					Usage: "Specify the message to use for notification",
				},
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	go handleSignal(cancel)

	os.Exit(HandleExit(app.RunContext(ctx, os.Args)))
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

func newConfig(ctx *cli.Context) (config.Config, error) {
	cfg := config.Config{}
	confPath, err := cfg.Find(ctx.String("config"))
	if err != nil {
		return cfg, err
	}
	if err := cfg.LoadFile(confPath); err != nil {
		return cfg, err
	}
	vars := ctx.StringSlice("var")
	vm := make(map[string]string, len(vars))
	if err := parseVarOpts(vars, vm); err != nil {
		return cfg, err
	}
	cfg.Vars = vm
	if err := cfg.Validation(); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func cmdPlan(ctx *cli.Context) error {
	cfg, err := newConfig(ctx)
	if err != nil {
		return err
	}

	// If when_destroy is not defined in configuration, tfnotify should not notify it
	warnDestroy := cfg.Terraform.Plan.WhenDestroy.Template != ""

	t := &tfnotify{
		config:                 cfg,
		context:                ctx,
		parser:                 terraform.NewPlanParser(),
		template:               terraform.NewPlanTemplate(cfg.Terraform.Plan.Template),
		destroyWarningTemplate: terraform.NewDestroyWarningTemplate(cfg.Terraform.Plan.WhenDestroy.Template),
		warnDestroy:            warnDestroy,
	}
	return t.Run(ctx.Context)
}

func cmdApply(ctx *cli.Context) error {
	cfg, err := newConfig(ctx)
	if err != nil {
		return err
	}
	t := &tfnotify{
		config:   cfg,
		context:  ctx,
		parser:   terraform.NewApplyParser(),
		template: terraform.NewApplyTemplate(cfg.Terraform.Apply.Template),
	}
	return t.Run(ctx.Context)
}
