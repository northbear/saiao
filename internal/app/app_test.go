package app

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	content := []byte(`
server:
  listen: ":0"
  shutdown_timeout_seconds: 10
logging:
  format: json
  level: info
identities:
  - name: test_identity
actions:
  - name: test_action
    type: shell
    identity: test_identity
    command_template: echo test
tool_groups:
  - name: test_group
    access_token_env: TEST_GROUP_TOKEN
    actions:
      - test_action
`)

	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("write temp config: %v", err)
	}

	t.Setenv("TEST_GROUP_TOKEN", "token")
	return path
}

func TestRunWithContextReturnsNilOnCanceledContext(t *testing.T) {
	cfgPath := writeTempConfig(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := RunWithContext(ctx, Options{ConfigPath: cfgPath})
	if err != nil {
		t.Fatalf("expected nil error on canceled context, got %v", err)
	}
}

func TestRunReturnsOnContextCancel(t *testing.T) {
	cfgPath := writeTempConfig(t)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)

	go func() {
		done <- RunWithContext(ctx, Options{ConfigPath: cfgPath})
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("RunWithContext did not return in time")
	}
}

func TestDefaultOptionsProvidesConfigPath(t *testing.T) {
	opts := DefaultOptions()
	if opts.ConfigPath == "" {
		t.Fatal("expected default config path to be set")
	}
}
