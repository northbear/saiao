package shell

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	saiaoerrors "saiao/internal/errors"
	"saiao/internal/models"
)

func TestExecutorExecuteRunsRenderedShellCommand(t *testing.T) {
	executor := &Executor{}
	action := models.Action{
		CommandTemplate: "printf 'hello %s' '{{name}}'",
	}

	output, exitCode, err := executor.Execute(context.Background(), action, models.Identity{}, map[string]any{
		"name": "world",
	})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if output != "hello world" {
		t.Fatalf("unexpected output %q", output)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
}

func TestExecutorExecuteUsesWorkingDirectoryAndEnvironment(t *testing.T) {
	dir := t.TempDir()
	executor := &Executor{}
	action := models.Action{
		CommandTemplate:  "printf '%s|%s' \"$PWD\" \"$GREETING\"",
		WorkingDirectory: dir,
		Environment: map[string]string{
			"GREETING": "hello {{name}}",
		},
	}

	output, exitCode, err := executor.Execute(context.Background(), action, models.Identity{}, map[string]any{
		"name": "ops",
	})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if output != dir+"|hello ops" {
		t.Fatalf("unexpected output %q", output)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
}

func TestExecutorExecuteReturnsExecutionFailedForNonZeroExit(t *testing.T) {
	executor := &Executor{}
	action := models.Action{
		CommandTemplate: "echo broken && exit 7",
	}

	output, exitCode, err := executor.Execute(context.Background(), action, models.Identity{}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, saiaoerrors.ErrExecutionFailed) {
		t.Fatalf("expected execution_failed error, got %v", err)
	}
	if exitCode != 7 {
		t.Fatalf("expected exit code 7, got %d", exitCode)
	}
	if !strings.Contains(output, "broken") {
		t.Fatalf("expected command output to be preserved, got %q", output)
	}
}

func TestExecutorExecuteReturnsTimeout(t *testing.T) {
	executor := &Executor{}
	action := models.Action{
		CommandTemplate: "sleep 2",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, exitCode, err := executor.Execute(ctx, action, models.Identity{}, nil)
	if err == nil {
		t.Fatal("expected timeout")
	}
	if !errors.Is(err, saiaoerrors.ErrTimeout) {
		t.Fatalf("expected timeout error, got %v", err)
	}
	if exitCode != -1 {
		t.Fatalf("expected exit code -1, got %d", exitCode)
	}
}

func TestExecutorExecuteReturnsInvalidInputForBadEnvironmentTemplate(t *testing.T) {
	executor := &Executor{}
	action := models.Action{
		CommandTemplate: "echo ok",
		Environment: map[string]string{
			"GREETING": "{{missing}}",
		},
	}

	_, _, err := executor.Execute(context.Background(), action, models.Identity{}, map[string]any{
		"name": "ops",
	})
	if err == nil {
		t.Fatal("expected invalid input error")
	}
	if !errors.Is(err, saiaoerrors.ErrInvalidInput) {
		t.Fatalf("expected invalid input error, got %v", err)
	}
}

func TestExecutorExecuteRespectsWorkingDirectoryFiles(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "message.txt")
	if err := os.WriteFile(filePath, []byte("from-file"), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	executor := &Executor{}
	action := models.Action{
		CommandTemplate:  "cat message.txt",
		WorkingDirectory: dir,
	}

	output, exitCode, err := executor.Execute(context.Background(), action, models.Identity{}, nil)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if output != "from-file" {
		t.Fatalf("unexpected output %q", output)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
}
