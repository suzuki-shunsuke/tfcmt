package controller_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/config"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/controller"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/terraform"
)

func TestControllerExportPlanResult(t *testing.T) {
	t.Run("no-op when ResultFile and GitHubOutput are disabled", func(t *testing.T) {
		ctrl := &controller.Controller{
			Config: config.Config{},
			Parser: terraform.NewPlanParser(),
		}
		if err := ctrl.ExportPlanResult("No changes.", 0); err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
	})

	t.Run("writes ResultFile with parsed plan attributes", func(t *testing.T) {
		tempDir := t.TempDir()
		resultPath := filepath.Join(tempDir, "plan-result.json")

		ctrl := &controller.Controller{
			Config: config.Config{
				Terraform: config.Terraform{
					Plan: config.Plan{
						ResultFile: resultPath,
					},
				},
				Vars: map[string]string{"target": "production"},
			},
			Parser: terraform.NewPlanParser(),
		}

		planOutput := `Terraform will perform the following actions:
  # aws_s3_bucket.data will be created
  + resource "aws_s3_bucket" "data" {}
Plan: 1 to add, 0 to change, 0 to destroy.`

		if err := ctrl.ExportPlanResult(planOutput, 2); err != nil {
			t.Fatalf("ExportPlanResult failed: %v", err)
		}

		data, err := os.ReadFile(resultPath)
		if err != nil {
			t.Fatalf("failed to read result file: %v", err)
		}

		var parsed terraform.PlanResult
		if err := json.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("failed to unmarshal result JSON: %v", err)
		}

		if parsed.ExitCode != 2 {
			t.Fatalf("expected ExitCode 2, got %d", parsed.ExitCode)
		}
		if parsed.HasNoChanges {
			t.Fatalf("expected HasNoChanges false, got true")
		}
		if !parsed.HasAddOrUpdateOnly {
			t.Fatalf("expected HasAddOrUpdateOnly true, got false")
		}
		if parsed.HasDestroy {
			t.Fatalf("expected HasDestroy false, got true")
		}
		if parsed.HasError {
			t.Fatalf("expected HasError false, got true")
		}
		if len(parsed.CreatedResources) != 1 || parsed.CreatedResources[0] != "aws_s3_bucket.data" {
			t.Fatalf("unexpected CreatedResources: %#v", parsed.CreatedResources)
		}
		if parsed.Vars["target"] != "production" {
			t.Fatalf("expected Vars[target] production, got %q", parsed.Vars["target"])
		}
	})

	t.Run("returns error when GitHubOutput enabled but GITHUB_OUTPUT unset", func(t *testing.T) {
		t.Setenv("GITHUB_OUTPUT", "")

		ctrl := &controller.Controller{
			Config: config.Config{
				Terraform: config.Terraform{
					Plan: config.Plan{
						GitHubOutput: true,
					},
				},
			},
			Parser: terraform.NewPlanParser(),
		}

		err := ctrl.ExportPlanResult("No changes.", 0)
		if err == nil {
			t.Fatalf("expected error when GITHUB_OUTPUT is unset, got nil")
		}
		if !strings.Contains(err.Error(), "GITHUB_OUTPUT") {
			t.Fatalf("expected error mentioning GITHUB_OUTPUT, got: %v", err)
		}
	})

	t.Run("appends to GITHUB_OUTPUT when enabled and set", func(t *testing.T) {
		tempDir := t.TempDir()
		githubOutputPath := filepath.Join(tempDir, "gh_out")
		t.Setenv("GITHUB_OUTPUT", githubOutputPath)

		ctrl := &controller.Controller{
			Config: config.Config{
				Terraform: config.Terraform{
					Plan: config.Plan{
						GitHubOutput: true,
					},
				},
			},
			Parser: terraform.NewPlanParser(),
		}

		planOutput := "No changes. Your infrastructure matches the configuration."
		if err := ctrl.ExportPlanResult(planOutput, 0); err != nil {
			t.Fatalf("ExportPlanResult failed: %v", err)
		}

		data, err := os.ReadFile(githubOutputPath)
		if err != nil {
			t.Fatalf("failed to read GITHUB_OUTPUT file: %v", err)
		}

		outStr := string(data)
		if !strings.Contains(outStr, "has_no_changes=true\n") {
			t.Fatalf("expected has_no_changes=true, got: %s", outStr)
		}
		if !strings.Contains(outStr, "exit_code=0\n") {
			t.Fatalf("expected exit_code=0, got: %s", outStr)
		}
		if !strings.Contains(outStr, "has_destroy=false\n") {
			t.Fatalf("expected has_destroy=false, got: %s", outStr)
		}
	})

	t.Run("writes both ResultFile and GITHUB_OUTPUT when both enabled", func(t *testing.T) {
		tempDir := t.TempDir()
		resultPath := filepath.Join(tempDir, "result.json")
		githubOutputPath := filepath.Join(tempDir, "github_output")
		t.Setenv("GITHUB_OUTPUT", githubOutputPath)

		ctrl := &controller.Controller{
			Config: config.Config{
				Terraform: config.Terraform{
					Plan: config.Plan{
						ResultFile:   resultPath,
						GitHubOutput: true,
					},
				},
			},
			Parser: terraform.NewPlanParser(),
		}

		planOutput := `Terraform will perform the following actions:
  # aws_instance.web will be destroyed
  - resource "aws_instance" "web" {}
Plan: 0 to add, 0 to change, 1 to destroy.`

		if err := ctrl.ExportPlanResult(planOutput, 2); err != nil {
			t.Fatalf("ExportPlanResult failed: %v", err)
		}

		if _, err := os.Stat(resultPath); err != nil {
			t.Fatalf("expected result file to exist: %v", err)
		}
		if _, err := os.Stat(githubOutputPath); err != nil {
			t.Fatalf("expected github output file to exist: %v", err)
		}

		ghData, _ := os.ReadFile(githubOutputPath)
		if !strings.Contains(string(ghData), "has_destroy=true\n") {
			t.Fatalf("expected has_destroy=true in github output, got: %s", string(ghData))
		}
	})
}
