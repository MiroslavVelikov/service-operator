package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestPostEcho(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		errorMsg       string
		expected       string
		expectedStatus int
	}{
		{
			name:           "should return success even if input message is empty",
			input:          `{"message":""}`,
			expectedStatus: http.StatusOK,
		}, {
			name:           "should echo the passed input",
			input:          `{"message":"Hello world"}`,
			expected:       `{"message":"Echoed Hello world"}`,
			expectedStatus: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			server := NewApiHandler()

			request, _ := http.NewRequest("POST", "/echo", bytes.NewBufferString(testCase.input))
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			assertContentType(t, response, CONTENT_TYPE)
			assertStatus(t, response.Code, testCase.expectedStatus)

			var expected string
			if testCase.errorMsg != "" {
				expected = testCase.errorMsg
			} else {
				expected = testCase.expected
			}

			actual := response.Body.String()
			if !strings.Contains(actual, expected) {
				t.Errorf("expected %s, actual %s", expected, actual)
			}
		})
	}
}

func assertStatus(t testing.TB, actual, expected int) {
	t.Helper()
	if actual != expected {
		t.Errorf("actual status %d, expected %d", actual, expected)
	}
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, expected string) {
	t.Helper()
	if response.Result().Header.Get("Content-Type") != expected {
		t.Errorf("response did not have content-type of %s, actual %v", expected, response.Result().Header)
	}
}
