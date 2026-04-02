# SAIAO Roadmap

## Current Baseline

- MVP runtime is implemented end to end
- manifest and invoke APIs are working, including request correlation IDs
- config loading, validation, secret resolution, and auth helpers are in place
- HTTP, SSH, shell, and email executors are implemented
- `/info` build metadata is injected at build time from the root `VERSION` file
- OpenAI-compatible manifest output is available through `GET /manifest?format=openai`
- binary startup and Docker / Docker Buildx verification are covered
- README, config examples, and integration documentation cover local use and common AI stacks
- `go test ./...` is green with direct coverage for runtime-critical internal packages

## Release Follow-Up

- improve structured audit logging while keeping request payloads and rendered content out by default
- harden shell and SSH execution policy boundaries beyond the current template validation rules
- add CI automation for `go test ./...`, `docker build`, and `docker buildx build`
- document release/build steps and reproducible container build expectations
- decide whether `Makefile` targets should default to Buildx-based image builds

## V2 Candidates

- add a richer audit/event model with stable event names and non-sensitive operator fields
- add optional detailed invocation logging with safe redaction controls
- tighten validation around secret-like fields and executor-specific sensitive settings
- define and enforce a clearer production policy for allowed shell and SSH command patterns
- add additional manifest adapter formats beyond the native SAIAO and OpenAI views where they provide real interoperability value
- normalize schemas further for stricter provider tool-definition compatibility when appropriate

## Runtime Safety

- add explicit redaction helpers for future detailed logging modes
- tighten validation around secret-like fields and executor-specific sensitive settings
- decide whether executor outputs should be truncated or normalized before logging
- define a clearer policy for allowed shell and SSH command patterns in production deployments

## API Evolution

- decide whether API error responses should carry stable request identifiers for support/debugging
- decide whether additional manifest response formats are worth supporting beyond `openai`
- decide whether stricter tool-definition compatibility should be exposed as an additional manifest mode

## Future

- optional detailed invocation logging with safe redaction controls
- richer audit/event model
- additional executor hardening and integration test depth
- more complete operator documentation and example deployment layouts
