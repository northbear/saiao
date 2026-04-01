package auth

import "testing"

func TestParseBearerToken(t *testing.T) {
	token, ok := ParseBearerToken("Bearer abc123")
	if !ok {
		t.Fatal("expected token to parse")
	}
	if token != "abc123" {
		t.Fatalf("expected abc123, got %q", token)
	}

	if _, ok := ParseBearerToken("Basic abc123"); ok {
		t.Fatal("did not expect basic auth to parse")
	}
}

func TestTokenStore(t *testing.T) {
	store := NewTokenStore()
	store.Register("token1", "group1")

	group, ok := store.GroupForToken("token1")
	if !ok {
		t.Fatal("expected token lookup to succeed")
	}
	if group != "group1" {
		t.Fatalf("expected group1, got %q", group)
	}

	if _, ok := store.GroupForToken("missing"); ok {
		t.Fatal("did not expect missing token to resolve")
	}
}
