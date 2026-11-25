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

func (c Chain) HandleFunc(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		c.Handle(handler).ServeHTTP(rw, req)
	}
}
