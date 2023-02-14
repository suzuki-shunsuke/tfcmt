package local_file

import (
	"context"

	"github.com/suzuki-shunsuke/tfcmt/pkg/terraform"
)

// Client is a fake API client for write to local file
type Client struct {
	Debug bool

	Config *Config

	common service

	Notify  *NotifyService
	Comment *LocalFileService
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
	UseRawOutput       bool
}

type service struct {
	client *Client
}

// NewClient returns Client initialized with Config
func NewClient(ctx context.Context, cfg *Config) (*Client, error) {
	c := &Client{
		Config: cfg,
	}

	c.common.client = c

	c.Notify = (*NotifyService)(&c.common)

	return c, nil
}
