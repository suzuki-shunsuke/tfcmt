package main

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/go-ci-env/cienv"
	"github.com/suzuki-shunsuke/tfcmt/pkg/apperr"
	"github.com/suzuki-shunsuke/tfcmt/pkg/config"
	"github.com/suzuki-shunsuke/tfcmt/pkg/constant"
	"github.com/suzuki-shunsuke/tfcmt/pkg/controller"
	"github.com/suzuki-shunsuke/tfcmt/pkg/signal"
	"github.com/suzuki-shunsuke/tfcmt/pkg/terraform"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "tfcmt"
	app.Usage = "Notify the execution result of terraform command"
	app.Version = constant.Version
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
	go signal.Handle(cancel)

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

	t := &controller.Controller{
		Config:                 cfg,
		Context:                ctx,
		Parser:                 terraform.NewPlanParser(),
		Template:               terraform.NewPlanTemplate(cfg.Terraform.Plan.Template),
		DestroyWarningTemplate: terraform.NewDestroyWarningTemplate(cfg.Terraform.Plan.WhenDestroy.Template),
		ParseErrorTemplate:     terraform.NewPlanParseErrorTemplate(cfg.Terraform.Plan.WhenParseError.Template),
	}
	return t.Run(ctx.Context)
}

func cmdApply(ctx *cli.Context) error {
	cfg, err := newConfig(ctx)
	if err != nil {
		return err
	}
	t := &controller.Controller{
		Config:             cfg,
		Context:            ctx,
		Parser:             terraform.NewApplyParser(),
		Template:           terraform.NewApplyTemplate(cfg.Terraform.Apply.Template),
		ParseErrorTemplate: terraform.NewApplyParseErrorTemplate(cfg.Terraform.Apply.WhenParseError.Template),
	}
	return t.Run(ctx.Context)
}
