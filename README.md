# tfcmt

[![Build Status](https://github.com/suzuki-shunsuke/tfcmt/workflows/test/badge.svg)](https://github.com/suzuki-shunsuke/tfcmt/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/suzuki-shunsuke/tfcmt)](https://goreportcard.com/report/github.com/suzuki-shunsuke/tfcmt)
[![GitHub last commit](https://img.shields.io/github/last-commit/suzuki-shunsuke/tfcmt.svg)](https://github.com/suzuki-shunsuke/tfcmt)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/tfcmt/master/LICENSE)

Fork of [mercari/tfnotify](https://github.com/mercari/tfnotify)

tfcmt parses Terraform commands' execution result and applies it to an arbitrary template and then notifies it to GitHub comments etc.

## Forked version

We forked [suzuki-shunsuke/tfnotify v1.3.3](https://github.com/suzuki-shunsuke/tfnotify/releases/tag/v1.3.3).

## Compared with tfnotify

Please see [Compared with tfnotify](COMPARED_WITH_TFNOTIFY.md).

**We recommend to read [Compared with tfnotify](COMPARED_WITH_TFNOTIFY.md) because there are some features which aren't described at README.**

## Motivation

There are commands such as `plan` and `apply` on Terraform command, but many developers think they would like to check if the execution of those commands succeeded.
Terraform commands are often executed via CI like CircleCI, but in that case you need to go to the CI page to check it.
This is very troublesome. It is very efficient if you can check it with GitHub comments.
You can do this by using this command.

<img src="./misc/images/1.png" width="600">

## Installation

Grab the binary from GitHub Releases (Recommended)

### What tfcmt does

1. Parse the execution result of Terraform
2. Bind parsed results to Go templates
3. Notify it to any platform (e.g. GitHub) as you like

Detailed specifications such as templates and notification destinations can be customized from the configuration files (described later).

## Usage

### Basic

tfcmt is just CLI command. So you can run it from your local after grabbing the binary.

tfcmt accpepts a command as arguments and run the command.

```console
$ tfcmt plan -- terraform plan -detailed-exitcode
```

```console
$ tfcmt apply -- terraform apply -auto-approve
```

### Configurations

When running tfcmt, you can specify the configuration path via `--config` option (if it's omitted, the configuration file `{.,}tfcmt.y{,a}ml` is searched from the current directory to the root directory).

The example settings of GitHub and GitHub Enterprise are as follows. Incidentally, there is no need to replace TOKEN string such as `$GITHUB_TOKEN` with the actual token. Instead, it must be defined as environment variables in CI settings.

[template](https://golang.org/pkg/text/template/) of Go can be used for `template`. The templates can be used in `tfcmt.yaml` are as follows:

Placeholder | Usage
---|---
`{{ .Result }}` | Matched result by parsing like `Plan: 1 to add` or `No changes`
`{{ .Body }}` | The entire of Terraform execution result
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
`{{ .ReplaecedResources }}` | a list of deleted resource paths. This variable can be used at only plan

On GitHub, tfcmt can also put a warning message if the plan result contains resource deletion (optional).

In the template, the [sprig template functions](http://masterminds.github.io/sprig/) can be used.

#### Template Examples

<details>
<summary>For GitHub</summary>

```yaml
---
ci: circleci
notifier:
  github:
    token: $GITHUB_TOKEN
    repository:
      owner: "suzuki-shunsuke"
      name: "tfcmt"
terraform:
  plan:
    template: |
      ## Plan Result <sup>[CI link]( {{ .Link }} )</sup>
      {{if .Result}}
      <pre><code>{{ .Result }}
      </pre></code>
      {{end}}
      <details><summary>Details (Click me)</summary>

      <pre><code>{{ .Body }}
      </pre></code></details>
  apply:
    template: |
      ## Apply Result
      {{if .Result}}
      <pre><code>{{ .Result }}
      </pre></code>
      {{end}}
      <details><summary>Details (Click me)</summary>

      <pre><code>{{ .Body }}
      </pre></code></details>
```

If you would like to let tfcmt warn the resource deletion, add `when_destroy` configuration as below.

```yaml
---
# ...
terraform:
  # ...
  plan:
    template: |
      ## Plan Result <sup>[CI link]( {{ .Link }} )</sup>
      {{if .Result}}
      <pre><code>{{ .Result }}
      </pre></code>
      {{end}}
      <details><summary>Details (Click me)</summary>

      <pre><code>{{ .Body }}
      </pre></code></details>
    when_destroy:
      template: |
        ## :warning: WARNING: Resource Deletion will happen :warning:

        This plan contains **resource deletion**. Please check the plan result very carefully!
  # ...
```

You can also let tfcmt add a label to PRs depending on the `terraform plan` output result. Currently, this feature is for Github labels only.

```yaml
---
# ...
terraform:
  # ...
  plan:
    template: |
      ## Plan Result <sup>[CI link]( {{ .Link }} )</sup>
      {{if .Result}}
      <pre><code>{{ .Result }}
      </pre></code>
      {{end}}
      <details><summary>Details (Click me)</summary>

      <pre><code>{{ .Body }}
      </pre></code></details>
    when_add_or_update_only:
      label: "add-or-update"
      label_color: "1d76db"  # blue
    when_destroy:
      label: "destroy"
      label_color: "d93f0b"  # red
    when_no_changes:
      label: "no-changes"
      label_color: "0e8a16"  # green
    when_plan_error:
      label: "error"
  # ...
```

Sometimes you may want not to HTML-escape Terraform command outputs.
For example, when you use code block to print command output, it's better to use raw characters instead of character references (e.g. `-/+` -> `-/&#43;`, `"` -> `&#34;`).

You can disable HTML escape by adding `use_raw_output: true` configuration.
With this configuration, Terraform doesn't HTML-escape any Terraform output.

~~~yaml
---
# ...
terraform:
  use_raw_output: true
  # ...
  plan:
    template: |
      ## Plan Result <sup>[CI link]( {{ .Link }} )</sup>
      {{if .Result}}
      ```
      {{ .Result }}
      ```
      {{end}}
      <details><summary>Details (Click me)</summary>

      ```
      {{ .Body }}
      ```
  # ...
~~~

</details>

<details>
<summary>For GitHub Enterprise</summary>

```yaml
---
ci: circleci
notifier:
  github:
    token: $GITHUB_TOKEN
    base_url: $GITHUB_BASE_URL # Example: https://github.example.com/api/v3
    repository:
      owner: "suzuki-shunsuke"
      name: "tfcmt"
terraform:
  plan:
    template: |
      ## Plan Result <sup>[CI link]( {{ .Link }} )</sup>
      {{if .Result}}
      <pre><code>{{ .Result }}
      </pre></code>
      {{end}}
      <details><summary>Details (Click me)</summary>

      <pre><code>{{ .Body }}
      </pre></code></details>
  apply:
    template: |
      ## Apply Result
      {{if .Result}}
      <pre><code>{{ .Result }}
      </pre></code>
      {{end}}
      <details><summary>Details (Click me)</summary>

      <pre><code>{{ .Body }}
      </pre></code></details>
```

</details>

## -var option

tfcmt supports to pass variables by `-var` option.
The format of the value should be `<name>:<value>`.

ex.

```
$ terraform plan | tfcmt -var label_prefix:foo -var header:hello plan
```

The variables can be referred in `template` and `label`.

```yaml
terraform:
  plan:
    template: |
      {{.Vars.header}}
      ...
    when_add_or_update_only:
      label: "{{.Vars.label_prefix}}/add-or-update"
```

### Supported CI

Currently, supported CI are here:

- CircleCI
- AWS CodeBuild
- GitHub Actions
- Google Cloud Build

### Private Repository Considerations

GitHub private repositories require the `repo` and `write:discussion` permissions.

### Google Cloud Build Considerations

- These environment variables are needed to be set using [substitutions](https://cloud.google.com/cloud-build/docs/configuring-builds/substitute-variable-values)
  - `COMMIT_SHA`
  - `BUILD_ID`
  - `PROJECT_ID`
  - `_PR_NUMBER`
- Recommended trigger events
  - `terraform plan`: Pull request
  - `terraform apply`: Push to branch

## License

### License of original code

This is a fork of [mercari/tfnotify](https://github.com/mercari/tfnotify), so about the origincal license, please see https://github.com/mercari/tfnotify#license .

Copyright 2018 Mercari, Inc.

Licensed under the MIT License.

### License of code which we wrote

MIT
