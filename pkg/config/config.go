package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/suzuki-shunsuke/go-findconfig/findconfig"
	"gopkg.in/yaml.v2"
)

// Config is for tfcmt config structure
type Config struct {
	CI                 CI `yaml:"-"`
	Terraform          Terraform
	Vars               map[string]string `yaml:"-"`
	EmbeddedVarNames   []string          `yaml:"embedded_var_names"`
	Templates          map[string]string
	Log                Log
	GHEBaseURL         string `yaml:"ghe_base_url"`
	GHEGraphQLEndpoint string `yaml:"ghe_graphql_endpoint"`
	GitHubToken        string `yaml:"-"`
	PlanPatch          bool   `yaml:"plan_patch"`
	RepoOwner          string `yaml:"repo_owner"`
	RepoName           string `yaml:"repo_name"`
	Output             string `yaml:"-"`
}

type CI struct {
	Name     string
	Owner    string
	Repo     string
	SHA      string
	Link     string
	PRNumber int
}

type Log struct {
	Level string
	// Format string
}

// Terraform represents terraform configurations
type Terraform struct {
	Plan         Plan
	Apply        Apply
	UseRawOutput bool `yaml:"use_raw_output"`
}

// Plan is a terraform plan config
type Plan struct {
	Template            string
	WhenAddOrUpdateOnly WhenAddOrUpdateOnly `yaml:"when_add_or_update_only"`
	WhenDestroy         WhenDestroy         `yaml:"when_destroy"`
	WhenNoChanges       WhenNoChanges       `yaml:"when_no_changes"`
	WhenPlanError       WhenPlanError       `yaml:"when_plan_error"`
	WhenParseError      WhenParseError      `yaml:"when_parse_error"`
	DisableLabel        bool                `yaml:"disable_label"`
}

// WhenAddOrUpdateOnly is a configuration to notify the plan result contains new or updated in place resources
type WhenAddOrUpdateOnly struct {
	Label        string
	Color        string `yaml:"label_color"`
	DisableLabel bool   `yaml:"disable_label"`
}

// WhenDestroy is a configuration to notify the plan result contains destroy operation
type WhenDestroy struct {
	Label        string
	Color        string `yaml:"label_color"`
	DisableLabel bool   `yaml:"disable_label"`
}

// WhenNoChanges is a configuration to add a label when the plan result contains no change
type WhenNoChanges struct {
	Label          string
	Color          string `yaml:"label_color"`
	DisableLabel   bool   `yaml:"disable_label"`
	DisableComment bool   `yaml:"disable_comment"`
}

// WhenPlanError is a configuration to notify the plan result returns an error
type WhenPlanError struct {
	Label        string
	Color        string `yaml:"label_color"`
	DisableLabel bool   `yaml:"disable_label"`
}

// WhenParseError is a configuration to notify the plan result returns an error
type WhenParseError struct {
	Template string
}

// Apply is a terraform apply config
type Apply struct {
	Template       string
	WhenParseError WhenParseError `yaml:"when_parse_error"`
}

// LoadFile binds the config file to Config structure
func (cfg *Config) LoadFile(path string) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("%s: no config file", path)
	}
	raw, _ := os.ReadFile(path)
	return yaml.Unmarshal(raw, cfg)
}

// Validate validates config file
func (cfg *Config) Validate() error {
	if cfg.Output != "" {
		return nil
	}
	if cfg.CI.Owner == "" {
		return errors.New("repository owner is missing")
	}

	if cfg.CI.Repo == "" {
		return errors.New("repository name is missing")
	}

	if cfg.CI.SHA == "" && cfg.CI.PRNumber <= 0 {
		return errors.New("pull request number or SHA (revision) is needed")
	}
	return nil
}

// Find returns config path
func (cfg *Config) Find(file string) (string, error) {
	if file != "" {
		if _, err := os.Stat(file); err == nil {
			return file, nil
		}
		return "", errors.New("config for tfcmt is not found at all")
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get a current directory path: %w", err)
	}
	if p := findconfig.Find(wd, findconfig.Exist, "tfcmt.yaml", "tfcmt.yml", ".tfcmt.yaml", ".tfcmt.yml"); p != "" {
		return p, nil
	}
	return "", nil
}
