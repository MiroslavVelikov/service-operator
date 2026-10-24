package main

import (
	"net/http"

	"github.com/service-operator/internal/api"

	log "github.com/sirupsen/logrus"
)

func main() {
	srv := api.NewApiHandler()
	log.Fatal(http.ListenAndServe(":8080", srv))
}
