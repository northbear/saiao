# Python Integration

This guide shows a compact Python integration pattern for SAIAO:

- wrap SAIAO HTTP calls in a small client
- fetch the OpenAI-formatted manifest from SAIAO
- pass manifest tools directly to an OpenAI Responses API request
- execute returned tool calls through SAIAO

OpenAI recommends the Responses API for new agentic projects, and its function-calling flow maps cleanly onto SAIAO tool execution.

## Install

```bash
pip install requests openai
```

## Token Source

The SAIAO client token is the bearer token configured for the target tool group.

For example, if the SAIAO config contains:

```yaml
tool_groups:
  - name: default_group
    access_token_env: SAIAO_TOKEN_DEFAULT
```

then the SAIAO server reads the token value from `SAIAO_TOKEN_DEFAULT` in its own process environment.

Your Python application must be provisioned with that same token value separately. The examples below use a client-side environment variable named `SAIAO_TOOL_GROUP_TOKEN` to make that explicit.

## SAIAO Client

```python
import os
import requests


class SAIAOClient:
    def __init__(self, base_url: str, group_name: str, token: str):
        self.base_url = base_url.rstrip("/")
        self.group_name = group_name
        self.session = requests.Session()
        self.session.headers.update({"Authorization": f"Bearer {token}"})

    def openai_manifest(self) -> dict:
        response = self.session.get(
            f"{self.base_url}/tool-groups/{self.group_name}/manifest",
            params={"format": "openai"},
            timeout=10,
        )
        response.raise_for_status()
        return response.json()

    def invoke(self, action_name: str, payload: dict) -> dict:
        response = self.session.post(
            f"{self.base_url}/tool-groups/{self.group_name}/actions/{action_name}/invoke",
            json=payload,
            timeout=30,
        )
        response.raise_for_status()
        return response.json()
```

## OpenAI Agent Loop

```python
import json
import os
from openai import OpenAI


saiao = SAIAOClient(
    "http://localhost:8080",
    "default_group",
    os.environ["SAIAO_TOOL_GROUP_TOKEN"],
)
manifest = saiao.openai_manifest()
tools = manifest["tools"]

client = OpenAI()

response = client.responses.create(
    model="gpt-5",
    instructions="Use SAIAO tools when they help answer the user request.",
    input="Ask the echo_message tool to repeat: hello from python agent",
    tools=tools,
    parallel_tool_calls=False,
)

tool_outputs = []

for item in response.output:
    if item.type != "function_call":
        continue

    args = json.loads(item.arguments)
    result = saiao.invoke(item.name, args)

    tool_outputs.append(
        {
            "type": "function_call_output",
            "call_id": item.call_id,
            "output": json.dumps(result["result"]),
        }
    )

final_response = client.responses.create(
    model="gpt-5",
    previous_response_id=response.id,
    input=tool_outputs,
)

print(final_response.output_text)
```

## Notes

- keep the SAIAO bearer token server-side
- source the token from deployment config that corresponds to `tool_groups[].access_token_env` in the SAIAO config
- log SAIAO `request_id` values from invoke responses when you need cross-system traceability
- the OpenAI-formatted manifest is requested with `?format=openai`
