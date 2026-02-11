package httpapi

import (
	"crypto-aggregator-service/internal/services"
	"net/http"

	"go.uber.org/zap"
)

type PollerController struct {
	poller *services.Poller
	logger *zap.SugaredLogger
}

// NewPollerController Creates a new instance
func NewPollerController(server *HTTPServer, poller *services.Poller) *PollerController {
	ac := &PollerController{
		poller: poller,
		logger: server.Logger,
	}

	// Loads routes
	server.Router.Get("/fetch", ac.handleFetch)

	return ac
}

func (pc *PollerController) handleFetch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	RenderJSON(ctx, w, http.StatusOK, pc.poller.Store.GetLayout())
}
