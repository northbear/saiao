package config

import (
	"errors"
	"os"

	"saiao/internal/models"
)

var ErrNotImplemented = errors.New("config loading not implemented")

func Load(path string) (*models.Config, error) {
	_, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &models.Config{
		Server: models.ServerConfig{
			Listen:                ":8080",
			ShutdownTimeoutSeconds: 10,
		},
		Logging: models.LoggingConfig{
			Format: "json",
			Level:  "info",
		},
	}

	return cfg, nil
}
