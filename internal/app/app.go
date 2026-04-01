package app

import (
	"context"
	"errors"
)

var ErrNotImplemented = errors.New("not implemented")

func Run(ctx context.Context) error {
	_ = ctx
	return ErrNotImplemented
}
