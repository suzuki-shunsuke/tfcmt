package localfile

import (
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/config"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/terraform"
)

// Client is a fake API client for write to local file
type Client struct {
	Debug bool

	Config *Config

	common service

	Notify *NotifyService
	Output *OutputService
}

// Config is a configuration for local file
type Config struct {
	OutputFile string
	Parser     terraform.Parser
	// Template is used for all Terraform command output
	Template           *terraform.Template
	ParseErrorTemplate *terraform.Template
	Vars               map[string]string
	EmbeddedVarNames   []string
	Templates          map[string]string
	CI                 string
	UseRawOutput       bool
	Masks              []*config.Mask

	// For labeling
	DisableLabel      bool
	GitHubLabelConfig *GitHubLabelConfig
}

type GitHubLabelConfig struct {
	BaseURL         string
	GraphQLEndpoint string
	Owner           string
	Repo            string
	PRNumber        int
	Revision        string
	Labels          interface{}
}

type service struct {
	client *Client
}

// NewClient returns Client initialized with Config
func NewClient(cfg *Config) (*Client, error) {
	c := &Client{
		Config: cfg,
	}

	c.common.client = c

	c.Notify = (*NotifyService)(&c.common)

	return c, nil
}
