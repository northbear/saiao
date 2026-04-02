# Authentication

Bearer token:

```text
Authorization: Bearer <token>
```

A bearer token grants access only to the tool group associated with that token.

# Endpoints

## GET /info

Returns service status and build/runtime metadata.

### Success response

```json
{
  "request_id": "4d3f9d1f7e2a4f2e9f0f0ccf7f4f1c57",
  "status": "ok",
  "service": "saiao",
  "version": "0.1.0",
  "commit": "abc1234"
}
```

---

## GET /tool-groups/{group_name}/manifest

Returns tools available to the caller.

### Query parameters

- `format` optional
- supported values: `saiao`, `openai`
- default: `saiao`

### Success response

```json
{
  "request_id": "4d3f9d1f7e2a4f2e9f0f0ccf7f4f1c57",
  "format": "saiao",
  "tools": [
    {
      "name": "action_name",
      "description": "optional description",
      "input_schema": {}
    }
  ]
}
```

`input_schema` follows the SAIAO JSON-Schema-compatible subset defined in `config-spec.md`.

### OpenAI format response

```json
{
  "request_id": "4d3f9d1f7e2a4f2e9f0f0ccf7f4f1c57",
  "format": "openai",
  "tools": [
    {
      "type": "function",
      "name": "action_name",
      "description": "optional description",
      "parameters": {},
      "strict": false
    }
  ]
}
```

---

## POST /tool-groups/{group_name}/actions/{action_name}/invoke

Executes the action.

### Request body

A JSON object matching the action `input_schema`.

### Success response

```json
{
  "request_id": "4d3f9d1f7e2a4f2e9f0f0ccf7f4f1c57",
  "status": "success",
  "result": {
    "output": "string",
    "exit_code": 0
  }
}
```

### Notes

* `output` is a generic textual result
* `exit_code` is mandatory in the normalized result model

# Error Format

```json
{
  "request_id": "4d3f9d1f7e2a4f2e9f0f0ccf7f4f1c57",
  "error": {
    "code": "invalid_input",
    "message": "description"
  }
}
```

Responses also include the `X-Request-Id` header. If the caller sends `X-Request-Id`, SAIAO reuses it when valid; otherwise SAIAO generates one.

## Usage Notes

Typical client flow:

1. Call `GET /info` to verify the service is reachable and inspect version metadata.
2. Call `GET /tool-groups/{group_name}/manifest` with a bearer token to discover available tools.
3. Call `POST /tool-groups/{group_name}/actions/{action_name}/invoke` with a JSON object matching the tool schema.

If you want to pass the manifest directly into OpenAI tool calling, request:

```text
GET /tool-groups/{group_name}/manifest?format=openai
```

Example:

```bash
curl http://localhost:8080/info

curl \
  -H "Authorization: Bearer secret" \
  http://localhost:8080/tool-groups/default_group/manifest

curl \
  -H "Authorization: Bearer secret" \
  "http://localhost:8080/tool-groups/default_group/manifest?format=openai"

curl -X POST \
  -H "Authorization: Bearer secret" \
  -H "Content-Type: application/json" \
  -d '{"message":"hello from saiao"}' \
  http://localhost:8080/tool-groups/default_group/actions/echo_message/invoke
```

## Error codes

* `invalid_input`
* `unauthorized`
* `not_found`
* `execution_failed`
* `timeout`
* `internal_error`
