package httpservermw

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"

	httpheader "github.com/gaiaz-iusipov/go-common/http/header"
)

type OTEL struct {
	Operation  string
	ServerName string
}

func (mw OTEL) Handler() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return Chain{
			otelhttp.NewMiddleware(mw.Operation,
				otelhttp.WithMetricRouteAttribute(),
				otelhttp.WithSpanNameFormatter(mw.spanNameFormatter),
				otelhttp.WithServerName(mw.ServerName),
			),
			mw.exportTraceID,
		}.Handle(next)
	}
}

func (OTEL) exportTraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if traceID := trace.SpanContextFromContext(req.Context()).TraceID(); traceID.IsValid() {
			rw.Header().Add(httpheader.TraceID, traceID.String())
		}
		next.ServeHTTP(rw, req)
	})
}

func (OTEL) spanNameFormatter(operation string, req *http.Request) string {
	if req.Pattern == "" {
		return operation
	}
	if operation == "" {
		return req.Pattern
	}
	return operation + " " + req.Pattern
}
