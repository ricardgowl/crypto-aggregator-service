package repositories

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// redirectTransport intercepts all requests and sends them to the test server.
type redirectTransport struct {
	target *httptest.Server
}

func (t *redirectTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Rewrite the URL to point to the test server, keeping path and query
	req.URL.Scheme = "http"
	req.URL.Host = t.target.Listener.Addr().String()
	return http.DefaultTransport.RoundTrip(req)
}

func newTestProvider(server *httptest.Server) *CryptoProvider {
	return &CryptoProvider{
		BaseURL: server.URL,
		Client:  &http.Client{Transport: &redirectTransport{target: server}},
	}
}

func TestCryptoProvider_Name(t *testing.T) {
	p := NewBitsoCryptoProvider(http.DefaultClient)
	assert.Equal(t, "bitso", p.Name())
}

func TestNewBitsoCryptoProvider(t *testing.T) {
	client := &http.Client{}
	p := NewBitsoCryptoProvider(client)

	assert.Equal(t, "https://bitso.com", p.BaseURL)
	assert.Same(t, client, p.Client)
}

func TestCryptoProvider_FetchBook_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "btc_mxn", r.URL.Query().Get("book"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"payload":{"last":"850000.50"}}`))
	}))
	defer server.Close()

	p := newTestProvider(server)
	price, err := p.fetchBook(context.Background(), "btc_mxn")

	require.NoError(t, err)
	assert.InDelta(t, 850000.50, price, 0.01)
}

func TestCryptoProvider_FetchBook_Non200Status(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	p := newTestProvider(server)
	_, err := p.fetchBook(context.Background(), "btc_mxn")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "bitso api status 503")
}

func TestCryptoProvider_FetchBook_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{invalid json`))
	}))
	defer server.Close()

	p := newTestProvider(server)
	_, err := p.fetchBook(context.Background(), "btc_mxn")

	assert.Error(t, err)
}

func TestCryptoProvider_FetchBook_InvalidPrice(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"payload":{"last":"not_a_number"}}`))
	}))
	defer server.Close()

	p := newTestProvider(server)
	_, err := p.fetchBook(context.Background(), "btc_mxn")

	assert.Error(t, err)
}

func TestCryptoProvider_FetchBook_ConnectionError(t *testing.T) {
	// Use a closed server to force connection error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Close()

	p := newTestProvider(server)
	_, err := p.fetchBook(context.Background(), "btc_mxn")

	assert.Error(t, err)
}

func TestCryptoProvider_GetPrice_BothBooksSucceed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		book := r.URL.Query().Get("book")
		switch book {
		case "btc_mxn":
			_, _ = w.Write([]byte(`{"payload":{"last":"850000.00"}}`))
		case "btc_usd":
			_, _ = w.Write([]byte(`{"payload":{"last":"50000.00"}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	p := newTestProvider(server)
	money, err := p.GetPrice(context.Background(), "BTC")

	require.NoError(t, err)
	assert.InDelta(t, 50000.00, money.USD, 0.01)
	assert.InDelta(t, 850000.00, money.MXN, 0.01)
}

func TestCryptoProvider_GetPrice_USDFallback(t *testing.T) {
	// USD book returns error, should fallback to MXN/20
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		book := r.URL.Query().Get("book")
		switch book {
		case "btc_mxn":
			_, _ = w.Write([]byte(`{"payload":{"last":"850000.00"}}`))
		case "btc_usd":
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	p := newTestProvider(server)
	money, err := p.GetPrice(context.Background(), "BTC")

	require.NoError(t, err)
	assert.InDelta(t, 850000.00, money.MXN, 0.01)
	// USD should be MXN / 20 as fallback
	assert.InDelta(t, 42500.00, money.USD, 0.01)
}

func TestCryptoProvider_GetPrice_MXNFails_ReturnsError(t *testing.T) {
	// If the MXN book fails, GetPrice should return an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	p := newTestProvider(server)
	_, err := p.GetPrice(context.Background(), "BTC")

	assert.Error(t, err)
}

func TestCryptoProvider_GetPrice_LowercasesSymbol(t *testing.T) {
	var receivedBooks []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedBooks = append(receivedBooks, r.URL.Query().Get("book"))
		_, _ = w.Write([]byte(`{"payload":{"last":"100.00"}}`))
	}))
	defer server.Close()

	p := newTestProvider(server)
	_, err := p.GetPrice(context.Background(), "ETH")

	require.NoError(t, err)
	assert.Contains(t, receivedBooks, "eth_mxn")
	assert.Contains(t, receivedBooks, "eth_usd")
}
