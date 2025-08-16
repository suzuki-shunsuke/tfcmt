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
brew install suzuki-shunsuke/tfcmt/tfcmt
```

## aqua

You can install tfcmt with [aqua](https://aquaproj.github.io/) too.

```sh
aqua g -i suzuki-shunsuke/tfcmt
```

## GitHub Releases

Grab the binary from [GitHub Releases](https://github.com/suzuki-shunsuke/tfcmt/releases)

### Verify downloaded binaries from GitHub Releases

You can verify downloaded binaries using some tools.

1. [Cosign](https://github.com/sigstore/cosign)
1. [slsa-verifier](https://github.com/slsa-framework/slsa-verifier)
1. [GitHub CLI](https://cli.github.com/)

#### 1. Cosign

You can install Cosign by aqua.

```sh
aqua g -i sigstore/cosign
```

```sh
# Download assets from GitHub Releases.
gh release download -R suzuki-shunsuke/tfcmt v4.14.0
# Verify a checksum file.
cosign verify-blob \
  --signature tfcmt_4.14.0_checksums.txt.sig \
  --certificate tfcmt_4.14.0_checksums.txt.pem \
  --certificate-identity-regexp 'https://github\.com/suzuki-shunsuke/go-release-workflow/\.github/workflows/release\.yaml@.*' \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  tfcmt_4.14.0_checksums.txt
```

Output:

```
Verified OK
```

After verifying the checksum, verify the artifact.

```sh
cat tfcmt_4.14.0_checksums.txt | sha256sum -c --ignore-missing
```

#### 2. slsa-verifier

You can install slsa-verifier by aqua.

```sh
aqua g -i slsa-framework/slsa-verifier
```

```sh
# Download assets from GitHub Releases.
gh release download -R suzuki-shunsuke/tfcmt v4.14.0
# Verify an asset.
slsa-verifier verify-artifact tfcmt_darwin_arm64.tar.gz \
  --provenance-path multiple.intoto.jsonl \
  --source-uri github.com/suzuki-shunsuke/tfcmt \
  --source-tag v4.14.0
```

Output:

```
Verified signature against tlog entry index 136685045 at URL: https://rekor.sigstore.dev/api/v1/log/entries/108e9186e8c5677a9b654937f69fcad5c5078be5a058025d612085e3f1befcae9b51fbcaca3edd08
Verified build using builder "https://github.com/slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@refs/tags/v2.0.0" at commit 13b3b64b1444d528db49d60a99310bcd45993a52
Verifying artifact tfcmt_darwin_arm64.tar.gz: PASSED
```

#### 3. GitHub CLI

You can install GitHub CLI by aqua.

```sh
aqua g -i cli/cli
```

```sh
# Download assets from GitHub Releases.
gh release download -R suzuki-shunsuke/tfcmt v4.14.0 -p tfcmt_darwin_arm64.tar.gz
# Verify an asset.
gh attestation verify tfcmt_darwin_arm64.tar.gz \
  -R suzuki-shunsuke/tfcmt \
  --signer-workflow suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml
```

Output:

```
Loaded digest sha256:5789ea2f3165b0448f84a46df6489b01d0c90802d2c95d3fa4b74de06177ced7 for file://tfcmt_darwin_arm64.tar.gz
Loaded 1 attestation from GitHub API
✓ Verification succeeded!

sha256:5789ea2f3165b0448f84a46df6489b01d0c90802d2c95d3fa4b74de06177ced7 was attested by:
REPO                                 PREDICATE_TYPE                  WORKFLOW                                                               
suzuki-shunsuke/go-release-workflow  https://slsa.dev/provenance/v1  .github/workflows/release.yaml@7f97a226912ee2978126019b1e95311d7d15c97a
```
