package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
	"saiao/internal/models"
)

var ErrNotImplemented = errors.New("config loading not implemented")

func Load(path string) (*models.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg models.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
