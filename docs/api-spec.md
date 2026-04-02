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

### Success response

```json
{
  "request_id": "4d3f9d1f7e2a4f2e9f0f0ccf7f4f1c57",
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

## Error codes

* `invalid_input`
* `unauthorized`
* `not_found`
* `execution_failed`
* `timeout`
* `internal_error`
