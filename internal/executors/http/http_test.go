package http

import (
	"context"
	"errors"
	"io"
	"net"
	stdhttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	saiaoerrors "saiao/internal/errors"
	"saiao/internal/models"
)

func TestExecutorExecuteSendsRenderedRequest(t *testing.T) {
	var gotMethod string
	var gotAuth string
	var gotBody string

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		if strings.Contains(err.Error(), "operation not permitted") {
			t.Skip("sandbox does not allow opening a local listener")
		}
		t.Fatalf("listen: %v", err)
	}

	server := httptest.NewUnstartedServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		gotMethod = r.Method
		gotAuth = r.Header.Get("Authorization")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		gotBody = string(body)
		w.WriteHeader(stdhttp.StatusCreated)
		_, _ = w.Write([]byte("created"))
	}))
	server.Listener = listener
	server.Start()
	defer server.Close()

	executor := &Executor{}
	action := models.Action{
		Method:       stdhttp.MethodPost,
		URL:          server.URL + "/issues/{{repo}}",
		Headers:      map[string]string{"Authorization": "Bearer {{identity.http_token}}"},
		BodyTemplate: `{"title":"{{title}}"}`,
	}
	identity := models.Identity{
		Secrets: map[string]string{"http_token": "secret-token"},
	}

	output, exitCode, err := executor.Execute(context.Background(), action, identity, map[string]any{
		"repo":  "demo",
		"title": "Hello",
	})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if output != "created" {
		t.Fatalf("expected created output, got %q", output)
	}
	if exitCode != stdhttp.StatusCreated {
		t.Fatalf("expected status code 201, got %d", exitCode)
	}
	if gotMethod != stdhttp.MethodPost {
		t.Fatalf("expected POST, got %s", gotMethod)
	}
	if gotAuth != "Bearer secret-token" {
		t.Fatalf("unexpected Authorization header %q", gotAuth)
	}
	if strings.TrimSpace(gotBody) != `{"title":"Hello"}` {
		t.Fatalf("unexpected body %q", gotBody)
	}
}

func TestExecutorExecuteReturnsExecutionFailedOnNon2xx(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		if strings.Contains(err.Error(), "operation not permitted") {
			t.Skip("sandbox does not allow opening a local listener")
		}
		t.Fatalf("listen: %v", err)
	}

	server := httptest.NewUnstartedServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		w.WriteHeader(stdhttp.StatusBadGateway)
		_, _ = w.Write([]byte("bad gateway"))
	}))
	server.Listener = listener
	server.Start()
	defer server.Close()

	executor := &Executor{}
	action := models.Action{
		Method: stdhttp.MethodGet,
		URL:    server.URL,
	}

	output, exitCode, err := executor.Execute(context.Background(), action, models.Identity{}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, saiaoerrors.ErrExecutionFailed) {
		t.Fatalf("expected execution failed error, got %v", err)
	}
	if output != "bad gateway" {
		t.Fatalf("unexpected output %q", output)
	}
	if exitCode != stdhttp.StatusBadGateway {
		t.Fatalf("unexpected exit code %d", exitCode)
	}
}

func TestExecutorExecuteReturnsTimeoutOnCanceledContext(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		if strings.Contains(err.Error(), "operation not permitted") {
			t.Skip("sandbox does not allow opening a local listener")
		}
		t.Fatalf("listen: %v", err)
	}

	server := httptest.NewUnstartedServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(stdhttp.StatusOK)
	}))
	server.Listener = listener
	server.Start()
	defer server.Close()

	executor := &Executor{}
	action := models.Action{
		Method: stdhttp.MethodGet,
		URL:    server.URL,
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
