# JavaScript And TypeScript Integration

This guide shows a structured Node.js or TypeScript integration:

- a small `SAIAOClient` wrapper
- fetching the OpenAI-formatted manifest
- an OpenAI Responses API tool loop

OpenAI recommends the Responses API for new agentic applications, and SAIAO fits naturally as a function-calling backend.

## Install

```bash
npm install openai
```

## Token Source

The SAIAO client token is the bearer token configured for the target tool group.

If the SAIAO config contains:

```yaml
tool_groups:
  - name: default_group
    access_token_env: SAIAO_TOKEN_DEFAULT
```

then the SAIAO server reads the token from `SAIAO_TOKEN_DEFAULT` in its own environment.

Your JavaScript or TypeScript application should receive the same token value through its own secure configuration path. The example below uses `process.env.SAIAO_TOOL_GROUP_TOKEN`.

## SAIAO Client

```ts
const saiaoToken = process.env.SAIAO_TOOL_GROUP_TOKEN!;

class SAIAOClient {
  constructor(
    private readonly baseUrl: string,
    private readonly groupName: string,
    private readonly token: string,
  ) {}

  async openAIManifest() {
    const response = await fetch(`${this.baseUrl}/tool-groups/${this.groupName}/manifest?format=openai`, {
      headers: { Authorization: `Bearer ${this.token}` },
    });

    if (!response.ok) throw new Error(`manifest failed: ${response.status}`);
    return response.json();
  }

  async invoke(actionName: string, payload: Record<string, unknown>) {
    const response = await fetch(
      `${this.baseUrl}/tool-groups/${this.groupName}/actions/${actionName}/invoke`,
      {
        method: "POST",
        headers: {
          Authorization: `Bearer ${this.token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      },
    );

    const body = await response.json();
    if (!response.ok) {
      throw new Error(`${body.error.code}: ${body.error.message}`);
    }

    return body;
  }
}
```

## OpenAI Agent Loop

```ts
import OpenAI from "openai";

const saiao = new SAIAOClient("http://localhost:8080", "default_group", saiaoToken);
const manifest = await saiao.openAIManifest();
const tools = manifest.tools;

const openai = new OpenAI();

const first = await openai.responses.create({
  model: "gpt-5",
  instructions: "Use SAIAO tools when they help complete the request.",
  input: "Use the echo_message tool to repeat: hello from javascript agent",
  tools,
  parallel_tool_calls: false,
});

const toolOutputs = [];

for (const item of first.output) {
  if (item.type !== "function_call") continue;

  const args = JSON.parse(item.arguments);
  const result = await saiao.invoke(item.name, args);

  toolOutputs.push({
    type: "function_call_output",
    call_id: item.call_id,
    output: JSON.stringify(result.result),
  });
}

const final = await openai.responses.create({
  model: "gpt-5",
  previous_response_id: first.id,
  input: toolOutputs,
});

console.log(final.output_text);
```

## Notes

- keep SAIAO tokens in backend configuration, not untrusted browser code
- source the token from deployment config that corresponds to `tool_groups[].access_token_env` in the SAIAO config
- capture SAIAO `request_id` in your application logs
- the OpenAI-formatted manifest is requested with `?format=openai`
