package terraform

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStripLines(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		pattern  string
		input    string
		expected string
	}{
		{
			name:     "empty input",
			pattern:  "anything",
			input:    "",
			expected: "",
		},
		{
			name:     "no matches keeps input intact",
			pattern:  "Refreshing state",
			input:    "line a\nline b\nline c",
			expected: "line a\nline b\nline c",
		},
		{
			name:    "matching lines removed, others kept",
			pattern: "Refreshing state",
			input: `data.terraform_remote_state.foo: Refreshing state...
google_project.my_project: Refreshing state... [id=my-project]
aws_iam_user.alice: Read complete after 0s [id=alice]

Plan: 1 to add, 0 to change, 0 to destroy.`,
			expected: `aws_iam_user.alice: Read complete after 0s [id=alice]

Plan: 1 to add, 0 to change, 0 to destroy.`,
		},
		{
			name:     "all lines match yields empty result",
			pattern:  "x",
			input:    "x\nxx\nxxx",
			expected: "",
		},
		{
			name:     "trailing newline preserved when last line is empty",
			pattern:  "drop",
			input:    "drop me\nkeep me\n",
			expected: "keep me\n",
		},
		{
			name:     "anchored pattern matches only at line start",
			pattern:  "^Refreshing",
			input:    "Refreshing state...\nfoo: Refreshing state...\nkeep",
			expected: "foo: Refreshing state...\nkeep",
		},
		{
			name:     "regex special characters work as Go regexp",
			pattern:  `Refreshing state\.\.\.`,
			input:    "a: Refreshing state...\nb: refreshing state (not dots)\nc",
			expected: "b: refreshing state (not dots)\nc",
		},
		{
			name:     "carriage returns preserved",
			pattern:  "drop",
			input:    "drop me\r\nkeep me\r\n",
			expected: "keep me\r\n",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			actual, err := stripLines(testCase.pattern, testCase.input)
			if err != nil {
				t.Fatalf("stripLines returned an unexpected error: %v", err)
			}
			if diff := cmp.Diff(actual, testCase.expected); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestStripLinesInvalidPattern(t *testing.T) {
	t.Parallel()
	if _, err := stripLines("(unclosed", "anything"); err == nil {
		t.Fatal("expected error for invalid regexp pattern, got nil")
	}
}

// TestStripLinesInTemplate verifies stripLines is registered in the html/template
// FuncMap (the default flow) and can be chained with wrapCode.
func TestStripLinesInTemplate(t *testing.T) {
	t.Parallel()
	tpl := NewPlanTemplate(`{{ stripLines "Refreshing state" .CombinedOutput | wrapCode }}`)
	tpl.SetValue(CommonTemplate{
		CombinedOutput: `data.terraform_remote_state.foo: Refreshing state...
aws_iam_user.alice: Refreshing state... [id=alice]
Plan: 0 to add, 1 to change, 0 to destroy.`,
	})

	actual, err := tpl.Execute()
	if err != nil {
		t.Fatalf("template execution failed: %v", err)
	}
	expected := "\n```hcl\nPlan: 0 to add, 1 to change, 0 to destroy.\n```\n"
	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Error(diff)
	}
}

// TestStripLinesInRawTemplate verifies stripLines is registered in the
// text/template FuncMap used when UseRawOutput is true.
func TestStripLinesInRawTemplate(t *testing.T) {
	t.Parallel()
	tpl := NewPlanTemplate(`{{ stripLines "noisy" .CombinedOutput }}`)
	tpl.SetValue(CommonTemplate{
		UseRawOutput:   true,
		CombinedOutput: "keep\nnoisy line\nkeep",
	})

	actual, err := tpl.Execute()
	if err != nil {
		t.Fatalf("template execution failed: %v", err)
	}
	expected := "keep\nkeep"
	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Error(diff)
	}
}
