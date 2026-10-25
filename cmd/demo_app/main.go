package main

import (
	"context"
	"time"

	"github.com/service-operator/internal/api"
	log "github.com/sirupsen/logrus"
)

const SHUTDOWN_TIMEOUT = time.Second * 5

func main() {
	if err := api.RunServer(context.Background(), SHUTDOWN_TIMEOUT); err != nil {
		log.Fatal(err)
	}
}
