package config

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"saiao/internal/models"
	"saiao/internal/render"
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
		if err := validateIdentitySecrets(cfg.Identities[i].Secrets, i); err != nil {
			return err
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
		cfg.Actions[i].Type = strings.TrimSpace(strings.ToLower(cfg.Actions[i].Type))
		if strings.TrimSpace(cfg.Actions[i].Identity) == "" {
			return fmt.Errorf("actions[%d].identity is required", i)
		}
		if _, ok := identityNames[cfg.Actions[i].Identity]; !ok {
			return fmt.Errorf("actions[%d].identity references unknown identity %q", i, cfg.Actions[i].Identity)
		}
		if err := validateAction(&cfg.Actions[i], i); err != nil {
			return err
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

func validateAction(action *models.Action, index int) error {
	switch action.Type {
	case "http":
		if strings.TrimSpace(action.Method) == "" {
			return fmt.Errorf("actions[%d].method is required for http actions", index)
		}
		action.Method = strings.ToUpper(strings.TrimSpace(action.Method))
		if strings.TrimSpace(action.URL) == "" {
			return fmt.Errorf("actions[%d].url is required for http actions", index)
		}
	case "ssh":
		if strings.TrimSpace(action.Host) == "" {
			return fmt.Errorf("actions[%d].host is required for ssh actions", index)
		}
		if strings.TrimSpace(action.CommandTemplate) == "" {
			return fmt.Errorf("actions[%d].command_template is required for ssh actions", index)
		}
	case "shell":
		if strings.TrimSpace(action.CommandTemplate) == "" {
			return fmt.Errorf("actions[%d].command_template is required for shell actions", index)
		}
	case "email":
		if len(action.To) == 0 {
			return fmt.Errorf("actions[%d].to is required for email actions", index)
		}
		if strings.TrimSpace(action.SubjectTemplate) == "" {
			return fmt.Errorf("actions[%d].subject_template is required for email actions", index)
		}
		if strings.TrimSpace(action.BodyTemplate) == "" {
			return fmt.Errorf("actions[%d].body_template is required for email actions", index)
		}
	default:
		return fmt.Errorf("actions[%d].type %q is not supported", index, action.Type)
	}

	if err := ValidateInputSchema(action.InputSchema); err != nil {
		return fmt.Errorf("actions[%d].input_schema: %w", index, err)
	}
	if err := validateActionTemplates(*action, index); err != nil {
		return err
	}

	return nil
}

func validateIdentitySecrets(secrets map[string]string, index int) error {
	if len(secrets) == 0 {
		return nil
	}

	sources := make(map[string]string, len(secrets))
	for key := range secrets {
		switch {
		case strings.HasSuffix(key, "_env"):
			base := strings.TrimSuffix(key, "_env")
			if previous, ok := sources[base]; ok {
				return fmt.Errorf("identities[%d].secrets defines multiple sources for %q: %s and %s", index, base, previous, key)
			}
			sources[base] = key
		case strings.HasSuffix(key, "_file"):
			base := strings.TrimSuffix(key, "_file")
			if previous, ok := sources[base]; ok {
				return fmt.Errorf("identities[%d].secrets defines multiple sources for %q: %s and %s", index, base, previous, key)
			}
			sources[base] = key
		}
	}

	return nil
}

func validateActionTemplates(action models.Action, index int) error {
	allowedInputs := allowedInputFields(action.InputSchema)

	check := func(fieldName, template string, allowIdentity bool) error {
		for _, placeholder := range render.Placeholders(template) {
			if strings.HasPrefix(placeholder, "identity.") {
				if !allowIdentity {
					return fmt.Errorf("actions[%d].%s must not reference identity placeholders", index, fieldName)
				}
				if placeholder == "identity" || len(strings.Split(placeholder, ".")) != 2 {
					return fmt.Errorf("actions[%d].%s references unsupported identity placeholder %q", index, fieldName, placeholder)
				}
				continue
			}

			if !slices.Contains(allowedInputs, placeholder) {
				return fmt.Errorf("actions[%d].%s references undeclared input %q", index, fieldName, placeholder)
			}
		}
		return nil
	}

	switch action.Type {
	case "http":
		if err := check("url", action.URL, false); err != nil {
			return err
		}
		if err := check("body_template", action.BodyTemplate, true); err != nil {
			return err
		}
		for headerName, value := range action.Headers {
			if err := check("headers."+headerName, value, true); err != nil {
				return err
			}
		}
	case "ssh", "shell":
		if err := check("command_template", action.CommandTemplate, false); err != nil {
			return err
		}
		for envName, value := range action.Environment {
			if err := check("environment."+envName, value, true); err != nil {
				return err
			}
		}
	case "email":
		if err := check("subject_template", action.SubjectTemplate, false); err != nil {
			return err
		}
		if err := check("body_template", action.BodyTemplate, false); err != nil {
			return err
		}
		for i, recipient := range action.To {
			if err := check(fmt.Sprintf("to[%d]", i), recipient, false); err != nil {
				return err
			}
		}
		for i, recipient := range action.Cc {
			if err := check(fmt.Sprintf("cc[%d]", i), recipient, false); err != nil {
				return err
			}
		}
		for i, recipient := range action.Bcc {
			if err := check(fmt.Sprintf("bcc[%d]", i), recipient, false); err != nil {
				return err
			}
		}
	}

	return nil
}

func allowedInputFields(schema map[string]any) []string {
	if len(schema) == 0 {
		return nil
	}

	properties, ok := schema["properties"].(map[string]any)
	if !ok || len(properties) == 0 {
		return nil
	}

	fields := make([]string, 0, len(properties))
	for name := range properties {
		fields = append(fields, name)
	}
	slices.Sort(fields)
	return fields
}
