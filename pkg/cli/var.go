package cli

import (
	"errors"
	"os"
	"strings"

	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/config"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/mask"
)

func parseVars(vars []string, envs []string, varsM map[string]string) error {
	parseVarEnvs(envs, varsM)
	return parseVarOpts(vars, varsM)
}

func parseVarOpts(vars []string, varsM map[string]string) error {
	for _, v := range vars {
		name, value, ok := strings.Cut(v, ":")
		if !ok {
			return errors.New("the value of var option is invalid. the format should be '<name>:<value>': " + v)
		}
		varsM[name] = value
	}
	return nil
}

func parseVarEnvs(envs []string, m map[string]string) {
	for _, kv := range envs {
		k, v, _ := strings.Cut(kv, "=")
		if a, ok := strings.CutPrefix(k, "TFCMT_VAR_"); ok {
			m[a] = v
		}
	}
}

func parseOptsPlan(args *PlanArgs, cfg *config.Config, envs []string) error { //nolint:cyclop
	if args.Owner != "" {
		cfg.CI.Owner = args.Owner
	}

	if args.Repo != "" {
		cfg.CI.Repo = args.Repo
	}

	if args.SHA != "" {
		cfg.CI.SHA = args.SHA
	}

	if args.PR != 0 {
		cfg.CI.PRNumber = args.PR
	}

	if args.PatchCount > 0 {
		cfg.PlanPatch = args.Patch
	}

	if args.BuildURL != "" {
		cfg.CI.Link = args.BuildURL
	}

	if args.Output != "" {
		cfg.Output = args.Output
	}

	if args.SkipNoChangesCount > 0 {
		cfg.Terraform.Plan.WhenNoChanges.DisableComment = args.SkipNoChanges
	}

	if args.IgnoreWarningCount > 0 {
		cfg.Terraform.Plan.IgnoreWarning = args.IgnoreWarning
	}

	vm := make(map[string]string, len(args.Var))
	if err := parseVars(args.Var, envs, vm); err != nil {
		return err
	}
	cfg.Vars = vm

	// Mask https://github.com/suzuki-shunsuke/tfcmt/discussions/1083
	masks, err := mask.ParseMasksFromEnv()
	if err != nil {
		return err
	}
	cfg.Masks = masks

	if args.DisableLabelCount > 0 {
		cfg.Terraform.Plan.DisableLabel = args.DisableLabel
	}

	if cfg.GHEBaseURL == "" {
		cfg.GHEBaseURL = os.Getenv("GITHUB_API_URL")
	}
	if cfg.GHEGraphQLEndpoint == "" {
		cfg.GHEGraphQLEndpoint = os.Getenv("GITHUB_GRAPHQL_URL")
	}

	return nil
}

func parseOptsApply(args *ApplyArgs, cfg *config.Config, envs []string) error { //nolint:cyclop
	if args.Owner != "" {
		cfg.CI.Owner = args.Owner
	}

	if args.Repo != "" {
		cfg.CI.Repo = args.Repo
	}

	if args.SHA != "" {
		cfg.CI.SHA = args.SHA
	}

	if args.PR != 0 {
		cfg.CI.PRNumber = args.PR
	}

	if args.BuildURL != "" {
		cfg.CI.Link = args.BuildURL
	}

	if args.Output != "" {
		cfg.Output = args.Output
	}

	vm := make(map[string]string, len(args.Var))
	if err := parseVars(args.Var, envs, vm); err != nil {
		return err
	}
	cfg.Vars = vm

	// Mask https://github.com/suzuki-shunsuke/tfcmt/discussions/1083
	masks, err := mask.ParseMasksFromEnv()
	if err != nil {
		return err
	}
	cfg.Masks = masks

	if cfg.GHEBaseURL == "" {
		cfg.GHEBaseURL = os.Getenv("GITHUB_API_URL")
	}
	if cfg.GHEGraphQLEndpoint == "" {
		cfg.GHEGraphQLEndpoint = os.Getenv("GITHUB_GRAPHQL_URL")
	}

	return nil
}
