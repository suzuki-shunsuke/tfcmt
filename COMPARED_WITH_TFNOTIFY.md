# Compared with tfnotify

tfcmt isn't compatible with tfnotify.

## Breaking Changes

* don't support platforms which we don't use
  * remove `notifier` option
* don't support `fmt` command
* configuration file name is changed
* command usage is changed
* don't remove duplicated comments

### don't support platforms which we don't use

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

### don't support `fmt` command

Because we don't use this command.
We notify the result of `terraform fmt` with [github-comment](https://github.com/suzuki-shunsuke/github-comment).

### configuration file name is changed

Not `{.,}tfnotify.y{,a}ml` but `{.,}tfcmt.y{,a}ml`.

### command usage is changed

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

### don't remove duplicated comments

tfnotify removes duplicated comments, but this feature isn't documented and confusing.
The link to the comment would be broken when the comment would be removed.

So this feature is removed from tfcmt.

## Features

* don't recreate labels
* support to configure label colors
* support template functions [sprig](http://masterminds.github.io/sprig/)
* support to pass variables by -var option
* support to find the configuration file recursively
* support --version option and `version` command
* support to post a comment when it failed to parse the result
* complement CI and GitHub Repository owner and name from environment variables

### don't recreate labels

If the label which tfnotify set is already set to a pull request, tfnotify removes the label from the pull request and re-adds the same label to the pull request.
This is meaningless.

So tfcmt doesn't recreate a label.

### support to configure label colors

tfcmt supports to configure label colors.
So we don't have to configure label colors manually.
This feature is useful especially for Monorepo.

### support to pass variables by -var option

tfcmt supports to pass variables to template by `-var <name>:<value>` options.
We can access the variable in the template by `{{.Vars.<variable name>}}`.

### support to find the configuration file recursively

tfcmt searches the configuration file from the current directory to the root directory recursively.

### support --version option and version command

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

### support to post a comment when it failed to parse the result

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

### complement CI and GitHub Repository owner and name from environment variables

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

## Others

* refactoring
* update urfave/cli to v2
