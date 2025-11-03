package httpclient

import (
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type config struct {
	otelOpts []otelhttp.Option
}

var defaultConfig = config{
	otelOpts: []otelhttp.Option{
		otelhttp.WithSpanNameFormatter(spanNameFormatter),
		otelhttp.WithClientTrace(clientTrace),
		otelhttp.WithMetricAttributesFn(metricAttributesFn),
	},
}

type Option func(cfg *config)

func WithOTELOptions(otelOpts ...otelhttp.Option) Option {
	return func(cfg *config) {
		cfg.otelOpts = append(cfg.otelOpts, otelOpts...)
	}
}
