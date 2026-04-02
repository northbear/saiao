# LangGraph Integration

This guide shows how to place SAIAO behind LangGraph tools using a small SAIAO client and the OpenAI-backed LangGraph tool loop.

## Why LangGraph Fits Well

LangGraph often needs a stable tool execution boundary:

- the graph decides when to use a tool
- SAIAO decides what the tool is allowed to do
- the executor credentials remain inside SAIAO

That is a good separation between orchestration and controlled action execution.

## Recommended Pattern

Use SAIAO as the execution backend and expose each manifest tool as a LangGraph-compatible callable.

High-level flow:

1. call `GET /tool-groups/{group_name}/manifest`
2. map each manifest tool into a LangGraph tool definition
3. on tool execution, call `POST /tool-groups/{group_name}/actions/{action_name}/invoke`
4. return `result.output` back to the graph

## Python Example Sketch

```python
import requests
from langchain_core.tools import tool

BASE_URL = "http://localhost:8080"
TOKEN = "secret"

The `TOKEN` value should come from the token configured for the target tool group in the SAIAO config. For example, if the group uses `access_token_env: SAIAO_TOKEN_DEFAULT`, SAIAO reads that token from its own environment and the LangGraph-side application must be provisioned with the same token value separately.

class SAIAOClient:
    def __init__(self, base_url: str, group_name: str, token: str):
        self.base_url = base_url.rstrip("/")
        self.group_name = group_name
        self.token = token

    def invoke(self, action_name: str, payload: dict) -> str:
        response = requests.post(
            f"{self.base_url}/tool-groups/{self.group_name}/actions/{action_name}/invoke",
            headers={"Authorization": f"Bearer {self.token}"},
            json=payload,
            timeout=30,
        )
        response.raise_for_status()
        body = response.json()
        return body["result"]["output"]

saiao = SAIAOClient(BASE_URL, "default_group", TOKEN)

@tool
def echo_message(message: str) -> str:
    """Echo a validated message through SAIAO."""
    return saiao.invoke("echo_message", {"message": message})
```

You can then register these tools in your LangGraph node structure exactly as you would any other tool.

If your LangGraph application already uses OpenAI chat or Responses-based models, SAIAO simply becomes the execution backend behind the tool functions that the graph exposes.

## Operational Notes

- cache the manifest rather than fetching it before every tool call
- keep SAIAO bearer tokens in server-side config, not prompt state
- propagate `request_id` into graph execution logs where possible
- let SAIAO remain the system of record for action definitions and input validation
