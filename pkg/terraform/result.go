package terraform

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

// PlanResult represents the machine-readable execution and parsed plan result.
type PlanResult struct {
	ExitCode           int               `json:"exit_code"`
	HasNoChanges       bool              `json:"has_no_changes"`
	HasAddOrUpdateOnly bool              `json:"has_add_or_update_only"`
	HasDestroy         bool              `json:"has_destroy"`
	HasError           bool              `json:"has_error"`
	CreatedResources   []string          `json:"created_resources"`
	UpdatedResources   []string          `json:"updated_resources"`
	DeletedResources   []string          `json:"deleted_resources"`
	ReplacedResources  []string          `json:"replaced_resources"`
	MovedResources     []*MovedResource  `json:"moved_resources"`
	ImportedResources  []string          `json:"imported_resources"`
	Vars               map[string]string `json:"vars"`
}

// NewPlanResult creates a PlanResult with empty non-nil slices and vars map.
func NewPlanResult(exitCode int, parseResult ParseResult, vars map[string]string) PlanResult {
	created := parseResult.CreatedResources
	if created == nil {
		created = []string{}
	}
	updated := parseResult.UpdatedResources
	if updated == nil {
		updated = []string{}
	}
	deleted := parseResult.DeletedResources
	if deleted == nil {
		deleted = []string{}
	}
	replaced := parseResult.ReplacedResources
	if replaced == nil {
		replaced = []string{}
	}
	moved := parseResult.MovedResources
	if moved == nil {
		moved = []*MovedResource{}
	}
	imported := parseResult.ImportedResources
	if imported == nil {
		imported = []string{}
	}
	v := vars
	if v == nil {
		v = map[string]string{}
	}

	return PlanResult{
		ExitCode:           exitCode,
		HasNoChanges:       parseResult.HasNoChanges,
		HasAddOrUpdateOnly: parseResult.HasAddOrUpdateOnly,
		HasDestroy:         parseResult.HasDestroy,
		HasError:           parseResult.HasError,
		CreatedResources:   created,
		UpdatedResources:   updated,
		DeletedResources:   deleted,
		ReplacedResources:  replaced,
		MovedResources:     moved,
		ImportedResources:  imported,
		Vars:               v,
	}
}

// WriteResultFile serializes PlanResult to a file formatted as pretty-printed JSON.
func WriteResultFile(filePath string, result PlanResult) error {
	dir := filepath.Dir(filePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create directory for result file: %w", err)
		}
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal plan result to JSON: %w", err)
	}
	data = append(data, '\n')

	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		return fmt.Errorf("write plan result file: %w", err)
	}
	return nil
}

// FormatGitHubOutput formats PlanResult as key=value lines for $GITHUB_OUTPUT.
func FormatGitHubOutput(result PlanResult) (string, error) {
	createdJSON, err := json.Marshal(result.CreatedResources)
	if err != nil {
		return "", fmt.Errorf("marshal created_resources: %w", err)
	}
	updatedJSON, err := json.Marshal(result.UpdatedResources)
	if err != nil {
		return "", fmt.Errorf("marshal updated_resources: %w", err)
	}
	deletedJSON, err := json.Marshal(result.DeletedResources)
	if err != nil {
		return "", fmt.Errorf("marshal deleted_resources: %w", err)
	}
	replacedJSON, err := json.Marshal(result.ReplacedResources)
	if err != nil {
		return "", fmt.Errorf("marshal replaced_resources: %w", err)
	}
	movedJSON, err := json.Marshal(result.MovedResources)
	if err != nil {
		return "", fmt.Errorf("marshal moved_resources: %w", err)
	}
	importedJSON, err := json.Marshal(result.ImportedResources)
	if err != nil {
		return "", fmt.Errorf("marshal imported_resources: %w", err)
	}
	varsJSON, err := json.Marshal(result.Vars)
	if err != nil {
		return "", fmt.Errorf("marshal vars: %w", err)
	}

	lines := fmt.Sprintf(
		"exit_code=%d\n"+
			"has_no_changes=%s\n"+
			"has_add_or_update_only=%s\n"+
			"has_destroy=%s\n"+
			"has_error=%s\n"+
			"created_resources=%s\n"+
			"updated_resources=%s\n"+
			"deleted_resources=%s\n"+
			"replaced_resources=%s\n"+
			"moved_resources=%s\n"+
			"imported_resources=%s\n"+
			"vars=%s\n",
		result.ExitCode,
		strconv.FormatBool(result.HasNoChanges),
		strconv.FormatBool(result.HasAddOrUpdateOnly),
		strconv.FormatBool(result.HasDestroy),
		strconv.FormatBool(result.HasError),
		string(createdJSON),
		string(updatedJSON),
		string(deletedJSON),
		string(replacedJSON),
		string(movedJSON),
		string(importedJSON),
		string(varsJSON),
	)

	return lines, nil
}

// ExportGitHubOutput appends formatted step outputs to the file specified by outputPath.
func ExportGitHubOutput(outputPath string, result PlanResult) error {
	content, err := FormatGitHubOutput(result)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open GITHUB_OUTPUT file: %w", err)
	}
	defer file.Close()

	if _, err := io.WriteString(file, content); err != nil {
		return fmt.Errorf("write to GITHUB_OUTPUT file: %w", err)
	}
	return nil
}
