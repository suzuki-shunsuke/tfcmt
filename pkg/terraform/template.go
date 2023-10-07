package terraform

import (
	"bytes"
	htmltemplate "html/template"
	"strings"
	texttemplate "text/template"

	tmpl "github.com/suzuki-shunsuke/tfcmt/v4/pkg/template"
)

const (
	// DefaultPlanTemplate is a default template for terraform plan
	DefaultPlanTemplate = `
{{template "plan_title" .}}

{{if .Link}}[CI link]({{.Link}}){{end}}

{{template "deletion_warning" .}}
{{template "result" .}}
{{template "updated_resources" .}}

{{template "changed_result" .}}
{{template "change_outside_terraform" .}}
{{template "warning" .}}
{{template "error_messages" .}}`

	// DefaultApplyTemplate is a default template for terraform apply
	DefaultApplyTemplate = `
{{template "apply_title" .}}

{{if .Link}}[CI link]({{.Link}}){{end}}

{{if ne .ExitCode 0}}{{template "guide_apply_failure" .}}{{end}}

{{template "result" .}}

<details><summary>Details (Click me)</summary>
{{wrapCode .CombinedOutput}}
</details>
{{template "error_messages" .}}`

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
{{template "apply_title" .}}

{{if .Link}}[CI link]({{.Link}}){{end}}

{{template "guide_apply_parse_error" .}}

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
	HasError               bool
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
	MovedResources         []*MovedResource
	ImportedResources      []string
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
	if len(text) > 60000 { //nolint:gomnd
		text = text[:20000] + `

# ...
# ... The maximum length of GitHub Comment is 65536, so the content is omitted by tfcmt.
# ...

` + text[len(text)-20000:]
	}
	if strings.Contains(text, "```") {
		return `<pre><code>` + text + `</code></pre>`
	}
	return htmltemplate.HTML("\n```hcl\n" + text + "\n```\n") //nolint:gosec
}

func generateOutput(kind, template string, data map[string]interface{}, useRawOutput bool) (string, error) {
	var b bytes.Buffer

	if useRawOutput {
		tpl, err := texttemplate.New(kind).Funcs(texttemplate.FuncMap{
			"avoidHTMLEscape": avoidHTMLEscape,
			"wrapCode":        wrapCode,
		}).Funcs(tmpl.TxtFuncMap()).Parse(template)
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
		}).Funcs(tmpl.FuncMap()).Parse(template)
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
		"HasError":               t.HasError,
		"ErrorMessages":          t.ErrorMessages,
		"CreatedResources":       t.CreatedResources,
		"UpdatedResources":       t.UpdatedResources,
		"DeletedResources":       t.DeletedResources,
		"ReplacedResources":      t.ReplacedResources,
		"MovedResources":         t.MovedResources,
		"ImportedResources":      t.ImportedResources,
		"HasDestroy":             t.HasDestroy,
	}

	templates := map[string]string{
		"plan_title":  "## {{if or (eq .ExitCode 1) .HasError}}:x: Plan Failed{{else}}Plan Result{{end}}{{if .Vars.target}} ({{.Vars.target}}){{end}}",
		"apply_title": "## {{if and (eq .ExitCode 0) (not .HasError)}}:white_check_mark: Apply Succeeded{{else}}:x: Apply Failed{{end}}{{if .Vars.target}} ({{.Vars.target}}){{end}}",
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
{{- end}}{{end}}{{if .ImportedResources}}
* Import
{{- range .ImportedResources}}
  * {{.}}
{{- end}}{{end}}{{if .MovedResources}}
* Move
{{- range .MovedResources}}
  * {{.Before}} => {{.After}}
{{- end}}{{end}}`,
		"deletion_warning": `{{if .HasDestroy}}
### :warning: Resource Deletion will happen :warning:
This plan contains resource delete operation. Please check the plan result very carefully!
{{end}}`,
		"changed_result": `{{if .ChangedResult}}
<details><summary>Change Result (Click me)</summary>
{{wrapCode .ChangedResult}}
</details>
{{end}}`,
		"change_outside_terraform": `{{if .ChangeOutsideTerraform}}
<details><summary>:information_source: Objects have changed outside of Terraform</summary>

_This feature was introduced from [Terraform v0.15.4](https://github.com/hashicorp/terraform/releases/tag/v0.15.4)._
{{wrapCode .ChangeOutsideTerraform}}
</details>
{{end}}`,
		"warning": `{{if .Warning}}
## :warning: Warnings :warning:
{{wrapCode .Warning}}
{{end}}`,
		"error_messages": `{{if .ErrorMessages}}
## :warning: Errors
{{range .ErrorMessages}}
* {{. -}}
{{- end}}{{end}}`,
		"guide_apply_failure":     "",
		"guide_apply_parse_error": "",
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
