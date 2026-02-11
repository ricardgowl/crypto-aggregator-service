package adapters

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockClient_Name(t *testing.T) {
	m := &MockClient{}
	assert.Equal(t, "mock", m.Name())
}

func TestMockClient_GetPrice_BTC(t *testing.T) {
	m := &MockClient{}
	price, err := m.GetPrice(context.Background(), "BTC")

	require.NoError(t, err)
	assert.GreaterOrEqual(t, price.USD, 10000.0)
	assert.Less(t, price.USD, 10001.0)
	assert.GreaterOrEqual(t, price.MXN, 200000.0)
	assert.Less(t, price.MXN, 200001.0)
}

func TestMockClient_GetPrice_ETH(t *testing.T) {
	m := &MockClient{}
	price, err := m.GetPrice(context.Background(), "ETH")

	require.NoError(t, err)
	assert.GreaterOrEqual(t, price.USD, 100.0)
	assert.Less(t, price.USD, 101.0)
}

func TestMockClient_GetPrice_XRP(t *testing.T) {
	m := &MockClient{}
	price, err := m.GetPrice(context.Background(), "XRP")

	require.NoError(t, err)
	assert.GreaterOrEqual(t, price.USD, 0.2)
	assert.Less(t, price.USD, 1.2)
}

func TestMockClient_GetPrice_DOGE(t *testing.T) {
	m := &MockClient{}
	price, err := m.GetPrice(context.Background(), "DOGE")

	require.NoError(t, err)
	assert.GreaterOrEqual(t, price.USD, 0.2)
	assert.Less(t, price.USD, 1.2)
}

func TestMockClient_GetPrice_UnknownSymbol(t *testing.T) {
	m := &MockClient{}
	price, err := m.GetPrice(context.Background(), "UNKNOWN")

	require.NoError(t, err)
	// Default base is 100.0
	assert.GreaterOrEqual(t, price.USD, 100.0)
	assert.Less(t, price.USD, 101.0)
}

func TestMockClient_GetPrice_MXNIsBaseTime20(t *testing.T) {
	m := &MockClient{}
	price, err := m.GetPrice(context.Background(), "BTC")

	require.NoError(t, err)
	// MXN should be approximately base*20 + rand
	assert.GreaterOrEqual(t, price.MXN, 200000.0)
}

func TestMockClient_ImplementsCryptoClient(t *testing.T) {
	// Compile-time check that MockClient implements CryptoClient
	m := &MockClient{}
	assert.Equal(t, "mock", m.Name())
	_, err := m.GetPrice(context.Background(), "BTC")
	assert.NoError(t, err)
}
