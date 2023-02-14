package tofile

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/tfcmt/pkg/notifier"
	"github.com/suzuki-shunsuke/tfcmt/pkg/terraform"
)

// Plan posts comment optimized for notifications
func (g *NotifyService) Plan(ctx context.Context, param *notifier.ParamExec) (int, error) {
	cfg := g.client.Config
	parser := g.client.Config.Parser
	template := g.client.Config.Template
	var errMsgs []string

	result := parser.Parse(param.CombinedOutput)
	result.ExitCode = param.ExitCode
	if result.HasParseError {
		template = g.client.Config.ParseErrorTemplate
	} else {
		if result.Error != nil {
			return result.ExitCode, result.Error
		}
		if result.Result == "" {
			return result.ExitCode, result.Error
		}
	}

	template.SetValue(terraform.CommonTemplate{
		Result:                 result.Result,
		ChangedResult:          result.ChangedResult,
		ChangeOutsideTerraform: result.OutsideTerraform,
		Warning:                result.Warning,
		HasDestroy:             result.HasDestroy,
		UseRawOutput:           cfg.UseRawOutput,
		Vars:                   cfg.Vars,
		Templates:              cfg.Templates,
		Stdout:                 param.Stdout,
		Stderr:                 param.Stderr,
		CombinedOutput:         param.CombinedOutput,
		ExitCode:               param.ExitCode,
		ErrorMessages:          errMsgs,
		CreatedResources:       result.CreatedResources,
		UpdatedResources:       result.UpdatedResources,
		DeletedResources:       result.DeletedResources,
		ReplacedResources:      result.ReplacedResources,
	})
	body, err := template.Execute()
	if err != nil {
		return result.ExitCode, err
	}

	logE := logrus.WithFields(logrus.Fields{
		"program": "tfcmt",
	})

	logE.Debug("write output plan to file")
	// WriteFile to Outputfile define in config (via command -output)
	if err := g.client.Comment.Post(ctx, body, cfg.OutputFile); err != nil {
		return result.ExitCode, err
	}
	return result.ExitCode, nil
}
