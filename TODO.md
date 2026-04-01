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

- [ ] Implement YAML config parsing
- [ ] Implement startup config validation
- [ ] Implement secret resolution from env vars
- [ ] Implement secret resolution from mounted files
- [ ] Implement bearer token parsing
- [ ] Implement token-to-group mapping
- [ ] Implement `/info` response wiring with build info source
- [ ] Implement tool manifest generation
- [ ] Implement invoke request flow
- [ ] Implement input schema validation subset
- [ ] Implement normalized error mapping
- [ ] Implement HTTP executor
- [ ] Implement SSH executor
- [ ] Implement Shell executor
- [ ] Implement Email executor

## Milestone 3: MVP hardening

- [ ] Add tests for config validation failures
- [ ] Add tests for secret resolution failures
- [ ] Add tests for unauthorized access
- [ ] Add tests for manifest filtering by tool group
- [ ] Add tests for invoke success and failure
- [ ] Add timeout handling tests
- [ ] Add integration test for binary startup
- [ ] Add container build verification
- [ ] Add README usage documentation

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
- the core MVP runtime behavior is still stubbed
- next work should focus on config parsing, auth, and action invocation
