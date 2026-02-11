package repositories

import (
	"context"
	"crypto-aggregator-service/internal/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
)

func NewBitsoCryptoProvider(client *http.Client) *CryptoProvider {
	return &CryptoProvider{
		BaseURL: "https://bitso.com",
		Client:  client,
	}
}

func (p *CryptoProvider) Name() string { return "bitso" }

type tickerResp struct {
	Success bool `json:"success"`
	Payload struct {
		Book      string `json:"book"`
		Last      string `json:"last"`
		CreatedAt string `json:"created_at"`
	} `json:"payload"`
}

/* func (p *CryptoProvider) GetQuote(ctx context.Context, ticker models.Ticker) (services.Quote, error) {
	// Bitso uses books like btc_mxn, eth_mxn, xrp_mxn
	book := strings.ToLower(string(ticker)) + "_mxn"

	u, err := url.Parse(p.BaseURL)
	if err != nil {
		return services.Quote{}, fmt.Errorf("parse base url: %w", err)
	}
	u.Path = "/api/v3/ticker"
	q := u.Query()
	q.Set("book", book)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return services.Quote{}, fmt.Errorf("new request: %w", err)
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		return services.Quote{}, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return services.Quote{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var body tickerResp
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return services.Quote{}, fmt.Errorf("decode: %w", err)
	}
	if !body.Success {
		return services.Quote{}, fmt.Errorf("bitso response success=false")
	}

	price, err := strconv.ParseFloat(body.Payload.Last, 64)
	if err != nil {
		return services.Quote{}, fmt.Errorf("parse last %q: %w", body.Payload.Last, err)
	}

	// created_at is ISO 8601; parse if possible
	ts := time.Now().UTC()
	if body.Payload.CreatedAt != "" {
		if parsed, err := time.Parse(time.RFC3339, body.Payload.CreatedAt); err == nil {
			ts = parsed.UTC()
		}
	}

	return services.Quote{
		Ticker: ticker,
		Time:   ts,
		Prices: map[string]float64{
			"MXN": price,
		},
	}, nil
} */

func (c *CryptoProvider) GetPrice(ctx context.Context, symbol string) (*models.Money, error) {
	// Bitso uses "btc_mxn" format
	// Parallel fetch could be done here too, but keeping it simple for now

	// Fetch MXN (Native to Bitso)
	mxnPrice, err := c.fetchBook(ctx, strings.ToLower(symbol)+"_mxn")
	if err != nil {
		return nil, err
	}

	// Fetch USD (Bitso has USD books for major coins)
	// If it fails (some coins don't have USD pairs on Bitso), calculate it or return 0
	usdPrice, err := c.fetchBook(ctx, strings.ToLower(symbol)+"_usd")
	if err != nil {
		// Fallback: Estimate USD based on a fixed rate or ignore
		usdPrice = mxnPrice / 20.0 // Rough fallback for demo
	}

	return &models.Money{
		USD: usdPrice,
		MXN: mxnPrice,
	}, nil
}

func (c *CryptoProvider) fetchBook(ctx context.Context, book string) (float64, error) {
	url := fmt.Sprintf("https://api.bitso.com/v3/ticker/?book=%s", book)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	resp, err := c.Client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("bitso api status %d", resp.StatusCode)
	}

	var result struct {
		Payload struct {
			Last string `json:"last"`
		} `json:"payload"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return strconv.ParseFloat(result.Payload.Last, 64)
}
