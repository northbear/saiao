package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"saiao/internal/models"
)

func Load(path string) (*models.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg models.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if err := ResolveSecrets(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
