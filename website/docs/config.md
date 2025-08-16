---
sidebar_position: 200
---

# Configuration

tfcmt provides the good default configuration and the configuration file is optional,
but you can customize the configuration with a configuration file.

## JSON Schema

[#1551](https://github.com/suzuki-shunsuke/tfcmt/pull/1551)

You can use JSON Schema of tfcmt's configuration file.

- https://github.com/suzuki-shunsuke/tfcmt/blob/main/json-schema/tfcmt.json
- https://raw.githubusercontent.com/suzuki-shunsuke/tfcmt/refs/heads/main/json-schema/tfcmt.json

If you look for a CLI tool to validate configuration with JSON Schema, [ajv-cli](https://ajv.js.org/packages/ajv-cli.html) is useful.

```sh
ajv --spec=draft2020 -s json-schema/tfcmt.json -d tfcmt.yaml
```

### Input Complementation by YAML Language Server

[Please see the comment too.](https://github.com/szksh-lab/.github/issues/67#issuecomment-2564960491)

Version: `main`

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/suzuki-shunsuke/tfcmt/main/json-schema/tfcmt.json
```

Or pinning version:

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/suzuki-shunsuke/tfcmt/v4.14.1/json-schema/tfcmt.json
```

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
`{{ .MovedResources }}` | a list of moved resource paths. This variable can be used at only plan
`{{ .ImportedResources }}` | a list of imported resources (`{"Before": "resource path", "After": "resource path"}`). This variable can be used at only plan

## Template Functions

In the template, the [sprig template functions](http://masterminds.github.io/sprig/) can be used.
And the following functions can be used.

* avoidHTMLEscape
* wrapCode

`avoidHTMLEscape` prevents the text from being HTML escaped.

`wrapCode` wraps a test with <code>\`\`\`</code> or `<pre><code>`.
If the text includes <code>\`\`\`</code>, the text wraps with `<pre><code>`, otherwise the text wraps with <code>\`\`\`</code> and the text isn't HTML escaped.

`wrapCode` omits too long text. ref. https://github.com/suzuki-shunsuke/tfcmt/releases/tag/v3.1.0

## Default Configuration

```yaml
embedded_var_names: []
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
    {{- end}}{{end}}{{if .ImportedResources}}
    * Import
    {{- range .ImportedResources}}
      * {{.}}
    {{- end}}{{end}}{{if .MovedResources}}
    * Move
    {{- range .MovedResources}}
      * {{.Before}} => {{.After}}
    {{- end}}{{end}}
  deletion_warning: |
    {{if .HasDestroy}}
    ### :warning: Resource Deletion will happen :warning:
    This plan contains resource delete operation. Please check the plan result very carefully!
    {{end}}
  changed_result: |
    {{if .ChangedResult}}
    <details><summary>Change Result (Click me)</summary>
    {{wrapCode .ChangedResult}}
    </details>
    {{end}}
  change_outside_terraform: |
    {{if .ChangeOutsideTerraform}}
    <details><summary>:information_source: Objects have changed outside of Terraform</summary>

    _This feature was introduced from [Terraform v0.15.4](https://github.com/hashicorp/terraform/releases/tag/v0.15.4)._
    {{wrapCode .ChangeOutsideTerraform}}
    </details>
    {{end}}
  warning: |
    {{if .Warning}}
    ## :warning: Warnings :warning:
    {{wrapCode .Warning}}
    {{end}}
  error_messages: |
    {{if .ErrorMessages}}
    ## :warning: Errors
    {{range .ErrorMessages}}
    * {{. -}}
    {{- end}}{{end}}
  guide_apply_failure: ""
  guide_apply_parse_error: ""
terraform:
  plan:
    disable_label: false
    ignore_warning: false # tfcmt >= v4.14.0
    template: |
      {{template "plan_title" .}}

      {{if .Link}}[CI link]({{.Link}}){{end}}

      {{template "deletion_warning" .}}
      {{template "result" .}}
      {{template "updated_resources" .}}

      {{template "changed_result" .}}
      {{template "change_outside_terraform" .}}
      {{template "warning" .}}
      {{template "error_messages" .}}
    when_add_or_update_only:
      label: "{{if .Vars.target}}{{.Vars.target}}/{{end}}add-or-update"
      label_color: 1d76db # blue
      # disable_label: false
    when_destroy:
      label: "{{if .Vars.target}}{{.Vars.target}}/{{end}}destroy"
      label_color: d93f0b # red
      # disable_label: false
    when_no_changes:
      label: "{{if .Vars.target}}{{.Vars.target}}/{{end}}no-changes"
      label_color: 0e8a16 # green
      # disable_label: false
      # disable_comment: false
    when_plan_error:
      label:
      label_color:
      # disable_label: false
    when_parse_error:
      template: |
        {{template "plan_title" .}}

        {{if .Link}}[CI link]({{.Link}}){{end}}

        It failed to parse the result.

        <details><summary>Details (Click me)</summary>
        {{wrapCode .CombinedOutput}}
        </details>
  apply:
    template: |
      {{template "apply_title" .}}

      {{if .Link}}[CI link]({{.Link}}){{end}}

      {{if ne .ExitCode 0}}{{template "guide_apply_failure" .}}{{end}}

      {{template "result" .}}

      <details><summary>Details (Click me)</summary>
      {{wrapCode .CombinedOutput}}
      </details>
      {{template "error_messages" .}}
    when_parse_error:
      template: |
        {{template "apply_title" .}}

        {{if .Link}}[CI link]({{.Link}}){{end}}

        {{template "guide_apply_parse_error" .}}

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

Please see [Custom Environment Variable Definition](environment-variable.md#custom-environment-variable-definition).

## Embed metadata in comments

Please see [Embed metadata in comments](embedded-metadata.md).

## Variables `ChangedResult`, `ChangeOutsideTerraform`, and `Warning`

[#103](https://github.com/suzuki-shunsuke/tfcmt/issues/103) [#107](https://github.com/suzuki-shunsuke/tfcmt/pull/107)

The following variables are added from tfcmt v1.1.0.

* ChangedResult
* ChangeOutsideTerraform
* Warning

Compared with `CombinedOutput`, you can make the result of `terraform plan` clear.

* You can exclude noisy `Refreshing state...` logs
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

## GitHub Enterprise Support

Please see [GitHub Enterprise Support](github-enterprise).

## Example configuration

```yaml
---
repo_owner: suzuki-shunsuke
repo_name: tfcmt
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
