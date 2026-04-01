package app

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

type stubServer struct {
	listenCalled  bool
	shutdownCalled bool
	listenErr     error
	shutdownErr   error
}

func (s *stubServer) ListenAndServe() error {
	s.listenCalled = true
	return s.listenErr
}

func (s *stubServer) Shutdown(ctx context.Context) error {
	_ = ctx
	s.shutdownCalled = true
	return s.shutdownErr
}

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

func TestRunWithOptionsReturnsNilOnCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := RunWithOptions(ctx, Options{})
	if err != nil {
		t.Fatalf("expected nil error on canceled context, got %v", err)
	}
}

func TestRunWithOptionsHandlesServerClosed(t *testing.T) {
	srv := &http.Server{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Fatalf("expected nil or server closed handling, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("RunWithOptions did not return in time")
	}
}
