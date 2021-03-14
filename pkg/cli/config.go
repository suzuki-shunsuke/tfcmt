package cli

import (
	"github.com/suzuki-shunsuke/go-ci-env/cienv"
	"github.com/suzuki-shunsuke/tfcmt/pkg/config"
	"github.com/urfave/cli/v2"
)

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
