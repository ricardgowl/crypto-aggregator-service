package httpapi

import (
	"crypto-aggregator-service/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func newTestServer() *HTTPServer {
	logger := zap.NewNop().Sugar()
	return NewHTTPServer(logger, config.ServerConfigurations{Port: 3000})
}

func TestHealthController_LivenessCheck(t *testing.T) {
	server := newTestServer()
	NewHealthController(server)

	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "ok", result["status"])
}

func TestHealthController_ReadinessCheck(t *testing.T) {
	server := newTestServer()
	NewHealthController(server)

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "ok", result["status"])
}

func TestHealthController_MetricsEndpoint(t *testing.T) {
	server := newTestServer()
	NewHealthController(server)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
