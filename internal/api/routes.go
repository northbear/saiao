package api

import (
	"fmt"
	"io"
	"net/http"

	"saiao/internal/auth"
	saiaoerrors "saiao/internal/errors"
	"saiao/internal/models"
)

func RegisterRoutes(mux *http.ServeMux, info BuildInfo, cfg *models.Config, store *auth.TokenStore) {
	mux.HandleFunc("GET /info", func(w http.ResponseWriter, r *http.Request) {
		_ = r
		writeInfo(w, info)
	})

	mux.HandleFunc("GET /tool-groups/{group_name}/manifest", func(w http.ResponseWriter, r *http.Request) {
		groupName := r.PathValue("group_name")
		if err := authorizeRequest(r, store, groupName); err != nil {
			writeMappedError(w, err)
			return
		}

		writeManifest(w, groupName, cfg)
	})

	mux.HandleFunc("POST /tool-groups/{group_name}/actions/{action_name}/invoke", func(w http.ResponseWriter, r *http.Request) {
		groupName := r.PathValue("group_name")
		actionName := r.PathValue("action_name")
		if err := authorizeRequest(r, store, groupName); err != nil {
			writeMappedError(w, err)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeMappedError(w, fmt.Errorf("%w: could not read request body", saiaoerrors.ErrInvalidInput))
			return
		}

		writeInvokeResult(w, groupName, actionName, body, cfg)
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
