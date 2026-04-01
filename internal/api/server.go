package api

import (
	"encoding/json"
	"net/http"

	"saiao/internal/models"
)

type BuildInfo struct {
	Service string
	Version string
	Commit  string
}

func NewServer(listenAddress string, info BuildInfo) *http.Server {
	mux := http.NewServeMux()
	RegisterRoutes(mux, info)

	return &http.Server{
		Addr:    listenAddress,
		Handler: mux,
	}
}

func RegisterRoutes(mux *http.ServeMux, info BuildInfo) {
	mux.HandleFunc("GET /info", func(w http.ResponseWriter, r *http.Request) {
		_ = r
		writeInfo(w, info)
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(v)
}

func writeInfo(w http.ResponseWriter, info BuildInfo) {
	writeJSON(w, http.StatusOK, models.InfoResponse{
		Status:  "ok",
		Service: info.Service,
		Version: info.Version,
		Commit:  info.Commit,
	})
}
