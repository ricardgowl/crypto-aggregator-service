package models

import (
	"errors"
	"fmt"
	"strings"
)

var ErrNoProviders = errors.New("no providers configured")

type ProvidersError struct {
	Ticker  string
	Details []error
}

func (e ProvidersError) Error() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("all providers failed for ticker=%s: ", e.Ticker))
	for i, err := range e.Details {
		if i > 0 {
			b.WriteString(" | ")
		}
		b.WriteString(err.Error())
	}
	return b.String()
}
