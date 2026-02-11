package httpapi

import (
	"crypto-aggregator-service/internal/services"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type AggregatorController struct {
	aggSvc *services.Aggregator
	poller *services.Poller
	logger *zap.SugaredLogger
}

// NewAggregatorController Creates a new instance
func NewAggregatorController(server *HTTPServer, aggSvc *services.Aggregator, poller *services.Poller) *AggregatorController {
	ac := &AggregatorController{
		aggSvc: aggSvc,
		poller: poller,
		logger: server.Logger,
	}

	// Loads routes
	server.Router.Get("/fetch", ac.handleFetch)

	return ac
}

func (ac *AggregatorController) handleFetch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	RenderJSON(ctx, w, http.StatusOK, ac.poller.Store.GetLayout())
}

func (ac *AggregatorController) handleAggregate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ly, err := ac.aggSvc.Execute(ctx)
	if err != nil {
		ac.logger.Error("aggregate failed", zap.Error(err))
		http.Error(w, `{"error":"failed to aggregate"}`, http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(ly); err != nil {
		ac.logger.Error("encode response failed", zap.Error(err))
		http.Error(w, `{"error":"encode failed"}`, http.StatusInternalServerError)
		return
	}
}
