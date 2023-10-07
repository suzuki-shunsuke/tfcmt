package cli

import (
	"errors"
	"strings"

	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/config"
	"github.com/urfave/cli/v2"
)

func parseVars(vars []string, envs []string, varsM map[string]string) error {
	parseVarEnvs(envs, varsM)
	return parseVarOpts(vars, varsM)
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

func parseVarEnvs(envs []string, m map[string]string) {
	for _, kv := range envs {
		k, v, _ := strings.Cut(kv, "=")
		if a := strings.TrimPrefix(k, "TFCMT_VAR_"); k != a {
			m[a] = v
		}
	}
}

func parseOpts(ctx *cli.Context, cfg *config.Config, envs []string) error {
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

	if ctx.IsSet("patch") {
		cfg.PlanPatch = ctx.Bool("patch")
	}

	if buildURL := ctx.String("build-url"); buildURL != "" {
		cfg.CI.Link = buildURL
	}

	if output := ctx.String("output"); output != "" {
		cfg.Output = output
	}

	if ctx.IsSet("skip-no-changes") {
		cfg.Terraform.Plan.WhenNoChanges.DisableComment = ctx.Bool("skip-no-changes")
	}

	vars := ctx.StringSlice("var")
	vm := make(map[string]string, len(vars))
	if err := parseVars(vars, envs, vm); err != nil {
		return err
	}
	cfg.Vars = vm

	return nil
}
