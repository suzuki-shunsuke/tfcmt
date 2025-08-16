---
sidebar_position: 100
---

# tfcmt

[![GitHub last commit](https://img.shields.io/github/last-commit/suzuki-shunsuke/tfcmt.svg)](https://github.com/suzuki-shunsuke/tfcmt)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/tfcmt/master/LICENSE)

[tfcmt](https://github.com/suzuki-shunsuke/tfcmt) is a fork of [mercari/tfnotify](https://github.com/mercari/tfnotify), enhancing tfnotify in many ways including Terraform >= v0.15 support and advanced formatting options.

tfcmt is a CLI tool to improve the experience of CI of Terraform.
By posting the result of `terraform plan` and `terraform apply` to GitHub Pull Requests as a comment,
you can know the result quickly without browsing the CI web page.

tfcmt enhances tfnotify in many ways, including Terraform >= v0.15 support and advanced formatting options.

[![image](https://user-images.githubusercontent.com/13323303/136236949-bac1a28d-4db2-4a08-900a-708a0a02311c.png)](https://github.com/suzuki-shunsuke/tfcmt/pull/132#issuecomment-936490121)

You can separate the changes outside of Terraform.

![image](https://user-images.githubusercontent.com/13323303/147385656-54cdbef1-a876-49dc-945c-39bcf443ca59.png)

You can exclude the log of refreshing state from the plan result.

![image](https://user-images.githubusercontent.com/13323303/136238225-1569f762-0087-4aae-a513-a63eb9701e05.png)

You can clarify the warning of Terraform.

![image](https://user-images.githubusercontent.com/13323303/136238685-be0bab01-f6cb-4b61-89fa-d94225e50ddb.png)

Combined with [github-comment](https://github.com/suzuki-shunsuke/github-comment), you can hide stale comments.

![image](https://user-images.githubusercontent.com/13323303/136240241-2f2e7455-8a2e-4fce-a91a-c8bab4d73510.png)

Instead of hiding stale comments and creating a new comment, you can update existing comment. This is useful to keep the pull request comments clean.
For detail, please see [Patch `tfcmt plan` comment](plan-patch.md).

![image](https://user-images.githubusercontent.com/13323303/164969354-02bdd49a-547e-4951-9262-033ec5b4db11.png)

--

![image](https://user-images.githubusercontent.com/13323303/164969385-355e801e-3d58-4b75-9657-0bcc10da8d12.png)

## Index

- [Install](install.md)
- [Getting Started](https://github.com/suzuki-shunsuke/tfcmt/tree/main/examples/getting-started)
- [Usage](usage.md)
- [Configuration](config.md)
- [Environment Variable](environment-variable.md)
- [Compared with tfnotify](compared-with-tfnotify.md)
- [Release Notes](https://github.com/suzuki-shunsuke/tfcmt/releases)
- [Blog](#blog)

## What tfcmt does

1. Parse the execution result of Terraform
2. Bind parsed results to Go templates
3. Update pull request labels
4. Post a comment to GitHub

## Blog

- [2021-12-26 tfcmt - Improve Terraform Workflow with PR Comment and Label](https://dev.to/suzukishunsuke/tfcmt-improve-terraform-workflow-with-pr-comment-and-label-1kh7)
- [2021-12-26 tfcmt で Terraform の CI/CD を改善する](https://zenn.dev/shunsuke_suzuki/articles/improve-terraform-cicd-with-tfcmt)

## Who uses tfcmt?

https://github.com/suzuki-shunsuke/tfcmt#who-uses-tfcmt

## License

[LICENSE](https://github.com/suzuki-shunsuke/tfcmt/blob/main/LICENSE)

### License of original code

This is a fork of [mercari/tfnotify](https://github.com/mercari/tfnotify), so about the origincal license, please see https://github.com/mercari/tfnotify#license .

Copyright 2018 Mercari, Inc.

Licensed under the MIT License.

### License of code which we wrote

MIT
