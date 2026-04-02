# Go Integration

This guide shows a compact Go pattern:

- wrap SAIAO in a reusable client
- fetch the OpenAI-formatted manifest once
- run a small tool loop

The SAIAO wrapper keeps the OpenAI-facing code short.

## Token Source

The SAIAO client token is the bearer token configured for the target tool group.

If the SAIAO config contains:

```yaml
tool_groups:
  - name: default_group
    access_token_env: SAIAO_TOKEN_DEFAULT
```

then the SAIAO server reads the token from `SAIAO_TOKEN_DEFAULT` in its own environment.

Your Go application should receive the same token value through its own configuration channel. The example below uses `os.Getenv("SAIAO_TOOL_GROUP_TOKEN")`.

## SAIAO Client

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type SAIAOClient struct {
	BaseURL string
	Group   string
	Token   string
	Client  *http.Client
}

func NewSAIAOClient(baseURL, group, token string) *SAIAOClient {
	return &SAIAOClient{
		BaseURL: baseURL,
		Group:   group,
		Token:   token,
		Client:  &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *SAIAOClient) OpenAIManifest() (map[string]any, error) {
	req, _ := http.NewRequest(http.MethodGet, c.BaseURL+"/tool-groups/"+c.Group+"/manifest?format=openai", nil)
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("manifest failed: %v", body)
	}
	return body, nil
}

func (c *SAIAOClient) Invoke(action string, payload map[string]any) (map[string]any, error) {
	raw, _ := json.Marshal(payload)
	req, _ := http.NewRequest(
		http.MethodPost,
		c.BaseURL+"/tool-groups/"+c.Group+"/actions/"+action+"/invoke",
		bytes.NewReader(raw),
	)
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("invoke failed: %v", body)
	}
	return body, nil
}
```

## OpenAI Agent Request

For Go, the smallest stable example is a direct HTTP call to the OpenAI Responses API.

```go
func OpenAIRequest(apiKey string, body map[string]any) (map[string]any, error) {
	raw, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/responses", bytes.NewReader(raw))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("openai request failed: %v", out)
	}
	return out, nil
}
```

## End-To-End Shape

```go
saiao := NewSAIAOClient("http://localhost:8080", "default_group", os.Getenv("SAIAO_TOOL_GROUP_TOKEN"))
manifest, _ := saiao.OpenAIManifest()
tools := manifest["tools"]

first, _ := OpenAIRequest(os.Getenv("OPENAI_API_KEY"), map[string]any{
	"model":               "gpt-5",
	"instructions":        "Use SAIAO tools when they help complete the request.",
	"input":               "Use the echo_message tool to repeat: hello from go agent",
	"tools":               tools,
	"parallel_tool_calls": false,
})

var toolOutputs []map[string]any
for _, raw := range first["output"].([]any) {
	item := raw.(map[string]any)
	if item["type"] != "function_call" {
		continue
	}

	var args map[string]any
	_ = json.Unmarshal([]byte(item["arguments"].(string)), &args)
	result, _ := saiao.Invoke(item["name"].(string), args)

	toolOutputs = append(toolOutputs, map[string]any{
		"type":    "function_call_output",
		"call_id": item["call_id"],
		"output":  result["result"],
	})
}
```

Then send `toolOutputs` back to `POST /v1/responses` with `previous_response_id` from the first response to get the model’s final answer.

## Notes

- this keeps SAIAO access in one reusable client
- source the token from deployment config that corresponds to `tool_groups[].access_token_env` in the SAIAO config
- the OpenAI-formatted manifest is requested with `?format=openai`
- forward SAIAO `request_id` into your own logs when you care about cross-service tracing
