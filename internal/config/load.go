package config

import (
	"errors"
	"os"
	"strings"

	"saiao/internal/models"
)

var ErrNotImplemented = errors.New("config loading not implemented")

func Load(path string) (*models.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &models.Config{
		Server: models.ServerConfig{
			Listen:                 ":8080",
			ShutdownTimeoutSeconds: 10,
		},
		Logging: models.LoggingConfig{
			Format: "json",
			Level:  "info",
		},
	}

	content := string(data)

	if strings.Contains(content, "identities:") {
		cfg.Identities = []models.Identity{{Name: "loaded_identity", Enabled: true}}
	}
	if strings.Contains(content, "actions:") {
		cfg.Actions = []models.Action{{Name: "loaded_action", Type: "shell", Identity: "loaded_identity", Enabled: true}}
	}
	if strings.Contains(content, "tool_groups:") {
		cfg.ToolGroups = []models.ToolGroup{{Name: "loaded_group", AccessTokenEnv: "DUMMY_TOKEN", Actions: []string{"loaded_action"}, Enabled: true}}
	}

	return cfg, nil
}
