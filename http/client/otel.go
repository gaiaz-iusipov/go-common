package httpclient

import (
	"context"
	"net/http"
	"net/http/httptrace"

	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/otel/attribute"
)

func spanNameFormatter(_ string, req *http.Request) string {
	if requestName := requestNameFromContext(req.Context()); requestName != "" {
		return requestName
	}
	return "HTTP " + req.Method
}

func clientTrace(ctx context.Context) *httptrace.ClientTrace {
	return otelhttptrace.NewClientTrace(ctx)
}

func metricAttributesFn(req *http.Request) []attribute.KeyValue {
	if requestName := requestNameFromContext(req.Context()); requestName != "" {
		return []attribute.KeyValue{attribute.String("http.request_name", requestName)}
	}
	return nil
}
