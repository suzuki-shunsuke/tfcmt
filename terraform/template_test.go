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

[CI link]()



<details><summary>Details (Click me)</summary>

` + "```" + `

` + "```" + `

</details>
`,
		},
		{
			name:     "case 1",
			template: DefaultPlanTemplate,
			value: CommonTemplate{
				Result:         "result",
				CombinedOutput: "body",
			},
			resp: `
## Plan Result

[CI link]()

<pre><code>result</code></pre>

<details><summary>Details (Click me)</summary>

` + "```" + `
body
` + "```" + `

</details>
`,
		},
		{
			name:     "case 2",
			template: DefaultPlanTemplate,
			value: CommonTemplate{
				Result:         "",
				CombinedOutput: "body",
			},
			resp: `
## Plan Result

[CI link]()



<details><summary>Details (Click me)</summary>

` + "```" + `
body
` + "```" + `

</details>
`,
		},
		{
			name:     "case 3",
			template: DefaultPlanTemplate,
			value: CommonTemplate{
				Result:         "",
				CombinedOutput: `This is a "body".`,
			},
			resp: `
## Plan Result

[CI link]()



<details><summary>Details (Click me)</summary>

` + "```" + `
This is a "body".
` + "```" + `

</details>
`,
		},
		{
			name:     "case 4",
			template: DefaultPlanTemplate,
			value: CommonTemplate{
				Result:         "",
				CombinedOutput: `This is a "body".`,
				UseRawOutput:   true,
			},
			resp: `
## Plan Result

[CI link]()



<details><summary>Details (Click me)</summary>

` + "```" + `
This is a "body".
` + "```" + `

</details>
`,
		},
		{
			name:     "case 5",
			template: "",
			value: CommonTemplate{
				Result:         "",
				CombinedOutput: "body",
			},
			resp: `
## Plan Result

[CI link]()



<details><summary>Details (Click me)</summary>

` + "```" + `
body
` + "```" + `

</details>
`,
		},
		{
			name:     "case 6",
			template: `{{ .Result }}-{{ .CombinedOutput }}`,
			value: CommonTemplate{
				Result:         "c",
				CombinedOutput: "d",
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
## Plan Result

[CI link]()

### :warning: Resource Deletion will happen :warning:
This plan contains resource delete operation. Please check the plan result very carefully!



<details><summary>Details (Click me)</summary>

` + "```" + `

` + "```" + `

</details>
`,
		},
		{
			name:     "case 1",
			template: DefaultDestroyWarningTemplate,
			value: CommonTemplate{
				Result: `This is a "result".`,
			},
			resp: `
## Plan Result

[CI link]()

### :warning: Resource Deletion will happen :warning:
This plan contains resource delete operation. Please check the plan result very carefully!
<pre><code>This is a &#34;result&#34;.</code></pre>


<details><summary>Details (Click me)</summary>

` + "```" + `

` + "```" + `

</details>
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
## Plan Result

[CI link]()

### :warning: Resource Deletion will happen :warning:
This plan contains resource delete operation. Please check the plan result very carefully!
<pre><code>This is a "result".</code></pre>


<details><summary>Details (Click me)</summary>

` + "```" + `

` + "```" + `

</details>
`,
		},
		{
			name:     "case 3",
			template: DefaultDestroyWarningTemplate,
			value: CommonTemplate{
				Result: "",
			},
			resp: `
## Plan Result

[CI link]()

### :warning: Resource Deletion will happen :warning:
This plan contains resource delete operation. Please check the plan result very carefully!



<details><summary>Details (Click me)</summary>

` + "```" + `

` + "```" + `

</details>
`,
		},
		{
			name:     "case 4",
			template: "",
			value: CommonTemplate{
				Result: "",
			},
			resp: `
## Plan Result

[CI link]()

### :warning: Resource Deletion will happen :warning:
This plan contains resource delete operation. Please check the plan result very carefully!



<details><summary>Details (Click me)</summary>

` + "```" + `

` + "```" + `

</details>
`,
		},
		{
			name:     "case 5",
			template: `{{ .Result }}-{{ .CombinedOutput }}`,
			value: CommonTemplate{
				Result:         "c",
				CombinedOutput: "d",
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
## :white_check_mark: Apply Result

[CI link]()



<details><summary>Details (Click me)</summary>

` + "```" + `

` + "```" + `

</details>
`,
		},
		{
			name:     "case 1",
			template: DefaultApplyTemplate,
			value: CommonTemplate{
				Result:         "result",
				CombinedOutput: "body",
			},
			resp: `
## :white_check_mark: Apply Result

[CI link]()

<pre><code>result</code></pre>

<details><summary>Details (Click me)</summary>

` + "```" + `
body
` + "```" + `

</details>
`,
		},
		{
			name:     "case 2",
			template: DefaultApplyTemplate,
			value: CommonTemplate{
				Result:         "",
				CombinedOutput: "body",
			},
			resp: `
## :white_check_mark: Apply Result

[CI link]()



<details><summary>Details (Click me)</summary>

` + "```" + `
body
` + "```" + `

</details>
`,
		},
		{
			name:     "case 3",
			template: "",
			value: CommonTemplate{
				Result:         "",
				CombinedOutput: "body",
			},
			resp: `
## :white_check_mark: Apply Result

[CI link]()



<details><summary>Details (Click me)</summary>

` + "```" + `
body
` + "```" + `

</details>
`,
		},
		{
			name:     "case 4",
			template: "",
			value: CommonTemplate{
				Result:         "",
				CombinedOutput: `This is a "body".`,
			},
			resp: `
## :white_check_mark: Apply Result

[CI link]()



<details><summary>Details (Click me)</summary>

` + "```" + `
This is a "body".
` + "```" + `

</details>
`,
		},
		{
			name:     "case 5",
			template: "",
			value: CommonTemplate{
				Result:         "",
				CombinedOutput: `This is a "body".`,
				UseRawOutput:   true,
			},
			resp: `
## :white_check_mark: Apply Result

[CI link]()



<details><summary>Details (Click me)</summary>

` + "```" + `
This is a "body".
` + "```" + `

</details>
`,
		},
		{
			name:     "case 6",
			template: `{{ .Result }}-{{ .CombinedOutput }}`,
			value: CommonTemplate{
				Result:         "c",
				CombinedOutput: "d",
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
