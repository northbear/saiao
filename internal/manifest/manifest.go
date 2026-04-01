package manifest

import (
	"fmt"

	saiaoerrors "saiao/internal/errors"
	"saiao/internal/models"
)

func Build(groupName string, cfg *models.Config) (models.ManifestResponse, error) {
	if cfg == nil {
		return models.ManifestResponse{}, fmt.Errorf("%w: config is nil", saiaoerrors.ErrInternal)
	}

	var group *models.ToolGroup
	for i := range cfg.ToolGroups {
		if cfg.ToolGroups[i].Name == groupName {
			group = &cfg.ToolGroups[i]
			break
		}
	}
	if group == nil {
		return models.ManifestResponse{}, fmt.Errorf("%w: tool group %q", saiaoerrors.ErrNotFound, groupName)
	}

	actionsByName := make(map[string]models.Action, len(cfg.Actions))
	for _, action := range cfg.Actions {
		actionsByName[action.Name] = action
	}

	tools := make([]models.ToolManifest, 0, len(group.Actions))
	for _, actionName := range group.Actions {
		action, ok := actionsByName[actionName]
		if !ok {
			return models.ManifestResponse{}, fmt.Errorf("%w: action %q", saiaoerrors.ErrNotFound, actionName)
		}

		tools = append(tools, models.ToolManifest{
			Name:        action.Name,
			Description: action.Description,
			InputSchema: action.InputSchema,
		})
	}

	return models.ManifestResponse{Tools: tools}, nil
}
