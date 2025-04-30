package localfile

import (
	"context"

	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/config"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/notifier/github"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/terraform"
)

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
	DisableLabel bool
}

type GitHubLabelConfig struct {
	BaseURL         string
	GraphQLEndpoint string
	Owner           string
	Repo            string
	PRNumber        int
	Revision        string
	Labels          github.ResultLabels
}

type Labeler interface {
	UpdateLabels(ctx context.Context, result terraform.ParseResult) []string
}
