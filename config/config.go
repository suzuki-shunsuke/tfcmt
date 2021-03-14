package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/suzuki-shunsuke/go-findconfig/findconfig"
	"gopkg.in/yaml.v2"
)

// Config is for tfcmt config structure
type Config struct {
	CI        string
	Notifier  Notifier
	Terraform Terraform
	Vars      map[string]string `yaml:"-"`
	Templates map[string]string
	Log       Log

	path string
}

type Log struct {
	Level string
	// Format string
}

// Notifier is a notification notifier
type Notifier struct {
	Github GithubNotifier
}

// GithubNotifier is a notifier for GitHub
type GithubNotifier struct {
	Token      string
	BaseURL    string `yaml:"base_url"`
	Repository Repository
}

// Repository represents a GitHub repository
type Repository struct {
	Owner string
	Name  string
}

// Terraform represents terraform configurations
type Terraform struct {
	Default      Default
	Plan         Plan
	Apply        Apply
	UseRawOutput bool `yaml:"use_raw_output"`
}

// Default is a default setting for terraform commands
type Default struct {
	Template string
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
	Label string
	Color string `yaml:"label_color"`
}

// WhenDestroy is a configuration to notify the plan result contains destroy operation
type WhenDestroy struct {
	Label    string
	Template string
	Color    string `yaml:"label_color"`
}

// WhenNoChange is a configuration to add a label when the plan result contains no change
type WhenNoChanges struct {
	Label string
	Color string `yaml:"label_color"`
}

// WhenPlanError is a configuration to notify the plan result returns an error
type WhenPlanError struct {
	Label string
	Color string `yaml:"label_color"`
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
	cfg.path = path
	if _, err := os.Stat(cfg.path); err != nil {
		return fmt.Errorf("%s: no config file", cfg.path)
	}
	raw, _ := ioutil.ReadFile(cfg.path)
	return yaml.Unmarshal(raw, cfg)
}

// Validation validates config file
func (cfg *Config) Validation() error {
	switch strings.ToLower(cfg.CI) {
	case "":
		break
	case "circleci":
		// ok pattern
	case "codebuild":
		// ok pattern
	case "github-actions":
		// ok pattern
	case "cloud-build", "cloudbuild":
		// ok pattern
	default:
		return fmt.Errorf("%s: not supported yet", cfg.CI)
	}
	if cfg.Notifier.Github.Repository.Owner == "" {
		return errors.New("repository owner is missing")
	}
	if cfg.Notifier.Github.Repository.Name == "" {
		return errors.New("repository name is missing")
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
