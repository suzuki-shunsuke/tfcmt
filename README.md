# tfnotify

[![Build Status](https://github.com/suzuki-shunsuke/tfnotify/workflows/test/badge.svg)](https://github.com/suzuki-shunsuke/tfnotify/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/suzuki-shunsuke/tfnotify)](https://goreportcard.com/report/github.com/suzuki-shunsuke/tfnotify)
[![GitHub last commit](https://img.shields.io/github/last-commit/suzuki-shunsuke/tfnotify.svg)](https://github.com/suzuki-shunsuke/tfnotify)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/tfnotify/master/LICENSE)

Fork of [mercari/tfnotify](https://github.com/mercari/tfnotify)

tfnotify parses Terraform commands' execution result and applies it to an arbitrary template and then notifies it to GitHub comments etc.

## Why do we fork mercari/tfnotify?

We have sent [some pull requests](https://github.com/mercari/tfnotify/pulls/suzuki-shunsuke) to mercari/tfnotify but they aren't merged yet.

We forked mercari/tfnotify v0.7.0 [fb178d8](https://github.com/mercari/tfnotify/tree/fb178d8a5a51f88a51b7fda93ed5443ff56dfc8f).

## Motivation

There are commands such as `plan` and `apply` on Terraform command, but many developers think they would like to check if the execution of those commands succeeded.
Terraform commands are often executed via CI like CircleCI, but in that case you need to go to the CI page to check it.
This is very troublesome. It is very efficient if you can check it with GitHub comments or Slack etc.
You can do this by using this command.

<img src="./misc/images/1.png" width="600">

<img src="./misc/images/2.png" width="500">

<img src="./misc/images/3.png" width="600">

## Installation

Grab the binary from GitHub Releases (Recommended)

### What tfnotify does

1. Parse the execution result of Terraform
2. Bind parsed results to Go templates
3. Notify it to any platform (e.g. GitHub) as you like

Detailed specifications such as templates and notification destinations can be customized from the configuration files (described later).

## Usage

### Basic

tfnotify is just CLI command. So you can run it from your local after grabbing the binary.

Basically tfnotify waits for the input from Stdin. So tfnotify needs to pipe the output of Terraform command like the following:

```console
$ terraform plan | tfnotify plan
```

For `plan` command, you also need to specify `plan` as the argument of tfnotify. In the case of `apply`, you need to do `apply`. Currently supported commands can be checked with `tfnotify --help`.

### Configurations

When running tfnotify, you can specify the configuration path via `--config` option (if it's omitted, the configuration file `{.,}tfnotify.y{,a}ml` is searched from the current directory to the root directory).

The example settings of GitHub and GitHub Enterprise, Slack, [Typetalk](https://www.typetalk.com/) are as follows. Incidentally, there is no need to replace TOKEN string such as `$GITHUB_TOKEN` with the actual token. Instead, it must be defined as environment variables in CI settings.

[template](https://golang.org/pkg/text/template/) of Go can be used for `template`. The templates can be used in `tfnotify.yaml` are as follows:

Placeholder | Usage
---|---
`{{ .Title }}` | Like `## Plan result`
`{{ .Message }}` | A string that can be set from CLI with `--message` option
`{{ .Result }}` | Matched result by parsing like `Plan: 1 to add` or `No changes`
`{{ .Body }}` | The entire of Terraform execution result
`{{ .Link }}` | The link of the build page on CI

On GitHub, tfnotify can also put a warning message if the plan result contains resource deletion (optional).

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
      name: "tfnotify"
terraform:
  fmt:
    template: |
      {{ .Title }}

      {{ .Message }}

      {{ .Result }}

      {{ .Body }}
  plan:
    template: |
      {{ .Title }} <sup>[CI link]( {{ .Link }} )</sup>
      {{ .Message }}
      {{if .Result}}
      <pre><code>{{ .Result }}
      </pre></code>
      {{end}}
      <details><summary>Details (Click me)</summary>

      <pre><code>{{ .Body }}
      </pre></code></details>
  apply:
    template: |
      {{ .Title }}
      {{ .Message }}
      {{if .Result}}
      <pre><code>{{ .Result }}
      </pre></code>
      {{end}}
      <details><summary>Details (Click me)</summary>

      <pre><code>{{ .Body }}
      </pre></code></details>
```

If you would like to let tfnotify warn the resource deletion, add `when_destroy` configuration as below.

```yaml
---
# ...
terraform:
  # ...
  plan:
    template: |
      {{ .Title }} <sup>[CI link]( {{ .Link }} )</sup>
      {{ .Message }}
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

You can also let tfnotify add a label to PRs depending on the `terraform plan` output result. Currently, this feature is for Github labels only.

```yaml
---
# ...
terraform:
  # ...
  plan:
    template: |
      {{ .Title }} <sup>[CI link]( {{ .Link }} )</sup>
      {{ .Message }}
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
      {{ .Title }} <sup>[CI link]( {{ .Link }} )</sup>
      {{ .Message }}
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
      name: "tfnotify"
terraform:
  fmt:
    template: |
      {{ .Title }}

      {{ .Message }}

      {{ .Result }}

      {{ .Body }}
  plan:
    template: |
      {{ .Title }} <sup>[CI link]( {{ .Link }} )</sup>
      {{ .Message }}
      {{if .Result}}
      <pre><code>{{ .Result }}
      </pre></code>
      {{end}}
      <details><summary>Details (Click me)</summary>

      <pre><code>{{ .Body }}
      </pre></code></details>
  apply:
    template: |
      {{ .Title }}
      {{ .Message }}
      {{if .Result}}
      <pre><code>{{ .Result }}
      </pre></code>
      {{end}}
      <details><summary>Details (Click me)</summary>

      <pre><code>{{ .Body }}
      </pre></code></details>
```

</details>

<details>
<summary>For GitLab</summary>

```yaml
---
ci: gitlabci
notifier:
  gitlab:
    token: $GITLAB_TOKEN
    base_url: $GITLAB_BASE_URL
    repository:
      owner: "suzuki-shunsuke"
      name: "tfnotify"
terraform:
  fmt:
    template: |
      {{ .Title }}

      {{ .Message }}

      {{ .Result }}

      {{ .Body }}
  plan:
    template: |
      {{ .Title }} <sup>[CI link]( {{ .Link }} )</sup>
      {{ .Message }}
      {{if .Result}}
      <pre><code> {{ .Result }}
      </pre></code>
      {{end}}
      <details><summary>Details (Click me)</summary>
      <pre><code> {{ .Body }}
      </pre></code></details>
  apply:
    template: |
      {{ .Title }}
      {{ .Message }}
      {{if .Result}}
      <pre><code> {{ .Result }}
      </pre></code>
      {{end}}
      <details><summary>Details (Click me)</summary>
      <pre><code> {{ .Body }}
      </pre></code></details>
```
</details>

<details>
<summary>For Slack</summary>

```yaml
---
ci: circleci
notifier:
  slack:
    token: $SLACK_TOKEN
    channel: $SLACK_CHANNEL_ID
    bot: $SLACK_BOT_NAME
terraform:
  plan:
    template: |
      {{ .Message }}
      {{if .Result}}
      ```
      {{ .Result }}
      ```
      {{end}}
      ```
      {{ .Body }}
      ```
```

</details>

<details>
<summary>For Typetalk</summary>

```yaml
---
ci: circleci
notifier:
  typetalk:
    token: $TYPETALK_TOKEN
    topic_id: $TYPETALK_TOPIC_ID
terraform:
  plan:
    template: |
      {{ .Message }}
      {{if .Result}}
      ```
      {{ .Result }}
      ```
      {{end}}
      ```
      {{ .Body }}
      ```
```

</details>

### Supported CI

Currently, supported CI are here:

- CircleCI
- Travis CI
- AWS CodeBuild
- TeamCity
- Drone
- Jenkins
- GitLab CI
- GitHub Actions
- Google Cloud Build

### Private Repository Considerations

GitHub private repositories require the `repo` and `write:discussion` permissions.

### Jenkins Considerations

- Plugin
  - [Git Plugin](https://wiki.jenkins.io/display/JENKINS/Git+Plugin)
- Environment Variable
  - `PULL_REQUEST_NUMBER` or `PULL_REQUEST_URL` are required to set by user for Pull Request Usage

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
