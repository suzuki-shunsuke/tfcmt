---
sidebar_position: 300
---

# Environment variable

- TFCMT_GITHUB_TOKEN (tfcmt >= [v4.8.0](https://github.com/suzuki-shunsuke/tfcmt/releases/tag/v4.8.0)), GITHUB_TOKEN
- TFCMT_REPO_OWNER
- TFCMT_REPO_NAME
- TFCMT_SHA
- TFCMT_PR_NUMBER
- TFCMT_CONFIG
- TFCMT_DISABLE_LABEL (tfcmt >= v4.10.0)
- [Native support of some CI platforms](#native-support-of-some-ci-platforms)
- [Custom Environment Variable Definition](#custom-environment-variable-definition)

## Native support of some CI platforms

Currently, supported CI are here:

- CircleCI
- Drone
- AWS CodeBuild
- GitHub Actions
- Google Cloud Build

On the supported CI platform, the following parameters are complemented by the built-in environment variables.

- `-owner`
- `-repo`
- `-pr`
- `-sha`
- `-build-url`

This feature is implemented by [go-ci-env](https://github.com/suzuki-shunsuke/go-ci-env).

:warning: You can also use tfcmt on other platform with CLI flags or Custom Environment Variable Definition.

## Google Cloud Build Support

tfcmt >= [v3.3.0](https://github.com/suzuki-shunsuke/tfcmt/releases/tag/v3.3.0)

[#376](https://github.com/suzuki-shunsuke/tfcmt/pull/376)

Set the environment variable `GOOGLE_CLOUD_BUILD`.

```sh
GOOGLE_CLOUD_BUILD=true
```

Set the following environment variables using [substitutions](https://cloud.google.com/cloud-build/docs/configuring-builds/substitute-variable-values).

* `COMMIT_SHA`
* `BUILD_ID`
* `PROJECT_ID`
* `_PR_NUMBER`
* `_REGION`

Specify the repository owner and name in `tfcmt.yaml`.

e.g.

tfcmt.yaml

```yaml
repo_owner: suzuki-shunsuke
repo_name: tfcmt
```

## Custom Environment Variable Definition

:::caution
This feature was removed from [v4.0.0](https://github.com/suzuki-shunsuke/tfcmt/releases/tag/v4.0.0) for security reason.
:::
