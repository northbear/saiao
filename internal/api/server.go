package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"saiao/internal/auth"
	saiaoerrors "saiao/internal/errors"
	"saiao/internal/invoke"
	"saiao/internal/manifest"
	"saiao/internal/models"
)

type BuildInfo struct {
	Service string
	Version string
	Commit  string
}

func NewServer(listenAddress string, info BuildInfo, cfg *models.Config, store *auth.TokenStore, logger *slog.Logger) *http.Server {
	mux := http.NewServeMux()
	RegisterRoutes(mux, info, cfg, store, defaultLogger(logger))

	var readTimeout time.Duration
	var writeTimeout time.Duration
	if cfg != nil {
		readTimeout = time.Duration(cfg.Server.ReadTimeoutSeconds) * time.Second
		writeTimeout = time.Duration(cfg.Server.WriteTimeoutSeconds) * time.Second
	}

	return &http.Server{
		Addr:         listenAddress,
		Handler:      withRequestID(mux),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}

func withRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := requestIDFromRequest(r)
		r.Header.Set(requestIDHeader, requestID)
		w.Header().Set(requestIDHeader, requestID)
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(v)
}

func writeInfo(w http.ResponseWriter, requestID string, info BuildInfo) {
	writeJSON(w, http.StatusOK, models.InfoResponse{
		RequestID: requestID,
		Status:    "ok",
		Service:   info.Service,
		Version:   info.Version,
		Commit:    info.Commit,
	})
}

func writeManifest(w http.ResponseWriter, requestID string, groupName string, cfg *models.Config) {
	response, err := manifest.Build(groupName, cfg)
	if err != nil {
		writeMappedError(w, requestID, err)
		return
	}
	response.RequestID = requestID
	writeJSON(w, http.StatusOK, response)
}

func writeInvokeResult(w http.ResponseWriter, requestID string, groupName, actionName string, body []byte, cfg *models.Config, logger *slog.Logger) {
	response, err := invoke.Execute(requestID, groupName, actionName, body, cfg, defaultLogger(logger))
	if err != nil {
		writeMappedError(w, requestID, err)
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func writeMappedError(w http.ResponseWriter, requestID string, err error) {
	statusCode := http.StatusInternalServerError
	code := "internal_error"

	switch {
	case errors.Is(err, saiaoerrors.ErrUnauthorized):
		statusCode = http.StatusUnauthorized
		code = "unauthorized"
	case errors.Is(err, saiaoerrors.ErrNotFound):
		statusCode = http.StatusNotFound
		code = "not_found"
	case errors.Is(err, saiaoerrors.ErrInvalidInput):
		statusCode = http.StatusBadRequest
		code = "invalid_input"
	case errors.Is(err, saiaoerrors.ErrExecutionFailed):
		statusCode = http.StatusBadGateway
		code = "execution_failed"
	case errors.Is(err, saiaoerrors.ErrTimeout):
		statusCode = http.StatusGatewayTimeout
		code = "timeout"
	case errors.Is(err, saiaoerrors.ErrInternal):
		statusCode = http.StatusInternalServerError
		code = "internal_error"
	}

	writeJSON(w, statusCode, models.NewErrorResponse(requestID, code, err.Error()))
}

func defaultLogger(logger *slog.Logger) *slog.Logger {
	if logger == nil {
		return slog.Default()
	}
	return logger
}
