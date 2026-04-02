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
	if got.Format != FormatSAIAO {
		t.Fatalf("expected format %q, got %q", FormatSAIAO, got.Format)
	}
}

func TestBuildOpenAIReturnsFunctionTools(t *testing.T) {
	cfg := &models.Config{
		Actions: []models.Action{
			{
				Name:        "a1",
				Description: "first",
				InputSchema: map[string]any{"type": "object"},
			},
			{Name: "a2", Description: "second"},
		},
		ToolGroups: []models.ToolGroup{
			{Name: "g1", Actions: []string{"a1", "a2"}},
		},
	}

	got, err := BuildOpenAI("g1", cfg)
	if err != nil {
		t.Fatalf("build openai manifest: %v", err)
	}

	if got.Format != FormatOpenAI {
		t.Fatalf("expected format %q, got %q", FormatOpenAI, got.Format)
	}
	if len(got.Tools) != 2 {
		t.Fatalf("expected 2 tools, got %d", len(got.Tools))
	}
	if got.Tools[0].Type != "function" {
		t.Fatalf("expected function tool type, got %q", got.Tools[0].Type)
	}
	if got.Tools[0].Name != "a1" {
		t.Fatalf("expected a1, got %q", got.Tools[0].Name)
	}
	if got.Tools[0].Strict {
		t.Fatal("expected strict false")
	}
	if got.Tools[1].Parameters["type"] != "object" {
		t.Fatalf("expected default object schema, got %#v", got.Tools[1].Parameters)
	}
}
