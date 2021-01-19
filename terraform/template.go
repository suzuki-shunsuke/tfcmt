package terraform

import (
	"bytes"
	htmltemplate "html/template"
	texttemplate "text/template"

	"github.com/Masterminds/sprig/v3"
)

const (
	// DefaultPlanTemplate is a default template for terraform plan
	DefaultPlanTemplate = `
## Plan Result

{{if .Result}}
<pre><code>{{ .Result }}
</code></pre>
{{end}}

<details><summary>Details (Click me)</summary>

<pre><code>{{ .Body }}
</code></pre></details>
{{if .ErrorMessages}}
## :warning: Errors
{{range .ErrorMessages}}
* {{. -}}
{{- end}}{{end}}`

	// DefaultApplyTemplate is a default template for terraform apply
	DefaultApplyTemplate = `
## Apply Result

{{if .Result}}
<pre><code>{{ .Result }}
</code></pre>
{{end}}

<details><summary>Details (Click me)</summary>

<pre><code>{{ .Body }}
</code></pre></details>
{{if .ErrorMessages}}
## :warning: Errors
{{range .ErrorMessages}}
* {{. -}}
{{- end}}{{end}}`

	// DefaultDestroyWarningTemplate is a default template for terraform plan
	DefaultDestroyWarningTemplate = `
## :warning: Plan Result: Resource Deletion will happen :warning:

This plan contains resource delete operation. Please check the plan result very carefully!

{{if .Result}}
<pre><code>{{ .Result }}
</code></pre>
{{end}}
`

	DefaultPlanParseErrorTemplate = `
## Plan Result

It failed to parse the result.

<details><summary>Details (Click me)</summary>

<pre><code>{{ .CombinedOutput }}
</code></pre></details>
`

	DefaultApplyParseErrorTemplate = `
## Apply Result

It failed to parse the result.

<details><summary>Details (Click me)</summary>

<pre><code>{{ .CombinedOutput }}
</code></pre></details>
`
)

// CommonTemplate represents template entities
type CommonTemplate struct {
	Result            string
	Body              string
	Link              string
	UseRawOutput      bool
	Vars              map[string]string
	Stdout            string
	Stderr            string
	CombinedOutput    string
	ExitCode          int
	ErrorMessages     []string
	CreatedResources  []string
	UpdatedResources  []string
	DeletedResources  []string
	ReplacedResources []string
}

// Template is a default template for terraform commands
type Template struct {
	Template string
	CommonTemplate
}

// NewPlanTemplate is PlanTemplate initializer
func NewPlanTemplate(template string) *Template {
	if template == "" {
		template = DefaultPlanTemplate
	}
	return &Template{
		Template: template,
	}
}

// NewDestroyWarningTemplate is DestroyWarningTemplate initializer
func NewDestroyWarningTemplate(template string) *Template {
	if template == "" {
		template = DefaultDestroyWarningTemplate
	}
	return &Template{
		Template: template,
	}
}

// NewApplyTemplate is ApplyTemplate initializer
func NewApplyTemplate(template string) *Template {
	if template == "" {
		template = DefaultApplyTemplate
	}
	return &Template{
		Template: template,
	}
}

func NewPlanParseErrorTemplate(template string) *Template {
	if template == "" {
		template = DefaultPlanParseErrorTemplate
	}
	return &Template{
		Template: template,
	}
}

func NewApplyParseErrorTemplate(template string) *Template {
	if template == "" {
		template = DefaultApplyParseErrorTemplate
	}
	return &Template{
		Template: template,
	}
}

func generateOutput(kind, template string, data map[string]interface{}, useRawOutput bool) (string, error) {
	var b bytes.Buffer

	if useRawOutput {
		tpl, err := texttemplate.New(kind).Funcs(sprig.TxtFuncMap()).Parse(template)
		if err != nil {
			return "", err
		}
		if err := tpl.Execute(&b, data); err != nil {
			return "", err
		}
	} else {
		tpl, err := htmltemplate.New(kind).Funcs(sprig.FuncMap()).Parse(template)
		if err != nil {
			return "", err
		}
		if err := tpl.Execute(&b, data); err != nil {
			return "", err
		}
	}

	return b.String(), nil
}

// Execute binds the execution result of terraform command into template
func (t *Template) Execute() (string, error) {
	data := map[string]interface{}{
		"Result":            t.Result,
		"Body":              t.Body,
		"Link":              t.Link,
		"Vars":              t.Vars,
		"Stdout":            t.Stdout,
		"Stderr":            t.Stderr,
		"CombinedOutput":    t.CombinedOutput,
		"ExitCode":          t.ExitCode,
		"ErrorMessages":     t.ErrorMessages,
		"CreatedResources":  t.CreatedResources,
		"UpdatedResources":  t.UpdatedResources,
		"DeletedResources":  t.DeletedResources,
		"ReplacedResources": t.ReplacedResources,
	}

	resp, err := generateOutput("default", t.Template, data, t.UseRawOutput)
	if err != nil {
		return "", err
	}

	return resp, nil
}

// SetValue sets template entities to CommonTemplate
func (t *Template) SetValue(ct CommonTemplate) {
	t.CommonTemplate = ct
}
