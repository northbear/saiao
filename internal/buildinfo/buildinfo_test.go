package buildinfo

import "testing"

func TestCurrentReturnsBuildMetadata(t *testing.T) {
	originalVersion := Version
	originalCommit := Commit
	t.Cleanup(func() {
		Version = originalVersion
		Commit = originalCommit
	})

	Version = "1.2.3"
	Commit = "abc1234"

	got := Current()
	if got.Service != Service {
		t.Fatalf("expected service %q, got %q", Service, got.Service)
	}
	if got.Version != "1.2.3" {
		t.Fatalf("expected version 1.2.3, got %q", got.Version)
	}
	if got.Commit != "abc1234" {
		t.Fatalf("expected commit abc1234, got %q", got.Commit)
	}
}
