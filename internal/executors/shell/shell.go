package shell

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	saiaoerrors "saiao/internal/errors"
	"saiao/internal/models"
	"saiao/internal/render"
)

type Executor struct{}

func (e *Executor) Execute(ctx context.Context, action models.Action, identity models.Identity, input map[string]any) (string, int, error) {
	command, err := render.Template(action.CommandTemplate, identity, input)
	if err != nil {
		return "", 0, fmt.Errorf("%w: %v", saiaoerrors.ErrInvalidInput, err)
	}

	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	if action.WorkingDirectory != "" {
		cmd.Dir = action.WorkingDirectory
	}

	if len(action.Environment) > 0 {
		cmd.Env = os.Environ()
		for key, value := range action.Environment {
			rendered, err := render.Template(value, identity, input)
			if err != nil {
				return "", 0, fmt.Errorf("%w: %v", saiaoerrors.ErrInvalidInput, err)
			}
			cmd.Env = append(cmd.Env, key+"="+rendered)
		}
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() != nil {
			return string(output), -1, fmt.Errorf("%w: %v", saiaoerrors.ErrTimeout, ctx.Err())
		}
		if exitErr, ok := err.(*exec.ExitError); ok {
			return string(output), exitErr.ExitCode(), fmt.Errorf("%w: %v", saiaoerrors.ErrExecutionFailed, err)
		}
		return string(output), -1, fmt.Errorf("%w: %v", saiaoerrors.ErrExecutionFailed, err)
	}

	return string(output), 0, nil
}
