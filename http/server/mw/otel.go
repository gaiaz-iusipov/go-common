package httpservermw

import (
	"net/http"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	semconv "go.opentelemetry.io/otel/semconv/v1.38.0"
	"go.opentelemetry.io/otel/trace"

	httpheader "github.com/gaiaz-iusipov/go-common/http/header"
)

type OTEL struct {
	ServerName string
}

func (mw OTEL) Handler() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return Chain{
			otelhttp.NewMiddleware("",
				otelhttp.WithSpanNameFormatter(mw.spanNameFormatter),
				otelhttp.WithServerName(mw.ServerName),
			),
			mw.routeMetrics,
			mw.exportTraceID,
		}.Handle(next)
	}
}

func (OTEL) routeMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if idx := strings.IndexByte(req.Pattern, '/'); idx >= 0 {
			labeler, _ := otelhttp.LabelerFromContext(req.Context())
			labeler.Add(semconv.HTTPRoute(req.Pattern[idx:]))
		}
		next.ServeHTTP(rw, req)
	})
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
