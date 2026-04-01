package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"saiao/internal/api"
)

var ErrNotImplemented = errors.New("not implemented")

type Options struct {
	ListenAddress string
	BuildInfo     BuildInfo
	Server        *http.Server
}

type BuildInfo struct {
	Service string
	Version string
	Commit  string
}

func DefaultOptions() Options {
	return Options{
		ListenAddress: ":8080",
		BuildInfo: BuildInfo{
			Service: "saiao",
			Version: "0.1.0",
			Commit:  "dev",
		},
	}
}

func Run(ctx context.Context) error {
	return RunWithOptions(ctx, DefaultOptions())
}

func RunWithOptions(ctx context.Context, opts Options) error {
	if ctx == nil {
		return errors.New("nil context")
	}

	select {
	case <-ctx.Done():
		return nil
	default:
	}

	srv := opts.Server
	if srv == nil {
		srv = api.NewServer(opts.ListenAddress, api.BuildInfo{
			Service: opts.BuildInfo.Service,
			Version: opts.BuildInfo.Version,
			Commit:  opts.BuildInfo.Commit,
		})
	}

	errCh := make(chan error, 1)

	go func() {
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown: %w", err)
		}
		return nil
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
}
