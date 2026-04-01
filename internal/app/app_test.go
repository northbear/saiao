package app

import (
	"context"
	"testing"
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
