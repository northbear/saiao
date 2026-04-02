# SAIAO Roadmap

## Current Baseline

- MVP runtime is implemented end to end
- manifest and invoke APIs are working
- config loading, validation, secret resolution, and auth helpers are in place
- HTTP, SSH, shell, and email executors are implemented
- binary startup and Docker / Docker Buildx verification are covered
- `go test ./...` is green with direct coverage for runtime-critical internal packages

## Next Priorities

- add request correlation IDs so API and invocation lifecycle logs can be tied together
- improve structured audit logging while keeping request payloads and rendered content out by default
- harden shell and SSH execution policy boundaries beyond the current template validation rules
- add CI automation for `go test ./...`, `docker build`, and `docker buildx build`
- improve README and operator-facing documentation for configuration, runtime model, and deployment

## Runtime Safety

- add explicit redaction helpers for future detailed logging modes
- tighten validation around secret-like fields and executor-specific sensitive settings
- decide whether executor outputs should be truncated or normalized before logging
- define a clearer policy for allowed shell and SSH command patterns in production deployments

## API Evolution

- add request IDs to API responses and logs where useful
- inject build metadata into `/info` from build-time variables instead of fixed defaults
- review manifest compatibility with common LLM tool-calling formats
- decide whether API error responses should carry stable request identifiers for support/debugging

## Delivery

- add CI workflows for tests and image builds
- document release/build steps and reproducible container build expectations
- decide whether `Makefile` targets should default to Buildx-based image builds

## Future

- optional detailed invocation logging with safe redaction controls
- richer audit/event model
- additional executor hardening and integration test depth
- more complete operator documentation and example deployment layouts
