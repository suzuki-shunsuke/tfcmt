package cli

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/suzuki-shunsuke/tfcmt/pkg/config"
	"github.com/urfave/cli/v2"
)

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

func parseOpts(ctx *cli.Context, cfg *config.Config) error {
	if owner := ctx.String("owner"); owner != "" {
		cfg.CI.Owner = owner
	}

	if repo := ctx.String("repo"); repo != "" {
		cfg.CI.Repo = repo
	}

	if sha := ctx.String("sha"); sha != "" {
		cfg.CI.SHA = sha
	}

	if pr := ctx.Int("pr"); pr != 0 {
		cfg.CI.PRNumber = pr
	}
	if cfg.CI.PRNumber == 0 {
		// support suzuki-shunsuke/ci-info
		if prS := os.Getenv("CI_INFO_PR_NUMBER"); prS != "" {
			a, err := strconv.Atoi(prS)
			if err != nil {
				return fmt.Errorf("parse CI_INFO_PR_NUMBER %s: %w", prS, err)
			}
			cfg.CI.PRNumber = a
		}
	}

	if buildURL := ctx.String("build-url"); buildURL != "" {
		cfg.CI.Link = buildURL
	}

	vars := ctx.StringSlice("var")
	vm := make(map[string]string, len(vars))
	if err := parseVarOpts(vars, vm); err != nil {
		return err
	}
	cfg.Vars = vm

	return nil
}
