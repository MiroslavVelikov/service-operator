package api

import "net/http"

func NewApiHandler() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("GET /version", getVersion)
	router.HandleFunc("POST /echo", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})

	return router
}

func getVersion(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"version": "1.0.0"}`))
}
