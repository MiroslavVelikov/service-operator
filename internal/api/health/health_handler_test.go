package health

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHealth(t *testing.T) {
	t.Run("should return success [OK, 200]", func(t *testing.T) {
		router := http.NewServeMux()
		SetupHealthHandler(router)

		request, _ := http.NewRequest("GET", "/healthz", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		if response.Code != http.StatusOK {
			t.Errorf("expected %d, actual %d", http.StatusOK, response.Code)
		}
		if response.Body.String() != "OK" {
			t.Errorf("expected %s, actual %s", "OK", response.Body.String())
		}
	})
}

func TestGetReady(t *testing.T) {
	t.Run("should return success [READY, 200]", func(t *testing.T) {
		router := http.NewServeMux()
		SetupHealthHandler(router)

		request, _ := http.NewRequest("GET", "/readyz", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		if response.Code != http.StatusOK {
			t.Errorf("expected %d, actual %d", http.StatusOK, response.Code)
		}
		if response.Body.String() != "READY" {
			t.Errorf("expected %s, actual %s", "READY", response.Body.String())
		}
	})
}
