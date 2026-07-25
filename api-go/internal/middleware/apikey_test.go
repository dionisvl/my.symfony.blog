package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"api-go/internal/middleware"
)

func TestAPIKeyAuthRejectsEmptyConfiguredKey(t *testing.T) {
	t.Parallel()

	called := false
	handler := middleware.APIKeyAuth("")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	req.Header.Set("X-API-Key", "secret")

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.False(t, called)
	assert.JSONEq(t, `{"error":"invalid API key"}`, recorder.Body.String())
}

func TestAPIKeyAuthAllowsMatchingKey(t *testing.T) {
	t.Parallel()

	called := false
	handler := middleware.APIKeyAuth("secret")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	req.Header.Set("X-API-Key", "secret")

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.True(t, called)
}
