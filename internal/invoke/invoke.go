package invoke

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	saiaoerrors "saiao/internal/errors"
	"saiao/internal/models"
)

func Execute(groupName string, actionName string, body []byte, cfg *models.Config) (models.SuccessResponse, error) {
	if cfg == nil {
		return models.SuccessResponse{}, fmt.Errorf("%w: config is nil", saiaoerrors.ErrInternal)
	}

	group, action, identity, err := resolve(groupName, actionName, cfg)
	if err != nil {
		return models.SuccessResponse{}, err
	}
	_ = group

	input, err := decodeInput(body)
	if err != nil {
		return models.SuccessResponse{}, err
	}

	if err := ValidateInput(action.InputSchema, input); err != nil {
		return models.SuccessResponse{}, err
	}

	ctx := context.Background()
	cancel := func() {}
	if action.TimeoutSeconds > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(action.TimeoutSeconds)*time.Second)
	}
	defer cancel()

	output, exitCode, err := Dispatch(ctx, action, identity, input)
	if err != nil {
		return models.SuccessResponse{}, err
	}

	return models.NewSuccessResponse(output, exitCode), nil
}

func resolve(groupName, actionName string, cfg *models.Config) (models.ToolGroup, models.Action, models.Identity, error) {
	var group models.ToolGroup
	foundGroup := false
	for _, candidate := range cfg.ToolGroups {
		if candidate.Name == groupName {
			group = candidate
			foundGroup = true
			break
		}
	}
	if !foundGroup {
		return models.ToolGroup{}, models.Action{}, models.Identity{}, fmt.Errorf("%w: tool group %q", saiaoerrors.ErrNotFound, groupName)
	}

	allowed := false
	for _, name := range group.Actions {
		if name == actionName {
			allowed = true
			break
		}
	}
	if !allowed {
		return models.ToolGroup{}, models.Action{}, models.Identity{}, fmt.Errorf("%w: action %q is not in tool group %q", saiaoerrors.ErrNotFound, actionName, groupName)
	}

	var action models.Action
	foundAction := false
	for _, candidate := range cfg.Actions {
		if candidate.Name == actionName {
			action = candidate
			foundAction = true
			break
		}
	}
	if !foundAction {
		return models.ToolGroup{}, models.Action{}, models.Identity{}, fmt.Errorf("%w: action %q", saiaoerrors.ErrNotFound, actionName)
	}

	var identity models.Identity
	foundIdentity := false
	for _, candidate := range cfg.Identities {
		if candidate.Name == action.Identity {
			identity = candidate
			foundIdentity = true
			break
		}
	}
	if !foundIdentity {
		return models.ToolGroup{}, models.Action{}, models.Identity{}, fmt.Errorf("%w: identity %q", saiaoerrors.ErrNotFound, action.Identity)
	}

	return group, action, identity, nil
}

func decodeInput(body []byte) (map[string]any, error) {
	if len(body) == 0 {
		return map[string]any{}, nil
	}

	var input map[string]any
	if err := json.Unmarshal(body, &input); err != nil {
		return nil, fmt.Errorf("%w: request body must be a JSON object", saiaoerrors.ErrInvalidInput)
	}
	if input == nil {
		input = map[string]any{}
	}

	return input, nil
}
