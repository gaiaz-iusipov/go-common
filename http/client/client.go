package httpclient

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// New wraps the provided [http.RoundTripper] and returns new [http.Client].
//
// If the provided [http.RoundTripper] is nil, [http.DefaultTransport] will be used
// as the base [http.RoundTripper].
func New(transport http.RoundTripper, opts ...Option) *http.Client {
	cfg := defaultConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	return &http.Client{
		Transport: otelhttp.NewTransport(transport, cfg.otelOpts...),
	}
}
