package app

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"saiao/internal/api"
)

func TestRunReturnsNotImplemented(t *testing.T) {
	err := Run(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if err != ErrNotImplemented {
		t.Fatalf("expected ErrNotImplemented, got %v", err)
	}
}

func TestRunWithOptionsRejectsNilContext(t *testing.T) {
	err := RunWithOptions(nil, DefaultOptions())
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "nil context" {
		t.Fatalf("expected nil context error, got %v", err)
	}
}

func TestRunWithOptionsShutsDownOnContextCancel(t *testing.T) {
	srv := api.NewServer(":0", api.BuildInfo{
		Service: "saiao",
		Version: "test",
		Commit:  "test",
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan error, 1)
	go func() {
		done <- RunWithOptions(ctx, Options{
			ListenAddress: ":0",
			BuildInfo: BuildInfo{
				Service: "saiao",
				Version: "test",
				Commit:  "test",
			},
			Server: srv,
		})
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("RunWithOptions did not return in time")
	}
}

func TestRunWithOptionsReturnsServerClosedAsNil(t *testing.T) {
	handler := http.NewServeMux()
	srv := &http.Server{Handler: handler}

	testListener := httptest.NewUnstartedServer(handler)
	testListener.Start()
	defer testListener.Close()

	srv.Addr = testListener.Listener.Addr().String()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- RunWithOptions(ctx, Options{
			ListenAddress: srv.Addr,
			BuildInfo: BuildInfo{
				Service: "saiao",
				Version: "test",
				Commit:  "test",
			},
			Server: srv,
		})
	}()

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Fatalf("expected nil or server closed handling, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("RunWithOptions did not return in time")
	}
}
