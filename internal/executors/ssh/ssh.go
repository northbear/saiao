package ssh

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"

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

	target := action.Host
	if identity.Username != "" {
		target = identity.Username + "@" + target
	}

	args := make([]string, 0, 4)
	if action.Port > 0 {
		args = append(args, "-p", strconv.Itoa(action.Port))
	}
	args = append(args, target, command)

	cmd := exec.CommandContext(ctx, "ssh", args...)
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
