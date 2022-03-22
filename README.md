# tfcmt

[![Build Status](https://github.com/suzuki-shunsuke/tfcmt/workflows/test/badge.svg)](https://github.com/suzuki-shunsuke/tfcmt/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/suzuki-shunsuke/tfcmt)](https://goreportcard.com/report/github.com/suzuki-shunsuke/tfcmt)
[![GitHub last commit](https://img.shields.io/github/last-commit/suzuki-shunsuke/tfcmt.svg)](https://github.com/suzuki-shunsuke/tfcmt)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/tfcmt/master/LICENSE)

Fork of [mercari/tfnotify](https://github.com/mercari/tfnotify)

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

## Index

- [Getting Started](examples/getting-started)
- [Usage](docs/USAGE.md)
- [Configuration](docs/CONFIGURATION.md)
- [Environment Variable](docs/ENVIRONMENT_VARIABLE.md)
- [Compared with tfnotify](docs/COMPARED_WITH_TFNOTIFY.md)
- [Release Notes](https://github.com/suzuki-shunsuke/tfcmt/releases)
- [Blog](#blog)

## Forked version

We forked [suzuki-shunsuke/tfnotify v1.3.3](https://github.com/suzuki-shunsuke/tfnotify/releases/tag/v1.3.3) ([mercari/tfnotify v0.7.0](https://github.com/mercari/tfnotify/releases/tag/v0.7.0)).

## Compared with tfnotify

Please see [Compared with tfnotify](docs/COMPARED_WITH_TFNOTIFY.md).

**We recommend to read this because there are some features which aren't described at README.**

## Install

Grab the binary from [GitHub Releases](https://github.com/suzuki-shunsuke/tfcmt/releases)

You can install tfcmt with [Homebrew](https://brew.sh/) too.

```console
$ brew install suzuki-shunsuke/tfcmt/tfcmt
```

You can install tfcmt with [aqua](https://aquaproj.github.io/) too.

```console
$ aqua g -i suzuki-shunsuke/tfcmt
```

## What tfcmt does

1. Parse the execution result of Terraform
2. Bind parsed results to Go templates
3. Update pull request labels
4. Post a comment to GitHub

## Getting Started

Please see [Getting Started](examples/getting-started).

## Usage

Please see [Command Usage](docs/USAGE.md).

## Configuration

Please see [Configuration](docs/CONFIGURATION.md).

## Blog

* [2021-12-26 tfcmt - Improve Terraform Workflow with PR Comment and Label](https://dev.to/suzukishunsuke/tfcmt-improve-terraform-workflow-with-pr-comment-and-label-1kh7)

### Japanese

* [2021-12-26 tfcmt で Terraform の CI/CD を改善する](https://zenn.dev/shunsuke_suzuki/articles/improve-terraform-cicd-with-tfcmt)

## Release Notes

Please see [GitHub Releases](https://github.com/suzuki-shunsuke/tfcmt/releases)

## License

### License of original code

This is a fork of [mercari/tfnotify](https://github.com/mercari/tfnotify), so about the origincal license, please see https://github.com/mercari/tfnotify#license .

Copyright 2018 Mercari, Inc.

Licensed under the MIT License.

### License of code which we wrote

MIT
