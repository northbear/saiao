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
		requestID := requestIDFromRequest(r)
		writeInfo(w, requestID, info)
	})

	mux.HandleFunc("GET /tool-groups/{group_name}/manifest", func(w http.ResponseWriter, r *http.Request) {
		requestID := requestIDFromRequest(r)
		groupName := r.PathValue("group_name")
		requestLogger := logger.With("request_id", requestID)
		if err := authorizeRequest(r, store, groupName); err != nil {
			requestLogger.Warn("manifest_request_denied",
				"tool_group", groupName,
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
				"error", err.Error(),
			)
			writeMappedError(w, requestID, err)
			return
		}

		requestLogger.Info("manifest_request",
			"tool_group", groupName,
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
		)
		writeManifest(w, requestID, groupName, cfg)
	})

	mux.HandleFunc("POST /tool-groups/{group_name}/actions/{action_name}/invoke", func(w http.ResponseWriter, r *http.Request) {
		requestID := requestIDFromRequest(r)
		groupName := r.PathValue("group_name")
		actionName := r.PathValue("action_name")
		requestLogger := logger.With("request_id", requestID)
		if err := authorizeRequest(r, store, groupName); err != nil {
			requestLogger.Warn("invoke_request_denied",
				"tool_group", groupName,
				"action", actionName,
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
				"error", err.Error(),
			)
			writeMappedError(w, requestID, err)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			requestLogger.Warn("invoke_request_invalid",
				"tool_group", groupName,
				"action", actionName,
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
				"error", err.Error(),
			)
			writeMappedError(w, requestID, fmt.Errorf("%w: could not read request body", saiaoerrors.ErrInvalidInput))
			return
		}

		requestLogger.Info("invoke_request",
			"tool_group", groupName,
			"action", actionName,
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
		)
		writeInvokeResult(w, requestID, groupName, actionName, body, cfg, requestLogger)
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
