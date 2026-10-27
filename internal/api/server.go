package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	PORT_KEY = "PORT"

	FALLBACK_PORT       = "8080"
	LAST_AVAILABLE_PORT = 65535
)

func RunServer(ctx context.Context, shutdownTimeout time.Duration) error {
	port, err := getPort()
	if err != nil {
		return err
	}

	srv := http.Server{
		Addr:    ":" + port,
		Handler: NewApiHandler(),
	}

	srvErr := make(chan error, 1)
	go func() {
		log.Infof("Listening on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			srvErr <- err
		}
		close(srvErr)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	select {
	case err := <-srvErr:
		return err
	case <-stop:
		log.Infof("Shutting down...")
	case <-ctx.Done():
		log.Infof("Context cancelled, shutting down...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		if closeErr := srv.Close(); closeErr != nil {
			return errors.Join(err, closeErr)
		}

		return err
	}

	log.Infof("Server shutdown complete")
	return nil
}

func getPort() (string, error) {
	port := os.Getenv(PORT_KEY)
	if port == "" {
		log.Warnf("PORT environment variable not set, using default: %s", FALLBACK_PORT)
		return FALLBACK_PORT, nil
	}
	iPort, err := strconv.Atoi(port)
	if err != nil {
		return "", fmt.Errorf("invalid PORT %q: must be an integer", port)
	} else if iPort < 1 || iPort > LAST_AVAILABLE_PORT {
		return "", fmt.Errorf("PORT out of range: %s, should be [1; %d]", port, LAST_AVAILABLE_PORT)
	}

	return port, nil
}
