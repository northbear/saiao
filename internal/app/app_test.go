package app

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"saiao/internal/api"
)

type immediateServer struct{}

func (s *immediateServer) ListenAndServe() error {
	return http.ErrServerClosed
}

func (s *immediateServer) Shutdown(context.Context) error {
	return nil
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

func TestRunWithOptionsReturnsServerClosedAsNil(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := RunWithOptions(ctx, Options{
		ListenAddress: ":0",
		BuildInfo: BuildInfo{
			Service: "saiao",
			Version: "test",
			Commit:  "test",
		},
		Server: &http.Server{},
	})

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		t.Fatalf("expected nil or server closed handling, got %v", err)
	}
}

func TestRunWithOptionsStopsOnCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	srv := api.NewServer(":0", api.BuildInfo{
		Service: "saiao",
		Version: "test",
		Commit:  "test",
	})

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
