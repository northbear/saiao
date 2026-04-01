package config

import (
	"strings"
	"testing"

	"saiao/internal/models"
)

func TestValidateRequiresSections(t *testing.T) {
	if err := Validate(nil); err == nil {
		t.Fatal("expected error for nil config")
	}

	if err := Validate(&models.Config{}); err == nil {
		t.Fatal("expected error for empty config")
	}
}

func TestValidateChecksReferencesAndUniqueness(t *testing.T) {
	t.Setenv("GROUP_TOKEN", "token")

	cfg := &models.Config{
		Identities: []models.Identity{
			{Name: "ops", Enabled: true},
			{Name: "ops", Enabled: true},
		},
		Actions: []models.Action{
			{Name: "action1", Type: "shell", Identity: "ops", Enabled: true},
			{Name: "action1", Type: "shell", Identity: "ops", Enabled: true},
		},
		ToolGroups: []models.ToolGroup{
			{Name: "group1", AccessTokenEnv: "GROUP_TOKEN", Actions: []string{"action1"}, Enabled: true},
		},
	}

	if err := Validate(cfg); err == nil {
		t.Fatal("expected validation error for duplicates")
	}
}

func TestValidateRequiresKnownReferences(t *testing.T) {
	t.Setenv("GROUP_TOKEN", "token")

	cfg := &models.Config{
		Identities: []models.Identity{{Name: "ops", Enabled: true}},
		Actions:    []models.Action{{Name: "action1", Type: "shell", Identity: "unknown", Enabled: true}},
		ToolGroups: []models.ToolGroup{{Name: "group1", AccessTokenEnv: "GROUP_TOKEN", Actions: []string{"action1"}, Enabled: true}},
	}

	if err := Validate(cfg); err == nil {
		t.Fatal("expected validation error for unknown identity")
	}
}

func TestValidateRequiresGroupTokenEnv(t *testing.T) {
	cfg := &models.Config{
		Identities: []models.Identity{{Name: "ops", Enabled: true}},
		Actions:    []models.Action{{Name: "action1", Type: "shell", Identity: "ops", Enabled: true}},
		ToolGroups: []models.ToolGroup{{Name: "group1", AccessTokenEnv: "MISSING_TOKEN", Actions: []string{"action1"}, Enabled: true}},
	}

	if err := Validate(cfg); err == nil {
		t.Fatal("expected validation error for missing token env")
	}
}

func TestValidateChecksActionSpecificFieldsAndSchema(t *testing.T) {
	t.Setenv("GROUP_TOKEN", "token")

	cfg := &models.Config{
		Identities: []models.Identity{{Name: "ops", Enabled: true}},
		Actions: []models.Action{
			{
				Name:     "action1",
				Type:     "shell",
				Identity: "ops",
				InputSchema: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"service": map[string]any{
							"type": "string",
							"enum": []any{"nginx", true},
						},
					},
				},
			},
		},
		ToolGroups: []models.ToolGroup{{Name: "group1", AccessTokenEnv: "GROUP_TOKEN", Actions: []string{"action1"}, Enabled: true}},
	}

	if err := Validate(cfg); err == nil {
		t.Fatal("expected validation error for missing command template and invalid enum schema")
	}
}

func TestValidateRejectsDuplicateSecretSources(t *testing.T) {
	t.Setenv("GROUP_TOKEN", "token")

	cfg := &models.Config{
		Identities: []models.Identity{
			{
				Name: "ops",
				Secrets: map[string]string{
					"http_token_env":  "HTTP_TOKEN",
					"http_token_file": "/run/secrets/http_token",
				},
			},
		},
		Actions: []models.Action{
			{
				Name:            "action1",
				Type:            "shell",
				Identity:        "ops",
				CommandTemplate: "echo ok",
			},
		},
		ToolGroups: []models.ToolGroup{
			{Name: "group1", AccessTokenEnv: "GROUP_TOKEN", Actions: []string{"action1"}, Enabled: true},
		},
	}

	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected duplicate secret source validation error")
	}
	if !strings.Contains(err.Error(), `defines multiple sources for "http_token"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateRejectsShellIdentityPlaceholders(t *testing.T) {
	t.Setenv("GROUP_TOKEN", "token")

	cfg := &models.Config{
		Identities: []models.Identity{{Name: "ops", Enabled: true}},
		Actions: []models.Action{
			{
				Name:            "action1",
				Type:            "shell",
				Identity:        "ops",
				CommandTemplate: "echo {{identity.http_token}}",
			},
		},
		ToolGroups: []models.ToolGroup{
			{Name: "group1", AccessTokenEnv: "GROUP_TOKEN", Actions: []string{"action1"}, Enabled: true},
		},
	}

	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected validation error for identity placeholder in shell command")
	}
	if !strings.Contains(err.Error(), "command_template must not reference identity placeholders") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateRejectsUndeclaredTemplateInput(t *testing.T) {
	t.Setenv("GROUP_TOKEN", "token")

	cfg := &models.Config{
		Identities: []models.Identity{{Name: "ops", Enabled: true}},
		Actions: []models.Action{
			{
				Name:         "action1",
				Type:         "http",
				Identity:     "ops",
				Method:       "POST",
				URL:          "https://example.com/repos/{{repo}}",
				BodyTemplate: `{"title":"{{title}}"}`,
				InputSchema: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"title": map[string]any{"type": "string"},
					},
				},
			},
		},
		ToolGroups: []models.ToolGroup{
			{Name: "group1", AccessTokenEnv: "GROUP_TOKEN", Actions: []string{"action1"}, Enabled: true},
		},
	}

	err := Validate(cfg)
	if err == nil {
		t.Fatal("expected validation error for undeclared template input")
	}
	if !strings.Contains(err.Error(), `references undeclared input "repo"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}
