package httpapi

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/goccy/go-json"
)

// RenderJSON Render A helper function to render a JSON response
func RenderJSON(ctx context.Context, w http.ResponseWriter, httpStatusCode int, payload interface{}) {
	// Headers
	w.Header().Set(middleware.RequestIDHeader, middleware.GetReqID(ctx))
	w.Header().Set("Content-Type", "application/json")

	js, err := json.Marshal(payload)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(httpStatusCode)
	_, _ = w.Write(js)
}

// RenderError Renders an error with some sane defaults.
// This function receives any type of error, but is recommended use a custom error
func RenderError(ctx context.Context, w http.ResponseWriter, err error) {
	var httpStatusCode int
	var code string
	var message string

	// All errors that not implement custom error will be parsed as a general internal error
	// with a default error message.
	httpStatusCode = http.StatusInternalServerError
	code = "UnsupportedError"
	message = "something went wrong...."

	payload := map[string]string{
		"code":    code,
		"message": message,
	}

	RenderJSON(ctx, w, httpStatusCode, payload)
}
