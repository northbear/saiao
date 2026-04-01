package ssh

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

func TestExecutorExecuteRunsSSHCommandViaSystemSSH(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "ssh-args.txt")
	sshPath := filepath.Join(dir, "ssh")
	script := "#!/bin/sh\nprintf '%s\n' \"$@\" >" + logPath + "\nprintf 'remote ok'\n"
	if err := os.WriteFile(sshPath, []byte(script), 0o755); err != nil {
		t.Fatalf("write ssh shim: %v", err)
	}

	originalPath := os.Getenv("PATH")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+originalPath)

	executor := &Executor{}
	action := models.Action{
		Host:            "example.com",
		Port:            2222,
		CommandTemplate: "echo {{service}}",
	}
	identity := models.Identity{Username: "ops"}

	output, exitCode, err := executor.Execute(context.Background(), action, identity, map[string]any{"service": "nginx"})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if output != "remote ok" {
		t.Fatalf("unexpected output %q", output)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	argsData, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read ssh args: %v", err)
	}
	gotArgs := strings.Split(strings.TrimSpace(string(argsData)), "\n")
	wantArgs := []string{"-p", "2222", "ops@example.com", "echo nginx"}
	if strings.Join(gotArgs, "|") != strings.Join(wantArgs, "|") {
		t.Fatalf("unexpected ssh args %v", gotArgs)
	}
}

func TestExecutorExecuteReturnsTimeout(t *testing.T) {
	dir := t.TempDir()
	sshPath := filepath.Join(dir, "ssh")
	script := "#!/bin/sh\nsleep 2\n"
	if err := os.WriteFile(sshPath, []byte(script), 0o755); err != nil {
		t.Fatalf("write ssh shim: %v", err)
	}

	originalPath := os.Getenv("PATH")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+originalPath)

	executor := &Executor{}
	action := models.Action{
		Host:            "example.com",
		CommandTemplate: "echo hello",
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
