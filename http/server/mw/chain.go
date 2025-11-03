package httpservermw

import (
	"net/http"
	"slices"
)

type Chain []func(next http.Handler) http.Handler

func (c Chain) Handle(handler http.Handler) http.Handler {
	for _, v := range slices.Backward(c) {
		handler = v(handler)
	}
	return handler
}

func (c Chain) HandleFunc(handler func(http.ResponseWriter, *http.Request)) http.Handler {
	return c.Handle(http.HandlerFunc(handler))
}
