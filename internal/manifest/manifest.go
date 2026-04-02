package manifest

import (
	"fmt"

	saiaoerrors "saiao/internal/errors"
	"saiao/internal/models"
)

const (
	FormatSAIAO  = "saiao"
	FormatOpenAI = "openai"
)

func Build(groupName string, cfg *models.Config) (models.ManifestResponse, error) {
	tools, err := buildTools(groupName, cfg)
	if err != nil {
		return models.ManifestResponse{}, err
	}

	return models.ManifestResponse{
		Format: FormatSAIAO,
		Tools:  tools,
	}, nil
}

func BuildOpenAI(groupName string, cfg *models.Config) (models.OpenAIManifestResponse, error) {
	tools, err := buildTools(groupName, cfg)
	if err != nil {
		return models.OpenAIManifestResponse{}, err
	}

	openAITools := make([]models.OpenAITool, 0, len(tools))
	for _, tool := range tools {
		parameters := tool.InputSchema
		if len(parameters) == 0 {
			parameters = map[string]any{
				"type":       "object",
				"properties": map[string]any{},
				"required":   []any{},
			}
		}

		openAITools = append(openAITools, models.OpenAITool{
			Type:        "function",
			Name:        tool.Name,
			Description: tool.Description,
			Parameters:  parameters,
			Strict:      false,
		})
	}

	return models.OpenAIManifestResponse{
		Format: FormatOpenAI,
		Tools:  openAITools,
	}, nil
}

func buildTools(groupName string, cfg *models.Config) ([]models.ToolManifest, error) {
	if cfg == nil {
		return nil, fmt.Errorf("%w: config is nil", saiaoerrors.ErrInternal)
	}

	var group *models.ToolGroup
	for i := range cfg.ToolGroups {
		if cfg.ToolGroups[i].Name == groupName {
			group = &cfg.ToolGroups[i]
			break
		}
	}
	if group == nil {
		return nil, fmt.Errorf("%w: tool group %q", saiaoerrors.ErrNotFound, groupName)
	}

	actionsByName := make(map[string]models.Action, len(cfg.Actions))
	for _, action := range cfg.Actions {
		actionsByName[action.Name] = action
	}

	tools := make([]models.ToolManifest, 0, len(group.Actions))
	for _, actionName := range group.Actions {
		action, ok := actionsByName[actionName]
		if !ok {
			return nil, fmt.Errorf("%w: action %q", saiaoerrors.ErrNotFound, actionName)
		}

		tools = append(tools, models.ToolManifest{
			Name:        action.Name,
			Description: action.Description,
			InputSchema: action.InputSchema,
		})
	}

	return tools, nil
}
