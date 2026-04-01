package invoke

import "testing"

func TestValidateInputChecksRequiredTypeAndEnum(t *testing.T) {
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"service": map[string]any{
				"type": "string",
				"enum": []any{"nginx", "redis"},
			},
			"force": map[string]any{
				"type": "boolean",
			},
		},
		"required": []any{"service"},
	}

	if err := ValidateInput(schema, map[string]any{"service": "nginx", "force": true}); err != nil {
		t.Fatalf("expected input to validate, got %v", err)
	}

	if err := ValidateInput(schema, map[string]any{}); err == nil {
		t.Fatal("expected missing required field to fail")
	}

	if err := ValidateInput(schema, map[string]any{"service": "postgres"}); err == nil {
		t.Fatal("expected enum mismatch to fail")
	}
}
