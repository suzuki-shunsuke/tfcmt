# Compared with tfnotify

tfcmt isn't compatible with tfnotify.

* Breaking Changes
  * [Don't support platforms which we don't use](#breaking-change-dont-support-platforms-which-we-dont-use)
    * Remove `notifier` option
  * [Remove `fmt` command](#breaking-change-remove-fmt-command)
  * [Configuration file name is changed](#breaking-change-configuration-file-name-is-changed)
  * [Command usage is changed](#breaking-change-command-usage-is-changed)
  * [Don't remove duplicate comments](#breaking-change-dont-remove-duplicate-comments)
  * [Change the behavior of deletion warning](#breaking-change-change-the-behavior-of-deletion-warning)
  * [Post a comment even if it failed to post a comment](#breaking-change-post-a-comment-even-if-it-failed-to-post-a-comment)
* Features
  * [Post a comment when it failed to parse the result](#feature-post-a-comment-when-it-failed-to-parse-the-result)
  * [Find the configuration file recursively](#feature-find-the-configuration-file-recursively)
  * [Complement CI and GitHub Repository owner and name from environment variables](#feature-complement-ci-and-github-repository-owner-and-name-from-environment-variables)
  * [Support to configure label colors](#feature-support-to-configure-label-colors)
  * Support template functions [sprig](http://masterminds.github.io/sprig/)
  * [Support to pass variables by -var option](#feature-support-to-pass-variables-by--var-option)
  * [Add template variables](#feature-add-template-variables)
  * [Don't recreate labels](#feature-dont-recreate-labels)
  * [--version option and `version` command](#feature---version-option-and-version-command)
* Others
  * refactoring
  * update urfave/cli to v2

## Breaking Change: Don't support platforms which we don't use

We support only the following platforms.

* CI
  * CircleCI
  * CodeBuild
  * CloudBuild
  * GitHub Actions
* Notifier
  * GitHub

We don't support the following platforms.

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

## Breaking Change: Remove `fmt` command

Because we don't use this command.
We notify the result of `terraform fmt` with [github-comment](https://github.com/suzuki-shunsuke/github-comment).

## Breaking Change: Configuration file name is changed

Not `{.,}tfnotify.y{,a}ml` but `{.,}tfcmt.y{,a}ml`.

## Breaking Change: Command usage is changed

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

## Breaking Change: Don't remove duplicate comments

tfnotify removes duplicate comments, but this feature isn't documented and confusing.
The link to the comment would be broken when the comment would be removed.

So this feature is removed from tfcmt.

## Breaking Change: Change the behavior of deletion warning

tfnotify posts a deletion warning comment as the other comment.
tfcmt posts only one comment whose template is `when_destroy.template`.

```yaml
    when_destroy:
      label: "destroy"
      label_color: "d93f0b"  # red
      template: |
        {{ .Title }}

        [CI link]( {{ .Link }} )

        This plan contains **resource deletion**. Please check the plan result very carefully!

        {{ .Message }}
        {{if .Result}}
        <pre><code>{{ .Result }}
        </pre></code>
        {{end}}
        <details><summary>Details (Click me)</summary>

        <pre><code>{{ .Body }}
        </pre></code></details>
```

And the default title of destroy warning is changed to `## :warning: Resource Deletion will happen :warning:`.

## Breaking Change: Post a comment even if it failed to post a comment

tfnotify doesn't post a comment when it failed to update labels.
For example, when the label length is too long, tfnotify failed to add the label and the comment isn't posted.
tfnotify exits immediately and the exit code isn't non zero.

On the other hand, tfcmt outputs the error log but the process continues even if it failed to update labels.
And the exit code is same as the result of terraform command.

### Feature: Post a comment when it failed to parse the result

tfnotify doesn't post a comment when it failed to parse the result.
tfcmt posts a comment when it failed to parse the result.

tfcmt supports to configure the template for the parse error.

```yaml
terraform:
  plan:
    when_parse_error:
      template: |
        {{ .Title }} <sup>[CI link]( {{ .Link }} )</sup>

        :warning: It failed to parse the result. :warning:

        {{ .Message }}

        <details><summary>Details (Click me)</summary>

        <pre><code>{{ .CombinedOutput }}
        </pre></code></details>
  apply:
    when_parse_error:
      template: |
        {{ .Title }} <sup>[CI link]( {{ .Link }} )</sup>

        :warning: It failed to parse the result. :warning:

        {{ .Message }}

        <details><summary>Details (Click me)</summary>

        <pre><code>{{ .CombinedOutput }}
        </pre></code></details>
```

## Feature: Find the configuration file recursively

tfcmt searches the configuration file from the current directory to the root directory recursively.

## Feature: Complement CI and GitHub Repository owner and name from environment variables

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

We can omit `ci` and `repository`.

```yaml
notifier:
  github:
    token: $GITHUB_TOKEN
```

## Feature: Support to configure label colors

tfcmt supports to configure label colors.
So we don't have to configure label colors manually.
This feature is useful especially for Monorepo.

## Feature: Support to pass variables by -var option

tfcmt supports to pass variables to template by `-var <name>:<value>` options.
We can access the variable in the template by `{{.Vars.<variable name>}}`.

## Feature: Add template variables

* Stdout: standard output of terraform command
* Stderr: standard error output of terraform command
* CombinedOutput: output of terraform command
* ExitCode: exit code of terraform command
* Vars: variables which are passed by `-var` option

## Feature: Don't recreate labels

If the label which tfnotify set is already set to a pull request, tfnotify removes the label from the pull request and re-adds the same label to the pull request.
This is meaningless.

So tfcmt doesn't recreate a label.

## Feature: --version option and version command

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
