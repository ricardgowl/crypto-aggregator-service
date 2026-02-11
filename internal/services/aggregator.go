package services

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"crypto-aggregator-service/internal/models"
)

type Aggregator struct {
	layoutLoader     LayoutLoader
	providers        []QuoteProvider
	perTickerTimeout time.Duration

	// Simple name dictionary (could be external later)
	names map[models.Ticker]string
}

func NewAggSvc(layoutLoader LayoutLoader, providers []QuoteProvider, perTickerTimeout time.Duration) *Aggregator {
	if perTickerTimeout <= 0 {
		perTickerTimeout = 2 * time.Second
	}

	return &Aggregator{
		layoutLoader:     layoutLoader,
		providers:        providers,
		perTickerTimeout: perTickerTimeout,
		names: map[models.Ticker]string{
			"BTC": "Bitcoin",
			"ETH": "Ethereum",
			"XRP": "XRP",
		},
	}
}

func (a *Aggregator) Execute(ctx context.Context) (models.Layout, error) {
	if len(a.providers) == 0 {
		return nil, models.ErrNoProviders
	}

	ly, err := a.layoutLoader.Load(ctx)
	if err != nil {
		return nil, fmt.Errorf("load layout: %w", err)
	}

	tickers := neededTickers(ly)

	type tickerResult struct {
		ticker models.Ticker
		model  models.Model
		err    error
	}

	out := make(chan tickerResult, len(tickers))
	var wg sync.WaitGroup

	for t := range tickers {
		wg.Add(1)
		go func(ticker models.Ticker) {
			defer wg.Done()

			tctx, cancel := context.WithTimeout(ctx, a.perTickerTimeout)
			defer cancel()

			model, err := a.fetchAndMerge(tctx, ticker)
			out <- tickerResult{ticker: ticker, model: model, err: err}
		}(t)
	}

	wg.Wait()
	close(out)

	quotes := make(map[models.Ticker]models.Model, len(tickers))
	var failed []error
	for r := range out {
		if r.err != nil {
			failed = append(failed, r.err)
			continue
		}
		quotes[r.ticker] = r.model
	}

	if len(failed) > 0 {
		return nil, fmt.Errorf("hydrate layout: %v", failed)
	}

	updated := make(models.Layout, 0, len(ly))
	for _, c := range ly {
		t, ok := tickerFromComponent(c.Component)
		if !ok {
			updated = append(updated, c)
			continue
		}
		c.Model = quotes[t]
		updated = append(updated, c)
	}

	return updated, nil
}

func (a *Aggregator) fetchAndMerge(ctx context.Context, ticker models.Ticker) (models.Model, error) {
	type pr struct {
		q   Quote
		err error
		p   string
	}

	ch := make(chan pr, len(a.providers))
	var wg sync.WaitGroup

	for _, p := range a.providers {
		wg.Add(1)
		go func(p QuoteProvider) {
			defer wg.Done()
			q, err := p.GetQuote(ctx, ticker)
			ch <- pr{q: q, err: err, p: p.Name()}
		}(p)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	// Merge policy:
	// - prefer first provider that supplies a given currency (USD/MXN)
	// - date: max(Time) among successful quotes
	// - name: provider name if supplied else our dictionary
	var (
		gotAny bool
		bestTS time.Time
		name   string
		usdSet bool
		mxnSet bool
		usd    float64
		mxn    float64
		errs   []error
	)

	for r := range ch {
		if r.err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", r.p, r.err))
			continue
		}
		gotAny = true
		if r.q.Time.After(bestTS) {
			bestTS = r.q.Time
		}
		if name == "" && r.q.Name != "" {
			name = r.q.Name
		}

		if !usdSet {
			if v, ok := r.q.Prices["USD"]; ok {
				usd = v
				usdSet = true
			}
		}
		if !mxnSet {
			if v, ok := r.q.Prices["MXN"]; ok {
				mxn = v
				mxnSet = true
			}
		}
	}

	if !gotAny {
		return models.Model{}, models.ProvidersError{Ticker: string(ticker), Details: errs}
	}

	if name == "" {
		name = a.names[ticker]
		if name == "" {
			name = string(ticker)
		}
	}

	// If one currency is missing, we treat as error for this exercise.
	// Alternative: allow partial response + errors array in component model.
	if !usdSet || !mxnSet {
		missing := []string{}
		if !usdSet {
			missing = append(missing, "USD")
		}
		if !mxnSet {
			missing = append(missing, "MXN")
		}
		return models.Model{}, fmt.Errorf("incomplete quote for %s (missing %v)", ticker, missing)
	}

	return models.Model{
		Date:         bestTS.UTC(),
		Name:         name,
		TickerSymbol: ticker,
		Price: models.Money{
			USD: usd,
			MXN: mxn,
		},
	}, nil
}

func neededTickers(ly models.Layout) map[models.Ticker]struct{} {
	out := make(map[models.Ticker]struct{}, len(ly))
	for _, c := range ly {
		if t, ok := tickerFromComponent(c.Component); ok {
			out[t] = struct{}{}
		}
	}
	return out
}

func tickerFromComponent(component models.ComponentType) (models.Ticker, bool) {
	s := string(component)
	if !strings.HasPrefix(s, "crypto_") {
		return "", false
	}
	t := strings.ToUpper(strings.TrimPrefix(s, "crypto_"))
	if t == "" {
		return "", false
	}
	return models.Ticker(t), true
}
