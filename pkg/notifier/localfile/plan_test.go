package localfile_test

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/notifier"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/notifier/localfile"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/terraform"
)

const planWithWarning = `
Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # terraform_data.a will be created
  + resource "terraform_data" "a" {
      + id = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.
╷
│ Warning: Resource targeting is in effect
│
│ You are creating a plan with the -target option, which means that the
│ result of this plan may not represent all of the changes requested by the
│ current configuration.
╵

─────────────────────────────────────────────────────────────────────────────
`

func TestNotifyService_Plan_ignoreWarning(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.DiscardHandler)
	for _, tc := range []struct {
		name          string
		ignoreWarning bool
		wantWarning   bool
	}{
		{name: "warnings section is written by default", ignoreWarning: false, wantWarning: true},
		{name: "warnings section is omitted with ignore_warning", ignoreWarning: true, wantWarning: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			outputFile := filepath.Join(t.TempDir(), "result.md")
			client, err := localfile.NewClient(&localfile.Config{
				OutputFile:         outputFile,
				Parser:             terraform.NewPlanParser(),
				Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
				ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate),
				IgnoreWarning:      tc.ignoreWarning,
				DisableLabel:       true,
			}, nil)
			if err != nil {
				t.Fatal(err)
			}
			if err := client.Notify.Plan(t.Context(), logger, &notifier.ParamExec{
				CombinedOutput: planWithWarning,
			}); err != nil {
				t.Fatal(err)
			}
			b, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatal(err)
			}
			if got := strings.Contains(string(b), "## :warning: Warnings"); got != tc.wantWarning {
				t.Fatalf("output contains the Warnings section = %v, want %v\n%s", got, tc.wantWarning, b)
			}
		})
	}
}
