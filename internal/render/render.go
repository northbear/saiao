package render

import (
	"fmt"
	"regexp"
	"strings"

	"saiao/internal/models"
)

var placeholderPattern = regexp.MustCompile(`{{\s*([a-zA-Z0-9_.]+)\s*}}`)

func Placeholders(template string) []string {
	matches := placeholderPattern.FindAllStringSubmatch(template, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(matches))
	placeholders := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) != 2 {
			continue
		}
		name := match[1]
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		placeholders = append(placeholders, name)
	}

	return placeholders
}

func Template(template string, identity models.Identity, input map[string]any) (string, error) {
	context := make(map[string]any, len(input)+1)
	for key, value := range input {
		context[key] = value
	}

	identityContext := make(map[string]any, len(identity.Secrets)+6)
	identityContext["name"] = identity.Name
	identityContext["username"] = identity.Username
	identityContext["email"] = identity.Email
	identityContext["display_name"] = identity.DisplayName
	identityContext["principal"] = identity.Principal
	identityContext["domain"] = identity.Domain
	for key, value := range identity.Secrets {
		identityContext[key] = value
	}
	context["identity"] = identityContext

	var renderErr error
	rendered := placeholderPattern.ReplaceAllStringFunc(template, func(match string) string {
		if renderErr != nil {
			return ""
		}

		parts := placeholderPattern.FindStringSubmatch(match)
		if len(parts) != 2 {
			renderErr = fmt.Errorf("invalid template placeholder %q", match)
			return ""
		}

		value, err := resolvePath(context, parts[1])
		if err != nil {
			renderErr = err
			return ""
		}

		return fmt.Sprint(value)
	})

	if renderErr != nil {
		return "", renderErr
	}

	return rendered, nil
}

func resolvePath(context map[string]any, path string) (any, error) {
	current := any(context)
	for _, segment := range strings.Split(path, ".") {
		node, ok := current.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("template path %q is invalid", path)
		}

		next, ok := node[segment]
		if !ok {
			return nil, fmt.Errorf("template path %q is undefined", path)
		}
		current = next
	}

	return current, nil
}
