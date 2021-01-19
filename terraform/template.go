package terraform

import (
	"bytes"
	htmltemplate "html/template"
	"strings"
	texttemplate "text/template"

	"github.com/Masterminds/sprig/v3"
)

const (
	// DefaultPlanTemplate is a default template for terraform plan
	DefaultPlanTemplate = `
## Plan Result

[CI link]({{ .Link }})

{{if .Result}}
<pre><code>{{ .Result }}
</code></pre>
{{end}}
{{if .CreatedResources}}
* Create
{{- range .CreatedResources}}
  * {{.}}
{{- end}}{{end}}{{if .UpdatedResources}}
* Update
{{- range .UpdatedResources}}
  * {{.}}
{{- end}}{{end}}
<details><summary>Details (Click me)</summary>
{{wrapCode .Body}}
</details>
{{if .ErrorMessages}}
## :warning: Errors
{{range .ErrorMessages}}
* {{. -}}
{{- end}}{{end}}`

	// DefaultApplyTemplate is a default template for terraform apply
	DefaultApplyTemplate = `
## Apply Result

[CI link]({{ .Link }})

{{if .Result}}
<pre><code>{{ .Result }}
</code></pre>
{{end}}

<details><summary>Details (Click me)</summary>
{{wrapCode .Body}}
</details>
{{if .ErrorMessages}}
## :warning: Errors
{{range .ErrorMessages}}
* {{. -}}
{{- end}}{{end}}`

	// DefaultDestroyWarningTemplate is a default template for terraform plan
	DefaultDestroyWarningTemplate = `
## :warning: Plan Result: Resource Deletion will happen :warning:

[CI link]({{ .Link }})

This plan contains resource delete operation. Please check the plan result very carefully!

{{if .Result}}
<pre><code>{{ .Result }}
</code></pre>
{{end}}
{{if .CreatedResources}}
* Create
{{- range .CreatedResources}}
  * {{.}}
{{- end}}{{end}}{{if .UpdatedResources}}
* Update
{{- range .UpdatedResources}}
  * {{.}}
{{- end}}{{end}}{{if .DeletedResources}}
* Delete
{{- range .DeletedResources}}
  * {{.}}
{{- end}}{{end}}{{if .ReplacedResources}}
* Replace
{{- range .ReplacedResources}}
  * {{.}}
{{- end}}{{end}}
<details><summary>Details (Click me)</summary>
{{wrapCode .Body}}
</details>
`

	DefaultPlanParseErrorTemplate = `
## Plan Result

[CI link]({{ .Link }})

It failed to parse the result.

<details><summary>Details (Click me)</summary>
{{wrapCode .CombinedOutput}}
</details>
`

	DefaultApplyParseErrorTemplate = `
## Apply Result

[CI link]({{ .Link }})

It failed to parse the result.

<details><summary>Details (Click me)</summary>
{{wrapCode .CombinedOutput}}
</details>
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

func avoidHTMLEscape(text string) htmltemplate.HTML {
	return htmltemplate.HTML(text) //nolint:gosec
}

func wrapCode(text string) interface{} {
	if strings.Contains(text, "```") {
		return `<pre><code>` + text + `</code></pre>`
	}
	return htmltemplate.HTML("\n```\n" + text + "\n```\n") //nolint:gosec
}

func generateOutput(kind, template string, data map[string]interface{}, useRawOutput bool) (string, error) {
	var b bytes.Buffer

	if useRawOutput {
		tpl, err := texttemplate.New(kind).Funcs(texttemplate.FuncMap{
			"avoidHTMLEscape": avoidHTMLEscape,
			"wrapCode":        wrapCode,
		}).Funcs(sprig.TxtFuncMap()).Parse(template)
		if err != nil {
			return "", err
		}
		if err := tpl.Execute(&b, data); err != nil {
			return "", err
		}
	} else {
		tpl, err := htmltemplate.New(kind).Funcs(htmltemplate.FuncMap{
			"avoidHTMLEscape": avoidHTMLEscape,
			"wrapCode":        wrapCode,
		}).Funcs(sprig.FuncMap()).Parse(template)
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
