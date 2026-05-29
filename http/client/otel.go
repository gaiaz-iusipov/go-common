package httpclient

import (
	"context"
	"net/http"
	"net/http/httptrace"

	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
)

func spanNameFormatter(_ string, req *http.Request) string {
	if requestName := requestNameFromContext(req.Context()); requestName != "" {
		return "HTTP " + req.Method + " " + requestName
	}
	return "HTTP " + req.Method
}

func clientTrace(ctx context.Context) *httptrace.ClientTrace {
	return otelhttptrace.NewClientTrace(ctx)
}
