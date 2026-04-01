package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"saiao/internal/api"
	"saiao/internal/auth"
	"saiao/internal/config"
	"saiao/internal/logging"
	"saiao/internal/models"
)

type Options struct {
	ConfigPath string
}

func DefaultOptions() Options {
	return Options{
		ConfigPath: "/etc/saiao/config.yaml",
	}
}

func Run(opts Options) error {
	return RunWithContext(context.Background(), opts)
}

func RunWithContext(ctx context.Context, opts Options) error {
	if ctx == nil {
		return errors.New("nil context")
	}

	if opts.ConfigPath == "" {
		opts = DefaultOptions()
	}

	cfg, err := config.Load(opts.ConfigPath)
	if err != nil {
		return err
	}

	if err := config.Validate(cfg); err != nil {
		return err
	}

	logger := logging.New(cfg.Logging.Format, cfg.Logging.Level)

	buildInfo := api.BuildInfo{
		Service: "saiao",
		Version: "0.1.0",
		Commit:  "dev",
	}

	srv := api.NewServer(cfg.Server.Listen, buildInfo, cfg, buildTokenStore(cfg), logger)

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Server.ShutdownTimeoutSeconds)*time.Second)
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

func buildTokenStore(cfg *models.Config) *auth.TokenStore {
	store := auth.NewTokenStore()
	if cfg == nil {
		return store
	}

	for _, group := range cfg.ToolGroups {
		token := os.Getenv(group.AccessTokenEnv)
		if token == "" {
			continue
		}
		store.Register(token, group.Name)
	}

	return store
}
