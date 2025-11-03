package httpclient

import "context"

type ctxKey struct{}

func WithRequestName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, ctxKey{}, name)
}

func requestNameFromContext(ctx context.Context) string {
	name, _ := ctx.Value(ctxKey{}).(string)
	return name
}
