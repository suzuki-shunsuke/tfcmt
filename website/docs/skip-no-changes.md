---
sidebar_position: 560
---

# Skip posting a comment if there is no change

tfcmt >= [v4.4.0](https://github.com/suzuki-shunsuke/tfcmt/releases/tag/v4.4.0) | [#773](https://github.com/suzuki-shunsuke/tfcmt/discussions/773) [#774](https://github.com/suzuki-shunsuke/tfcmt/pull/774)

You can skip posting a comment if there is no change using the command line option `-skip-no-changes` or configuration field `disable_comment`.

e.g.

```sh
tfcmt plan -skip-no-changes -- terraform plan
```

tfcmt.yaml

```yaml
terraform:
  plan:
    when_no_changes:
      disable_comment: true
```

If the option is set, `tfcmt plan` adds or updates a pull request label but doesn't post a comment if the result of `terraform plan` has no change and no warning.

Even if there are no comment, the pull request label lets you know the result.
This feature is useful when you want to keep pull request comments clean.

`tfcmt apply` accepts the same flag and a separate yaml key. Apply has no result label, so the option only omits the comment. A zero-resource apply is the summary line `Apply complete! Resources: 0 added, 0 changed, 0 destroyed.`

e.g.

```sh
tfcmt apply -skip-no-changes -- terraform apply
```

tfcmt.yaml

```yaml
terraform:
  apply:
    when_no_changes:
      disable_comment: true
```

`terraform.plan.when_no_changes.disable_comment` does not skip apply comments.
