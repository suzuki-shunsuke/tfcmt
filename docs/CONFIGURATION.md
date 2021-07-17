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
`{{ .ChangedResult }}` |
`{{ .ChangeOutsideTerraform }}` |
`{{ .Warning }}` |
`{{ .Link }}` | The link of the build page on CI
`{{ .Vars }}` | The variables which are passed by `-var` option
`{{ .Stdout }}` | The standard output of terraform command
`{{ .Stderr }}` | The standard error output of terraform command
`{{ .CombinedOutput }}` | The output of terraform command
`{{ .ExitCode }}` | The exit code of terraform command
`{{ .HasDestroy }}` | Whether there are destroyed resources
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
ci:
  pr: []
  owner: []
  repo: []
  sha: []
  link: []
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

      {{if .Link}}[CI link]({{.Link}}){{end}}

      {{if .HasDestroy}}{{template "deletion_warning" .}}{{end}}
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
      label_color: 1d76db # blue
    when_destroy:
      label: "{{if .Vars.target}}{{.Vars.target}}/{{end}}destroy"
      label_color: d93f0b # red
    when_no_changes:
      label: "{{if .Vars.target}}{{.Vars.target}}/{{end}}no-changes"
      label_color: 0e8a16 # green
    when_plan_error:
      label:
      label_color:
    when_parse_error:
      label:
      label_color:
  apply:
    template: |
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
      {{- end}}{{end}}
    when_parse_error:
      template: |
        ## Apply Result{{if .Vars.target}} ({{.Vars.target}}){{end}}

        {{if .Link}}[CI link]({{.Link}}){{end}}

        It failed to parse the result.

        <details><summary>Details (Click me)</summary>
        {{wrapCode .CombinedOutput}}
        </details>
```

If you don't want to update labels, please set `terraform.plan.disable_label: true`.

```yaml
terraform:
  plan:
    disable_label: true
```

## Custom Environment Variable Definition

Please see [Custom Environment Variable Definition](ENVIRONMENT_VARIABLE.md#custom-environment-variable-definition).

## Embed metadata in comments

Please see [Embed metadata in comments](EMBED_METADATA.md).

## Variables `ChangedResult`, `ChangeOutsideTerraform`, and `Warning`

[#103](https://github.com/suzuki-shunsuke/tfcmt/issues/103) [#107](https://github.com/suzuki-shunsuke/tfcmt/pull/107)

The following variables are added from tfcmt v1.1.0.

* ChangedResult
* ChangeOutsideTerraform
* Warning

Compared with `CombinedOutput`, we can make the result of `terraform plan` clear.

* We can exclude noisy `Refreshing state...` logs
* Separate `changes made outside of Terraform since the last "terraform apply":`
* Make the warning easy to see

### Variable: ChangedResult

```
{{if .ChangedResult}}
<details><summary>Change Result (Click me)</summary>
{{wrapCode .ChangedResult}}
</details>
{{end}}
```

![image](https://user-images.githubusercontent.com/13323303/126020688-dc3c64be-bf01-4ee9-9693-39f85bc67442.png)

### Variable: ChangeOutsideTerraform

```
{{if .ChangeOutsideTerraform}}
<details><summary>:warning: Note: Objects have changed outside of Terraform</summary>
{{wrapCode .ChangeOutsideTerraform}}
</details>
{{end}}
```

![image](https://user-images.githubusercontent.com/13323303/126021350-be037a55-2d83-48a3-a76d-7f9da23fde29.png)

### Variable: Warning

```
{{if .Warning}}
## :warning: Warnings :warning:
{{wrapCode .Warning}}
{{end}}
```

![image](https://user-images.githubusercontent.com/13323303/126020762-68f99375-f860-4c66-964e-6dd0c9578cb1.png)

### Example configuration

```yaml
---
templates:
  changed_result: |
    {{if .ChangedResult}}
    <details><summary>Change Result (Click me)</summary>
    {{wrapCode .ChangedResult}}
    </details>
    {{end}}
  change_outside_terraform: |
    {{if .ChangeOutsideTerraform}}
    <details><summary>:warning: Note: Objects have changed outside of Terraform</summary>
    {{wrapCode .ChangeOutsideTerraform}}
    </details>
    {{end}}
  warning: |
    {{if .Warning}}
    ## :warning: Warnings :warning:
    {{wrapCode .Warning}}
    {{end}}
  error_message: |
    {{if .ErrorMessages}}
    ## :warning: Errors
    {{range .ErrorMessages}}
    * {{. -}}
    {{- end}}{{end}}

terraform:
  plan:
    template: |
      {{template "plan_title" .}}

      {{if .Link}}[CI link]({{.Link}}){{end}}

      {{template "result" .}}
      {{template "updated_resources" .}}

      {{template "changed_result" .}}
      {{template "change_outside_terraform" .}}
      {{template "warning" .}}
      {{template "error_message" .}}
    when_destroy:
      template: |
        {{template "plan_title" .}}

        {{if .Link}}[CI link]({{.Link}}){{end}}

        {{template "deletion_warning" .}}

        {{template "result" .}}
        {{template "updated_resources" .}}

        {{template "changed_result" .}}
        {{template "change_outside_terraform" .}}
        {{template "warning" .}}
        {{template "error_message" .}}
```
