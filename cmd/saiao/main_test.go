package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestBinaryStartsAndServesInfo(t *testing.T) {
	dir := t.TempDir()
	binaryPath := filepath.Join(dir, "saiao")
	cacheDir := filepath.Join(dir, "gocache")

	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/saiao")
	buildCmd.Dir = filepath.Join("..", "..")
	buildCmd.Env = append(os.Environ(), "GOCACHE="+cacheDir)
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build binary: %v\n%s", err, string(output))
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		if strings.Contains(err.Error(), "operation not permitted") {
			t.Skip("sandbox does not allow reserving a local listener")
		}
		t.Fatalf("reserve port: %v", err)
	}
	address := listener.Addr().String()
	listener.Close()

	configPath := filepath.Join(dir, "config.yaml")
	config := []byte(`
server:
  listen: "` + address + `"
  shutdown_timeout_seconds: 1
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
	if err := os.WriteFile(configPath, config, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runCmd := exec.CommandContext(ctx, binaryPath, "-config", configPath)
	runCmd.Env = append(os.Environ(), "TEST_GROUP_TOKEN=token")
	runCmd.Dir = dir
	stderr, err := runCmd.StderrPipe()
	if err != nil {
		t.Fatalf("stderr pipe: %v", err)
	}

	if err := runCmd.Start(); err != nil {
		if strings.Contains(err.Error(), "operation not permitted") {
			t.Skip("sandbox does not allow starting the binary listener")
		}
		t.Fatalf("start binary: %v", err)
	}
	defer func() {
		cancel()
		_ = runCmd.Process.Kill()
		_ = runCmd.Wait()
	}()

	errCh := make(chan string, 1)
	go func() {
		data, _ := io.ReadAll(stderr)
		errCh <- string(data)
	}()

	client := &http.Client{Timeout: 200 * time.Millisecond}
	infoURL := "http://" + address + "/info"

	var lastErr error
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := client.Get(infoURL)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected 200 from /info, got %d", resp.StatusCode)
			}

			var payload map[string]any
			if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
				t.Fatalf("decode /info: %v", err)
			}
			if payload["service"] != "saiao" {
				t.Fatalf("expected service saiao, got %v", payload["service"])
			}
			return
		}

		lastErr = err
		select {
		case stderrOutput := <-errCh:
			if strings.Contains(stderrOutput, "operation not permitted") {
				t.Skip("sandbox does not allow binding the binary listener")
			}
			if stderrOutput != "" {
				t.Fatalf("binary exited before serving /info: %s", stderrOutput)
			}
		default:
		}
		time.Sleep(50 * time.Millisecond)
	}

	if lastErr != nil && errors.Is(lastErr, context.DeadlineExceeded) {
		t.Fatalf("binary did not become ready: %v", lastErr)
	}
	t.Fatalf("binary did not become ready: %v", lastErr)
}
