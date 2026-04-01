package app

import (
	"context"
	"testing"
	"time"
)

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

func TestRunReturnsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)

	go func() {
		done <- Run(ctx)
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return in time")
	}
}
