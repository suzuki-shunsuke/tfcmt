---
sidebar_position: 120
---

# Install

- [Homebrew](#homebrew)
- [aqua](#aqua)
- [GitHub Releases](#github-releases)

## Homebrew

You can install tfcmt with [Homebrew](https://brew.sh/).

```sh
brew install tfcmt
```

Or

```sh
brew install --cask suzuki-shunsuke/tfcmt/tfcmt
```

## aqua

You can install tfcmt with [aqua](https://aquaproj.github.io/) too.

```sh
aqua g -i suzuki-shunsuke/tfcmt
```

## GitHub Releases

Grab the binary from [GitHub Releases](https://github.com/suzuki-shunsuke/tfcmt/releases)

### Verify downloaded binaries from GitHub Releases

You can verify downloaded assets using some tools.

1. [GitHub CLI](https://cli.github.com/)
1. [Cosign](https://github.com/sigstore/cosign)
1. [Packslip](https://packslip.dev/) `tfcmt >= v4.14.20`
1. [slsa-verifier](https://github.com/slsa-framework/slsa-verifier) `tfcmt <= v4.14.19`

#### 1. GitHub CLI

You can install GitHub CLI by aqua.

```sh
aqua g -i cli/cli
```

```sh
version=v4.14.20
asset=tfcmt_darwin_arm64.tar.gz
gh release download -R suzuki-shunsuke/tfcmt "$version" -p "$asset"
gh attestation verify "$asset" \
  -R suzuki-shunsuke/tfcmt \
  --signer-workflow suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml
```

#### 2. Cosign

You can install Cosign by aqua.

```sh
aqua g -i sigstore/cosign
```

From tfcmt v4.14.16, the bundle file is `tfcmt_<version>_checksums.txt.sigstore.json`.
Cosign v2.4.2 or later is required.

```sh
version=v4.14.20
checksum_file="tfcmt_${version#v}_checksums.txt"
asset=tfcmt_darwin_arm64.tar.gz
gh release download "$version" \
  -R suzuki-shunsuke/tfcmt \
  -p "$asset" \
  -p "$checksum_file" \
  -p "${checksum_file}.sigstore.json"
cosign verify-blob \
  "$checksum_file" \
  --bundle "${checksum_file}.sigstore.json" \
  --certificate-identity-regexp 'https://github\.com/suzuki-shunsuke/go-release-workflow/\.github/workflows/release\.yaml@.*' \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com"
cat "$checksum_file" | sha256sum -c --ignore-missing -
```

For tfcmt v4.14.15 or older, `*.pem` and `*.sig` are available instead.

```sh
version=v4.14.15
checksum_file="tfcmt_${version#v}_checksums.txt"
asset=tfcmt_darwin_arm64.tar.gz
gh release download "$version" \
  -R suzuki-shunsuke/tfcmt \
  -p "$asset" \
  -p "$checksum_file" \
  -p "${checksum_file}.pem" \
  -p "${checksum_file}.sig"
cosign verify-blob \
  --signature "${checksum_file}.sig" \
  --certificate "${checksum_file}.pem" \
  --certificate-identity-regexp 'https://github\.com/suzuki-shunsuke/go-release-workflow/\.github/workflows/release\.yaml@.*' \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  "$checksum_file"
cat "$checksum_file" | sha256sum -c --ignore-missing -
```

#### 3. Packslip

`tfcmt >= v4.14.20`

You can install packslip by aqua.

```sh
aqua g -i jdx/packslip
```

```sh
version=v4.14.20
asset=tfcmt_darwin_arm64.tar.gz
gh release download -R suzuki-shunsuke/tfcmt "$version" -p "$asset" -p packslip.sigstore.json
packslip verify packslip.sigstore.json \
    --pin ps1_baqonzlpvy53hygidrvusbl2nq \
    --identity-prefix https://github.com/suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml@ \
    --artifact "$asset"
```

- `--pin` requires the release to be signed from the repository `suzuki-shunsuke/tfcmt`, identified by its repository ID rather than its name
- `--identity-prefix` requires the release to be signed by the release workflow of [suzuki-shunsuke/go-release-workflow](https://github.com/suzuki-shunsuke/go-release-workflow)

On success, packslip prints a line starting with `ok:`.
Please make sure it shows the project and version you meant, e.g. `github.com/suzuki-shunsuke/tfcmt 4.14.20`.

#### 4. slsa-verifier

:::caution
As of tfcmt v4.14.20, `multiple.intoto.jsonl` hasn't been published anymore.
:::

You can install slsa-verifier by aqua.

```sh
aqua g -i slsa-framework/slsa-verifier
```

```sh
version=v4.14.19
asset=tfcmt_darwin_arm64.tar.gz
gh release download -R suzuki-shunsuke/tfcmt "$version" -p "$asset" -p multiple.intoto.jsonl
slsa-verifier verify-artifact "$asset" \
  --provenance-path multiple.intoto.jsonl \
  --source-uri github.com/suzuki-shunsuke/tfcmt \
  --source-tag "$version"
```
