package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"saiao/internal/auth"
	"saiao/internal/logging"
	"saiao/internal/models"
)

func testConfig() *models.Config {
	return &models.Config{
		Identities: []models.Identity{
			{
				Name:    "ops",
				Secrets: map[string]string{},
			},
		},
		Actions: []models.Action{
			{
				Name:            "echo",
				Type:            "shell",
				Identity:        "ops",
				Description:     "Echo a validated message",
				CommandTemplate: "printf 'hello %s' '{{message}}'",
				InputSchema: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"message": map[string]any{
							"type": "string",
						},
					},
					"required": []any{"message"},
				},
			},
			{
				Name:            "other",
				Type:            "shell",
				Identity:        "ops",
				Description:     "A second action in another group",
				CommandTemplate: "printf other",
			},
		},
		ToolGroups: []models.ToolGroup{
			{
				Name:           "ops_group",
				AccessTokenEnv: "OPS_GROUP_TOKEN",
				Actions:        []string{"echo"},
			},
			{
				Name:           "other_group",
				AccessTokenEnv: "OTHER_GROUP_TOKEN",
				Actions:        []string{"other"},
			},
		},
	}
}

func TestNewServerSetsAddress(t *testing.T) {
	srv := NewServer(":1234", BuildInfo{
		Service: "saiao",
		Version: "1.0.0",
		Commit:  "abc1234",
	}, testConfig(), auth.NewTokenStore(), slog.Default())
	if srv.Addr != ":1234" {
		t.Fatalf("expected addr :1234, got %s", srv.Addr)
	}
	if srv.Handler == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestInfoEndpointReturnsInfoJSON(t *testing.T) {
	srv := NewServer(":0", BuildInfo{
		Service: "saiao",
		Version: "1.0.0",
		Commit:  "abc1234",
	}, testConfig(), auth.NewTokenStore(), slog.Default())

	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	rr := httptest.NewRecorder()

	srv.Handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var got map[string]any
	body, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}

	if got["status"] != "ok" {
		t.Fatalf("expected status ok, got %v", got["status"])
	}
	if got["request_id"] == "" {
		t.Fatalf("expected request_id, got %v", got["request_id"])
	}
	if rr.Header().Get(requestIDHeader) != got["request_id"] {
		t.Fatalf("expected header request id %v, got %q", got["request_id"], rr.Header().Get(requestIDHeader))
	}
	if got["service"] != "saiao" {
		t.Fatalf("expected service saiao, got %v", got["service"])
	}
	if got["version"] != "1.0.0" {
		t.Fatalf("expected version 1.0.0, got %v", got["version"])
	}
	if got["commit"] != "abc1234" {
		t.Fatalf("expected commit abc1234, got %v", got["commit"])
	}
}

func TestManifestEndpointRequiresMatchingToken(t *testing.T) {
	store := auth.NewTokenStore()
	store.Register("secret", "ops_group")

	srv := NewServer(":0", BuildInfo{Service: "saiao"}, testConfig(), store, slog.Default())

	req := httptest.NewRequest(http.MethodGet, "/tool-groups/ops_group/manifest", nil)
	req.Header.Set("Authorization", "Bearer secret")
	req.Header.Set(requestIDHeader, "req-manifest-1")
	rr := httptest.NewRecorder()

	srv.Handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var response models.ManifestResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(response.Tools) != 1 || response.Tools[0].Name != "echo" {
		t.Fatalf("unexpected manifest: %#v", response.Tools)
	}
	if response.RequestID != "req-manifest-1" {
		t.Fatalf("expected request id req-manifest-1, got %q", response.RequestID)
	}
	if rr.Header().Get(requestIDHeader) != "req-manifest-1" {
		t.Fatalf("expected response header request id req-manifest-1, got %q", rr.Header().Get(requestIDHeader))
	}
}

func TestManifestEndpointRejectsWrongToolGroupToken(t *testing.T) {
	store := auth.NewTokenStore()
	store.Register("secret", "other_group")

	srv := NewServer(":0", BuildInfo{Service: "saiao"}, testConfig(), store, slog.Default())

	req := httptest.NewRequest(http.MethodGet, "/tool-groups/ops_group/manifest", nil)
	req.Header.Set("Authorization", "Bearer secret")
	rr := httptest.NewRecorder()

	srv.Handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}

	var response models.ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.RequestID == "" {
		t.Fatal("expected request id on error response")
	}
}

func TestInvokeEndpointExecutesAction(t *testing.T) {
	store := auth.NewTokenStore()
	store.Register("secret", "ops_group")

	srv := NewServer(":0", BuildInfo{Service: "saiao"}, testConfig(), store, slog.Default())

	req := httptest.NewRequest(http.MethodPost, "/tool-groups/ops_group/actions/echo/invoke", bytes.NewBufferString(`{"message":"world"}`))
	req.Header.Set("Authorization", "Bearer secret")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(requestIDHeader, "req-invoke-1")
	rr := httptest.NewRecorder()

	srv.Handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
	}

	var response models.SuccessResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Result.Output != "hello world" {
		t.Fatalf("unexpected output %q", response.Result.Output)
	}
	if response.RequestID != "req-invoke-1" {
		t.Fatalf("expected request id req-invoke-1, got %q", response.RequestID)
	}
}

func TestInvokeEndpointReturnsInvalidInput(t *testing.T) {
	store := auth.NewTokenStore()
	store.Register("secret", "ops_group")

	srv := NewServer(":0", BuildInfo{Service: "saiao"}, testConfig(), store, slog.Default())

	req := httptest.NewRequest(http.MethodPost, "/tool-groups/ops_group/actions/echo/invoke", bytes.NewBufferString(`{}`))
	req.Header.Set("Authorization", "Bearer secret")
	rr := httptest.NewRecorder()

	srv.Handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}

	var response models.ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Error.Code != "invalid_input" {
		t.Fatalf("expected invalid_input, got %q", response.Error.Code)
	}
	if response.RequestID == "" {
		t.Fatal("expected request id on error response")
	}
}

func TestInvokeLogsShareRequestIDAcrossAPIAndLifecycle(t *testing.T) {
	store := auth.NewTokenStore()
	store.Register("secret", "ops_group")

	var buf bytes.Buffer
	logger := logging.NewWithWriter(&buf, "json", "info")
	srv := NewServer(":0", BuildInfo{Service: "saiao"}, testConfig(), store, logger)

	req := httptest.NewRequest(http.MethodPost, "/tool-groups/ops_group/actions/echo/invoke", bytes.NewBufferString(`{"message":"world"}`))
	req.Header.Set("Authorization", "Bearer secret")
	req.Header.Set(requestIDHeader, "req-log-1")
	rr := httptest.NewRecorder()

	srv.Handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	logOutput := buf.String()
	if !strings.Contains(logOutput, `"msg":"invoke_request"`) {
		t.Fatalf("expected invoke_request log, got %s", logOutput)
	}
	if !strings.Contains(logOutput, `"msg":"invoke_started"`) {
		t.Fatalf("expected invoke_started log, got %s", logOutput)
	}
	if !strings.Contains(logOutput, `"msg":"invoke_finished"`) {
		t.Fatalf("expected invoke_finished log, got %s", logOutput)
	}
	if count := strings.Count(logOutput, `"request_id":"req-log-1"`); count < 3 {
		t.Fatalf("expected request_id to appear in correlated logs, got %d entries in %s", count, logOutput)
	}
}
