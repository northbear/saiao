package manifest

import "testing"

import "saiao/internal/models"

func TestBuildReturnsToolsForRequestedGroup(t *testing.T) {
	cfg := &models.Config{
		Actions: []models.Action{
			{Name: "a1", Description: "first", InputSchema: map[string]any{"type": "object"}},
			{Name: "a2", Description: "second"},
		},
		ToolGroups: []models.ToolGroup{
			{Name: "g1", Actions: []string{"a1"}},
		},
	}

	got, err := Build("g1", cfg)
	if err != nil {
		t.Fatalf("build manifest: %v", err)
	}

	if len(got.Tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(got.Tools))
	}
	if got.Tools[0].Name != "a1" {
		t.Fatalf("expected a1, got %q", got.Tools[0].Name)
	}
}
