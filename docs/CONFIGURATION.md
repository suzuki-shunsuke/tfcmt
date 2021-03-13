# Configuration

tfcmt provides the good default configuration and the configuration file is optional,
but we can customize the configuration with a configuration file.

## Configuration File Path

When running tfcmt, you can specify the configuration path via `--config` option (if it's omitted, the configuration file `{.,}tfcmt.y{,a}ml` is searched from the current directory to the root directory).
Configuration file is optional. If `--config` option isn't used and the configuration file isn't found, tfcmt works with the default configuration.

## Template Engine

The template is rendered with Go's [template](https://golang.org/pkg/html/template/).

## Template variables

Placeholder | Usage
---|---
`{{ .Result }}` | Matched result by parsing like `Plan: 1 to add` or `No changes`
`{{ .Link }}` | The link of the build page on CI
`{{ .Vars }}` | The variables which are passed by `-var` option
`{{ .Stdout }}` | The standard output of terraform command
`{{ .Stderr }}` | The standard error output of terraform command
`{{ .CombinedOutput }}` | The output of terraform command
`{{ .ExitCode }}` | The exit code of terraform command
`{{ .ErrorMessages }}` | a list of error messages which occur in tfcmt
`{{ .CreatedResources }}` | a list of created resource paths. This variable can be used at only plan
`{{ .UpdatedResources }}` | a list of updated resource paths. This variable can be used at only plan
`{{ .DeletedResources }}` | a list of deleted resource paths. This variable can be used at only plan
`{{ .ReplacedResources }}` | a list of deleted resource paths. This variable can be used at only plan

## Template Functions

In the template, the [sprig template functions](http://masterminds.github.io/sprig/) can be used.
And the following functions can be used.

* avoidHTMLEscape
* wrapCode

`avoidHTMLEscape` prevents the text from being HTML escaped.

`wrapCode` wraps a test with <code>\`\`\`</code> or `<pre><code>`.
If the text includes <code>\`\`\`</code>, the text wraps with `<pre><code>`, otherwise the text wraps with <code>\`\`\`</code> and the text isn't HTML escaped.

## Default Configuration

```yaml
vars: {}
templates:
  plan_title: "## {{if eq .ExitCode 1}}:x: {{end}}Plan Result{{if .Vars.target}} ({{.Vars.target}}){{end}}"
  apply_title: "## :{{if eq .ExitCode 0}}white_check_mark{{else}}x{{end}}: Apply Result{{if .Vars.target}} ({{.Vars.target}}){{end}}"

  result: "{{if .Result}}<pre><code>{{ .Result }}</code></pre>{{end}}"
  updated_resources: |
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
  deletion_warning: |
    ### :warning: Resource Deletion will happen :warning:
    This plan contains resource delete operation. Please check the plan result very carefully!
terraform:
  plan:
    disable_label: false
    template: |
      {{template "plan_title" .}}

      [CI link]({{ .Link }})

      {{template "result" .}}
      {{template "updated_resources" .}}
      <details><summary>Details (Click me)</summary>
      {{wrapCode .CombinedOutput}}
      </details>
      {{if .ErrorMessages}}
      ## :warning: Errors
      {{range .ErrorMessages}}
      * {{. -}}
      {{- end}}{{end}}
    when_add_or_update_only:
      label: "{{if .Vars.target}}{{.Vars.target}}/{{end}}add-or-update"
      color: 1d76db # blue
    when_destroy:
      label: "{{if .Vars.target}}{{.Vars.target}}/{{end}}destroy"
      color: d93f0b # red
      template: |
        {{template "plan_title" .}}

        [CI link]({{ .Link }})

        {{template "deletion_warning" .}}
        {{template "result" .}}

        {{template "updated_resources" .}}
        <details><summary>Details (Click me)</summary>
        {{wrapCode .CombinedOutput}}
        </details>
    when_no_changes:
      label: "{{if .Vars.target}}{{.Vars.target}}/{{end}}no-changes"
      color: 0e8a16 # green
    when_plan_error:
      label:
      color:
    when_parse_error:
      label:
      color:
  apply:
    template: |
      {{template "apply_title" .}}

      [CI link]({{ .Link }})

      {{template "result" .}}

      <details><summary>Details (Click me)</summary>
      {{wrapCode .CombinedOutput}}
      </details>
      {{if .ErrorMessages}}
      ## :warning: Errors
      {{range .ErrorMessages}}
      * {{. -}}
      {{- end}}{{end}}
    when_parse_error:
      template: |
        ## Apply Result{{if .Vars.target}} ({{.Vars.target}}){{end}}

        [CI link]({{ .Link }})

        It failed to parse the result.

        <details><summary>Details (Click me)</summary>
        {{wrapCode .CombinedOutput}}
        </details>
```

If the plan contains resource deletion, the template of `when_destroy` is used.

If you don't want to update labels, please set `terraform.plan.disable_label: true`.

```yaml
terraform:
  plan:
    disable_label: true
```
