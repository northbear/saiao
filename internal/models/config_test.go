package models

import "testing"

func TestConfigHelpers(t *testing.T) {
	var cfg Config

	if cfg.HasIdentities() {
		t.Fatal("expected no identities")
	}
	if cfg.HasActions() {
		t.Fatal("expected no actions")
	}
	if cfg.HasToolGroups() {
		t.Fatal("expected no tool groups")
	}

	cfg.Identities = []Identity{{Name: "id1"}}
	cfg.Actions = []Action{{Name: "a1"}}
	cfg.ToolGroups = []ToolGroup{{Name: "tg1"}}

	if !cfg.HasIdentities() {
		t.Fatal("expected identities")
	}
	if !cfg.HasActions() {
		t.Fatal("expected actions")
	}
	if !cfg.HasToolGroups() {
		t.Fatal("expected tool groups")
	}
}
