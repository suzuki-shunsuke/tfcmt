package terraform

import (
	"bytes"
	htmltemplate "html/template"
	texttemplate "text/template"

	"github.com/Masterminds/sprig/v3"
)

const (
	// DefaultPlanTitle is a default title for terraform plan
	DefaultPlanTitle = "## Plan result"
	// DefaultDestroyWarningTitle is a default title of destroy warning
	DefaultDestroyWarningTitle = "## :warning: Resource Deletion will happen :warning:"
	// DefaultApplyTitle is a default title for terraform apply
	DefaultApplyTitle = "## Apply result"

	defaultTemplate = `
{{ .Title }}

{{ .Message }}

{{if .Result}}
<pre><code>{{ .Result }}
</code></pre>
{{end}}

<details><summary>Details (Click me)</summary>

<pre><code>{{ .Body }}
</code></pre></details>
`

	// DefaultPlanTemplate is a default template for terraform plan
	DefaultPlanTemplate = defaultTemplate
	// DefaultApplyTemplate is a default template for terraform apply
	DefaultApplyTemplate = defaultTemplate

	// DefaultDestroyWarningTemplate is a default template for terraform plan
	DefaultDestroyWarningTemplate = `
{{ .Title }}

This plan contains resource delete operation. Please check the plan result very carefully!

{{if .Result}}
<pre><code>{{ .Result }}
</code></pre>
{{end}}
`

	DefaultParseErrorTemplate = `
{{ .Title }}

It failed to parse the result.

{{ .Message }}

<details><summary>Details (Click me)</summary>

<pre><code>{{ .CombinedOutput }}
</code></pre></details>
`
)

// CommonTemplate represents template entities
type CommonTemplate struct {
	Title          string
	Message        string
	Result         string
	Body           string
	Link           string
	UseRawOutput   bool
	Vars           map[string]string
	Stdout         string
	Stderr         string
	CombinedOutput string
	ExitCode       int
}

// Template is a default template for terraform commands
type Template struct {
	Template     string
	defaultTitle string

	CommonTemplate
}

// NewPlanTemplate is PlanTemplate initializer
func NewPlanTemplate(template string) *Template {
	if template == "" {
		template = DefaultPlanTemplate
	}
	return &Template{
		Template:     template,
		defaultTitle: DefaultPlanTitle,
	}
}

// NewDestroyWarningTemplate is DestroyWarningTemplate initializer
func NewDestroyWarningTemplate(template string) *Template {
	if template == "" {
		template = DefaultDestroyWarningTemplate
	}
	return &Template{
		Template:     template,
		defaultTitle: DefaultDestroyWarningTitle,
	}
}

// NewApplyTemplate is ApplyTemplate initializer
func NewApplyTemplate(template string) *Template {
	if template == "" {
		template = DefaultApplyTemplate
	}
	return &Template{
		Template:     template,
		defaultTitle: DefaultApplyTitle,
	}
}

func NewPlanParseErrorTemplate(template string) *Template {
	if template == "" {
		template = DefaultParseErrorTemplate
	}
	return &Template{
		Template:     template,
		defaultTitle: DefaultPlanTitle,
	}
}

func NewApplyParseErrorTemplate(template string) *Template {
	if template == "" {
		template = DefaultParseErrorTemplate
	}
	return &Template{
		Template:     template,
		defaultTitle: DefaultApplyTitle,
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
		"Title":          t.Title,
		"Message":        t.Message,
		"Result":         t.Result,
		"Body":           t.Body,
		"Link":           t.Link,
		"Vars":           t.Vars,
		"Stdout":         t.Stdout,
		"Stderr":         t.Stderr,
		"CombinedOutput": t.CombinedOutput,
		"ExitCode":       t.ExitCode,
	}

	resp, err := generateOutput("default", t.Template, data, t.UseRawOutput)
	if err != nil {
		return "", err
	}

	return resp, nil
}

// SetValue sets template entities to CommonTemplate
func (t *Template) SetValue(ct CommonTemplate) {
	if ct.Title == "" {
		ct.Title = t.defaultTitle
	}
	t.CommonTemplate = ct
}
