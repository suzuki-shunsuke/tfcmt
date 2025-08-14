---
sidebar_position: 560
---

# Mask sensitive data

[#1083](https://github.com/suzuki-shunsuke/tfcmt/discussions/1083) [#1115](https://github.com/suzuki-shunsuke/tfcmt/pull/1115) `tfcmt >= v4.9.0`

You can mask sensitive data in outputs of terraform.
This feature prevents the leak of sensitive data.

The following outputs are masked.

- Standard output of terraform command
- Standard error output of terraform command
- Pull request comment of `tfcmt plan` and `tfcmt apply`
- [local files created by `--output` option](output-file.md)

:::caution
Even if you maske secrets using this feature, secrets are still stored in Terraform States.
Please see also [Sensitive Data in State](https://developer.hashicorp.com/terraform/language/state/sensitive-data).
:::

You can use environment variables `TFCMT_MASKS` and `TFCMT_MASKS_SEPARATOR`.

- `TFCMT_MASKS`: A list of masks. Masks are joined by `TFCMT_MASKS_SEPARATOR`
- `TFCMT_MASKS_SEPARATOR`: A separator of masks. The default value is `,`

The format of each mask is `${type}:${value}`.
`${type}` must be either `env` or `regexp`.
If `${type}` is `env`, `${value}` is a masked environment variable name.
If `${type}` is `regexp`, `${value}` is a masked regular expression.

e.g. Mask GitHub access tokens and the environment variable `DATADOG_API_KEY`.

```sh
export TFCMT_MASKS='env:GITHUB_TOKEN,env:DATADOG_API_KEY,regexp:ghp_[^ ]+'
tfcmt plan -- terraform plan
```

e.g. Change the separator to `/`.

```sh
export TFCMT_MASKS_SEPARATOR=/
export TFCMT_MASKS='env:GITHUB_TOKEN/env:DATADOG_API_KEY/regexp:ghp_[^ ]+'
```

All matching strings are replaced with `***`.
Replacements are done in order of `TFCMT_MASKS`, so the result depends on the order of `TFCMT_MASKS`.
For example, if `TFCMT_MASKS` is `regexp:foo,regexp:foo.*`, `regexp:foo.*` has no meaning because all `foo` are replaced with `***` before replacing `foo.*` with `***` so `foo.*` doesn't match with anything.

## Example

This example creates a resource [google_cloudbuild_trigger](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloudbuild_trigger).
This resource has a GitHub Access token as a field `substitutions._GH_TOKEN`.

main.tf

```tf
resource "google_cloudbuild_trigger" "filename_trigger" {
  location = "us-central1"

  trigger_template {
    branch_name = "main"
    repo_name   = "my-repo"
  }

  substitutions = {
    _GH_TOKEN = var.gh_token # Secret
  }

  filename = "cloudbuild.yaml"
}

variable "gh_token" {
  type        = string
  description = "GitHub Access token"
}

terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "5.13.0"
    }
  }
}
```

If you run `terraform plan` without masking, the secret would be leaked.
To prevent the leak, let's mask the secret.

```sh
export TFCMT_MASKS=env:TF_VAR_gh_token # Mask the environment variable TF_VAR_gh_token
```

Please see `_GH_TOKEN` in the output of `tfcmt plan` and the pull request comment.
You can confirm `_GH_TOKEN` is masked as `***` properly.

```console
$ tfcmt plan -- terraform plan
tfcmt plan -- terraform plan

Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # google_cloudbuild_trigger.filename_trigger will be created
  + resource "google_cloudbuild_trigger" "filename_trigger" {
      + create_time   = (known after apply)
      + filename      = "cloudbuild.yaml"
      + id            = (known after apply)
      + location      = "us-central1"
      + name          = (known after apply)
      + project       = "hello"
      + substitutions = {
          + "_GH_TOKEN" = "***"
        }
      + trigger_id    = (known after apply)

      + trigger_template {
          + branch_name = "main"
          + project_id  = (known after apply)
          + repo_name   = "my-repo"
        }
    }

Plan: 1 to add, 0 to change, 0 to destroy.

─────────────────────────────────────────────────────────────────────────────

Note: You didn't use the -out option to save this plan, so Terraform can't
guarantee to take exactly these actions if you run "terraform apply" now.
```

![image](https://github.com/suzuki-shunsuke/tfcmt-docs/assets/13323303/7b79481b-923c-40cf-8bbb-f955b0685d1f)

## Terraform sensitive input variables and outputs and sensitive function

Terraform itself has features to prevent sensitive data from being leaked.

- https://developer.hashicorp.com/terraform/tutorials/configuration-language/sensitive-variables
- https://developer.hashicorp.com/terraform/language/functions/sensitive
- https://developer.hashicorp.com/terraform/language/values/outputs#sensitive-suppressing-values-in-cli-output
- https://developer.hashicorp.com/terraform/language/values/variables#suppressing-values-in-cli-output
- https://www.hashicorp.com/blog/terraform-0-14-adds-the-ability-to-redact-sensitive-values-in-console-output
- https://www.hashicorp.com/blog/announcing-hashicorp-terraform-0-15-general-availability

So first you should use these features.
But even if these features are available, it still makes sense for tfcmt to mask sensitive data.
Please imagine the situation that platform engineers manage Terraform workflows and product teams manage Terraform codes in a Monorepo.
Then platform engineers need to prevent sensitive data from being leaked, but if product teams forget to protect them with `sensitive` flags, sensitive data would be leaked.
By protecting sensitive data using tfcmt, platform engineers can prevent sensitive data from being leaked while delegating the management of Terraform codes to product teams.
tfcmt's masking feature works as a guardrail.
