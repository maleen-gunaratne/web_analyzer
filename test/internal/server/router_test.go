package server

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"web-analyzer/internal/server"
)

func TestSetupRouter(t *testing.T) {
	t.Run("Router Configuration", func(t *testing.T) {

		router := server.SetupRouter()

		ts := httptest.NewServer(router)
		defer ts.Close()

		resp, err := http.Get(ts.URL + "/health")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		resp, err = http.Get(ts.URL + "/metrics")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		resp, err = http.Get(ts.URL + "/api/v1/analyze")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		resp, err = http.Get(ts.URL + "/health")
		assert.NoError(t, err)
	})
}

func TestMiddleware(t *testing.T) {
	t.Run("RequestID Middleware", func(t *testing.T) {
		router := server.SetupRouter()

		ts := httptest.NewServer(router)
		defer ts.Close()

		resp, err := http.Get(ts.URL + "/health")
		assert.NoError(t, err)

		requestID := resp.Header.Get("X-Request-ID")
		assert.NotEmpty(t, requestID)
	})

}
