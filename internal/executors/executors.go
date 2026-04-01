package executors

import (
	"context"
	"fmt"
	"strings"

	emailExecutor "saiao/internal/executors/email"
	httpExecutor "saiao/internal/executors/http"
	shellExecutor "saiao/internal/executors/shell"
	sshExecutor "saiao/internal/executors/ssh"
	"saiao/internal/models"
)

type Executor interface {
	Execute(ctx context.Context, action models.Action, identity models.Identity, input map[string]any) (string, int, error)
}

func ForActionType(actionType string) (Executor, error) {
	switch strings.ToLower(strings.TrimSpace(actionType)) {
	case "http":
		return &httpExecutor.Executor{}, nil
	case "ssh":
		return &sshExecutor.Executor{}, nil
	case "shell":
		return &shellExecutor.Executor{}, nil
	case "email":
		return &emailExecutor.Executor{}, nil
	default:
		return nil, fmt.Errorf("unsupported executor type %q", actionType)
	}
}
