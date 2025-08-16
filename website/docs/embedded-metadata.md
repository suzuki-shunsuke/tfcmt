---
sidebar_position: 400
---

# Embed metadata in comments

[#67](https://github.com/suzuki-shunsuke/tfcmt/pull/67)

tfcmt embeds metadata into comment with [github-comment-metadata](https://github.com/suzuki-shunsuke/github-comment-metadata).
tfcmt itself doesn't support hiding old comments, but you can hide comments with [github-comment's hide command](https://github.com/suzuki-shunsuke/github-comment#hide).

## embedded_var_names

[#108](https://github.com/suzuki-shunsuke/tfcmt/issues/108) [#115](https://github.com/suzuki-shunsuke/tfcmt/pull/115)

If you want to embed variables passed by `-var` option, you have to specify variable names by `embedded_var_names` in the configuration.

```yaml
embedded_var_names:
- name
```
