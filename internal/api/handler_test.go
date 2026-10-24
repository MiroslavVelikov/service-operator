package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetVersion(t *testing.T) {
	server := NewApiHandler()

	t.Run("should return current api version", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "/version", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertContentType(t, response, "application/json")
		assertStatus(t, response.Code, http.StatusOK)

		expected := `{"version": "1.0.0"}`
		actual := response.Body.String()

		if actual != expected {
			t.Errorf("expected %s, actual %s", expected, actual)
		}
	})
}

func assertStatus(t testing.TB, actual, expected int) {
	t.Helper()
	if actual != expected {
		t.Errorf("actual status %d, expected %d", actual, expected)
	}
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, expected string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != expected {
		t.Errorf("response did not have content-type of %s, got %v", expected, response.Result().Header)
	}
}
