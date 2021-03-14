package cli

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/suzuki-shunsuke/tfcmt/pkg/platform"
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

func parseOpts(ctx *cli.Context, ci *platform.CI) error {
	if sha := ctx.String("sha"); sha != "" {
		ci.PR.Revision = sha
	}
	if pr := ctx.Int("pr"); pr != 0 {
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
	if buildURL := ctx.String("build-url"); buildURL != "" {
		ci.URL = buildURL
	}
	return nil
}
