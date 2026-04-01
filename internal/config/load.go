package config

import (
	"errors"

	"saiao/internal/models"
)

var ErrNotImplemented = errors.New("config loading not implemented")

func Load(_ string) (*models.Config, error) {
	return nil, ErrNotImplemented
}
