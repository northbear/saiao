# SAIAO TODO

## Milestone 1: Minimal runnable binary

- [x] Initialize Go module
- [x] Create application entrypoint
- [x] Add internal app wiring
- [x] Add config models
- [x] Add config loader stub
- [x] Add config validation stub
- [x] Add secret resolution stub
- [x] Add auth package stub
- [x] Add API server stub
- [x] Add manifest package stub
- [x] Add invoke package stub
- [x] Add executor interface and stubs
- [x] Add logging package stub
- [x] Add error mapping stub
- [x] Add basic tests
- [x] Verify project builds with `go test ./...`
- [x] Add Dockerfile
- [x] Add example config
- [x] Add Makefile for tests/build/docker-build/clean

## Milestone 2: Core MVP runtime behavior

- [x] Implement YAML config parsing
- [x] Implement startup config validation
- [x] Implement secret resolution from env vars
- [x] Implement secret resolution from mounted files
- [x] Implement bearer token parsing
- [x] Implement token-to-group mapping
- [x] Implement `/info` response wiring with build info source
- [x] Implement tool manifest generation
- [x] Implement invoke request flow
- [x] Implement input schema validation subset
- [x] Implement normalized error mapping
- [x] Implement HTTP executor
- [x] Implement SSH executor
- [x] Implement Shell executor
- [x] Implement Email executor

## Milestone 3: MVP hardening

- [x] Add tests for config validation failures
- [x] Add tests for secret resolution failures
- [x] Add tests for unauthorized access
- [x] Add tests for manifest filtering by tool group
- [x] Add tests for invoke success and failure
- [x] Add timeout handling tests
- [x] Add integration test for binary startup
- [x] Add container build verification
- [x] Add README usage documentation

## Notes

Completed today:
- go module initialized
- CLI entrypoint added
- app startup path wired
- config-backed runnable binary created
- /info endpoint added
- Dockerfile created
- Makefile created
- example config added
- unit tests passing

Current state:
- the binary is runnable
- Milestone 2 runtime behavior is implemented end-to-end for manifest retrieval and action invocation
- Milestone 3 verification coverage now includes binary startup and Docker build validation
- the next meaningful gaps are broader executor-specific integration coverage and runtime polish
