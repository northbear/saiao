package config

import (
	"fmt"
	"os"
	"strings"

	"saiao/internal/models"
)

func ResolveSecrets(cfg *models.Config) error {
	if cfg == nil {
		return nil
	}

	for i := range cfg.Identities {
		if cfg.Identities[i].Secrets == nil {
			continue
		}

		resolved := make(map[string]string, len(cfg.Identities[i].Secrets))
		for key, value := range cfg.Identities[i].Secrets {
			switch {
			case strings.HasSuffix(key, "_env"):
				envValue, ok := os.LookupEnv(value)
				if !ok {
					return fmt.Errorf("missing required environment variable %q", value)
				}
				if strings.TrimSpace(envValue) == "" {
					return fmt.Errorf("environment variable %q is empty", value)
				}
				resolved[strings.TrimSuffix(key, "_env")] = envValue

			case strings.HasSuffix(key, "_file"):
				data, err := os.ReadFile(value)
				if err != nil {
					return fmt.Errorf("read secret file %q: %w", value, err)
				}
				secret := strings.TrimSpace(string(data))
				if secret == "" {
					return fmt.Errorf("secret file %q is empty", value)
				}
				resolved[strings.TrimSuffix(key, "_file")] = secret

			default:
				resolved[key] = value
			}
		}
		cfg.Identities[i].Secrets = resolved
	}

	return nil
}
