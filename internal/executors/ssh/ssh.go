package ssh

import "saiao/internal/models"

type Executor struct{}

func (e *Executor) Execute(action models.Action, input map[string]any) (string, int, error) {
	_ = e
	_ = action
	_ = input
	return "", 0, nil
}
