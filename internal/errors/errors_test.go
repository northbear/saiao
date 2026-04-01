package errors

import (
	stdErrors "errors"
	"testing"
)

func TestCanonicalErrorsMatchByIs(t *testing.T) {
	tests := []struct {
		name string
		err  error
		ok   func(error) bool
	}{
		{name: "unauthorized", err: ErrUnauthorized, ok: IsUnauthorized},
		{name: "not found", err: ErrNotFound, ok: IsNotFound},
		{name: "invalid input", err: ErrInvalidInput, ok: IsInvalidInput},
		{name: "execution failed", err: ErrExecutionFailed, ok: IsExecutionFailed},
		{name: "timeout", err: ErrTimeout, ok: IsTimeout},
		{name: "internal", err: ErrInternal, ok: IsInternal},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.ok(tc.err) {
				t.Fatalf("expected helper to match %v", tc.err)
			}

			wrapped := stdErrors.New("wrap: " + tc.err.Error())
			if tc.ok(wrapped) {
				t.Fatalf("did not expect helper to match unrelated wrapped error")
			}

			if !stdErrors.Is(tc.err, tc.err) {
				t.Fatalf("expected errors.Is to match identical error")
			}
		})
	}
}
