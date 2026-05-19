package terraform

import (
	"strings"
	"testing"
)

func TestStripLines(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		pattern string
		input   string
		want    string
	}{
		{
			name:    "empty input",
			pattern: "anything",
			input:   "",
			want:    "",
		},
		{
			name:    "no matches keeps input intact",
			pattern: "Refreshing state",
			input:   "line a\nline b\nline c",
			want:    "line a\nline b\nline c",
		},
		{
			name:    "matching lines removed, others kept",
			pattern: "Refreshing state",
			input: `data.terraform_remote_state.foo: Refreshing state...
google_project.my_project: Refreshing state... [id=my-project]
aws_iam_user.alice: Read complete after 0s [id=alice]

Plan: 1 to add, 0 to change, 0 to destroy.`,
			want: `aws_iam_user.alice: Read complete after 0s [id=alice]

Plan: 1 to add, 0 to change, 0 to destroy.`,
		},
		{
			name:    "all lines match yields empty result",
			pattern: "x",
			input:   "x\nxx\nxxx",
			want:    "",
		},
		{
			name:    "trailing newline preserved when last line is empty",
			pattern: "drop",
			input:   "drop me\nkeep me\n",
			want:    "keep me\n",
		},
		{
			name:    "anchored pattern matches only at line start",
			pattern: "^Refreshing",
			input:   "Refreshing state...\nfoo: Refreshing state...\nkeep",
			want:    "foo: Refreshing state...\nkeep",
		},
		{
			name:    "regex special characters work as Go regexp",
			pattern: `Refreshing state\.\.\.`,
			input:   "a: Refreshing state...\nb: refreshing state (not dots)\nc",
			want:    "b: refreshing state (not dots)\nc",
		},
		{
			name:    "carriage returns preserved",
			pattern: "drop",
			input:   "drop me\r\nkeep me\r\n",
			want:    "keep me\r\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := stripLines(tc.pattern, tc.input)
			if err != nil {
				t.Fatalf("stripLines returned error: %v", err)
			}
			if got != tc.want {
				t.Errorf("stripLines(%q, %q):\nwant %q\n got %q", tc.pattern, tc.input, tc.want, got)
			}
		})
	}
}

func TestStripLinesInvalidPattern(t *testing.T) {
	t.Parallel()
	_, err := stripLines("(unclosed", "anything")
	if err == nil {
		t.Fatal("expected error for invalid regexp pattern, got nil")
	}
}

// TestStripLinesInTemplate verifies the helper is registered in the template
// FuncMap and can be chained with wrapCode in a realistic plan template.
func TestStripLinesInTemplate(t *testing.T) {
	t.Parallel()
	tpl := NewPlanTemplate(`{{ stripLines "Refreshing state" .CombinedOutput | wrapCode }}`)
	tpl.SetValue(CommonTemplate{
		CombinedOutput: `data.terraform_remote_state.foo: Refreshing state...
aws_iam_user.alice: Refreshing state... [id=alice]
Plan: 0 to add, 1 to change, 0 to destroy.`,
	})

	got, err := tpl.Execute()
	if err != nil {
		t.Fatalf("template execution failed: %v", err)
	}

	if strings.Contains(got, "Refreshing state") {
		t.Errorf("expected `Refreshing state` lines to be stripped, got:\n%s", got)
	}
	if !strings.Contains(got, "Plan: 0 to add, 1 to change, 0 to destroy.") {
		t.Errorf("expected Plan summary line to survive, got:\n%s", got)
	}
	if !strings.Contains(got, "```hcl") {
		t.Errorf("expected wrapCode to fence the result with ```hcl, got:\n%s", got)
	}
}

// TestStripLinesInRawTemplate verifies the helper works in the text/template
// flow used when UseRawOutput is true.
func TestStripLinesInRawTemplate(t *testing.T) {
	t.Parallel()
	tpl := NewPlanTemplate(`{{ stripLines "noisy" .CombinedOutput }}`)
	tpl.SetValue(CommonTemplate{
		UseRawOutput:   true,
		CombinedOutput: "keep\nnoisy line\nkeep",
	})

	got, err := tpl.Execute()
	if err != nil {
		t.Fatalf("template execution failed: %v", err)
	}
	if strings.Contains(got, "noisy line") {
		t.Errorf("expected noisy line to be stripped, got:\n%s", got)
	}
	if !strings.Contains(got, "keep\nkeep") {
		t.Errorf("expected kept lines to be adjacent after strip, got:\n%s", got)
	}
}
