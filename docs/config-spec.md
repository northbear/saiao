# docs/config-spec.md

## Purpose

This document defines the YAML configuration contract for SAIAO.

It specifies:
- supported top-level sections
- field meanings
- required and optional fields
- secret source conventions
- startup validation behavior

SAIAO is designed to run primarily as a containerized application.
Therefore, configuration is split between:
- YAML files mounted into the container
- secret values injected through environment variables or mounted secret files

The configuration contract is intentionally explicit.
SAIAO prefers clarity over magical behavior.

## Configuration File Model

SAIAO reads one primary YAML configuration file at startup.

The file should usually be mounted into the container, for example under:

```text
/etc/saiao/config.yaml
```

The configuration file contains non-sensitive structure and metadata, while sensitive values should be supplied indirectly.

### Minimal runnable example

```yaml
server:
  listen: ":8080"

logging:
  format: json
  level: info

identities:
  - name: local_shell
    enabled: true

actions:
  - name: echo_message
    type: shell
    identity: local_shell
    description: Echo a validated message
    command_template: "printf '%s' '{{message}}'"
    input_schema:
      type: object
      properties:
        message:
          type: string
      required:
        - message

tool_groups:
  - name: default_group
    access_token_env: SAIAO_TOKEN_DEFAULT
    actions:
      - echo_message
```

This example is also available as [`configs/example.yaml`](/home/space/devel/aikvn/saiao/configs/example.yaml).

## Top-Level Sections

The configuration file may contain the following top-level sections:

- `server`
- `logging`
- `identities`
- `actions`
- `tool_groups`

### Required top-level sections for MVP

The following sections are required:
- `identities`
- `actions`
- `tool_groups`

The following sections are optional:
- `server`
- `logging`

---

## Secret Source Guidance

### General rule

Sensitive values must not be hardcoded in committed YAML files.

Sensitive values should be provided through:
- environment variable references
- mounted secret files

Non-sensitive values may be stored directly in YAML.

### What usually belongs in YAML

Typical non-sensitive values:
- object names
- usernames
- email addresses
- display names
- domains or realms
- hostnames
- ports
- action descriptions
- URLs that are not secret
- action membership in tool groups

### What usually belongs in environment variables or files

Typical sensitive values:
- access tokens
- API tokens
- passwords
- SMTP passwords
- SSH private keys
- optional SSH passphrases

---

## Secret Reference Naming Convention

SAIAO uses explicit suffix-based secret references.

### `_env` fields
A field ending with `_env` contains the **name of an environment variable**.

Example:

```yaml
http_token_env: GITHUB_TOKEN
```

This means SAIAO must read the value of the environment variable `GITHUB_TOKEN` during startup.

### `_file` fields
A field ending with `_file` contains a **file path**.

Example:

```yaml
ssh_private_key_file: /run/secrets/saiao_ssh_key
```

This means SAIAO must read the secret value from that file during startup.

### Resolution rules
During startup:
- `_env` fields are resolved from process environment variables
- `_file` fields are resolved from file contents

### Failure rules
Startup must fail fast if:
- a referenced environment variable does not exist
- a referenced file does not exist
- a referenced file cannot be read
- a required resolved value is empty unless explicitly allowed by the field definition

### Inline secrets
Direct inline secrets in YAML are discouraged for v1 and should generally be treated as invalid for sensitive fields.

Example of discouraged pattern:

```yaml
http_token: "hardcoded-secret"
```

### Source exclusivity
If both `_env` and `_file` variants are defined for the same logical secret field, configuration validation should fail unless the field explicitly allows multiple sources.

For MVP, the recommended rule is:
- one secret field → one source only

---

## `server` Section

The `server` section defines HTTP server runtime settings.

### Fields

#### `listen`
- Type: string
- Required: no
- Default: `:8080`
- Meaning: bind address and port for the SAIAO HTTP server

#### `read_timeout_seconds`
- Type: integer
- Required: no
- Default: implementation-defined reasonable default
- Meaning: maximum request read timeout

#### `write_timeout_seconds`
- Type: integer
- Required: no
- Default: implementation-defined reasonable default
- Meaning: maximum response write timeout

#### `shutdown_timeout_seconds`
- Type: integer
- Required: no
- Default: implementation-defined reasonable default
- Meaning: graceful shutdown timeout

Example:

```yaml
server:
  listen: ":8080"
  read_timeout_seconds: 10
  write_timeout_seconds: 30
  shutdown_timeout_seconds: 10
```

---

## `logging` Section

The `logging` section defines logging behavior.

SAIAO is container-first, so default logging should go to stdout/stderr.

### Fields

#### `format`
- Type: string
- Required: no
- Allowed values: `json`, `text`
- Default: `json`

#### `level`
- Type: string
- Required: no
- Allowed values: `debug`, `info`, `warn`, `error`
- Default: `info`

Example:

```yaml
logging:
  format: json
  level: info
```

---

## Template Context

SAIAO uses Go template syntax for all templated fields.

At runtime, templates are evaluated with the following context:

### Available variables

- `identity.*` — resolved identity fields and secrets
- input fields — provided at invocation (top-level)

### Identity secret resolution

Fields defined with `_env` or `_file` are resolved at startup and exposed without suffix.

Example:

```yaml
http_token_env: GITHUB_TOKEN
```

Becomes available in templates as:

```text
{{identity.http_token}}
```

### Example usage

```yaml
Authorization: Bearer {{identity.http_token}}
```

### Notes

* No nested `identity.secrets.*` structure is exposed
* Only explicitly defined fields are available
* Missing variables must cause execution failure

---

## `identities` Section

The `identities` section defines internal execution principals.

An identity is a single structured principal object that may own several secret values and account-like attributes.

Each identity must have a unique `name`.

### Identity Fields

#### `name`
- Type: string
- Required: yes
- Meaning: unique internal identity name referenced by actions

#### `username`
- Type: string
- Required: no
- Meaning: login or account name

#### `email`
- Type: string
- Required: no
- Meaning: email address associated with the identity

#### `display_name`
- Type: string
- Required: no
- Meaning: human-readable label

#### `principal`
- Type: string
- Required: no
- Meaning: enterprise-style principal or service account name

#### `domain`
- Type: string
- Required: no
- Meaning: domain or realm context

#### `groups`
- Type: array of strings
- Required: no
- Meaning: descriptive organizational groups or labels

#### `enabled`
- Type: boolean
- Required: no
- Default: `true`
- Meaning: whether this identity is enabled for use

#### `secrets`
- Type: object
- Required: no, but usually needed for executable actions
- Meaning: secret references and secret-adjacent connection values used by actions

### `identities[].secrets` Fields

The `secrets` object may contain fields such as:

#### HTTP/API-related
- `http_token_env`
- `http_token_file`
- `basic_auth_password_env`
- `basic_auth_password_file`

#### Unix-related
- `unix_password_env`
- `unix_password_file`

#### SSH-related
- `ssh_private_key_file`
- `ssh_private_key_env`
- `ssh_passphrase_env`
- `ssh_passphrase_file`

#### Email-related
- `smtp_host`
- `smtp_port`
- `smtp_username`
- `smtp_password_env`
- `smtp_password_file`
- `from_email`

### Identity Validation Rules

- `name` must be unique across all identities
- disabled identities must not be used by enabled actions
- action references to unknown identities are invalid
- identity secrets must follow `_env` / `_file` rules for sensitive fields
- if both `_env` and `_file` are set for the same logical secret, validation should fail

### Example

```yaml
identities:
  - name: ops_service_account
    username: svc_ops
    email: svc_ops@example.com
    display_name: Operations Service Account
    principal: svc_ops@example.com
    domain: example.com
    groups:
      - ops
      - automation
    enabled: true
    secrets:
      http_token_env: GITHUB_TOKEN
      unix_password_env: OPS_UNIX_PASSWORD
      ssh_private_key_file: /run/secrets/saiao_ssh_key
      smtp_host: smtp.example.com
      smtp_port: 587
      smtp_username: svc_ops@example.com
      smtp_password_env: SMTP_PASSWORD
      from_email: svc_ops@example.com
```

---

## `actions` Section

The `actions` section defines predefined executable operations.

Each action references exactly one identity in MVP.
Each action must have a unique `name`.

### Common Action Fields

#### `name`
- Type: string
- Required: yes
- Meaning: globally unique action name

#### `type`
- Type: string
- Required: yes
- Allowed values: `http`, `ssh`, `shell`, `email`

#### `identity`
- Type: string
- Required: yes
- Meaning: name of the referenced identity

#### `description`
- Type: string
- Required: no
- Meaning: human-readable description used in docs/manifests

#### `enabled`
- Type: boolean
- Required: no
- Default: `true`

#### `timeout_seconds`
- Type: integer
- Required: no
- Meaning: maximum execution time for the action

#### `input_schema`
- Type: object
- Required: no
- Meaning: limited JSON-schema-like input validation contract

### Type-Specific Fields

#### HTTP action
Typical fields:
- `method`
- `url`
- `headers`
- `body_template`

Required minimum:
- `method`
- `url`

#### SSH action
Typical fields:
- `host`
- `port`
- `command_template`
- `working_directory`

Required minimum:
- `host`
- `command_template`

#### Shell action
Typical fields:
- `command_template`
- `working_directory`
- `environment`

Required minimum:
- `command_template`

#### Email action
Typical fields:
- `to`
- `cc`
- `bcc`
- `subject_template`
- `body_template`

Required minimum:
- `to`
- `subject_template`
- `body_template`

### Action Validation Rules

- `name` must be unique across all actions
- `type` must be supported
- `identity` must reference an existing enabled identity
- disabled actions must not be exposed in manifests or invocation
- action type-specific required fields must be present
- `input_schema` must be structurally valid if defined
- command templates must not rely on arbitrary unchecked user input

### Example

```yaml
actions:
  - name: create_github_issue
    type: http
    identity: ops_service_account
    description: Create GitHub issue
    method: POST
    url: https://api.github.com/repos/org/repo/issues
    headers:
      Accept: application/vnd.github+json
      Authorization: Bearer {{identity.http_token}}
    body_template: |
      {
        "title": "{{title}}",
        "body": "{{body}}"
      }
    input_schema:
      type: object
      properties:
        title:
          type: string
        body:
          type: string
      required: [title]

  - name: restart_service_remote
    type: ssh
    identity: ops_service_account
    description: Restart an allowed remote service
    host: app01.example.com
    command_template: sudo systemctl restart {{service}}
    input_schema:
      type: object
      properties:
        service:
          type: string
          enum: [nginx, postgresql, redis]
      required: [service]
```

---

## `tool_groups` Section

The `tool_groups` section defines the sets of actions exposed to AI applications.

Each tool group must have a unique `name` and exactly one access token source in MVP.

### Tool Group Fields

#### `name`
- Type: string
- Required: yes
- Meaning: unique tool group name

#### `description`
- Type: string
- Required: no
- Meaning: human-readable description

#### `enabled`
- Type: boolean
- Required: no
- Default: `true`

#### `access_token_env`
- Type: string
- Required: yes for MVP unless a future alternate token source is added
- Meaning: name of environment variable containing the bearer token for this tool group

#### `actions`
- Type: array of strings
- Required: yes
- Meaning: list of action names available in the tool group

### Tool Group Validation Rules

- `name` must be unique across all tool groups
- `access_token_env` must be present in MVP
- referenced environment variable must exist at startup
- `actions` must not be empty
- every listed action must exist and be enabled
- disabled tool groups must not be exposed through runtime APIs

### Example

```yaml
tool_groups:
  - name: ops_basic
    description: Basic operations tool set
    enabled: true
    access_token_env: SAIAO_TOKEN_OPS_BASIC
    actions:
      - restart_service_remote
      - create_github_issue
```

---

## Input Schema Subset

SAIAO MVP supports a limited schema subset for `input_schema`.

Supported concepts:
- `type: object`
- `properties`
- `required`
- scalar property types such as `string`, `integer`, `boolean`
- `enum`

Unsupported or out of scope for MVP unless explicitly implemented:
- deep schema composition
- custom validators
- advanced array constraints
- conditional schema logic

---

## Startup Validation Behavior

SAIAO should validate the full configuration at startup and fail fast on invalid configuration.

### Validation should include at least:
- YAML syntax validity
- presence of required sections
- uniqueness of identity/action/tool group names
- valid cross-references
- presence of required env vars
- readability of required files
- action type required fields
- tool group token availability
- structural validity of input schemas

### Startup failure examples
- missing `SAIAO_TOKEN_OPS_BASIC`
- unknown identity referenced by an action
- unknown action referenced by a tool group
- unreadable SSH private key file
- duplicate action names

---

## Container Runtime Guidance

SAIAO is intended to run as a containerized application.

Recommended runtime model:
- mount configuration file or config directory into the container
- inject secret values through environment variables
- mount secret files such as SSH keys into the container filesystem
- log to stdout/stderr

Typical runtime layout:
- config file: `/etc/saiao/config.yaml`
- mounted secret file: `/run/secrets/saiao_ssh_key`

Example Docker run pattern:

```bash
docker run -p 8080:8080 \
  -v $(pwd)/config:/etc/saiao \
  -v $(pwd)/secrets:/run/secrets \
  -e SAIAO_TOKEN_OPS_BASIC=secret \
  -e GITHUB_TOKEN=ghp_xxx \
  -e SMTP_PASSWORD=xxx \
  saiao:latest
```

---

## Full Example Configuration

```yaml
server:
  listen: ":8080"
  read_timeout_seconds: 10
  write_timeout_seconds: 30
  shutdown_timeout_seconds: 10

logging:
  format: json
  level: info

identities:
  - name: ops_service_account
    username: svc_ops
    email: svc_ops@example.com
    display_name: Operations Service Account
    principal: svc_ops@example.com
    domain: example.com
    groups:
      - ops
      - automation
    enabled: true
    secrets:
      http_token_env: GITHUB_TOKEN
      unix_password_env: OPS_UNIX_PASSWORD
      ssh_private_key_file: /run/secrets/saiao_ssh_key
      smtp_host: smtp.example.com
      smtp_port: 587
      smtp_username: svc_ops@example.com
      smtp_password_env: SMTP_PASSWORD
      from_email: svc_ops@example.com

actions:
  - name: create_github_issue
    type: http
    identity: ops_service_account
    description: Create GitHub issue
    method: POST
    url: https://api.github.com/repos/org/repo/issues
    headers:
      Accept: application/vnd.github+json
      Authorization: Bearer {{identity.http_token}}
    body_template: |
      {
        "title": "{{title}}",
        "body": "{{body}}"
      }
    input_schema:
      type: object
      properties:
        title:
          type: string
        body:
          type: string
      required: [title]

  - name: restart_service_remote
    type: ssh
    identity: ops_service_account
    description: Restart an allowed remote service
    host: app01.example.com
    command_template: sudo systemctl restart {{service}}
    input_schema:
      type: object
      properties:
        service:
          type: string
          enum: [nginx, postgresql, redis]
      required: [service]

  - name: check_disk_local
    type: shell
    identity: ops_service_account
    description: Check root filesystem usage
    command_template: df -h /

  - name: send_alert_email
    type: email
    identity: ops_service_account
    description: Send operations alert email
    to:
      - ops@example.com
    subject_template: "Alert: {{subject}}"
    body_template: |
      {{message}}
    input_schema:
      type: object
      properties:
        subject:
          type: string
        message:
          type: string
      required: [subject, message]

tool_groups:
  - name: ops_basic
    description: Basic operations tool set
    enabled: true
    access_token_env: SAIAO_TOKEN_OPS_BASIC
    actions:
      - restart_service_remote
      - create_github_issue
      - check_disk_local
      - send_alert_email
```

---

## Summary

The SAIAO configuration contract is designed to be:
- explicit
- validation-friendly
- container-ready
- safe for open-source operational use

It intentionally separates:
- declarative YAML structure
- sensitive runtime secrets
- executable action definitions
