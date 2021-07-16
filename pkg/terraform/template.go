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
{{template "plan_title" .}}

{{if .Link}}[CI link]({{.Link}}){{end}}

{{template "result" .}}
{{template "updated_resources" .}}
<details><summary>Details (Click me)</summary>
{{wrapCode .CombinedOutput}}
</details>
{{if .ErrorMessages}}
## :warning: Errors
{{range .ErrorMessages}}
* {{. -}}
{{- end}}{{end}}`

	// DefaultApplyTemplate is a default template for terraform apply
	DefaultApplyTemplate = `
{{template "apply_title" .}}

{{if .Link}}[CI link]({{.Link}}){{end}}

{{template "result" .}}

<details><summary>Details (Click me)</summary>
{{wrapCode .CombinedOutput}}
</details>
{{if .ErrorMessages}}
## :warning: Errors
{{range .ErrorMessages}}
* {{. -}}
{{- end}}{{end}}`

	// DefaultPlanParseErrorTemplate is a default template for terraform plan parse error
	DefaultPlanParseErrorTemplate = `
{{template "plan_title" .}}

{{if .Link}}[CI link]({{.Link}}){{end}}

It failed to parse the result.

<details><summary>Details (Click me)</summary>
{{wrapCode .CombinedOutput}}
</details>
`

	// DefaultApplyParseErrorTemplate  is a default template for terraform apply parse error
	DefaultApplyParseErrorTemplate = `
## Apply Result{{if .Vars.target}} ({{.Vars.target}}){{end}}

{{if .Link}}[CI link]({{.Link}}){{end}}

It failed to parse the result.

<details><summary>Details (Click me)</summary>
{{wrapCode .CombinedOutput}}
</details>
`
)

// CommonTemplate represents template entities
type CommonTemplate struct {
	Result                 string
	ChangedResult          string
	ChangeOutsideTerraform string
	Warning                string
	Link                   string
	UseRawOutput           bool
	HasDestroy             bool
	Vars                   map[string]string
	Templates              map[string]string
	Stdout                 string
	Stderr                 string
	CombinedOutput         string
	ExitCode               int
	ErrorMessages          []string
	CreatedResources       []string
	UpdatedResources       []string
	DeletedResources       []string
	ReplacedResources      []string
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
		"Result":                 t.Result,
		"ChangedResult":          t.ChangedResult,
		"ChangeOutsideTerraform": t.ChangeOutsideTerraform,
		"Warning":                t.Warning,
		"Link":                   t.Link,
		"Vars":                   t.Vars,
		"Stdout":                 t.Stdout,
		"Stderr":                 t.Stderr,
		"CombinedOutput":         t.CombinedOutput,
		"ExitCode":               t.ExitCode,
		"ErrorMessages":          t.ErrorMessages,
		"CreatedResources":       t.CreatedResources,
		"UpdatedResources":       t.UpdatedResources,
		"DeletedResources":       t.DeletedResources,
		"ReplacedResources":      t.ReplacedResources,
		"HasDestroy":             t.HasDestroy,
	}

	templates := map[string]string{
		"plan_title":  "## {{if eq .ExitCode 1}}:x: {{end}}Plan Result{{if .Vars.target}} ({{.Vars.target}}){{end}}",
		"apply_title": "## :{{if eq .ExitCode 0}}white_check_mark{{else}}x{{end}}: Apply Result{{if .Vars.target}} ({{.Vars.target}}){{end}}",
		"result":      "{{if .Result}}<pre><code>{{ .Result }}</code></pre>{{end}}",
		"updated_resources": `{{if .CreatedResources}}
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
{{- end}}{{end}}`,
		"deletion_warning": `### :warning: Resource Deletion will happen :warning:
This plan contains resource delete operation. Please check the plan result very carefully!`,
	}

	for k, v := range t.Templates {
		templates[k] = v
	}

	resp, err := generateOutput("default", addTemplates(t.Template, templates), data, t.UseRawOutput)
	if err != nil {
		return "", err
	}

	return resp, nil
}

// SetValue sets template entities to CommonTemplate
func (t *Template) SetValue(ct CommonTemplate) {
	t.CommonTemplate = ct
}

func addTemplates(tpl string, templates map[string]string) string {
	for k, v := range templates {
		tpl += `{{define "` + k + `"}}` + v + "{{end}}"
	}
	return tpl
}
