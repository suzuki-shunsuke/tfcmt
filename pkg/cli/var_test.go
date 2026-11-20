package cli

import (
	"testing"

	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/config"
)

func TestParseOptsPlan(t *testing.T) {
	t.Parallel()

	t.Run("sets ResultFile and GitHubOutput from args", func(t *testing.T) {
		t.Parallel()
		args := &PlanArgs{
			GlobalArgs:        &GlobalArgs{},
			ResultFile:        "result.json",
			GitHubOutput:      true,
			GitHubOutputCount: 1,
		}
		var cfg config.Config

		if err := parseOptsPlan(args, &cfg, nil); err != nil {
			t.Fatalf("parseOptsPlan failed: %v", err)
		}

		if cfg.Terraform.Plan.ResultFile != "result.json" {
			t.Fatalf("expected ResultFile result.json, got %q", cfg.Terraform.Plan.ResultFile)
		}
		if !cfg.Terraform.Plan.GitHubOutput {
			t.Fatalf("expected GitHubOutput true, got false")
		}
	})

	t.Run("preserves config defaults when args not supplied", func(t *testing.T) {
		t.Parallel()
		args := &PlanArgs{
			GlobalArgs: &GlobalArgs{},
		}
		cfg := config.Config{
			Terraform: config.Terraform{
				Plan: config.Plan{
					ResultFile:   "default.json",
					GitHubOutput: true,
				},
			},
		}

		if err := parseOptsPlan(args, &cfg, nil); err != nil {
			t.Fatalf("parseOptsPlan failed: %v", err)
		}

		if cfg.Terraform.Plan.ResultFile != "default.json" {
			t.Fatalf("expected ResultFile to remain default.json, got %q", cfg.Terraform.Plan.ResultFile)
		}
		if !cfg.Terraform.Plan.GitHubOutput {
			t.Fatalf("expected GitHubOutput to remain true, got false")
		}
	})

	t.Run("parses variables from args and environment", func(t *testing.T) {
		t.Parallel()
		args := &PlanArgs{
			GlobalArgs: &GlobalArgs{
				Var: []string{"target:production", "tier:frontend"},
			},
		}
		var cfg config.Config
		envs := []string{"TFCMT_VAR_cluster=us-east-1", "OTHER_ENV=ignore"}

		if err := parseOptsPlan(args, &cfg, envs); err != nil {
			t.Fatalf("parseOptsPlan failed: %v", err)
		}

		if cfg.Vars["target"] != "production" {
			t.Fatalf("expected target=production, got %q", cfg.Vars["target"])
		}
		if cfg.Vars["tier"] != "frontend" {
			t.Fatalf("expected tier=frontend, got %q", cfg.Vars["tier"])
		}
		if cfg.Vars["cluster"] != "us-east-1" {
			t.Fatalf("expected cluster=us-east-1, got %q", cfg.Vars["cluster"])
		}
	})
}
