# Compared with tfnotify

tfcmt isn't compatible with tfnotify.

## Breaking Changes

* don't support platforms which we don't use
  * remove `notifier` option
* don't support `fmt` command
* configuration file name is changed
* command usage is changed

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

Not`{.,}tfnotify.y{,a}ml` but `{.,}tfcmt.y{,a}ml`.

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

## Features

* don't recreate labels
* support to configure label colors
* support template functions [sprig](http://masterminds.github.io/sprig/)
* support to pass variables by -var option
* support to find the configuration file recursively
* support --version option and add `version` command

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

### support --version option

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

## Others

* refactoring
* update urfave/cli to v2
