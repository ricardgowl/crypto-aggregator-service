package models

import (
	"time"
)

type Ticker string

type Money struct {
	USD float64 `json:"usd"`
	MXN float64 `json:"mxn"`
}
type Model struct {
	Date         time.Time `json:"date"`
	Name         string    `json:"name"`
	TickerSymbol Ticker    `json:"ticker_symbol"`
	Price        Money     `json:"price"`
}
