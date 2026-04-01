package executors

import "saiao/internal/models"

type Executor interface {
	Execute(action models.Action, input map[string]any) (string, int, error)
}
