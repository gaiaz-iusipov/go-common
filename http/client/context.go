package httpclient

import (
	"context"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
)

type ctxKey struct{}

func WithRequestName(ctx context.Context, requestName string) context.Context {
	ctx = context.WithValue(ctx, ctxKey{}, requestName)

	labeler, found := otelhttp.LabelerFromContext(ctx)
	if !found {
		ctx = otelhttp.ContextWithLabeler(ctx, labeler)
	}
	labeler.Add(attribute.String("http.request_name", requestName))

	return ctx
}

func requestNameFromContext(ctx context.Context) string {
	name, _ := ctx.Value(ctxKey{}).(string)
	return name
}
