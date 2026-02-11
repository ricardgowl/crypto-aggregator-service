package config

import (
	"context"
	"crypto-aggregator-service/internal/models"
)

type StaticLayoutLoader struct{}

func NewStaticLayoutLoader() *StaticLayoutLoader { return &StaticLayoutLoader{} }

func (l *StaticLayoutLoader) Load(ctx context.Context) (models.Layout, error) {
	_ = ctx

	return models.Layout{
		{ID: 1, Component: "crypto_btc", Model: map[string]any{}},
		{ID: 2, Component: "crypto_eth", Model: map[string]any{}},
		{ID: 3, Component: "crypto_xrp", Model: map[string]any{}},
	}, nil
}
