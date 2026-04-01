package config

import (
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
		Actions: []models.Action{{Name: "action1", Type: "shell", Identity: "unknown", Enabled: true}},
		ToolGroups: []models.ToolGroup{{Name: "group1", AccessTokenEnv: "GROUP_TOKEN", Actions: []string{"action1"}, Enabled: true}},
	}

	if err := Validate(cfg); err == nil {
		t.Fatal("expected validation error for unknown identity")
	}
}

func TestValidateRequiresGroupTokenEnv(t *testing.T) {
	cfg := &models.Config{
		Identities: []models.Identity{{Name: "ops", Enabled: true}},
		Actions: []models.Action{{Name: "action1", Type: "shell", Identity: "ops", Enabled: true}},
		ToolGroups: []models.ToolGroup{{Name: "group1", AccessTokenEnv: "MISSING_TOKEN", Actions: []string{"action1"}, Enabled: true}},
	}

	if err := Validate(cfg); err == nil {
		t.Fatal("expected validation error for missing token env")
	}
}
