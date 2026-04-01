package invoke

import (
	"context"

	"saiao/internal/executors"
	"saiao/internal/models"
)

func Dispatch(ctx context.Context, action models.Action, identity models.Identity, input map[string]any) (string, int, error) {
	executor, err := executors.ForActionType(action.Type)
	if err != nil {
		return "", 0, err
	}

	return executor.Execute(ctx, action, identity, input)
}
