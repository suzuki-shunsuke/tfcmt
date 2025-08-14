---
sidebar_position: 500
---

# Compared with tfnotify v0.7.0

:::caution
This page isn't fully maintained.
We stopped maintaining this page so there are some differences not described in this page.
:::

tfcmt is a fork of [mercari/tfnotify v0.7.0](https://github.com/mercari/tfnotify/releases/tag/v0.7.0).

tfcmt isn't compatible with tfnotify.

* Breaking Changes
  * [Don't support platforms which we don't use](#breaking-change-dont-support-platforms-which-we-dont-use)
    * Remove `notifier` option
  * [Remove `fmt` command](#breaking-change-remove-fmt-command)
  * [Configuration file name is changed](#breaking-change-configuration-file-name-is-changed)
  * [Configuration file structure is change](#breaking-change-configuration-file-structure-is-changed)
  * [Command usage is changed](#breaking-change-command-usage-is-changed)
  * [template variable Body is removed](#breaking-change-template-variable-body-is-removed)
  * [Remove --message and --destroy-warning-message option and template variable .Message](#breaking-change-remove---message-and---destroy-warning-message-option-and-template-variable-message)
  * [Remove --title and --destroy-warning-title options and template variable .Title](#breaking-change-remove---title-and---destroy-warning-title-options-and-template-variable-title)
  * [Don't remove duplicate comments](#breaking-change-dont-remove-duplicate-comments)
  * [Embed metadata into comment](#breaking-change-embed-metadata-into-comment)
  * [Change the behavior of deletion warning](#breaking-change-change-the-behavior-of-deletion-warning)
  * [Update labels by default](#breaking-change-update-pull-request-labels-by-default)
* Features
  * [Support Terraform >= v0.15](#feature-support-terraform--v015)
  * [Support patching the comment of `tfcmt plan`](#feature-support-patching-the-comment-of-tfcmt-plan)
  * [Add template variables](#feature-add-template-variables)
  * [Post a comment when it failed to parse the result](#feature-post-a-comment-when-it-failed-to-parse-the-result)
  * [Find the configuration file recursively](#feature-find-the-configuration-file-recursively)
  * [Make a configuration file optional](#feature-make-a-configuration-file-optional)
    * [Complement CI and GitHub Repository owner and name from environment variables](#feature-complement-ci-and-github-repository-owner-and-name-from-environment-variables)
    * [Get GitHub Token from the environment variable "GITHUB_TOKEN" by default](#feature-get-github-token-from-the-environment-variable-github_token-by-default)
  * [Support Custom Environment Variable Definition](#feature-custom-environment-variable-definition)
  * [Syntax Highlight](#feature-syntax-highlight)
  * [Support configuring label colors](#feature-support-configuring-label-colors)
  * Support template functions [sprig](http://masterminds.github.io/sprig/)
  * [Support passing variables by -var option](#feature-support-passing-variables-by--var-option)
  * [Add templates configuration](#feature-add-templates-configuration)
  * [Add template functions](#feature-add-template-functions)
  * [Add command-line options about CI](#feature-add-command-line-options-about-ci)
  * [Get pull request number from CI_INFO_PR_NUMBER](#feature-get-pull-request-number-from-ci_info_pr_number)
  * [Add --log-level option and log.level configuration and output structured log with logrus](#feature-add---log-level-option-and-loglevel-configuration-and-output-structured-log-with-logrus)
  * [Don't recreate labels](#feature-dont-recreate-labels)
  * [--version option and `version` command](#feature---version-option-and-version-command)
  * [Omit too long text in the template function `wrapCode` automatically](https://github.com/suzuki-shunsuke/tfcmt/releases/tag/v3.1.0)
* Fixes
  * [Post a comment even if it failed to update labels](#fix-post-a-comment-even-if-it-failed-to-update-labels)
* Others
  * refactoring
  * update urfave/cli to v2

## Breaking Change: Don't support platforms which we don't use

[#4](https://github.com/suzuki-shunsuke/tfcmt/pull/4)

tfcmt supports only the following platforms.

* CI
  * CircleCI
  * CodeBuild
  * CloudBuild
  * GitHub Actions
* Notifier
  * GitHub

tfcmt doesn't support the following platforms.

* CI
  * Jenkins
  * Travis
  * GitLab
  * Drone
  * TeamCity
* Notification
  * Slack
  * TypeTalk
  * GitLab

Because we don't use these platforms and it is hard to maintain them.
By removing them, the code makes simple.

:::info
For GitLab Users, please check [hirosassa/tfcmt-gitlab](https://github.com/hirosassa/tfcmt-gitlab), which is a fork of tfcmt for GitLab.
:::

## Breaking Change: Remove `fmt` command

[#5](https://github.com/suzuki-shunsuke/tfcmt/pull/5)

Because we don't use this command.
We notify the result of `terraform fmt` with [github-comment](https://github.com/suzuki-shunsuke/github-comment).

## Breaking Change: Configuration file name is changed

[#6](https://github.com/suzuki-shunsuke/tfcmt/pull/6)

Not `{.,}tfnotify.y{,a}ml` but `{.,}tfcmt.y{,a}ml`.

## Breaking Change: Configuration file structure is changed

[#79](https://github.com/suzuki-shunsuke/tfcmt/pull/79) [#80](https://github.com/suzuki-shunsuke/tfcmt/pull/80)

* `notifier` was removed
* structure of `ci` was changed

Please see [Configuration](config.md) too.

## Breaking Change: Command usage is changed

[#7](https://github.com/suzuki-shunsuke/tfcmt/pull/7)

AS IS

```
terraform plan | tfnotify plan
terraform apply | tfnotify apply
```

TO BE

```
tfcmt plan -- terraform plan
tfcmt apply -- terraform apply
```

By this change, tfcmt can handle the standard error output and exit code of the terraform command.

## Breaking Change: template variable Body is removed

template variable `Body` is removed. Replace `Body` to `CombinedOutput`.
`CombinedOutput` includes both standard output and standard error output.

## Breaking Change: Remove --message and --destroy-warning-message option and template variable .Message

[#40](https://github.com/suzuki-shunsuke/tfcmt/pull/40)

We introduced more general option `-var` and template variable `.Vars`,
so the `--message` and `--destroy-warning-message` options aren't needed.

## Breaking Change: Remove --title and --destroy-warning-title options and template variable .Title

[#41](https://github.com/suzuki-shunsuke/tfcmt/pull/41)

We introduced more general option `-var` and template variable `.Vars`,
so the `--title` and `--destroy-warning-title` options aren't needed.

## Breaking Change: Don't remove duplicate comments

[#14](https://github.com/suzuki-shunsuke/tfcmt/pull/14)

tfnotify removes duplicate comments, but this feature isn't documented and confusing.
The link to the comment would be broken when the comment would be removed.

So this feature is removed from tfcmt.

## Breaking Change: Embed metadata into comment

[#67](https://github.com/suzuki-shunsuke/tfcmt/pull/67)

Instead of removing duplicate comments, tfcmt embeds metadata into comment with [github-comment-metadata](https://github.com/suzuki-shunsuke/github-comment-metadata).
tfcmt itself doesn't support hiding old comments, but you can hide comments with [github-comment's hide command](https://github.com/suzuki-shunsuke/github-comment#hide).

## Breaking Change: Change the behavior of deletion warning

[#32](https://github.com/suzuki-shunsuke/tfcmt/pull/32)

tfnotify posts a deletion warning comment as the other comment.
tfcmt posts only one comment whose template is `when_destroy.template`.

```yaml
    when_destroy:
      label: "destroy"
      label_color: "d93f0b"  # red
      template: |
        ## Plan Result

        [CI link]( {{ .Link }} )

        This plan contains **resource deletion**. Please check the plan result very carefully!

        {{if .Result}}
        <pre><code>{{ .Result }}
        </pre></code>
        {{end}}
        <details><summary>Details (Click me)</summary>

        <pre><code>{{ .CombinedOutput }}
        </pre></code></details>
```

## Breaking Change: Update pull request labels by default

[#44](https://github.com/suzuki-shunsuke/tfcmt/pull/44)

tfcmt updates pull request labels by default using default label name and color.

* no-changes, green
* add-or-update, blue
* destroy, red

If you don't want to update labels, please configure `disable_label: true`.

```yaml
terraform:
  plan:
    disable_label: true
```

## Feature: Support Terraform >= v0.15

[#90](https://github.com/suzuki-shunsuke/tfcmt/pull/90) [#91](https://github.com/suzuki-shunsuke/tfcmt/pull/91)

From Terraform v0.15, the message which terraform plan has no change was changed.

AS IS

```
No changes. Infrastructure is up-to-date.
```

TO BE

```
No changes. Your infrastructure matches the configuration.
```

tfcmt supports both messages.

## Feature: Support patching the comment of `tfcmt plan`

Please see [Patch `tfcmt plan` comment](plan-patch.md).

## Feature: Add template variables

* `Stdout`: standard output of terraform command
* `Stderr`: standard error output of terraform command
* `CombinedOutput`: output of terraform command
* `ExitCode`: exit code of terraform command
* `Vars`: variables which are passed by `-var` option
* `ErrorMessages`: a list of error messages which occur in tfcmt
* [#39](https://github.com/suzuki-shunsuke/tfcmt/pull/39) `CreatedResources`: a list of created resource paths ([]string)
* [#39](https://github.com/suzuki-shunsuke/tfcmt/pull/39) `UpdatedResources`: a list of updated resource paths ([]string)
* [#39](https://github.com/suzuki-shunsuke/tfcmt/pull/39) `DeletedResources`: a list of deleted resource paths ([]string)
* [#39](https://github.com/suzuki-shunsuke/tfcmt/pull/39) `ReplacedResources`: a list of replaced resource paths ([]string)
* [#103](https://github.com/suzuki-shunsuke/tfcmt/pull/103) [#107](https://github.com/suzuki-shunsuke/tfcmt/pull/107) `ChangedResult`
* [#103](https://github.com/suzuki-shunsuke/tfcmt/pull/103) [#107](https://github.com/suzuki-shunsuke/tfcmt/pull/107) `ChangeOutsideTerraform`
* [#103](https://github.com/suzuki-shunsuke/tfcmt/pull/103) [#107](https://github.com/suzuki-shunsuke/tfcmt/pull/107) `Warning`

### Feature: Add template variables of changed resource paths

[#39](https://github.com/suzuki-shunsuke/tfcmt/pull/39)

As a summary of the result of `terraform plan`, it is convenient to show a list of resource paths.
So the following template variables are added.

* `CreatedResources`: a list of created resource paths ([]string)
* `UpdatedResources`: a list of updated resource paths ([]string)
* `DeletedResources`: a list of deleted resource paths ([]string)
* `ReplacedResources`: a list of replaced resource paths ([]string)

For example,

```
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
```

```
* Create
  * null_resource.foo
* Update
  * mysql_database.bar
* Delete
  * null_resource.bar
* Replace
  * mysql_database.foo
```

## Feature: Post a comment when it failed to parse the result

[#21](https://github.com/suzuki-shunsuke/tfcmt/pull/21)

tfnotify doesn't post a comment when it failed to parse the result.
tfcmt posts a comment when it failed to parse the result.

tfcmt supports configuring the template for the parse error.

```yaml
terraform:
  plan:
    when_parse_error:
      template: |
        ## Plan Result <sup>[CI link]( {{ .Link }} )</sup>

        :warning: It failed to parse the result. :warning:

        <details><summary>Details (Click me)</summary>

        <pre><code>{{ .CombinedOutput }}
        </pre></code></details>
  apply:
    when_parse_error:
      template: |
        ## Apply Result <sup>[CI link]( {{ .Link }} )</sup>

        :warning: It failed to parse the result. :warning:

        <details><summary>Details (Click me)</summary>

        <pre><code>{{ .CombinedOutput }}
        </pre></code></details>
```

## Feature: Find the configuration file recursively

[suzuki-shunsuke/tfnotify#19](https://github.com/suzuki-shunsuke/tfnotify/pull/19)

tfcmt searches the configuration file from the current directory to the root directory recursively.

## Feature: Make a configuration file optional

[#45](https://github.com/suzuki-shunsuke/tfcmt/pull/45)

When a configuration file path isn't specified and it isn't found, tfcmt works with the default configuration.

### Feature: Complement CI and GitHub Repository owner and name from environment variables

[#25](https://github.com/suzuki-shunsuke/tfcmt/pull/25)

tfcmt complement the configuration CI and GitHub Repository owner and name from CI builtin environment variables.
tfcmt uses [suzuki-shunsuke/go-ci-env](https://github.com/suzuki-shunsuke/go-ci-env) for this feature.
So currently, this feature doesn't support Google CloudBuild for now.

AS IS

```yaml
ci: circleci
notifier:
  github:
    token: $GITHUB_TOKEN
    repository:
      owner: suzuki-shunsuke
      name: tfcmt
```

You can omit `ci` and `repository`.

```yaml
notifier:
  github:
    token: $GITHUB_TOKEN
```

## Feature: Get GitHub Token from the environment variable "GITHUB_TOKEN" by default

You can omit the configuration `notifier.github.token`.

## Feature: Custom Environment Variable Definition

You can complement the parameters like `pr` and `repo` on the other platform like Travis CI and Jenkins with Custom Environment Variable Definition.

Please see [here](environment-variable.md#custom-environment-variable-definition).

## Feature: Syntax Highlight

[#146](https://github.com/suzuki-shunsuke/tfcmt/pull/146)

Use HCL Syntax Hightlit. Please see [Release Note](https://github.com/suzuki-shunsuke/tfcmt/releases/tag/v2.1.0)

## Feature: Support configuring label colors

[98547135a6d37b11b641feb399eec17721fe0bc0](https://github.com/suzuki-shunsuke/tfnotify/commit/98547135a6d37b11b641feb399eec17721fe0bc0)
[49ea5c3a8c01e53cac6d3b529bd5d9907c41e9d3](https://github.com/suzuki-shunsuke/tfnotify/commit/49ea5c3a8c01e53cac6d3b529bd5d9907c41e9d3)

tfcmt supports configuring label colors.
So you don't have to configure label colors manually.
This feature is useful especially for Monorepo.

## Feature: Support passing variables by -var option

[suzuki-shunsuke/tfnotify#29](https://github.com/suzuki-shunsuke/tfnotify/pull/29)

tfcmt supports passing variables to template by `-var <name>:<value>` options.
You can access the variable in the template by `{{.Vars.<variable name>}}`.

The variable `target` has a special meaning.
This variable is used at the default template and default label name.
This is useful for Monorepo. By setting `target`, you can distinguish the comment and label of each service.
When this variable isn't set, this is just ignored.

## Feature: Add templates configuration and builtin templates

[#50](https://github.com/suzuki-shunsuke/tfcmt/issues/50) [#51](https://github.com/suzuki-shunsuke/tfcmt/pull/51)

e.g.

```yaml
templates:
  title: "## Plan Result ({{.Vars.target}})"
terraform:
  plan:
    template: |
      {{template "title" .}}
```

The following builtin templates are defined. You can override them.

* plan_title
* apply_title
* result
* updated_resources
* deletion_warning

## Feature: Add template functions

[#42](https://github.com/suzuki-shunsuke/tfcmt/pull/42)

* avoidHTMLEscape
* wrapCode

`avoidHTMLEscape` prevents the text from being HTML escaped.

`wrapCode` wraps a test with <code>\`\`\`</code> or `<pre><code>`.
If the text includes <code>\`\`\`</code>, the text wraps with `<pre><code>`, otherwise the text wraps with <code>\`\`\`</code> and the text isn't HTML escaped.

## Feature: Add command-line options about CI

* -owner
* -repo
* -pr
* -sha
* -build-url

mercari/tfnotify gets these parameters from only environment variables, so you can't use mercari/tfnotify on the platform which mercari/tfnotify doesn't support.
On the other hand, tfcmt supports specifying these parameters by command-line options, so you can use tfcmt anywhere.

e.g.

```console
$ tfcmt -owner suzuki-shunsuke -repo tfcmt -pr 3 -- terraform plan
```

## Feature: Get pull request number from commit hash

[v3.2.2](https://github.com/suzuki-shunsuke/tfcmt/releases/tag/v3.2.2) [#288](https://github.com/suzuki-shunsuke/tfcmt/issues/288)

## Feature: Create comments to pull request instead of commit

[v3.4.0](https://github.com/suzuki-shunsuke/tfcmt/releases/tag/v3.4.0) [#387](https://github.com/suzuki-shunsuke/tfcmt/issues/387) [#390](https://github.com/suzuki-shunsuke/tfcmt/issues/390)

## Feature: Get pull request number from CI_INFO_PR_NUMBER

[ci-info](https://github.com/suzuki-shunsuke/ci-info) is a CLI tool to get CI related information, and the environment variable `CI_INFO_PR_NUMBER` is set via ci-info by default.
If the pull request number can't bet gotten but `CI_INFO_PR_NUMBER` is being set, `CI_INFO_PR_NUMBER` is used.

## Feature: Add --log-level option and log.level configuration and output structured log with logrus

[#59](https://github.com/suzuki-shunsuke/tfcmt/pull/59)

e.g.

```
$ tfcmt --log-level debug plan -- terraform plan
```

```yaml
---
log:
  level: debug
```

## Feature: Don't recreate labels

[suzuki-shunsuke/tfnotify#32](https://github.com/suzuki-shunsuke/tfnotify/pull/32)

If the label which tfnotify set is already set to a pull request, tfnotify removes the label from the pull request and re-adds the same label to the pull request.
This is meaningless.

So tfcmt doesn't recreate a label.

## Feature: --version option and version command

[suzuki-shunsuke/tfnotify#4](https://github.com/suzuki-shunsuke/tfnotify/pull/4)
[#11](https://github.com/suzuki-shunsuke/tfcmt/pull/11)

AS IS

```
$ tfnotify --version
tfnotify version unset
```

TO BE

```
$ tfcmt --version
tfcmt version 0.1.0

$ tfcmt version
tfcmt version 0.1.0
```

## Fix: Post a comment even if it failed to update labels

[#35](https://github.com/suzuki-shunsuke/tfcmt/pull/35)

tfnotify doesn't post a comment when it failed to update labels.
For example, when the label length is too long, tfnotify failed to add the label and the comment isn't posted.

On the other hand, tfcmt outputs the error log but the process continues even if it failed to update labels.
