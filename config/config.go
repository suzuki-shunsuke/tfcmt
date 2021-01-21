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
	CI        string            `yaml:"ci"`
	Notifier  Notifier          `yaml:"notifier"`
	Terraform Terraform         `yaml:"terraform"`
	Vars      map[string]string `yaml:"-"`
	Templates map[string]string

	path string
}

// Notifier is a notification notifier
type Notifier struct {
	Github GithubNotifier `yaml:"github"`
}

// GithubNotifier is a notifier for GitHub
type GithubNotifier struct {
	Token      string     `yaml:"token"`
	BaseURL    string     `yaml:"base_url"`
	Repository Repository `yaml:"repository"`
}

// Repository represents a GitHub repository
type Repository struct {
	Owner string `yaml:"owner"`
	Name  string `yaml:"name"`
}

// Terraform represents terraform configurations
type Terraform struct {
	Default      Default `yaml:"default"`
	Plan         Plan    `yaml:"plan"`
	Apply        Apply   `yaml:"apply"`
	UseRawOutput bool    `yaml:"use_raw_output,omitempty"`
}

// Default is a default setting for terraform commands
type Default struct {
	Template string `yaml:"template"`
}

// Plan is a terraform plan config
type Plan struct {
	Template            string              `yaml:"template"`
	WhenAddOrUpdateOnly WhenAddOrUpdateOnly `yaml:"when_add_or_update_only,omitempty"`
	WhenDestroy         WhenDestroy         `yaml:"when_destroy,omitempty"`
	WhenNoChanges       WhenNoChanges       `yaml:"when_no_changes,omitempty"`
	WhenPlanError       WhenPlanError       `yaml:"when_plan_error,omitempty"`
	WhenParseError      WhenParseError      `yaml:"when_parse_error,omitempty"`
	DisableLabel        bool                `yaml:"disable_label,omitempty"`
}

// WhenAddOrUpdateOnly is a configuration to notify the plan result contains new or updated in place resources
type WhenAddOrUpdateOnly struct {
	Label string `yaml:"label,omitempty"`
	Color string `yaml:"label_color,omitempty"`
}

// WhenDestroy is a configuration to notify the plan result contains destroy operation
type WhenDestroy struct {
	Label    string `yaml:"label,omitempty"`
	Template string `yaml:"template,omitempty"`
	Color    string `yaml:"label_color,omitempty"`
}

// WhenNoChange is a configuration to add a label when the plan result contains no change
type WhenNoChanges struct {
	Label string `yaml:"label,omitempty"`
	Color string `yaml:"label_color,omitempty"`
}

// WhenPlanError is a configuration to notify the plan result returns an error
type WhenPlanError struct {
	Label string `yaml:"label,omitempty"`
	Color string `yaml:"label_color,omitempty"`
}

// WhenParseError is a configuration to notify the plan result returns an error
type WhenParseError struct {
	Template string `yaml:"template"`
}

// Apply is a terraform apply config
type Apply struct {
	Template       string         `yaml:"template"`
	WhenParseError WhenParseError `yaml:"when_parse_error,omitempty"`
}

// LoadFile binds the config file to Config structure
func (cfg *Config) LoadFile(path string) error {
	cfg.path = path
	_, err := os.Stat(cfg.path)
	if err != nil {
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
