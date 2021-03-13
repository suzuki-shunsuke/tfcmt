# Command Usage

```console
$ tfcmt --version
tfcmt version 0.7.0
```

```console
$ tfcmt help
NAME:
   tfcmt - Notify the execution result of terraform command

USAGE:
   tfcmt [global options] command [command options] [arguments...]

VERSION:
   0.7.0

COMMANDS:
   plan     Run terraform plan and post a comment to GitHub commit or pull request
   apply    Run terraform apply and post a comment to GitHub commit or pull request
   version  Show version
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --ci value         name of CI to run tfcmt
   --owner value      GitHub Repository owner name
   --repo value       GitHub Repository name
   --sha value        commit SHA (revision)
   --build-url value  build url
   --log-level value  log level
   --pr value         pull request number (default: 0)
   --config value     config path
   --var value        template variables. The format of value is '<name>:<value>'
   --help, -h         show help (default: false)
   --version, -v      print the version (default: false)
```

### -var option

tfcmt supports to pass variables by `-var` option.
The format of the value should be `<name>:<value>`.

```console
$ tfcmt -var name:foo plan -- terraform plan
```

The variables can be referred in `template` and `label`.

```yaml
terraform:
  plan:
    template: |
      {{.Vars.name}}
      ...
    when_add_or_update_only:
      label: "{{.Vars.name}}/add-or-update"
```

## tfcmt plan

```console
$ tfcmt help plan
NAME:
   tfcmt plan - Run terraform plan and post a comment to GitHub commit or pull request

USAGE:
   tfcmt plan [arguments...]
```

e.g.

```console
$ tfcmt plan -- terraform plan
```

## tfcmt apply

```console
$ tfcmt help apply
NAME:
   tfcmt apply - Run terraform apply and post a comment to GitHub commit or pull request

USAGE:
   tfcmt apply [arguments...]
```

e.g.

```console
$ tfcmt apply -- terraform apply -auto-approve
```
