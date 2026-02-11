package adapters

import (
	"context"
	"crypto-aggregator-service/internal/models"
	"math/rand"
)

type MockClient struct{}

func (m *MockClient) Name() string { return "mock" }

func (m *MockClient) GetPrice(ctx context.Context, symbol string) (*models.Money, error) {
	// Simulate random fluctuation
	base := 100.0
	if symbol == "DOGE" {
		base = 0.20
	} else if symbol == "BTC" {
		base = 10000.0
	} else if symbol == "ETH" {
		base = 100.0
	} else if symbol == "XRP" {
		base = 0.2
	}

	return &models.Money{
		USD: base + rand.Float64(),
		MXN: (base * 20) + rand.Float64(),
	}, nil
}
