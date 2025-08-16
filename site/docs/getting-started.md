---
sidebar_position: 150
---

# Getting Started

In this getting started, you can understand tfcmt's primary feature.

## Requirements

- Terraform
- [tfcmt](install)
- GitHub Access Token

GitHub Access Token requires the following permissions:

- Pull Requests: Write - To post comments to pull requests
- Issues: Write - To create GitHub Issue Labels

## Setup

```sh
git clone https://github.com/suzuki-shunsuke/tfcmt
cd tfcmt/examples/getting-started
terraform init
terraform validate
terraform plan
```

```sh
export GITHUB_TOKEN=xxx # your personal access token
```

Open an issue or pull request at your repository to post comments with tfcmt.
**Please change values properly**.

```sh
PR_NUMBER=70
OWNER=suzuki-shunsuke
REPO=tfcmt
```

## tfcmt plan

By `tfcmt plan` command, you can post the result of `terraform plan` as a comment.

```console
$ tfcmt -owner "$OWNER" -repo "$REPO" -pr "$PR_NUMBER" plan -- terraform plan

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # null_resource.foo will be created
  + resource "null_resource" "foo" {
      + id = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

------------------------------------------------------------------------

Note: You didn't specify an "-out" parameter to save this plan, so Terraform
can't guarantee that exactly these actions will be performed if
"terraform apply" is subsequently run.
```

https://github.com/suzuki-shunsuke/tfcmt/pull/70#issuecomment-797854184

![image](https://user-images.githubusercontent.com/13323303/111016701-b6f89200-83f2-11eb-9fed-35d8249c9ba0.png)

1. By the following message, you can know the number of added, changed, and destroyed resources.

```
Plan: 1 to add, 0 to change, 0 to destroy.
```

2. By the following message, you can know which resources are added, changed, and destroyed, and recreated.

```markdown
* Create
  * null_resource.foo
```

Please see [List of changed resources](#list-of-changed-resources) too.

3. By opening `Details`, you can confirm the full output of `terraform plan`.

![image](https://user-images.githubusercontent.com/13323303/111022026-7fe6a880-8413-11eb-84db-3159402d42f3.png)

### Pull Request label

By `tfcmt plan`, pull request labels are set.
It makes easy to understand the result of `terraform plan`.

![image](https://user-images.githubusercontent.com/13323303/111016806-31291680-83f3-11eb-94f0-be22585aae64.png)

The following labels are set according to the result of `terraform plan`

- no-changes: there is no resource to be changed
- add-or-update: there are resources to be created or updated but there is no resource to be destroyed or recreated
- destroy: there are resources to be destroyed or recreated

The label color is configured automatically.

- no-changes: green
- add-or-update: blue
- destroy: red

### Configuration file is optional

You don't have to prepare the configuration file for tfcmt.
The configuration file is optional.
tfcmt provides the good default configuration.
You can also customize the configuration with configuration file if needed.

Please see [Configuration](config.md) too.

### List of changed resources

You can know which resources are changed without looking the output of `terraform plan`.

Please look at the comment closer.

https://github.com/suzuki-shunsuke/tfcmt/pull/70#issuecomment-797854184

![image](https://user-images.githubusercontent.com/13323303/111016959-1014f580-83f4-11eb-8a6f-2f5ee7bd9607.png)

The following resources are listed.

* Create
  * created resource 1
  * ...
* Update
  * updated resource 1
  * ...
* Delete
  * deleted resource 1
  * ...
* Replace
  * recreated resource 1
  * ...
* Import
  * imported resource 1
  * ...
* Move
  * old resource path 1 => new resource path 1
  * ...

## tfcmt apply

By `tfcmt apply` command, you can post the result of `terraform apply`

```console
$ tfcmt -owner "$OWNER" -repo "$REPO" -pr "$PR_NUMBER" apply -- terraform apply -auto-approve
null_resource.foo: Creating...
null_resource.foo: Creation complete after 0s [id=459501600381334523]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

https://github.com/suzuki-shunsuke/tfcmt/pull/70#issuecomment-797856135

![image](https://user-images.githubusercontent.com/13323303/111017196-3c7d4180-83f5-11eb-9e4b-ba350e07adf8.png)

## Deletion warning

Let's remove the resource `null_resource.foo` and run `terraform plan`.

```sh
vi main.tf
```

```tf
# comment out
# resource "null_resource" "foo" {}
```

```console
$ tfcmt -owner "$OWNER" -repo "$REPO" -pr "$PR_NUMBER" plan -- terraform plan
null_resource.foo: Refreshing state... [id=459501600381334523]

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # null_resource.foo will be destroyed
  - resource "null_resource" "foo" {
      - id = "459501600381334523" -> null
    }

Plan: 0 to add, 0 to change, 1 to destroy.

------------------------------------------------------------------------

Note: You didn't specify an "-out" parameter to save this plan, so Terraform
can't guarantee that exactly these actions will be performed if
"terraform apply" is subsequently run.
```

Then in the comment of `tfcmt plan`, the deletion warning is included.

https://github.com/suzuki-shunsuke/tfcmt/pull/70#issuecomment-797856650

![image](https://user-images.githubusercontent.com/13323303/111017250-97169d80-83f5-11eb-9144-039dc4b62b37.png)

You can also find the pull request label is changed from `add-or-update` to `destroy`.

It is helpful to prevent unexpected resource deletion.

## Support of CI platforms

In the above commands, you specify the repository owner, name, and pull request number as the command line arguments.

```console
$ tfcmt -owner "$OWNER" -repo "$REPO" -pr "$PR_NUMBER"
```

But on the following CI platform, tfcmt gets these parameters from the built in environment variables so you don't have to specify these arguments.

- AWS CodeBuild
- CircleCI
- Drone
- GitHub Actions
- [Google Cloud Build](environment-variable.md#google-cloud-build-support)

AS IS

```sh
tfcmt -owner "$OWNER" -repo "$REPO" -pr "$PR_NUMBER" plan -- terraform plan
```

TO BE

```sh
tfcmt plan -- terraform plan
```

Note that if tfcmt can't get the pull request number from environment variables you have to complement it.

## Hide old comments with github-comment

When running CI of the same pull request at many times,
it is convenient to hide old comments posted by tfcmt.

tfcmt itself doesn't support to hide old comments, but you can hide old comments with [github-comment](https://github.com/suzuki-shunsuke/github-comment).
tfcmt embeds metadata in a comment as HTML comment.
Please check comments posted by tfcmt.

![image](https://user-images.githubusercontent.com/13323303/111018042-20c86a00-83fa-11eb-9f85-491649411005.png)

![image](https://user-images.githubusercontent.com/13323303/111018071-40f82900-83fa-11eb-8583-1601ea3af484.png)

```html
<!-- github-comment: {"Command":"plan","PRNumber":70,"Program":"tfcmt","SHA1":"","Vars":{}} -->
```

So you can hide comments with [github-comment hide](https://github.com/suzuki-shunsuke/github-comment#hide) command.

## Monorepo support: target variable

Let's assume that the repository is Monorepo and there are multiple Terraform states in the repository.

For example,

```
foo/
  main.tf
  ...
bar/
  main.tf
  ...
```

In the above case, you have to distinguish comments for the state `foo` and `bar`.
By specifying the special variable `target` by `-var` argument, you can do it.

```sh
vi main.tf
```

```tf
resource "null_resource" "foo" {}
```

```console
$ tfcmt -owner "$OWNER" -repo "$REPO" -pr "$PR_NUMBER" -var "target:foo" plan -- terraform plan
null_resource.foo: Refreshing state... [id=459501600381334523]

No changes. Infrastructure is up-to-date.

This means that Terraform did not detect any differences between your
configuration and real physical resources that exist. As a result, no
actions need to be performed.
```

https://github.com/suzuki-shunsuke/tfcmt/pull/70#issuecomment-797861332

![image](https://user-images.githubusercontent.com/13323303/111018399-ea8bea00-83fb-11eb-8efe-205ab8c996b7.png)

We can find

- the target name is included in the comment title `Plan Result (foo)`
- the target name is included in the pull request label `foo/no-changes`
