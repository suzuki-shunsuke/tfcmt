---
sidebar_position: 580
---

# Use tfcmt with `terragrunt run-all`

[#843](https://github.com/suzuki-shunsuke/tfcmt/discussions/843)

[Terragrunt](https://terragrunt.gruntwork.io/) is a thin wrapper that provides extra tools for keeping your configurations DRY, working with multiple Terraform modules, and managing remote state.

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

2. Run `terragrunt run-all` with `--terragrunt-tfpath`

```sh
terragrunt run-all plan --terragrunt-tfpath "<absolute path of tfwrapper.sh>"
```

Then the result of `terraform plan` is posted per module.
