package config

import (
	"errors"
	"fmt"
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

	identityNames := make(map[string]struct{}, len(cfg.Identities))
	for i := range cfg.Identities {
		name := strings.TrimSpace(cfg.Identities[i].Name)
		if name == "" {
			return fmt.Errorf("identities[%d].name is required", i)
		}
		if _, exists := identityNames[name]; exists {
			return fmt.Errorf("duplicate identity name %q", name)
		}
		identityNames[name] = struct{}{}
		cfg.Identities[i].Name = name
	}

	actionNames := make(map[string]struct{}, len(cfg.Actions))
	for i := range cfg.Actions {
		name := strings.TrimSpace(cfg.Actions[i].Name)
		if name == "" {
			return fmt.Errorf("actions[%d].name is required", i)
		}
		if _, exists := actionNames[name]; exists {
			return fmt.Errorf("duplicate action name %q", name)
		}
		actionNames[name] = struct{}{}
		cfg.Actions[i].Name = name

		if strings.TrimSpace(cfg.Actions[i].Type) == "" {
			return fmt.Errorf("actions[%d].type is required", i)
		}
		if strings.TrimSpace(cfg.Actions[i].Identity) == "" {
			return fmt.Errorf("actions[%d].identity is required", i)
		}
		if _, ok := identityNames[cfg.Actions[i].Identity]; !ok {
			return fmt.Errorf("actions[%d].identity references unknown identity %q", i, cfg.Actions[i].Identity)
		}
	}

	toolGroupNames := make(map[string]struct{}, len(cfg.ToolGroups))
	for i := range cfg.ToolGroups {
		name := strings.TrimSpace(cfg.ToolGroups[i].Name)
		if name == "" {
			return fmt.Errorf("tool_groups[%d].name is required", i)
		}
		if _, exists := toolGroupNames[name]; exists {
			return fmt.Errorf("duplicate tool group name %q", name)
		}
		toolGroupNames[name] = struct{}{}
		cfg.ToolGroups[i].Name = name

		if strings.TrimSpace(cfg.ToolGroups[i].AccessTokenEnv) == "" {
			return fmt.Errorf("tool_groups[%d].access_token_env is required", i)
		}
		if _, ok := os.LookupEnv(cfg.ToolGroups[i].AccessTokenEnv); !ok {
			return fmt.Errorf("missing required environment variable %q", cfg.ToolGroups[i].AccessTokenEnv)
		}
		if len(cfg.ToolGroups[i].Actions) == 0 {
			return fmt.Errorf("tool_groups[%d].actions must not be empty", i)
		}
		for _, actionName := range cfg.ToolGroups[i].Actions {
			if _, ok := actionNames[actionName]; !ok {
				return fmt.Errorf("tool_groups[%d] references unknown action %q", i, actionName)
			}
		}
	}

	return nil
}
