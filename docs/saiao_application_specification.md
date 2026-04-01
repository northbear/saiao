# SAIAO Application Specification

## 1. Purpose

**SAIAO** (Simple AI Action Orchestrator) is a minimal application that exposes predefined executable actions to AI applications through secured tool groups.

Its purpose is to provide a small, practical, and auditable execution layer between an AI application and real-world operations such as HTTP calls, remote SSH commands, local shell commands, and email notifications.

SAIAO is intentionally limited in scope:
- it does not let models invent arbitrary integrations
- it does not expose execution credentials to AI applications or models
- it does not attempt to be a full workflow engine, secrets manager, or IAM platform

The system is designed to be easy to configure, easy to run, and ready for practical use without heavy integration work.

---

## 2. Design Goals

### 2.1 Primary goals
- Provide a **minimal execution layer** for AI applications
- Expose actions in an **LLM-friendly tools format**
- Support **controlled invocation** of predefined actions
- Separate **client access authorization** from **execution credentials**
- Keep configuration simple and human-readable using **YAML**
- Be implementation-friendly for a small standalone open-source project

### 2.2 Non-goals for v1
- Arbitrary command execution
- Dynamic action creation through API
- Multi-user RBAC
- Built-in secrets vault
- Workflow DAG engine
- Plugin marketplace
- Full policy engine
- OIDC/JWT identity federation

---

## 3. Core Concepts

SAIAO defines four core objects.

### 3.1 Identity
An **Identity** is a named execution principal used internally by SAIAO.

In SAIAO v1, an identity is modeled as a **single account-like object** that may own several secret values and connection-related attributes.

The identity should be represented in a form that resembles how technical identities are commonly received from enterprise authentication and directory systems such as:
- Active Directory
- LDAP
- IdP / organizational authentication services

That means an identity is not just a raw token or a single credential file. It is a structured principal-like object that may include fields such as:
- username / login name
- email address
- display name
- principal name / account name
- domain or realm
- optional groups or labels
- associated secrets and credential material

An identity may contain multiple secret kinds at once, for example:
- API token
- Unix password
- SSH private key
- SMTP password

Identities are not exposed to AI applications.

### 3.2 Action
An **Action** is a predefined executable operation.

Each action:
- has a unique name
- has a type
- references one identity
- defines its execution configuration
- may define an input schema for validation
- may define a manifest-visible tool description for LLMs

Supported v1 action types:
- `http`
- `ssh`
- `shell`
- `email`

### 3.3 Tool Group
A **Tool Group** is a named collection of actions exposed together to an AI application.

A tool group defines the subset of actions that a specific AI application is allowed to discover and invoke.

### 3.4 Tool Group Access Token
A **Tool Group Access Token** is a bearer token used by an AI application to:
- retrieve the manifest of tools in a specific tool group
- invoke actions belonging to that tool group

This token authorizes access to the AI-facing interface only.
It is distinct from identities used for execution.

---

## 4. Security Model

### 4.1 Separation of concerns
SAIAO separates two kinds of secrets:

1. **Execution credentials**
   - held by a structured internal identity
   - an identity may own several secrets at once
   - used internally by SAIAO to execute actions
   - never exposed to the AI application or model

2. **Client access credentials**
   - held as tool group access tokens
   - used by AI applications to access tool manifests and invoke actions
   - do not grant direct access to execution credentials

### 4.2 Security principles
- AI applications authenticate using a tool group token
- an AI application can only access the specific tool group bound to its token
- identities are internal and never returned by API
- configuration mutation is out of scope for the runtime API
- arbitrary shell and SSH command construction is discouraged and restricted
- all action invocations should be logged

### 4.3 Recommended v1 restrictions
- one static bearer token per tool group
- tokens loaded from environment variables or configuration
- no token introspection API
- no token issuance API
- no config editing API
- no secrets retrieval API

---

## 5. Action Types

## 5.1 HTTP Action
An HTTP action performs a predefined HTTP request.

Typical use cases:
- calling REST APIs
- creating issues/tickets
- triggering webhooks
- querying internal services

Typical fields:
- method
- url
- headers template
- body template
- timeout
- input schema

Execution notes:
- authentication data may be injected from identity
- request body may be rendered from validated input

## 5.2 SSH Action
An SSH action executes a predefined command or script on a remote host over SSH.

Typical use cases:
- restarting a service
- checking system state on a remote server
- running predefined maintenance scripts

Typical fields:
- host
- port
- username or identity-provided user
- command template or script path
- timeout
- input schema

Security notes:
- arbitrary commands should not be accepted from the model
- preferred approach is fixed commands with validated parameters
- better yet, execute predefined scripts with restricted arguments

## 5.3 Shell Action
A shell action executes a predefined local Linux shell command or script.

Typical use cases:
- collecting diagnostics
- running local maintenance tasks
- interacting with local services or files

Typical fields:
- command template or script path
- working directory
- environment variables
- timeout
- input schema

Security notes:
- arbitrary local shell command execution must not be the default mode
- prefer fixed commands or predefined scripts with validated arguments

## 5.4 Email Action
An email action sends an email using predefined SMTP or relay credentials.

Typical use cases:
- alerting
- sending reports
- sending operator notifications

Typical fields:
- SMTP host/port
- sender address
- recipients
- subject template
- body template
- optional CC/BCC
- input schema

---

## 6. Identity Model

## 6.1 General definition
In SAIAO v1, an identity is a **single structured principal object** containing:
- account-style descriptive fields
- optional directory/organization-related attributes
- one or more secret values used for execution

This identity model is intentionally closer to how service accounts or technical users are commonly represented in real organizations through AD, LDAP, or IdP-backed environments.

So instead of treating identity as a narrow per-protocol credential container, SAIAO treats it as a reusable execution principal that may participate in several action types.

Each action references exactly one identity in v1.

## 6.2 Identity shape
A v1 identity should support ordinary account-like fields such as:
- `name` — internal SAIAO identity name
- `username` — login or account name
- `email` — email address associated with the identity
- `display_name` — human-readable label
- `principal` — optional enterprise-style principal name
- `domain` or `realm` — optional domain context
- `groups` — optional labels or group names
- `enabled` — optional administrative flag

These fields are descriptive and do not by themselves expose secrets.

## 6.3 Secret ownership
An identity may own multiple secrets and credential values simultaneously.

Supported v1 credential content may include:

### API credentials
- bearer token
- optional basic auth username/password

### Unix credentials
- username
- optional password

### SSH credentials
- username
- private key path or inline private key reference
- optional passphrase source

Note: public key generation may be supported as a helper feature, but pre-installation of the public key on target systems is assumed.

### Email credentials
- SMTP host
- SMTP port
- username
- password
- from address
- optional TLS settings

The same identity may provide the credential material needed for different actions, for example:
- SSH action using `username` + `ssh_private_key`
- email action using `email` + SMTP credentials
- HTTP action using bearer token

## 6.4 Identity loading
Identities may load values from:
- literal config values for non-sensitive fields
- environment variable references
- file paths

Recommended v1 convention:
- prefer environment variables and secret files for sensitive values
- avoid storing raw secrets directly in committed YAML files

## 6.5 Identity exposure rules
Identity objects are internal to SAIAO.

The runtime API must not expose:
- raw secret values
- full identity objects
- internal connection material

Tool manifests may only expose action-level information needed by the AI application.

---

## 7. Action Definition Model

Each action must define:
- `name`
- `type`
- `identity`
- execution-specific configuration

Optional fields:
- `description`
- `input_schema`
- `timeout`
- `manifest`
- `enabled`

### 7.1 Naming requirements
- action names must be unique globally within a SAIAO instance
- names should be stable and API-safe
- recommended format: lowercase, underscore-separated

### 7.2 Input schema
Actions may define an input schema used to validate invocation input.

The schema should be simple and LLM-friendly, close to JSON Schema concepts:
- object type
- properties
- required fields
- enum values
- basic scalar types

The v1 implementation does not need full JSON Schema compliance.
A limited subset is sufficient.

### 7.3 Templates
Actions may use templates to render:
- HTTP headers
- HTTP request body
- shell commands
- SSH commands
- email subject/body

Templates may reference validated input values.

Templates must not expose raw identity secrets in response payloads.

---

## 8. Tool Group Model

A tool group defines a set of actions that can be discovered and invoked together.

Each tool group must define:
- `name`
- access token source
- list of allowed action names

Optional fields:
- `description`
- `enabled`

### 8.1 Tool group responsibilities
A tool group is responsible for:
- grouping actions for a client use case
- limiting the visible action surface
- limiting the invokable action surface
- associating the group with one access token

### 8.2 v1 access model
In v1:
- one tool group has one static bearer token
- the token grants access only to that group
- one action may belong to multiple tool groups if desired

---

## 9. Runtime API

SAIAO should expose a small HTTP API for AI applications.

## 9.1 Authentication
Authentication is performed using:

`Authorization: Bearer <tool-group-token>`

The token determines which tool group is accessible.

## 9.2 Required endpoints

### 9.2.1 Retrieve tool manifest
`GET /tool-groups/{group_name}/manifest`

Purpose:
- return LLM-friendly tool definitions for the specified group

Behavior:
- validate bearer token
- ensure token matches the requested group
- return only manifest-visible fields

### 9.2.2 Invoke action
`POST /tool-groups/{group_name}/actions/{action_name}/invoke`

Purpose:
- invoke an allowed action with validated input

Behavior:
- validate bearer token
- ensure action belongs to requested group
- validate input payload against action schema
- execute action using referenced identity
- return structured result

## 9.3 Optional endpoints for v1

### Runtime info endpoint
`GET /info`

This endpoint may combine the purposes of health and version reporting in a single response.

Recommended response content:
- service status
- application name
- version
- build or commit identifier if available
- current time or uptime if desired

Example conceptual response:

```json
{
  "status": "ok",
  "service": "saiao",
  "version": "0.1.0",
  "commit": "abc1234"
}
```

A separate health-only endpoint is optional, but a single runtime info endpoint is sufficient for v1.

---

## 10. Tool Manifest Format

The manifest returned to AI applications should be compatible with common LLM tool/function calling formats.

At minimum, each tool manifest should include:
- `name`
- `description`
- `parameters`

Example conceptual format:

```json
{
  "name": "send_alert_email",
  "description": "Send an alert email to operations",
  "parameters": {
    "type": "object",
    "properties": {
      "subject": {"type": "string"},
      "message": {"type": "string"}
    },
    "required": ["subject", "message"]
  }
}
```

The manifest must not contain:
- raw secrets
- internal identity details
- hidden headers
- hidden connection data

---

## 11. Invocation Flow

A typical flow is:

1. AI application obtains or is configured with a tool group access token
2. AI application requests the manifest for its allowed tool group
3. SAIAO returns tool definitions in LLM-friendly format
4. The AI model selects a tool and provides arguments
5. AI application calls the invoke endpoint with those arguments
6. SAIAO validates input
7. SAIAO resolves the referenced action and identity
8. SAIAO executes the action
9. SAIAO returns structured execution result
10. SAIAO writes an audit/log entry

---

## 12. Logging and Audit

v1 should include basic structured logging.

Each action invocation should log at least:
- timestamp
- tool group name
- action name
- invocation status
- execution duration
- high-level caller identity derived from token context if available

The log should avoid writing secrets.

### 12.1 Default runtime logging model
SAIAO is intended to run primarily as a containerized application.

Therefore, the default logging behavior should be:
- write logs to `stdout`
- write fatal/runtime diagnostics to `stderr` when appropriate
- avoid default file-based logging inside the container

This matches ordinary container and Kubernetes operational practice.

Optional alternate logging outputs may be added later, but stdout/stderr should be the default and documented runtime mode.

---

## 13. Error Handling

SAIAO should return structured errors.

Recommended categories:
- authentication error
- authorization error
- unknown tool group
- unknown action
- invalid input
- execution error
- timeout
- internal server error

Error responses should:
- be machine-readable
- avoid leaking secrets or internal sensitive config

---

## 14. Configuration

SAIAO should be configured primarily via YAML.

Configuration should define:
- identities
- actions
- tool groups
- server settings
- logging settings

Sensitive values should support indirection via:
- environment variable references
- file path references

### 14.1 Container-oriented configuration model
SAIAO is intended to run in a container by default.

Therefore the recommended runtime model is:
- secrets are provided through environment variables and/or mounted secret files
- configuration files are provided through a mounted volume or mounted directory
- the application container reads configuration from a known path, for example `/etc/saiao/`

Recommended v1 conventions:
- mount one configuration file or a configuration directory into the container
- inject secret values through environment variables when possible
- allow sensitive material such as SSH private keys to be mounted as files
- avoid baking site-specific configuration or secrets into the container image

The final build artifact should be a ready-to-run Docker/OCI container image suitable for:
- local Docker or Podman execution
- Kubernetes deployment
- simple container-based demos

### 14.2 Example configuration

```yaml
server:
  listen: ":8080"
  read_timeout_seconds: 10
  write_timeout_seconds: 30

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
    secrets:
      http_token_env: GITHUB_TOKEN
      unix_password_env: OPS_UNIX_PASSWORD
      ssh_private_key_file: /etc/saiao/keys/svc_ops_id_ed25519
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
      Authorization: "Bearer {{identity.http_token}}"
      Accept: application/vnd.github+json
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
    access_token_env: SAIAO_TOKEN_OPS_BASIC
    actions:
      - restart_service_remote
      - check_disk_local
      - send_alert_email
```

---

## 15. Execution Safety Requirements

The v1 implementation should apply the following safety rules:

### 15.1 Input validation first
- every invocation input must be validated before rendering templates or executing actions

### 15.2 No arbitrary shell by default
- shell and SSH actions should favor fixed commands or restricted templates
- arguments should be validated using enums or strict typing when possible

### 15.3 Secret isolation
- secrets must never be returned in API responses
- secrets must not be written to logs
- manifests must not expose internal execution fields not needed by the AI application

### 15.4 Timeouts
- every action must run with a bounded timeout

### 15.5 Explicit group scoping
- action invocation must be rejected if the token’s group does not include the requested action

---

## 16. Suggested Internal Architecture

A simple implementation may consist of the following modules:

- **config loader**
- **identity resolver**
- **action registry**
- **tool group registry**
- **manifest generator**
- **auth middleware**
- **input validator**
- **executors**
  - HTTP executor
  - SSH executor
  - shell executor
  - email executor
- **logger/audit module**
- **HTTP server**

This architecture is descriptive, not mandatory.

---

## 17. Minimal MVP Scope

The recommended MVP should include:
- single binary application
- YAML configuration
- environment variable secret support
- mounted-file support for config and key material
- static bearer-token auth per tool group
- manifest endpoint
- invoke endpoint
- combined runtime info endpoint
- four action types: HTTP, SSH, shell, email
- basic input schema validation
- structured logging to stdout
- Docker/OCI image build output

That is sufficient for a first usable open-source release.

---

## 18. Packaging and Distribution

## 18.1 Final build product
The final build product of SAIAO should be a Docker/OCI container image ready to run.

The image should:
- launch the SAIAO server directly
- accept config path and listen settings via flags or environment variables
- work cleanly in Docker, Podman, and Kubernetes-style environments

## 18.2 Runtime assumptions
The runtime model should assume:
- immutable container image
- mounted configuration
- externally injected environment variables
- stdout/stderr logging

This should be reflected in both implementation and documentation.

## 18.3 Documentation requirements
The project documentation should include at least:
- a `README.md`
- quick-start instructions for running the container
- explanation of config mount layout
- explanation of required environment variables
- example YAML configuration files
- example requests to retrieve manifests and invoke actions
- troubleshooting notes for common startup/config errors

## 18.4 Client integration examples
The documentation should provide example client code for interacting with SAIAO from popular AI-related programming languages.

At minimum, examples should be provided for:
- Python
- JavaScript / TypeScript
- Go

Optional additional examples may include:
- Bash / curl
- Java
- C#

The examples should cover at least:
- fetching a tool group manifest
- invoking an action
- handling success and error responses
- passing the bearer token correctly

## 18.5 README scope
The `README.md` should be sufficient for a user to:
- start the container
- mount configuration
- inject secrets
- retrieve a manifest
- invoke a test action

In other words, the README should describe a runnable path, not only architecture.

---

## 19. Future Extensions

Possible future extensions after MVP:
- multiple tokens per tool group
- audit persistence backend
- richer schema validation
- action output schema
- script-based action packaging
- per-action allowlists for hosts/commands
- secret backends
- per-token caller labels
- rate limits
- tool group metadata
- OpenAI-compatible combined tools endpoint

These are explicitly outside the required MVP.

---

## 20. Summary Definition

SAIAO is a minimal AI-facing action execution service that:
- organizes predefined actions into tool groups
- secures AI application access with tool group bearer tokens
- executes actions using internal identities
- supports HTTP, SSH, local shell, and email actions
- exposes tool manifests in an LLM-friendly format
- keeps execution credentials isolated from models and client applications

Its defining property is controlled, practical execution with minimal moving parts.

