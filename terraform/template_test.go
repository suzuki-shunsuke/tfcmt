package terraform

import (
	"testing"
)

func TestPlanTemplateExecute(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		template string
		value    CommonTemplate
		resp     string
	}{
		{
			name:     "case 0",
			template: DefaultPlanTemplate,
			value:    CommonTemplate{},
			resp: `
## Plan Result



<details><summary>Details (Click me)</summary>

<pre><code>
</code></pre></details>
`,
		},
		{
			name:     "case 1",
			template: DefaultPlanTemplate,
			value: CommonTemplate{
				Result: "result",
				Body:   "body",
			},
			resp: `
## Plan Result


<pre><code>result
</code></pre>


<details><summary>Details (Click me)</summary>

<pre><code>body
</code></pre></details>
`,
		},
		{
			name:     "case 2",
			template: DefaultPlanTemplate,
			value: CommonTemplate{
				Result: "",
				Body:   "body",
			},
			resp: `
## Plan Result



<details><summary>Details (Click me)</summary>

<pre><code>body
</code></pre></details>
`,
		},
		{
			name:     "case 3",
			template: DefaultPlanTemplate,
			value: CommonTemplate{
				Result: "",
				Body:   `This is a "body".`,
			},
			resp: `
## Plan Result



<details><summary>Details (Click me)</summary>

<pre><code>This is a &#34;body&#34;.
</code></pre></details>
`,
		},
		{
			name:     "case 4",
			template: DefaultPlanTemplate,
			value: CommonTemplate{
				Result:       "",
				Body:         `This is a "body".`,
				UseRawOutput: true,
			},
			resp: `
## Plan Result



<details><summary>Details (Click me)</summary>

<pre><code>This is a "body".
</code></pre></details>
`,
		},
		{
			name:     "case 5",
			template: "",
			value: CommonTemplate{
				Result: "",
				Body:   "body",
			},
			resp: `
## Plan Result



<details><summary>Details (Click me)</summary>

<pre><code>body
</code></pre></details>
`,
		},
		{
			name:     "case 6",
			template: `{{ .Result }}-{{ .Body }}`,
			value: CommonTemplate{
				Result: "c",
				Body:   "d",
			},
			resp: `c-d`,
		},
	}
	for i, testCase := range testCases {
		testCase := testCase
		if testCase.name == "" {
			t.Fatalf("testCase.name is required: index %d", i)
		}
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			template := NewPlanTemplate(testCase.template)
			template.SetValue(testCase.value)
			resp, err := template.Execute()
			if err != nil {
				t.Fatal(err)
			}
			if resp != testCase.resp {
				t.Errorf("got %s but want %s", resp, testCase.resp)
			}
		})
	}
}

func TestDestroyWarningTemplateExecute(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		template string
		value    CommonTemplate
		resp     string
	}{
		{
			name:     "case 0",
			template: DefaultDestroyWarningTemplate,
			value:    CommonTemplate{},
			resp: `
## :warning: Plan Result: Resource Deletion will happen :warning:

This plan contains resource delete operation. Please check the plan result very carefully!


`,
		},
		{
			name:     "case 1",
			template: DefaultDestroyWarningTemplate,
			value: CommonTemplate{
				Result: `This is a "result".`,
			},
			resp: `
## :warning: Plan Result: Resource Deletion will happen :warning:

This plan contains resource delete operation. Please check the plan result very carefully!


<pre><code>This is a &#34;result&#34;.
</code></pre>

`,
		},
		{
			name:     "case 2",
			template: DefaultDestroyWarningTemplate,
			value: CommonTemplate{
				Result:       `This is a "result".`,
				UseRawOutput: true,
			},
			resp: `
## :warning: Plan Result: Resource Deletion will happen :warning:

This plan contains resource delete operation. Please check the plan result very carefully!


<pre><code>This is a "result".
</code></pre>

`,
		},
		{
			name:     "case 3",
			template: DefaultDestroyWarningTemplate,
			value: CommonTemplate{
				Result: "",
			},
			resp: `
## :warning: Plan Result: Resource Deletion will happen :warning:

This plan contains resource delete operation. Please check the plan result very carefully!


`,
		},
		{
			name:     "case 4",
			template: "",
			value: CommonTemplate{
				Result: "",
			},
			resp: `
## :warning: Plan Result: Resource Deletion will happen :warning:

This plan contains resource delete operation. Please check the plan result very carefully!


`,
		},
		{
			name:     "case 5",
			template: `{{ .Result }}-{{ .Body }}`,
			value: CommonTemplate{
				Result: "c",
				Body:   "d",
			},
			resp: `c-d`,
		},
	}
	for i, testCase := range testCases {
		testCase := testCase
		if testCase.name == "" {
			t.Fatalf("testCase.name is required: index %d", i)
		}
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			template := NewDestroyWarningTemplate(testCase.template)
			template.SetValue(testCase.value)
			resp, err := template.Execute()
			if err != nil {
				t.Fatal(err)
			}
			if resp != testCase.resp {
				t.Errorf("got %s but want %s", resp, testCase.resp)
			}
		})
	}
}

func TestApplyTemplateExecute(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		template string
		value    CommonTemplate
		resp     string
	}{
		{
			name:     "case 0",
			template: DefaultApplyTemplate,
			value:    CommonTemplate{},
			resp: `
## Apply Result



<details><summary>Details (Click me)</summary>

<pre><code>
</code></pre></details>
`,
		},
		{
			name:     "case 1",
			template: DefaultApplyTemplate,
			value: CommonTemplate{
				Result: "result",
				Body:   "body",
			},
			resp: `
## Apply Result


<pre><code>result
</code></pre>


<details><summary>Details (Click me)</summary>

<pre><code>body
</code></pre></details>
`,
		},
		{
			name:     "case 2",
			template: DefaultApplyTemplate,
			value: CommonTemplate{
				Result: "",
				Body:   "body",
			},
			resp: `
## Apply Result



<details><summary>Details (Click me)</summary>

<pre><code>body
</code></pre></details>
`,
		},
		{
			name:     "case 3",
			template: "",
			value: CommonTemplate{
				Result: "",
				Body:   "body",
			},
			resp: `
## Apply Result



<details><summary>Details (Click me)</summary>

<pre><code>body
</code></pre></details>
`,
		},
		{
			name:     "case 4",
			template: "",
			value: CommonTemplate{
				Result: "",
				Body:   `This is a "body".`,
			},
			resp: `
## Apply Result



<details><summary>Details (Click me)</summary>

<pre><code>This is a &#34;body&#34;.
</code></pre></details>
`,
		},
		{
			name:     "case 5",
			template: "",
			value: CommonTemplate{
				Result:       "",
				Body:         `This is a "body".`,
				UseRawOutput: true,
			},
			resp: `
## Apply Result



<details><summary>Details (Click me)</summary>

<pre><code>This is a "body".
</code></pre></details>
`,
		},
		{
			name:     "case 6",
			template: `{{ .Result }}-{{ .Body }}`,
			value: CommonTemplate{
				Result: "c",
				Body:   "d",
			},
			resp: `c-d`,
		},
	}
	for i, testCase := range testCases {
		testCase := testCase
		if testCase.name == "" {
			t.Fatalf("testCase.name is required: index %d", i)
		}
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			template := NewApplyTemplate(testCase.template)
			template.SetValue(testCase.value)
			resp, err := template.Execute()
			if err != nil {
				t.Error(err)
			}
			if resp != testCase.resp {
				t.Errorf("got %s but want %s", resp, testCase.resp)
			}
		})
	}
}
