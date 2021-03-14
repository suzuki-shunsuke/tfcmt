# tfcmt

[![Build Status](https://github.com/suzuki-shunsuke/tfcmt/workflows/test/badge.svg)](https://github.com/suzuki-shunsuke/tfcmt/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/suzuki-shunsuke/tfcmt)](https://goreportcard.com/report/github.com/suzuki-shunsuke/tfcmt)
[![GitHub last commit](https://img.shields.io/github/last-commit/suzuki-shunsuke/tfcmt.svg)](https://github.com/suzuki-shunsuke/tfcmt)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/tfcmt/master/LICENSE)

Fork of [mercari/tfnotify](https://github.com/mercari/tfnotify)

tfcmt is a CLI tool to improve the experience of CI of Terraform.
By posting the result of `terraform plan` and `terraform apply` to GitHub Pull Requests as a comment,
we can know the result quickly without browsing the CI web page.

![image](https://user-images.githubusercontent.com/13323303/111016701-b6f89200-83f2-11eb-9fed-35d8249c9ba0.png)

https://github.com/suzuki-shunsuke/tfcmt/pull/70#issuecomment-797854184

## Index

- [Getting Started](examples/getting-started)
- [Usage](docs/USAGE.md)
- [Configuration](docs/CONFIGURATION.md)
- [Environment Variable](docs/ENVIRONMENT_VARIABLE.md)
- [Compared with tfnotify](docs/COMPARED_WITH_TFNOTIFY.md)
- [Release Notes](https://github.com/suzuki-shunsuke/tfcmt/releases)

## Forked version

We forked [suzuki-shunsuke/tfnotify v1.3.3](https://github.com/suzuki-shunsuke/tfnotify/releases/tag/v1.3.3).

## Compared with tfnotify

Please see [Compared with tfnotify](docs/COMPARED_WITH_TFNOTIFY.md).

**We recommend to read this because there are some features which aren't described at README.**

## Install

Grab the binary from [GitHub Releases](https://github.com/suzuki-shunsuke/tfcmt/releases)

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

## Release Notes

Please see [GitHub Releases](https://github.com/suzuki-shunsuke/tfcmt/releases)

## License

### License of original code

This is a fork of [mercari/tfnotify](https://github.com/mercari/tfnotify), so about the origincal license, please see https://github.com/mercari/tfnotify#license .

Copyright 2018 Mercari, Inc.

Licensed under the MIT License.

### License of code which we wrote

MIT
