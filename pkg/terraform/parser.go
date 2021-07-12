package terraform

import (
	"errors"
	"regexp"
	"strings"
)

// Parser is an interface for parsing terraform execution result
type Parser interface {
	Parse(body string) ParseResult
}

// ParseResult represents the result of parsed terraform execution
type ParseResult struct {
	Result             string
	OutsideTerraform   string
	ChangeResult       string
	Warnings           string
	HasAddOrUpdateOnly bool
	HasDestroy         bool
	HasNoChanges       bool
	HasPlanError       bool
	HasParseError      bool
	ExitCode           int
	Error              error
	CreatedResources   []string
	UpdatedResources   []string
	DeletedResources   []string
	ReplacedResources  []string
}

// DefaultParser is a parser for terraform commands
type DefaultParser struct{}

// PlanParser is a parser for terraform plan
type PlanParser struct {
	Pass         *regexp.Regexp
	Fail         *regexp.Regexp
	HasDestroy   *regexp.Regexp
	HasNoChanges *regexp.Regexp
	Create       *regexp.Regexp
	Update       *regexp.Regexp
	Delete       *regexp.Regexp
	Replace      *regexp.Regexp
}

// ApplyParser is a parser for terraform apply
type ApplyParser struct {
	Pass *regexp.Regexp
	Fail *regexp.Regexp
}

// NewDefaultParser is DefaultParser initializer
func NewDefaultParser() *DefaultParser {
	return &DefaultParser{}
}

// NewPlanParser is PlanParser initialized with its Regexp
func NewPlanParser() *PlanParser {
	return &PlanParser{
		Pass: regexp.MustCompile(`(?m)^(Plan: \d|No changes.)`),
		Fail: regexp.MustCompile(`(?m)^(Error: )`),
		// "0 to destroy" should be treated as "no destroy"
		HasDestroy:   regexp.MustCompile(`(?m)([1-9][0-9]* to destroy.)`),
		HasNoChanges: regexp.MustCompile(`(?m)^(No changes.)`),
		Create:       regexp.MustCompile(`^ *# (.*) will be created$`),
		Update:       regexp.MustCompile(`^ *# (.*) will be updated in-place$`),
		Delete:       regexp.MustCompile(`^ *# (.*) will be destroyed$`),
		Replace:      regexp.MustCompile(`^ *# (.*) must be replaced$`),
	}
}

// NewApplyParser is ApplyParser initialized with its Regexp
func NewApplyParser() *ApplyParser {
	return &ApplyParser{
		Pass: regexp.MustCompile(`(?m)^(Apply complete!)`),
		Fail: regexp.MustCompile(`(?m)^(Error: )`),
	}
}

// Parse returns ParseResult related with terraform commands
func (p *DefaultParser) Parse(body string) ParseResult {
	return ParseResult{
		Result:   body,
		ExitCode: ExitPass,
		Error:    nil,
	}
}

func extractResource(pattern *regexp.Regexp, line string) string {
	if arr := pattern.FindStringSubmatch(line); len(arr) == 2 { //nolint:gomnd
		return arr[1]
	}
	return ""
}

// Parse returns ParseResult related with terraform plan
func (p *PlanParser) Parse(body string) ParseResult { //nolint:cyclop
	var exitCode int
	switch {
	case p.Pass.MatchString(body):
		exitCode = ExitPass
	case p.Fail.MatchString(body):
		exitCode = ExitFail
	default:
		return ParseResult{
			Result:        "",
			HasParseError: true,
			ExitCode:      ExitFail,
			Error:         errors.New("cannot parse plan result"),
		}
	}
	lines := strings.Split(body, "\n")
	firstMatchLineIndex := -1
	var result, firstMatchLine string
	var createdResources, updatedResources, deletedResources, replacedResources []string
	startOutsideTerraform := -1
	endOutsideTerraform := -1
	startChangeOutput := -1
	endChangeOutput := -1
	startWarning := -1
	endWarning := -1
	for i, line := range lines {
		if line == "Note: Objects have changed outside of Terraform" { // https://github.com/hashicorp/terraform/blob/332045a4e4b1d256c45f98aac74e31102ace7af7/internal/command/views/plan.go#L403
			startOutsideTerraform = i + 1
		}
		if startOutsideTerraform != -1 && endOutsideTerraform == -1 && strings.HasPrefix(line, "Unless you have made equivalent changes to your configuration") { // https://github.com/hashicorp/terraform/blob/332045a4e4b1d256c45f98aac74e31102ace7af7/internal/command/views/plan.go#L110
			endOutsideTerraform = i + 1
		}
		if line == "Terraform will perform the following actions:" { // https://github.com/hashicorp/terraform/blob/332045a4e4b1d256c45f98aac74e31102ace7af7/internal/command/views/plan.go#L252
			startChangeOutput = i + 1
		}
		if startChangeOutput != -1 && endChangeOutput == -1 && strings.HasPrefix(line, "Plan: ") { // https://github.com/hashicorp/terraform/blob/dfc12a6a9e1cff323829026d51873c1b80200757/internal/command/views/plan.go#L306
			endChangeOutput = i + 1
		}
		if strings.HasPrefix(line, "Warning:") && startWarning == -1 {
			startWarning = i
		}
		if strings.HasPrefix(line, "─────") && startWarning != -1 && endWarning == -1 {
			endWarning = i
		}
		if firstMatchLineIndex == -1 {
			if p.Pass.MatchString(line) || p.Fail.MatchString(line) {
				firstMatchLineIndex = i
				firstMatchLine = line
			}
		}
		if rsc := extractResource(p.Create, line); rsc != "" {
			createdResources = append(createdResources, rsc)
		} else if rsc := extractResource(p.Update, line); rsc != "" {
			updatedResources = append(updatedResources, rsc)
		} else if rsc := extractResource(p.Delete, line); rsc != "" {
			deletedResources = append(deletedResources, rsc)
		} else if rsc := extractResource(p.Replace, line); rsc != "" {
			replacedResources = append(replacedResources, rsc)
		}
	}
	var hasPlanError bool
	switch {
	case p.Pass.MatchString(firstMatchLine):
		result = lines[firstMatchLineIndex]
	case p.Fail.MatchString(firstMatchLine):
		hasPlanError = true
		result = strings.Join(trimLastNewline(lines[firstMatchLineIndex:]), "\n")
	}

	hasDestroy := p.HasDestroy.MatchString(firstMatchLine)
	hasNoChanges := p.HasNoChanges.MatchString(firstMatchLine)
	HasAddOrUpdateOnly := !hasNoChanges && !hasDestroy && !hasPlanError

	outsideTerraform := ""
	if startOutsideTerraform != -1 {
		outsideTerraform = strings.Join(lines[startOutsideTerraform:endOutsideTerraform], "\n")
	}

	changeResult := ""
	if startChangeOutput != -1 {
		changeResult = strings.Join(lines[startChangeOutput:endChangeOutput], "\n")
	}

	warnings := ""
	if startWarning != -1 {
		if endWarning == -1 {
			warnings = strings.Join(lines[startWarning:], "\n")
		} else {
			warnings = strings.Join(lines[startWarning:endWarning], "\n")
		}
	}

	return ParseResult{
		Result:             result,
		ChangeResult:       changeResult,
		OutsideTerraform:   outsideTerraform,
		Warnings:           warnings,
		HasAddOrUpdateOnly: HasAddOrUpdateOnly,
		HasDestroy:         hasDestroy,
		HasNoChanges:       hasNoChanges,
		HasPlanError:       hasPlanError,
		ExitCode:           exitCode,
		Error:              nil,
		CreatedResources:   createdResources,
		UpdatedResources:   updatedResources,
		DeletedResources:   deletedResources,
		ReplacedResources:  replacedResources,
	}
}

// Parse returns ParseResult related with terraform apply
func (p *ApplyParser) Parse(body string) ParseResult {
	var exitCode int
	switch {
	case p.Pass.MatchString(body):
		exitCode = ExitPass
	case p.Fail.MatchString(body):
		exitCode = ExitFail
	default:
		return ParseResult{
			Result:        "",
			ExitCode:      ExitFail,
			HasParseError: true,
			Error:         errors.New("cannot parse apply result"),
		}
	}
	lines := strings.Split(body, "\n")
	var i int
	var result, line string
	for i, line = range lines {
		if p.Pass.MatchString(line) || p.Fail.MatchString(line) {
			break
		}
	}
	switch {
	case p.Pass.MatchString(line):
		result = lines[i]
	case p.Fail.MatchString(line):
		result = strings.Join(trimLastNewline(lines[i:]), "\n")
	}
	return ParseResult{
		Result:   result,
		ExitCode: exitCode,
		Error:    nil,
	}
}

func trimLastNewline(s []string) []string {
	if len(s) == 0 {
		return s
	}
	last := len(s) - 1
	if s[last] == "" {
		return s[:last]
	}
	return s
}
