package api

import (
	"encoding/json"
	"net/http"

	"saiao/internal/models"
)

func NewServer(listenAddress string, info appInfo) *http.Server {
	mux := http.NewServeMux()
	RegisterRoutes(mux, info)

	return &http.Server{
		Addr:    listenAddress,
		Handler: mux,
	}
}

type appInfo struct {
	Service string
	Version string
	Commit  string
}

func writeJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(v)
}

func writeInfo(w http.ResponseWriter, info appInfo) {
	writeJSON(w, http.StatusOK, models.InfoResponse{
		Status:  "ok",
		Service: info.Service,
		Version: info.Version,
		Commit:  info.Commit,
	})
}
