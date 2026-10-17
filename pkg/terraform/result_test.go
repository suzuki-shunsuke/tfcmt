package terraform_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/terraform"
)

func TestNewPlanResult(t *testing.T) {
	t.Parallel()

	t.Run("nil slices and vars are initialized as non-nil empty collections", func(t *testing.T) {
		t.Parallel()
		parseResult := terraform.ParseResult{
			HasNoChanges:       true,
			HasAddOrUpdateOnly: false,
			HasDestroy:         false,
			HasError:           false,
		}

		result := terraform.NewPlanResult(0, parseResult, nil)

		if result.ExitCode != 0 {
			t.Fatalf("expected ExitCode 0, got %d", result.ExitCode)
		}
		if !result.HasNoChanges {
			t.Fatalf("expected HasNoChanges true, got false")
		}
		if result.HasAddOrUpdateOnly || result.HasDestroy || result.HasError {
			t.Fatalf("unexpected boolean flags set")
		}
		if result.CreatedResources == nil || len(result.CreatedResources) != 0 {
			t.Fatalf("expected non-nil empty CreatedResources, got %#v", result.CreatedResources)
		}
		if result.UpdatedResources == nil || len(result.UpdatedResources) != 0 {
			t.Fatalf("expected non-nil empty UpdatedResources, got %#v", result.UpdatedResources)
		}
		if result.DeletedResources == nil || len(result.DeletedResources) != 0 {
			t.Fatalf("expected non-nil empty DeletedResources, got %#v", result.DeletedResources)
		}
		if result.ReplacedResources == nil || len(result.ReplacedResources) != 0 {
			t.Fatalf("expected non-nil empty ReplacedResources, got %#v", result.ReplacedResources)
		}
		if result.MovedResources == nil || len(result.MovedResources) != 0 {
			t.Fatalf("expected non-nil empty MovedResources, got %#v", result.MovedResources)
		}
		if result.ImportedResources == nil || len(result.ImportedResources) != 0 {
			t.Fatalf("expected non-nil empty ImportedResources, got %#v", result.ImportedResources)
		}
		if result.Vars == nil || len(result.Vars) != 0 {
			t.Fatalf("expected non-nil empty Vars, got %#v", result.Vars)
		}

		data, err := json.Marshal(result)
		if err != nil {
			t.Fatalf("failed to marshal result: %v", err)
		}
		jsonStr := string(data)
		expectedSubstrings := []string{
			`"created_resources":[]`,
			`"updated_resources":[]`,
			`"deleted_resources":[]`,
			`"replaced_resources":[]`,
			`"moved_resources":[]`,
			`"imported_resources":[]`,
			`"vars":{}`,
		}
		for _, expected := range expectedSubstrings {
			if !strings.Contains(jsonStr, expected) {
				t.Fatalf("expected JSON to contain %q, got: %s", expected, jsonStr)
			}
		}
	})

	t.Run("populated parse result is preserved", func(t *testing.T) {
		t.Parallel()
		parseResult := terraform.ParseResult{
			HasNoChanges:       false,
			HasAddOrUpdateOnly: true,
			HasDestroy:         false,
			HasError:           false,
			CreatedResources:   []string{"aws_s3_bucket.example"},
			UpdatedResources:   []string{"aws_iam_role.app"},
			DeletedResources:   []string{"aws_sqs_queue.old"},
			ReplacedResources:  []string{"aws_instance.worker"},
			MovedResources: []*terraform.MovedResource{
				{Before: "aws_s3_bucket.old", After: "aws_s3_bucket.new"},
			},
			ImportedResources: []string{"aws_route53_record.primary"},
		}
		vars := map[string]string{"target": "dev", "env": "stage"}

		result := terraform.NewPlanResult(2, parseResult, vars)

		if result.ExitCode != 2 {
			t.Fatalf("expected ExitCode 2, got %d", result.ExitCode)
		}
		if result.HasNoChanges {
			t.Fatalf("expected HasNoChanges false, got true")
		}
		if !result.HasAddOrUpdateOnly {
			t.Fatalf("expected HasAddOrUpdateOnly true, got false")
		}
		if len(result.CreatedResources) != 1 || result.CreatedResources[0] != "aws_s3_bucket.example" {
			t.Fatalf("unexpected CreatedResources: %#v", result.CreatedResources)
		}
		if len(result.MovedResources) != 1 || result.MovedResources[0].Before != "aws_s3_bucket.old" || result.MovedResources[0].After != "aws_s3_bucket.new" {
			t.Fatalf("unexpected MovedResources: %#v", result.MovedResources)
		}
		if result.Vars["target"] != "dev" || result.Vars["env"] != "stage" {
			t.Fatalf("unexpected Vars: %#v", result.Vars)
		}

		data, err := json.Marshal(result)
		if err != nil {
			t.Fatalf("failed to marshal result: %v", err)
		}
		jsonStr := string(data)
		expectedSubstrings := []string{
			`"exit_code":2`,
			`"has_no_changes":false`,
			`"has_add_or_update_only":true`,
			`"created_resources":["aws_s3_bucket.example"]`,
			`"moved_resources":[{"before":"aws_s3_bucket.old","after":"aws_s3_bucket.new"}]`,
			`"target":"dev"`,
		}
		for _, expected := range expectedSubstrings {
			if !strings.Contains(jsonStr, expected) {
				t.Fatalf("expected JSON to contain %q, got: %s", expected, jsonStr)
			}
		}
	})
}

func TestWriteResultFile(t *testing.T) {
	t.Parallel()

	t.Run("successfully writes formatted JSON to new and nested file path", func(t *testing.T) {
		t.Parallel()
		tempDir := t.TempDir()
		resultPath := filepath.Join(tempDir, "nested", "dir", "result.json")

		planResult := terraform.PlanResult{
			ExitCode:           2,
			HasNoChanges:       false,
			HasAddOrUpdateOnly: true,
			HasDestroy:         false,
			HasError:           false,
			CreatedResources:   []string{"aws_s3_bucket.example"},
			UpdatedResources:   []string{},
			DeletedResources:   []string{},
			ReplacedResources:  []string{},
			MovedResources:     []*terraform.MovedResource{},
			ImportedResources:  []string{},
			Vars:               map[string]string{"target": "dev"},
		}

		if err := terraform.WriteResultFile(resultPath, planResult); err != nil {
			t.Fatalf("WriteResultFile failed: %v", err)
		}

		content, err := os.ReadFile(resultPath)
		if err != nil {
			t.Fatalf("failed to read result file: %v", err)
		}

		if !strings.HasSuffix(string(content), "\n") {
			t.Fatalf("expected trailing newline in result file")
		}

		var parsed terraform.PlanResult
		if err := json.Unmarshal(content, &parsed); err != nil {
			t.Fatalf("failed to unmarshal written JSON: %v", err)
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
		if len(parsed.CreatedResources) != 1 || parsed.CreatedResources[0] != "aws_s3_bucket.example" {
			t.Fatalf("unexpected CreatedResources: %#v", parsed.CreatedResources)
		}
		if parsed.Vars["target"] != "dev" {
			t.Fatalf("unexpected Vars: %#v", parsed.Vars)
		}
	})

	t.Run("overwrites existing file cleanly", func(t *testing.T) {
		t.Parallel()
		tempDir := t.TempDir()
		resultPath := filepath.Join(tempDir, "result.json")

		if err := os.WriteFile(resultPath, []byte("stale-content-that-should-be-truncated"), 0o644); err != nil {
			t.Fatalf("failed to write initial file: %v", err)
		}

		planResult := terraform.PlanResult{
			ExitCode:           0,
			HasNoChanges:       true,
			HasAddOrUpdateOnly: false,
			HasDestroy:         false,
			HasError:           false,
			CreatedResources:   []string{},
			UpdatedResources:   []string{},
			DeletedResources:   []string{},
			ReplacedResources:  []string{},
			MovedResources:     []*terraform.MovedResource{},
			ImportedResources:  []string{},
			Vars:               map[string]string{},
		}

		if err := terraform.WriteResultFile(resultPath, planResult); err != nil {
			t.Fatalf("WriteResultFile failed: %v", err)
		}

		content, err := os.ReadFile(resultPath)
		if err != nil {
			t.Fatalf("failed to read result file: %v", err)
		}

		if strings.Contains(string(content), "stale-content") {
			t.Fatalf("expected stale content to be overwritten, got: %s", string(content))
		}
	})

	t.Run("returns error when directory path is invalid", func(t *testing.T) {
		t.Parallel()
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "file.txt")
		if err := os.WriteFile(filePath, []byte("regular file"), 0o644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}

		invalidPath := filepath.Join(filePath, "cannot_be_child", "result.json")
		planResult := terraform.PlanResult{}
		if err := terraform.WriteResultFile(invalidPath, planResult); err == nil {
			t.Fatalf("expected error for invalid directory path, got nil")
		}
	})
}

func TestFormatGitHubOutput(t *testing.T) {
	t.Parallel()

	planResult := terraform.PlanResult{
		ExitCode:           2,
		HasNoChanges:       false,
		HasAddOrUpdateOnly: true,
		HasDestroy:         false,
		HasError:           false,
		CreatedResources:   []string{"aws_s3_bucket.example"},
		UpdatedResources:   []string{},
		DeletedResources:   []string{},
		ReplacedResources:  []string{},
		MovedResources: []*terraform.MovedResource{
			{Before: "aws_s3_bucket.old", After: "aws_s3_bucket.new"},
		},
		ImportedResources: []string{},
		Vars:              map[string]string{"target": "dev"},
	}

	formatted, err := terraform.FormatGitHubOutput(planResult)
	if err != nil {
		t.Fatalf("FormatGitHubOutput failed: %v", err)
	}

	lines := strings.Split(strings.TrimSuffix(formatted, "\n"), "\n")
	kv := make(map[string]string, len(lines))
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			t.Fatalf("malformed output line: %q", line)
		}
		kv[parts[0]] = parts[1]
	}

	expectedValues := map[string]string{
		"exit_code":              "2",
		"has_no_changes":         "false",
		"has_add_or_update_only": "true",
		"has_destroy":            "false",
		"has_error":              "false",
		"created_resources":      `["aws_s3_bucket.example"]`,
		"updated_resources":      `[]`,
		"deleted_resources":      `[]`,
		"replaced_resources":     `[]`,
		"moved_resources":        `[{"before":"aws_s3_bucket.old","after":"aws_s3_bucket.new"}]`,
		"imported_resources":     `[]`,
		"vars":                   `{"target":"dev"}`,
	}

	for k, expected := range expectedValues {
		actual, ok := kv[k]
		if !ok {
			t.Fatalf("missing key %q in formatted output", k)
		}
		if actual != expected {
			t.Fatalf("key %q: expected %q, got %q", k, expected, actual)
		}
	}
}

func TestExportGitHubOutput(t *testing.T) {
	t.Parallel()

	t.Run("appends output to file", func(t *testing.T) {
		t.Parallel()
		tempDir := t.TempDir()
		outputPath := filepath.Join(tempDir, "github_output")

		if err := os.WriteFile(outputPath, []byte("existing_key=existing_value\n"), 0o644); err != nil {
			t.Fatalf("failed to create existing file: %v", err)
		}

		planResult := terraform.PlanResult{
			ExitCode:           0,
			HasNoChanges:       true,
			HasAddOrUpdateOnly: false,
			HasDestroy:         false,
			HasError:           false,
			CreatedResources:   []string{},
			UpdatedResources:   []string{},
			DeletedResources:   []string{},
			ReplacedResources:  []string{},
			MovedResources:     []*terraform.MovedResource{},
			ImportedResources:  []string{},
			Vars:               map[string]string{},
		}

		if err := terraform.ExportGitHubOutput(outputPath, planResult); err != nil {
			t.Fatalf("ExportGitHubOutput failed: %v", err)
		}

		content, err := os.ReadFile(outputPath)
		if err != nil {
			t.Fatalf("failed to read output file: %v", err)
		}

		str := string(content)
		if !strings.HasPrefix(str, "existing_key=existing_value\n") {
			t.Fatalf("expected file to retain existing content, got: %s", str)
		}
		if !strings.Contains(str, "has_no_changes=true\n") {
			t.Fatalf("expected file to contain has_no_changes=true, got: %s", str)
		}
		if !strings.Contains(str, "exit_code=0\n") {
			t.Fatalf("expected file to contain exit_code=0, got: %s", str)
		}
	})

	t.Run("returns error on invalid output path", func(t *testing.T) {
		t.Parallel()
		invalidPath := filepath.Join("/nonexistent_dir", "sub", "output")
		planResult := terraform.PlanResult{}
		if err := terraform.ExportGitHubOutput(invalidPath, planResult); err == nil {
			t.Fatalf("expected error for invalid output path, got nil")
		}
	})
}
