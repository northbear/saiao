package invoke

import (
	"bytes"
	stdErrors "errors"
	"strings"
	"testing"

	saiaoerrors "saiao/internal/errors"
	"saiao/internal/logging"
	"saiao/internal/models"
)

func TestExecuteRunsShellAction(t *testing.T) {
	cfg := &models.Config{
		Identities: []models.Identity{
			{Name: "ops", Secrets: map[string]string{}},
		},
		Actions: []models.Action{
			{
				Name:            "echo",
				Type:            "shell",
				Identity:        "ops",
				CommandTemplate: "printf 'hello %s' '{{message}}'",
				InputSchema: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"message": map[string]any{"type": "string"},
					},
					"required": []any{"message"},
				},
			},
		},
		ToolGroups: []models.ToolGroup{
			{Name: "ops_group", Actions: []string{"echo"}},
		},
	}

	got, err := Execute("req-123", "ops_group", "echo", []byte(`{"message":"world"}`), cfg, logging.Discard())
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if got.Result.Output != "hello world" {
		t.Fatalf("expected hello world, got %q", got.Result.Output)
	}
	if got.Result.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", got.Result.ExitCode)
	}
	if got.RequestID != "req-123" {
		t.Fatalf("expected request id req-123, got %q", got.RequestID)
	}
}

func TestExecuteReturnsTimeoutWhenActionExceedsDeadline(t *testing.T) {
	cfg := &models.Config{
		Identities: []models.Identity{
			{Name: "ops", Secrets: map[string]string{}},
		},
		Actions: []models.Action{
			{
				Name:            "slow",
				Type:            "shell",
				Identity:        "ops",
				CommandTemplate: "sleep 2",
				TimeoutSeconds:  1,
			},
		},
		ToolGroups: []models.ToolGroup{
			{Name: "ops_group", Actions: []string{"slow"}},
		},
	}

	_, err := Execute("req-123", "ops_group", "slow", nil, cfg, logging.Discard())
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if !stdErrors.Is(err, saiaoerrors.ErrTimeout) {
		t.Fatalf("expected timeout error, got %v", err)
	}
}

func TestExecuteLogsLifecycle(t *testing.T) {
	cfg := &models.Config{
		Identities: []models.Identity{
			{Name: "ops", Secrets: map[string]string{}},
		},
		Actions: []models.Action{
			{
				Name:            "echo",
				Type:            "shell",
				Identity:        "ops",
				CommandTemplate: "printf hi",
			},
		},
		ToolGroups: []models.ToolGroup{
			{Name: "ops_group", Actions: []string{"echo"}},
		},
	}

	var buf bytes.Buffer
	logger := logging.NewWithWriter(&buf, "json", "info")

	_, err := Execute("req-123", "ops_group", "echo", nil, cfg, logger)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	logOutput := buf.String()
	if !strings.Contains(logOutput, `"msg":"invoke_started"`) {
		t.Fatalf("expected invoke_started log, got %s", logOutput)
	}
	if !strings.Contains(logOutput, `"msg":"invoke_finished"`) {
		t.Fatalf("expected invoke_finished log, got %s", logOutput)
	}
	if !strings.Contains(logOutput, `"status":"success"`) {
		t.Fatalf("expected success status log, got %s", logOutput)
	}
	if !strings.Contains(logOutput, `"action":"echo"`) {
		t.Fatalf("expected action field log, got %s", logOutput)
	}
	if !strings.Contains(logOutput, `"request_id":"req-123"`) {
		t.Fatalf("expected request_id field log, got %s", logOutput)
	}
}
