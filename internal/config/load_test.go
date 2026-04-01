package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadParsesYAMLAndResolvesSecrets(t *testing.T) {
	dir := t.TempDir()

	secretPath := filepath.Join(dir, "secret.txt")
	if err := os.WriteFile(secretPath, []byte("top-secret\n"), 0o600); err != nil {
		t.Fatalf("write secret: %v", err)
	}

	t.Setenv("TEST_GROUP_TOKEN", "group-token")
	t.Setenv("HTTP_TOKEN", "api-token")

	cfgPath := filepath.Join(dir, "config.yaml")
	content := []byte(`
server:
  listen: ":9090"
  shutdown_timeout_seconds: 12
logging:
  format: text
  level: debug
identities:
  - name: ops
    enabled: true
    secrets:
      http_token_env: HTTP_TOKEN
      ssh_private_key_file: ` + secretPath + `
actions:
  - name: ping
    type: shell
    identity: ops
    enabled: true
    command_template: echo ok
tool_groups:
  - name: ops_group
    enabled: true
    access_token_env: TEST_GROUP_TOKEN
    actions:
      - ping
`)

	if err := os.WriteFile(cfgPath, content, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if err := Validate(cfg); err != nil {
		t.Fatalf("validate config: %v", err)
	}

	if got := cfg.Server.Listen; got != ":9090" {
		t.Fatalf("expected listen :9090, got %q", got)
	}
	if got := cfg.Identities[0].Secrets["http_token"]; got != "api-token" {
		t.Fatalf("expected resolved env secret, got %q", got)
	}
	if got := cfg.Identities[0].Secrets["ssh_private_key"]; got != "top-secret" {
		t.Fatalf("expected resolved file secret, got %q", got)
	}
}

func TestLoadReturnsErrorForMissingFile(t *testing.T) {
	_, err := Load("test.yaml")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadReturnsErrorForMissingSecretEnv(t *testing.T) {
	t.Setenv("TEST_GROUP_TOKEN", "group-token")

	cfgPath := filepath.Join(t.TempDir(), "config.yaml")
	content := []byte(`
identities:
  - name: ops
    secrets:
      http_token_env: MISSING_HTTP_TOKEN
actions:
  - name: ping
    type: shell
    identity: ops
    command_template: echo ok
tool_groups:
  - name: ops_group
    access_token_env: TEST_GROUP_TOKEN
    actions:
      - ping
`)

	if err := os.WriteFile(cfgPath, content, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	_, err := Load(cfgPath)
	if err == nil {
		t.Fatal("expected missing env secret to fail")
	}
	if !strings.Contains(err.Error(), `missing required environment variable "MISSING_HTTP_TOKEN"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadReturnsErrorForMissingSecretFile(t *testing.T) {
	t.Setenv("TEST_GROUP_TOKEN", "group-token")

	cfgPath := filepath.Join(t.TempDir(), "config.yaml")
	content := []byte(`
identities:
  - name: ops
    secrets:
      ssh_private_key_file: /tmp/does-not-exist
actions:
  - name: ping
    type: shell
    identity: ops
    command_template: echo ok
tool_groups:
  - name: ops_group
    access_token_env: TEST_GROUP_TOKEN
    actions:
      - ping
`)

	if err := os.WriteFile(cfgPath, content, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	_, err := Load(cfgPath)
	if err == nil {
		t.Fatal("expected missing secret file to fail")
	}
	if !strings.Contains(err.Error(), `read secret file "/tmp/does-not-exist"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadReturnsErrorForEmptySecretFile(t *testing.T) {
	t.Setenv("TEST_GROUP_TOKEN", "group-token")

	dir := t.TempDir()
	secretPath := filepath.Join(dir, "empty-secret.txt")
	if err := os.WriteFile(secretPath, []byte("\n"), 0o600); err != nil {
		t.Fatalf("write secret: %v", err)
	}

	cfgPath := filepath.Join(dir, "config.yaml")
	content := []byte(`
identities:
  - name: ops
    secrets:
      ssh_private_key_file: ` + secretPath + `
actions:
  - name: ping
    type: shell
    identity: ops
    command_template: echo ok
tool_groups:
  - name: ops_group
    access_token_env: TEST_GROUP_TOKEN
    actions:
      - ping
`)

	if err := os.WriteFile(cfgPath, content, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	_, err := Load(cfgPath)
	if err == nil {
		t.Fatal("expected empty secret file to fail")
	}
	if !strings.Contains(err.Error(), `secret file "`) || !strings.Contains(err.Error(), `is empty`) {
		t.Fatalf("unexpected error: %v", err)
	}
}
