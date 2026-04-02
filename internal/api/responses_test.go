package api

import "testing"

func TestSuccessResponse(t *testing.T) {
	got := SuccessResponse("req-123", "hello", 0)

	if got.Status != "success" {
		t.Fatalf("expected success status, got %q", got.Status)
	}
	if got.RequestID != "req-123" {
		t.Fatalf("expected request id req-123, got %q", got.RequestID)
	}
	if got.Result.Output != "hello" {
		t.Fatalf("expected output hello, got %q", got.Result.Output)
	}
	if got.Result.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", got.Result.ExitCode)
	}
}
