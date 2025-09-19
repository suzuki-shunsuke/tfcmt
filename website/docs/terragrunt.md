---
sidebar_position: 575
---

# Terragurnt Support

[Terragrunt](https://terragrunt.gruntwork.io/) is a thin wrapper that provides extra tools for keeping your configurations DRY, working with multiple Terraform modules, and managing remote state.

## Failed to parse the output of Terragrunt

[#1972](https://github.com/suzuki-shunsuke/tfcmt/issues/1972)
[#1541](https://github.com/suzuki-shunsuke/tfcmt/issues/1541)

By default, tfcmt can't parse the output of Terragrunt.
This isn't a bug.
The log of terragrunt includes the prefix `timestamp STDOUT terraform:`.

e.g.

```
09:32:46.963 STDOUT terraform: data.aws_caller_identity.current: Reading...
09:32:46.963 STDOUT terraform: module.cwlogs1.aws_cloudwatch_log_group.this: Refreshing state... [id=test-2]
09:32:47.088 STDOUT terraform:
```

Due to the prefix, tfcmt can't parse the output.
You can suppress the preifx by the environment variable `TERRAGRUNT_LOG_DISABLE=true, then tfcmt can parse the output.

## terragrunt run-all

[#843](https://github.com/suzuki-shunsuke/tfcmt/discussions/843)
[#1923](https://github.com/suzuki-shunsuke/tfcmt/issues/1923)

Terragrunt supports deploying multiple Terraform modules in a single command by [terragrunt run-all command](https://terragrunt.gruntwork.io/docs/features/execute-terraform-commands-on-multiple-modules-at-once/#the-run-all-command), but tfcmt doesn't support the output.

:x: This doesn't work.

```sh
tfcmt plan -- terragrunt run-all plan
```

You can solve the issue by Terragrunt's [--terragrunt-tfpath](https://terragrunt.gruntwork.io/docs/reference/cli-options/#terragrunt-tfpath) option.

1. Create a wrapper script of `terraform` and make it executable

```sh
vi tfwrapper.sh
chmod a+x tfwrapper.sh
```

tfwrapper.sh

```sh
#!/bin/bash

set -euo pipefail

command=$1

base_dir=$(git rev-parse --show-toplevel) # Please fix if necessary
target=${PWD#"$base_dir"/}

if [ "$command" == "plan" ]; then
  tfcmt -var "target:${target}" plan -- terraform "$@"
elif [ "$command" == "apply" ]; then
  tfcmt -var "target:${target}" apply -- terraform "$@"
else
  terraform "$@"
fi
```

2. Run `terragrunt plan --all` with `--tf-path`

Latest Terragrunt (`>= 0.85.0`)

```sh
terragrunt plan --all --tf-path "<absolute path of tfwrapper.sh>"
```

Old Terragrunt: (`< 0.85.0`)

```sh
terragrunt run-all plan --terragrunt-tfpath "<absolute path of tfwrapper.sh>"
```

Then the result of `terraform plan` is posted per module.
