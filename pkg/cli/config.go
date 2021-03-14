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
		cfg.CI.Owner = owner
	}
	if repo := ctx.String("repo"); repo != "" {
		cfg.CI.Repo = repo
	}

	platform := cienv.Get()
	if platform != nil {
		if cfg.CI.Owner == "" {
			cfg.CI.Owner = platform.RepoOwner()
		}
		if cfg.CI.Repo == "" {
			cfg.CI.Repo = platform.RepoName()
		}
	}

	if err := cfg.Validation(); err != nil {
		return cfg, err
	}
	return cfg, nil
}
