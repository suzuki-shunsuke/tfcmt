package terraform

import (
	"bytes"
	htmltemplate "html/template"
	texttemplate "text/template"

	"github.com/Masterminds/sprig"
)

const (
	// DefaultPlanTitle is a default title for terraform plan
	DefaultPlanTitle = "## Plan result"
	// DefaultDestroyWarningTitle is a default title of destroy warning
	DefaultDestroyWarningTitle = "## WARNING: Resource Deletion will happen"
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
)

// Template is an template interface for parsed terraform execution result
type Template interface {
	Execute() (resp string, err error)
	SetValue(template CommonTemplate)
	GetValue() CommonTemplate
}

// CommonTemplate represents template entities
type CommonTemplate struct {
	Title        string
	Message      string
	Result       string
	Body         string
	Link         string
	UseRawOutput bool
	Vars         map[string]string
}

// DefaultTemplate is a default template for terraform commands
type DefaultTemplate struct {
	Template     string
	defaultTitle string

	CommonTemplate
}

// NewPlanTemplate is PlanTemplate initializer
func NewPlanTemplate(template string) *DefaultTemplate {
	if template == "" {
		template = DefaultPlanTemplate
	}
	return &DefaultTemplate{
		Template:     template,
		defaultTitle: DefaultPlanTitle,
	}
}

// NewDestroyWarningTemplate is DestroyWarningTemplate initializer
func NewDestroyWarningTemplate(template string) *DefaultTemplate {
	if template == "" {
		template = DefaultDestroyWarningTemplate
	}
	return &DefaultTemplate{
		Template:     template,
		defaultTitle: DefaultDestroyWarningTitle,
	}
}

// NewApplyTemplate is ApplyTemplate initializer
func NewApplyTemplate(template string) *DefaultTemplate {
	if template == "" {
		template = DefaultApplyTemplate
	}
	return &DefaultTemplate{
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
func (t *DefaultTemplate) Execute() (string, error) {
	data := map[string]interface{}{
		"Title":   t.Title,
		"Message": t.Message,
		"Result":  t.Result,
		"Body":    t.Body,
		"Link":    t.Link,
		"Vars":    t.Vars,
	}

	resp, err := generateOutput("default", t.Template, data, t.UseRawOutput)
	if err != nil {
		return "", err
	}

	return resp, nil
}

// SetValue sets template entities to CommonTemplate
func (t *DefaultTemplate) SetValue(ct CommonTemplate) {
	if ct.Title == "" {
		ct.Title = t.defaultTitle
	}
	t.CommonTemplate = ct
}

// GetValue gets template entities
func (t *DefaultTemplate) GetValue() CommonTemplate {
	return t.CommonTemplate
}
