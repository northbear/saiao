# Version Management

This document describes a simple versioning policy that can be reused across small service-style applications, not only SAIAO.

## Goals

The policy aims to keep versioning:

- explicit in the repository
- reproducible in local and CI builds
- visible at runtime
- aligned between binaries and container images

## Single Source Of Truth

The canonical application version is stored in a root `VERSION` file.

Example:

```text
0.1.0
```

This value represents the intended application version independently of a specific build environment.

## Embedded Build Metadata

Builds should embed at least:

- `version`
- `commit`

Recommended meaning:

- `version`: the effective application version derived from `VERSION`
- `commit`: the git revision used to produce the build artifact

Embedding this metadata makes runtime version reporting trustworthy and easy to inspect through an info endpoint, CLI flag, or startup log.

## Branch Policy

The effective version is derived from the repository version plus branch context:

- builds from `master` use the plain version from `VERSION`
- builds from non-`master` branches append `-dev`

Examples:

- `VERSION=0.1.0` on `master` becomes `0.1.0`
- `VERSION=0.1.0` on `feature/x` becomes `0.1.0-dev`

This keeps release-line builds and development builds clearly distinguishable.

Detached HEAD builds should usually be treated as development builds unless the release pipeline applies a stricter rule.

## Build Injection

The version should be injected at build time rather than hardcoded in application logic.

Typical pattern:

1. Read `VERSION`.
2. Determine the current branch.
3. Compute the effective version.
4. Read the current commit SHA.
5. Inject both into the build.

For Go services this commonly means using `-ldflags` to set package variables.

## Container Tag Policy

The default Docker image tag should match the effective application version.

Examples:

- release build on `master`: `0.1.0`
- development build on another branch: `0.1.0-dev`

This keeps the runtime-reported version and the image tag aligned.

Additional tags such as `latest` may still be published intentionally, but they should not be the primary version identity.

## Release Discipline

Recommended release rules:

- bump `VERSION` intentionally in source control
- ensure release builds use the plain version, not a `-dev` suffix
- if git tags are used, require the release tag to match `VERSION`

This gives one explicit repository version and one exact source revision for every artifact.

## Why This Policy Works Well

This approach avoids several common problems:

- no version string hidden deep in source code
- no dependency on runtime environment variables
- no ambiguity between app version and image tag
- no requirement that every build happen from a tagged checkout

It works well for:

- local builds
- CI builds
- container images
- source-based deployments outside containers

## SAIAO Example

SAIAO applies this policy as follows:

- root [`VERSION`](/home/space/devel/aikvn/saiao/VERSION) is the canonical version
- [`Makefile`](/home/space/devel/aikvn/saiao/Makefile) computes the effective version and injects it into the Go build
- [`Dockerfile`](/home/space/devel/aikvn/saiao/Dockerfile) applies the same logic for container builds
- the runtime exposes embedded build metadata through [`docs/api-spec.md`](/home/space/devel/aikvn/saiao/docs/api-spec.md) `GET /info`

