package config

import (
	"context"
	"crypto-aggregator-service/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppConfigurations_ToModel(t *testing.T) {
	app := AppConfigurations{
		Layout: []ItemConfig{
			{ID: 1, Component: "crypto_btc", Vendor: "bitso"},
			{ID: 2, Component: "crypto_eth", Vendor: "mock"},
		},
	}

	result := app.ToModel()

	require.Len(t, result, 2)
	assert.Equal(t, 1, result[0].ID)
	assert.Equal(t, models.ComponentType("crypto_btc"), result[0].Component)
	assert.Nil(t, result[0].Model)
	assert.Equal(t, 2, result[1].ID)
	assert.Equal(t, models.ComponentType("crypto_eth"), result[1].Component)
}

func TestAppConfigurations_ToModel_Empty(t *testing.T) {
	app := AppConfigurations{Layout: nil}
	result := app.ToModel()
	assert.Empty(t, result)
}

func TestAppConfigurations_GetVendorMap(t *testing.T) {
	app := AppConfigurations{
		Layout: []ItemConfig{
			{ID: 1, Component: "crypto_btc", Vendor: "bitso"},
			{ID: 2, Component: "crypto_eth", Vendor: "mock"},
			{ID: 3, Component: "crypto_xrp", Vendor: "bitso"},
		},
	}

	result := app.GetVendorMap()

	require.Len(t, result, 3)
	assert.Equal(t, "bitso", result[1])
	assert.Equal(t, "mock", result[2])
	assert.Equal(t, "bitso", result[3])
}

func TestAppConfigurations_GetVendorMap_Empty(t *testing.T) {
	app := AppConfigurations{Layout: nil}
	result := app.GetVendorMap()
	assert.Empty(t, result)
}

func TestStaticLayoutLoader_Load(t *testing.T) {
	loader := NewStaticLayoutLoader()
	require.NotNil(t, loader)

	layout, err := loader.Load(context.Background())
	require.NoError(t, err)
	require.Len(t, layout, 3)

	assert.Equal(t, 1, layout[0].ID)
	assert.Equal(t, models.ComponentType("crypto_btc"), layout[0].Component)
	assert.Equal(t, 2, layout[1].ID)
	assert.Equal(t, models.ComponentType("crypto_eth"), layout[1].Component)
	assert.Equal(t, 3, layout[2].ID)
	assert.Equal(t, models.ComponentType("crypto_xrp"), layout[2].Component)
}

func TestConfigurations_Structures(t *testing.T) {
	cfg := Configurations{
		Server: ServerConfigurations{Port: 3000, RefreshInterval: 10},
		App: AppConfigurations{
			Layout: []ItemConfig{{ID: 1, Component: "test", Vendor: "mock"}},
		},
		Keys: KeysConfigurations{Public: "test-key"},
	}

	assert.Equal(t, 3000, cfg.Server.Port)
	assert.Equal(t, 10, cfg.Server.RefreshInterval)
	assert.Equal(t, "test-key", cfg.Keys.Public)
	assert.Len(t, cfg.App.Layout, 1)
}
