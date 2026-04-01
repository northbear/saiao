package invoke

import (
	stdErrors "errors"
	"testing"

	saiaoerrors "saiao/internal/errors"
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

	got, err := Execute("ops_group", "echo", []byte(`{"message":"world"}`), cfg)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if got.Result.Output != "hello world" {
		t.Fatalf("expected hello world, got %q", got.Result.Output)
	}
	if got.Result.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", got.Result.ExitCode)
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

	_, err := Execute("ops_group", "slow", nil, cfg)
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if !stdErrors.Is(err, saiaoerrors.ErrTimeout) {
		t.Fatalf("expected timeout error, got %v", err)
	}
}
