# SAIAO — Simple AI Action Orchestrator

SAIAO is a simple application that implements the controlled execution layer approach for AI systems.

It is a minimal, container-ready service that demonstrates how AI applications can discover and invoke predefined actions through secured tool groups, while execution credentials remain internal to the service.

SAIAO is designed as a practical reference for building controlled AI-to-action integrations with clear boundaries, predictable behavior, and minimal operational complexity.

## Features

SAIAO is a minimal execution layer that allows AI applications to safely trigger predefined actions such as HTTP calls, SSH commands, local shell operations, and email notifications.

It provides:
- controlled execution
- strict boundaries
- no credential exposure

## Why SAIAO

Most AI systems stop at generating answers.
SAIAO enables AI to **act**, safely and predictably.

## Core Concepts

- Identity — execution credentials (internal only)
- Action — predefined operation
- Tool Group — allowed actions set
- Access Token — AI application authorization

## Quick Start

```bash
docker run -p 8080:8080 \
  -v $(pwd)/config:/etc/saiao \
  -e SAIAO_TOKEN_OPS_BASIC=secret \
  saiao:latest
```

## Example

Fetch manifest:

```bash
curl -H "Authorization: Bearer secret" \
  http://localhost:8080/tool-groups/ops_basic/manifest
```

Invoke action:

```bash
curl -X POST \
  -H "Authorization: Bearer secret" \
  -H "Content-Type: application/json" \
  -d '{"subject":"Test","message":"Hello"}' \
  http://localhost:8080/tool-groups/ops_basic/actions/send_alert_email/invoke
```
