package health

import "net/http"

func SetupHealthHandler(router *http.ServeMux) {
	router.HandleFunc("GET /healthz", getHealth)
	router.HandleFunc("GET /readyz", getReadiness)
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func getReadiness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("READY"))
}
