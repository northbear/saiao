package config

import (
	"os"

	"saiao/internal/models"
)

func Load(path string) (*models.Config, error) {
	_, err := os.ReadFile(path)
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
		Identities: []models.Identity{
			{Name: "default_identity", Enabled: true},
		},
		Actions: []models.Action{
			{Name: "default_action", Type: "shell", Identity: "default_identity", Enabled: true},
		},
		ToolGroups: []models.ToolGroup{
			{Name: "default_group", AccessTokenEnv: "SAIAO_TOKEN_DEFAULT", Actions: []string{"default_action"}, Enabled: true},
		},
	}

	return cfg, nil
}
