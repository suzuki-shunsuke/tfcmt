---
sidebar_position: 550
---

# Patch `tfcmt plan` comment

tfcmt >= [v3.2.0](https://github.com/suzuki-shunsuke/tfcmt/releases/tag/v3.2.0) | [#199](https://github.com/suzuki-shunsuke/tfcmt/issues/199) [#245](https://github.com/suzuki-shunsuke/tfcmt/issues/245) [#248](https://github.com/suzuki-shunsuke/tfcmt/issues/248) [#249](https://github.com/suzuki-shunsuke/tfcmt/issues/249)

Instead of creating a new comment, you can update existing comment. This is useful to keep the pull request comments clean.

![image](https://user-images.githubusercontent.com/13323303/164969354-02bdd49a-547e-4951-9262-033ec5b4db11.png)

--

![image](https://user-images.githubusercontent.com/13323303/164969385-355e801e-3d58-4b75-9657-0bcc10da8d12.png)

The option `-patch` has been added to `tfcmt plan` command.

```sh
tfcmt plan -patch -- terraform plan -no-color
```

And the configuration option `plan_patch` has been added.

```yaml
plan_patch: true
```

The command line option `-patch` takes precedence over configuration file option `plan_patch`.

If you want to disable patching although `plan_patch` is true, please set `-patch=false`.

```
tfcmt plan -patch=false -- terraform plan -no-color
```

### Motivation

By patching the comment instead of creating a new comment, you can keep the pull request comments clean.

### Using `-patch` with monorepos containing multiple root modules (tfstates)

You can specify the `target` variable to instruct tfcmt which comments should be updated:

```sh
cd /path/to/root-modules/dev
tfcmt -var 'target:dev' plan -patch -- terraform plan -no-color

cd /path/to/root-modules/prd
tfcmt -var 'target:prd' plan -patch -- terraform plan -no-color
```

See also [Monorepo support: target variable](getting-started#monorepo-support-target-variable).

### Trouble shooting

If the comment isn't patched expectedly, please set `-log-level=debug`.

```sh
tfcmt -log-level=debug plan -patch -- terraform plan -no-color
```

### :warning: Note to use  tfcmt plan's patch option with github-comment hide

If you hide comments by [github-comment hide](https://suzuki-shunsuke.github.io/github-comment/hide) and enable tfcmt plan's patch option,
you should be careful not to hide tfcmt plan's comments.

There are some ways to fix the problem.

1. Stop using `github-comment hide`
1. Fix github-comment hide's condition and exclude tfcmt's comments from the target ofgithub-comment hide
1. Run github-comment hide after tfcmt
