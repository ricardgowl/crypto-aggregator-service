package services

import (
	"context"
	"crypto-aggregator-service/internal/models"
	"time"
)

type LayoutLoader interface {
	Load(ctx context.Context) (models.Layout, error)
}

// Quote is a vendor-agnostic partial quote.
// A provider can return only USD, only MXN, or both.
type Quote struct {
	Ticker models.Ticker
	Time   time.Time

	// Optional fields
	Name string

	// Prices; keys: "USD", "MXN" (extendable)
	Prices map[string]float64
}

type QuoteProvider interface {
	Name() string
	GetQuote(ctx context.Context, ticker models.Ticker) (Quote, error)
}
