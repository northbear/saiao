package config

import (
	"errors"
	"strings"

	"saiao/internal/models"
)

func Validate(cfg *models.Config) error {
	if cfg == nil {
		return errors.New("config is nil")
	}

	if strings.TrimSpace(cfg.Server.Listen) == "" {
		cfg.Server.Listen = ":8080"
	}

	if cfg.Server.ShutdownTimeoutSeconds <= 0 {
		cfg.Server.ShutdownTimeoutSeconds = 10
	}

	if cfg.Logging.Format == "" {
		cfg.Logging.Format = "json"
	}
	if cfg.Logging.Level == "" {
		cfg.Logging.Level = "info"
	}

	if len(cfg.Identities) == 0 {
		return errors.New("identities section is required")
	}
	if len(cfg.Actions) == 0 {
		return errors.New("actions section is required")
	}
	if len(cfg.ToolGroups) == 0 {
		return errors.New("tool_groups section is required")
	}

	return nil
}
