package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	saiaoerrors "saiao/internal/errors"
	"saiao/internal/models"
	"saiao/internal/render"
)

type Executor struct{}

func (e *Executor) Execute(ctx context.Context, action models.Action, identity models.Identity, input map[string]any) (string, int, error) {
	url, err := render.Template(action.URL, identity, input)
	if err != nil {
		return "", 0, fmt.Errorf("%w: %v", saiaoerrors.ErrInvalidInput, err)
	}

	var bodyReader io.Reader
	if strings.TrimSpace(action.BodyTemplate) != "" {
		body, err := render.Template(action.BodyTemplate, identity, input)
		if err != nil {
			return "", 0, fmt.Errorf("%w: %v", saiaoerrors.ErrInvalidInput, err)
		}
		bodyReader = strings.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, action.Method, url, bodyReader)
	if err != nil {
		return "", 0, fmt.Errorf("%w: %v", saiaoerrors.ErrExecutionFailed, err)
	}

	for key, value := range action.Headers {
		rendered, err := render.Template(value, identity, input)
		if err != nil {
			return "", 0, fmt.Errorf("%w: %v", saiaoerrors.ErrInvalidInput, err)
		}
		req.Header.Set(key, rendered)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return "", -1, fmt.Errorf("%w: %v", saiaoerrors.ErrTimeout, ctx.Err())
		}
		return "", -1, fmt.Errorf("%w: %v", saiaoerrors.ErrExecutionFailed, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, fmt.Errorf("%w: %v", saiaoerrors.ErrExecutionFailed, err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return string(data), resp.StatusCode, fmt.Errorf("%w: remote returned status %d", saiaoerrors.ErrExecutionFailed, resp.StatusCode)
	}

	return string(data), resp.StatusCode, nil
}
