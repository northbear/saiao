package api

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"saiao/internal/auth"
	saiaoerrors "saiao/internal/errors"
	"saiao/internal/models"
)

func RegisterRoutes(mux *http.ServeMux, info BuildInfo, cfg *models.Config, store *auth.TokenStore, logger *slog.Logger) {
	mux.HandleFunc("GET /info", func(w http.ResponseWriter, r *http.Request) {
		_ = r
		writeInfo(w, info)
	})

	mux.HandleFunc("GET /tool-groups/{group_name}/manifest", func(w http.ResponseWriter, r *http.Request) {
		groupName := r.PathValue("group_name")
		if err := authorizeRequest(r, store, groupName); err != nil {
			logger.Warn("manifest_request_denied",
				"tool_group", groupName,
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
				"error", err.Error(),
			)
			writeMappedError(w, err)
			return
		}

		logger.Info("manifest_request",
			"tool_group", groupName,
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
		)
		writeManifest(w, groupName, cfg)
	})

	mux.HandleFunc("POST /tool-groups/{group_name}/actions/{action_name}/invoke", func(w http.ResponseWriter, r *http.Request) {
		groupName := r.PathValue("group_name")
		actionName := r.PathValue("action_name")
		if err := authorizeRequest(r, store, groupName); err != nil {
			logger.Warn("invoke_request_denied",
				"tool_group", groupName,
				"action", actionName,
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
				"error", err.Error(),
			)
			writeMappedError(w, err)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Warn("invoke_request_invalid",
				"tool_group", groupName,
				"action", actionName,
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
				"error", err.Error(),
			)
			writeMappedError(w, fmt.Errorf("%w: could not read request body", saiaoerrors.ErrInvalidInput))
			return
		}

		logger.Info("invoke_request",
			"tool_group", groupName,
			"action", actionName,
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
		)
		writeInvokeResult(w, groupName, actionName, body, cfg, logger)
	})
}

func authorizeRequest(r *http.Request, store *auth.TokenStore, groupName string) error {
	token, ok := auth.ParseBearerToken(r.Header.Get("Authorization"))
	if !ok {
		return fmt.Errorf("%w: missing or invalid bearer token", saiaoerrors.ErrUnauthorized)
	}

	authorizedGroup, ok := store.GroupForToken(token)
	if !ok || authorizedGroup != groupName {
		return fmt.Errorf("%w: token does not grant access to tool group %q", saiaoerrors.ErrUnauthorized, groupName)
	}

	return nil
}
