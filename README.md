# SAIAO: Simple AI Action Orchestrator

SAIAO is a small HTTP service that lets AI-facing applications discover and invoke predefined actions without exposing execution credentials to the caller.

It is designed for the "controlled execution layer" pattern:

- the AI application gets a scoped bearer token for a tool group
- SAIAO owns the actual execution credentials and runtime boundaries
- actions are predefined in YAML rather than assembled dynamically by the caller

SAIAO currently supports:

- `shell` actions
- `ssh` actions
- `http` actions
- `email` actions

## Why SAIAO

SAIAO separates three concerns that are often mixed together in AI-to-system integrations:

- authorization for what the caller may invoke
- identity and credentials used to perform the work
- execution of the underlying operation

That separation makes the runtime easier to reason about, test, and operate.

## Core Model

- `identity`: execution-side identity and secret material used by an action
- `action`: a predefined operation with an executor type and optional input schema
- `tool_group`: the subset of actions exposed to a specific caller
- bearer token: authorization for a single tool group

## Quick Start

1. Create a config file from [`configs/example.yaml`](/home/space/devel/aikvn/saiao/configs/example.yaml).
2. Export the access token environment variable referenced by the tool group.
3. Start SAIAO.
4. Fetch the manifest and invoke an action.

In the example config, the tool group sets:

```yaml
tool_groups:
  - name: default_group
    access_token_env: SAIAO_TOKEN_DEFAULT
```

That means SAIAO reads the bearer token value from the server process environment variable `SAIAO_TOKEN_DEFAULT`. Client applications that call `default_group` must be provisioned with that same token value out of band.

### Local Run

```bash
export SAIAO_TOKEN_DEFAULT=secret
go run ./cmd/saiao -config ./configs/example.yaml
```

### Docker Run

```bash
make docker-build

docker run --rm -p 8080:8080 \
  -v "$(pwd)/configs/example.yaml:/etc/saiao/config.yaml:ro" \
  -e SAIAO_TOKEN_DEFAULT=secret \
  saiao:0.1.0
```

The default image tag matches the effective application version. On `master` that is the plain `VERSION`, and on other branches it gets a `-dev` suffix.

## End-to-End Example

The example config exposes one shell action named `echo_message` through the tool group `default_group`.

### 1. Get service metadata

```bash
curl http://localhost:8080/info
```

Example response:

```json
{
  "request_id": "4d3f9d1f7e2a4f2e9f0f0ccf7f4f1c57",
  "status": "ok",
  "service": "saiao",
  "version": "0.1.0",
  "commit": "abc1234"
}
```

### 2. Fetch the manifest

```bash
curl \
  -H "Authorization: Bearer secret" \
  http://localhost:8080/tool-groups/default_group/manifest
```

Example response:

```json
{
  "request_id": "4d3f9d1f7e2a4f2e9f0f0ccf7f4f1c57",
  "tools": [
    {
      "name": "echo_message",
      "description": "Echo a validated message",
      "input_schema": {
        "type": "object",
        "properties": {
          "message": {
            "type": "string"
          }
        },
        "required": [
          "message"
        ]
      }
    }
  ]
}
```

### 3. Invoke the action

```bash
curl -X POST \
  -H "Authorization: Bearer secret" \
  -H "Content-Type: application/json" \
  -d '{"message":"hello from saiao"}' \
  http://localhost:8080/tool-groups/default_group/actions/echo_message/invoke
```

Example response:

```json
{
  "request_id": "4d3f9d1f7e2a4f2e9f0f0ccf7f4f1c57",
  "status": "success",
  "result": {
    "output": "hello from saiao",
    "exit_code": 0
  }
}
```

## Configuration

SAIAO reads a single YAML config file at startup. By default the binary looks for:

```text
/etc/saiao/config.yaml
```

The main config sections are:

- `server`
- `logging`
- `identities`
- `actions`
- `tool_groups`

Secret values should not be committed inline. Instead, reference them through environment variables or files and let SAIAO resolve them at startup.

See:

- [`configs/example.yaml`](/home/space/devel/aikvn/saiao/configs/example.yaml)
- [`docs/config-spec.md`](/home/space/devel/aikvn/saiao/docs/config-spec.md)

## API Summary

SAIAO exposes three HTTP endpoints:

- `GET /info`
- `GET /tool-groups/{group_name}/manifest`
- `POST /tool-groups/{group_name}/actions/{action_name}/invoke`

All manifest and invoke calls require `Authorization: Bearer <token>`.

Responses also include `X-Request-Id`. If the client sends a valid `X-Request-Id`, SAIAO reuses it; otherwise SAIAO generates one.

See [`docs/api-spec.md`](/home/space/devel/aikvn/saiao/docs/api-spec.md).

For OpenAI tool calling, SAIAO can return an OpenAI-adapted manifest directly from:

```text
GET /tool-groups/{group_name}/manifest?format=openai
```

The bearer token used by clients is the tool-group token configured through `tool_groups[].access_token_env` in the SAIAO config, not an executor credential and not an OpenAI API secret.

## Integration Guides

Standalone integration guides are available for common AI application stacks:

- [Python](/home/space/devel/aikvn/saiao/docs/integrations/python.md)
- [JavaScript / TypeScript](/home/space/devel/aikvn/saiao/docs/integrations/javascript.md)
- [Go](/home/space/devel/aikvn/saiao/docs/integrations/go.md)
- [LangGraph](/home/space/devel/aikvn/saiao/docs/integrations/langgraph.md)

## Logging

Logging is written to stdout and supports:

- `json`
- `text`

For `json` format, each log record is emitted as one JSON object per line, which works well with line-oriented tooling and log collectors.

Invocation logs share the same `request_id` so a single request can be followed across API and execution lifecycle events.

## Versioning And Builds

The canonical application version lives in the root [`VERSION`](/home/space/devel/aikvn/saiao/VERSION) file.

Builds embed:

- `version` from `VERSION`
- `commit` from the current git revision

Branch behavior:

- builds from `master` use the plain version, for example `0.1.0`
- builds from non-`master` branches use a `-dev` suffix, for example `0.1.0-dev`

The default Docker image tag follows the same effective version.

The full policy is documented in [`docs/version-management.md`](/home/space/devel/aikvn/saiao/docs/version-management.md).

Useful commands:

```bash
make test
make build
make docker-build
```

## Current Scope

SAIAO is intentionally small. It focuses on:

- explicit configuration
- predictable executor behavior
- minimal runtime surface

It does not yet try to be a full workflow engine, policy engine, or secrets platform.

## License

This project is licensed under the MIT License. See [`LICENSE`](/home/space/devel/aikvn/saiao/LICENSE).
