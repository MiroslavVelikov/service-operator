package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/service-operator/internal/api/health"
	"github.com/service-operator/internal/config"
)

const (
	MAX_POST_SIZE = 1048576

	CONTENT_TYPE = "application/json"
)

func NewApiHandler() *http.ServeMux {
	router := http.NewServeMux()

	health.SetupHealthHandler(router)

	router.HandleFunc("GET /version", getVersion)
	router.HandleFunc("POST /echo", replyEcho)

	return router
}

func getVersion(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", CONTENT_TYPE)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"version": "1.0.0"}`))
}

func replyEcho(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", CONTENT_TYPE)
	req.Body = http.MaxBytesReader(w, req.Body, MAX_POST_SIZE)

	var input config.EchoInput

	dec := json.NewDecoder(req.Body)
	err := dec.Decode(&input)
	if err != nil {
		if err == io.EOF {
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		http.Error(w, "Request body must contain only the required information", http.StatusBadRequest)
		return
	}

	if input.Message != "" {
		output := config.EchoInput{Message: "Echoed " + input.Message}
		err := json.NewEncoder(w).Encode(output)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	w.WriteHeader(http.StatusOK)
}
