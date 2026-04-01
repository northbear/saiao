package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux, info BuildInfo) {
	mux.HandleFunc("GET /info", func(w http.ResponseWriter, r *http.Request) {
		_ = r
		writeInfo(w, info)
	})
}
