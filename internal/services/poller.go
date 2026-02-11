package services

import (
	"context"
	"crypto-aggregator-service/internal/models"
	"crypto-aggregator-service/internal/repositories"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Poller struct {
	Store     *repositories.LayoutStore
	vendors   map[string]repositories.CryptoClient
	vendorMap map[int]string // LOOKUP: ComponentID -> VendorName
	logger    *zap.SugaredLogger
}

func NewPoller(s *repositories.LayoutStore, v map[string]repositories.CryptoClient, vendorMap map[int]string, l *zap.SugaredLogger) *Poller {
	return &Poller{Store: s, vendors: v, vendorMap: vendorMap, logger: l}
}

func (p *Poller) Start(ctx context.Context, interval time.Duration) {
	p.logger.Info("Starting poller service", zap.Duration("interval", interval))

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Initial fetch immediately
	p.refresh(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.refresh(ctx)
		}
	}
}

/*func (p *Poller) refresh(ctx context.Context) {
	layout := p.Store.GetLayout()
	var wg sync.WaitGroup

	p.logger.Info("Refreshing data for components", zap.Int("count", len(layout)))

	for i, comp := range layout {
		wg.Add(1)

		go func(index int, c models.Component) {
			defer wg.Done()

			// Parse symbol (e.g., "crypto_btc" -> "BTC")
			parts := strings.Split(string(c.Component), "_")
			if len(parts) < 2 {
				return
			}
			symbol := strings.ToUpper(parts[1])

			// Select Vendor
			client, exists := p.vendors[c.Component]
			if !exists {
				// Use Mock as fallback or log error
				p.logger.Warn("Vendor not found, using mock", zap.String("vendor", c.Vendor))
				client = p.vendors["mock"]
			}

			// Fetch Data (with separate timeout)
			// This ensures one slow request doesn't block the whole batch beyond 10s
			reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			price, err := client.GetPrice(reqCtx, symbol)

			model := models.Model{
				Date:         time.Now(),
				Name:         symbol, // Could map BTC -> Bitcoin here
				TickerSymbol: models.Ticker(symbol),
			}

			if err != nil {
				p.logger.Error("Failed to fetch price",
					zap.String("symbol", symbol),
					zap.String("vendor", client.Name()),
					zap.Error(err))
			} else {
				model.Price = *price
			}

			// Update State
			p.Store.UpdateModel(index, model)

		}(i, comp)
	}

	wg.Wait()
}*/

func (p *Poller) refresh(ctx context.Context) {
	layout := p.Store.GetLayout()
	var wg sync.WaitGroup

	for i, comp := range layout {
		// 1. Lookup Vendor for this ID
		vendorName, ok := p.vendorMap[comp.ID]
		if !ok {
			p.logger.Warn("No vendor configured for component", zap.Int("id", comp.ID))
			continue
		}

		// 2. Lookup the actual Client (Bitso/Binance)
		client, ok := p.vendors[vendorName]
		if !ok {
			// Fallback to mock or skip
			client = p.vendors["mock"]
		}

		wg.Add(1)

		go func(index int, c models.Component, vClient repositories.CryptoClient) {
			defer wg.Done()

			// ... (Same fetching logic as before) ...

			// Logic to extract symbol from "crypto_btc"
			parts := strings.Split(string(c.Component), "_")
			symbol := "BTC"
			if len(parts) > 1 {
				symbol = strings.ToUpper(parts[1])
			}

			price, err := vClient.GetPrice(ctx, symbol)

			// ... (Create model and update store) ...
			// p.store.UpdateModel(index, model)

			model := models.Model{
				Date:         time.Now(),
				Name:         symbol, // Could map BTC -> Bitcoin here
				TickerSymbol: models.Ticker(symbol),
			}

			if err != nil {
				p.logger.Error("Failed to fetch price",
					zap.String("symbol", symbol),
					zap.String("vendor", client.Name()),
					zap.Error(err))
			} else {
				model.Price = *price
			}

			// Update State
			p.Store.UpdateModel(index, model)

		}(i, comp, client)
	}
	wg.Wait()
}
