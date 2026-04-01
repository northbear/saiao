# Definition

SAIAO is a minimal service that exposes predefined actions to AI systems via secured tool groups.

# Key Properties

- No credential exposure
- Deterministic execution
- Config-driven
- Container-first

# Action Types

- HTTP
- SSH
- Shell
- Email

# Security Model

- identities are internal
- tokens grant group access only
- a bearer token maps to exactly one tool group in MVP
- no secrets in responses

# Execution Result Model

All actions return a normalized result structure.

```json
{
  "status": "success",
  "result": {
    "output": "string",
    "exit_code": 0
  }
}
```

## Behavior by action type

* HTTP:

  * `output` = response body or response summary as string
  * `exit_code` = 0 on successful execution, non-zero on execution failure

* SSH / Shell:

  * `output` = stdout, with stderr optionally appended or logged separately
  * `exit_code` = process exit code

* Email:

  * `output` = "sent" or provider-specific message identifier
  * `exit_code` = 0 on successful execution, non-zero on execution failure

## Failure

On failure:

```json
{
  "error": {
    "code": "execution_failed",
    "message": "description"
  }
}
```

## Error Codes

* `invalid_input`
* `unauthorized`
* `not_found`
* `execution_failed`
* `timeout`
* `internal_error`

# Timeout Behavior

If an action exceeds `timeout_seconds`:

* execution must be terminated
* response must be:

```json
{
  "error": {
    "code": "timeout",
    "message": "execution exceeded timeout"
  }
}
```

`timeout_seconds` is configured per action in the YAML configuration.

## Caller Identity

SAIAO does not authenticate end-users directly.

The caller identity is derived from the tool group associated with the provided access token.

For logging and auditing:

* caller = tool group name

## Identity Selection in MVP

SAIAO supports multiple identities in configuration.

In MVP, each action references exactly one configured identity, and that identity is used as the execution identity for the action.
