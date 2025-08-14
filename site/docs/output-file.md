---
sidebar_position: 570
---

# Output the result to a local file

[#194](https://github.com/suzuki-shunsuke/tfcmt/issues/194) [#654](https://github.com/suzuki-shunsuke/tfcmt/pull/654) `tfcmt >= v4.2.0`

tfcmt normally posts the result of `terraform plan` and `terraform apply` to GitHub Pull Requests as a comment.
But tfcmt also supports outputting the result to a local file by `--output` option.

tfcmt plan:

```sh
tfcmt --output plan.md plan -- terraform plan
```

tfcmt apply:

```sh
tfcmt --output apply.md apply -- terraform apply
```

If a specified file doesn't exist, the file is created.
If the file already exist, the file content is appended.

:::tip
If you want to overwrite the file content instead of appending, please make [the file empty](https://www.tecmint.com/empty-delete-file-content-linux/) before running tfcmt.

e.g.

```sh
: > plan.md # Make the file empty
tfcmt --output plan.md plan -- terraform plan
```
:::

[Metadata](embedded-metadata.md) isn't embedded.
