package render

import (
	"strings"
	"testing"

	"saiao/internal/models"
)

func TestTemplateRendersInputAndIdentityFields(t *testing.T) {
	identity := models.Identity{
		Name:     "ops",
		Username: "svc_ops",
		Secrets: map[string]string{
			"http_token": "secret-token",
		},
	}

	rendered, err := Template("user={{identity.username}} token={{identity.http_token}} repo={{repo}}", identity, map[string]any{
		"repo": "saiao",
	})
	if err != nil {
		t.Fatalf("template render: %v", err)
	}

	if rendered != "user=svc_ops token=secret-token repo=saiao" {
		t.Fatalf("unexpected rendered template %q", rendered)
	}
}

func TestTemplateReturnsErrorForUnknownPlaceholder(t *testing.T) {
	_, err := Template("{{missing}}", models.Identity{}, map[string]any{"known": "value"})
	if err == nil {
		t.Fatal("expected missing placeholder to fail")
	}
	if !strings.Contains(err.Error(), `template path "missing" is undefined`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPlaceholdersReturnsUniqueNames(t *testing.T) {
	got := Placeholders("{{name}} {{identity.http_token}} {{name}} {{identity.http_token}} {{title}}")

	want := []string{"name", "identity.http_token", "title"}
	if len(got) != len(want) {
		t.Fatalf("expected %d placeholders, got %d: %v", len(want), len(got), got)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected placeholder %q at index %d, got %q", want[i], i, got[i])
		}
	}
}
