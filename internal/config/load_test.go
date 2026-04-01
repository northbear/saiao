package config

import "testing"

func TestLoadReturnsNotImplemented(t *testing.T) {
	_, err := Load("test.yaml")
	if err == nil {
		t.Fatal("expected error")
	}
}
