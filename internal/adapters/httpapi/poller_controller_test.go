package httpapi

import (
	"crypto-aggregator-service/config"
	"crypto-aggregator-service/internal/models"
	"crypto-aggregator-service/internal/repositories"
	"crypto-aggregator-service/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestPollerController_Fetch_EmptyLayout(t *testing.T) {
	logger := zap.NewNop().Sugar()
	server := NewHTTPServer(logger, config.ServerConfigurations{Port: 3000})

	store := repositories.NewLayoutStore(nil)
	poller := services.NewPoller(store, nil, nil, logger)
	NewPollerController(server, poller)

	req := httptest.NewRequest(http.MethodGet, "/fetch", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result []models.Component
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestPollerController_Fetch_WithData(t *testing.T) {
	logger := zap.NewNop().Sugar()
	server := NewHTTPServer(logger, config.ServerConfigurations{Port: 3000})

	components := []models.Component{
		{ID: 1, Component: "crypto_btc", Model: map[string]any{"price": 50000}},
		{ID: 2, Component: "crypto_eth", Model: map[string]any{"price": 3000}},
	}
	store := repositories.NewLayoutStore(components)
	poller := services.NewPoller(store, nil, nil, logger)
	NewPollerController(server, poller)

	req := httptest.NewRequest(http.MethodGet, "/fetch", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var result []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Len(t, result, 2)
}
