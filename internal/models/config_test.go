package models

import "testing"

func TestConfigTypesExist(t *testing.T) {
	cfg := Config{
		Server: ServerConfig{Listen: ":8080"},
		Logging: LoggingConfig{Format: "json", Level: "info"},
	}
	if cfg.Server.Listen != ":8080" {
		t.Fatal("unexpected listen value")
	}
}
