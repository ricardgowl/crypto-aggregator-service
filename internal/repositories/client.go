package repositories

import (
	"context"
	"crypto-aggregator-service/internal/models"
	"net/http"
)

type CryptoProvider struct {
	BaseURL string
	Client  *http.Client
}

// CryptoClient is the interface all vendors must implement.
type CryptoClient interface {
	GetPrice(ctx context.Context, symbol string) (*models.Money, error)
	Name() string
}
