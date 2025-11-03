package slogmulti

import (
	"context"
	"errors"
	"log/slog"
)

var _ slog.Handler = (*Handler)(nil)

type Handler []slog.Handler

func (h Handler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h Handler) Handle(ctx context.Context, record slog.Record) error {
	var errs error
	for _, handler := range h {
		if handler.Enabled(ctx, record.Level) {
			errs = errors.Join(errs, handler.Handle(ctx, record))
		}
	}
	return errs
}

func (h Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make(Handler, len(h))
	for i, handler := range h {
		handlers[i] = handler.WithAttrs(attrs)
	}
	return handlers
}

func (h Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	handlers := make(Handler, len(h))
	for i, handler := range h {
		handlers[i] = handler.WithGroup(name)
	}
	return handlers
}
