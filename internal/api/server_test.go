package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewServerSetsAddress(t *testing.T) {
	srv := NewServer(":1234", appInfo{
		Service: "saiao",
		Version: "1.0.0",
		Commit:  "abc1234",
	})
	if srv.Addr != ":1234" {
		t.Fatalf("expected addr :1234, got %s", srv.Addr)
	}
	if srv.Handler == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestInfoEndpointReturnsInfoJSON(t *testing.T) {
	srv := NewServer(":0", appInfo{
		Service: "saiao",
		Version: "1.0.0",
		Commit:  "abc1234",
	})

	req := httptest.NewRequest(http.MethodGet, "/info", nil)
	rr := httptest.NewRecorder()

	srv.Handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var got map[string]any
	body, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}

	if got["status"] != "ok" {
		t.Fatalf("expected status ok, got %v", got["status"])
	}
	if got["service"] != "saiao" {
		t.Fatalf("expected service saiao, got %v", got["service"])
	}
	if got["version"] != "1.0.0" {
		t.Fatalf("expected version 1.0.0, got %v", got["version"])
	}
	if got["commit"] != "abc1234" {
		t.Fatalf("expected commit abc1234, got %v", got["commit"])
	}
}
