package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMoney_JSONSerialization(t *testing.T) {
	m := Money{USD: 50000.50, MXN: 900000.75}

	data, err := json.Marshal(m)
	require.NoError(t, err)

	var decoded Money
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.InDelta(t, m.USD, decoded.USD, 0.01)
	assert.InDelta(t, m.MXN, decoded.MXN, 0.01)
}

func TestModel_JSONSerialization(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	model := Model{
		Date:         now,
		Name:         "Bitcoin",
		TickerSymbol: Ticker("BTC"),
		Price:        Money{USD: 50000.0, MXN: 900000.0},
	}

	data, err := json.Marshal(model)
	require.NoError(t, err)

	var decoded Model
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, model.Name, decoded.Name)
	assert.Equal(t, model.TickerSymbol, decoded.TickerSymbol)
	assert.InDelta(t, model.Price.USD, decoded.Price.USD, 0.01)
	assert.InDelta(t, model.Price.MXN, decoded.Price.MXN, 0.01)
}

func TestTicker_StringConversion(t *testing.T) {
	ticker := Ticker("BTC")
	assert.Equal(t, "BTC", string(ticker))
}
