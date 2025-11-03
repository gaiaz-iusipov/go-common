package slogmulti_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"

	slogmulti "github.com/gaiaz-iusipov/go-common/slog/multi"
)

func TestHandler_Enabled(t *testing.T) {
	tests := [...]struct {
		name    string
		handler slogmulti.Handler
		level   slog.Level
		want    bool
	}{
		{
			name:  "empty",
			level: slog.LevelDebug,
		},
		{
			name: "disabled",
			handler: slogmulti.Handler{
				slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo}),
				slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo}),
			},
			level: slog.LevelDebug,
		},
		{
			name: "enabled",
			handler: slogmulti.Handler{
				slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo}),
				slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}),
			},
			level: slog.LevelDebug,
			want:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.handler.Enabled(t.Context(), test.level)
			if got != test.want {
				t.Errorf("Enabled() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestHandler_Handle(t *testing.T) {
	tests := [...]struct {
		name      string
		handlerFn func(writer io.Writer) slogmulti.Handler
		level     slog.Level
		wantLogs  string
		wantErr   string
	}{
		{
			name: "empty",
			handlerFn: func(writer io.Writer) slogmulti.Handler {
				return nil
			},
			level: slog.LevelInfo,
		},
		{
			name: "disabled",
			handlerFn: func(writer io.Writer) slogmulti.Handler {
				return slogmulti.Handler{
					slog.NewTextHandler(writer, &slog.HandlerOptions{Level: slog.LevelInfo}),
					slog.NewTextHandler(writer, &slog.HandlerOptions{Level: slog.LevelInfo}),
				}
			},
			level: slog.LevelDebug,
		},
		{
			name: "enabled",
			handlerFn: func(writer io.Writer) slogmulti.Handler {
				return slogmulti.Handler{
					slog.NewTextHandler(writer, &slog.HandlerOptions{Level: slog.LevelDebug}),
					slog.NewTextHandler(writer, &slog.HandlerOptions{Level: slog.LevelInfo}),
				}
			},
			level:    slog.LevelDebug,
			wantLogs: "time=2025-01-10T20:30:40.000Z level=DEBUG msg=message\n",
		},
		{
			name: "with errors",
			handlerFn: func(writer io.Writer) slogmulti.Handler {
				return slogmulti.Handler{
					errorHandler("error 1"),
					errorHandler("error 2"),
				}
			},
			level:   slog.LevelDebug,
			wantErr: "error 1\nerror 2",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			handler := test.handlerFn(buf)
			record := slog.NewRecord(ts, test.level, "message", 0)

			gotErr := handler.Handle(t.Context(), record)

			assert.EqualError(t, gotErr, test.wantErr)
			assert.Equal(t, test.wantLogs, buf.String())
		})
	}
}

func TestHandler_WithAttrs(t *testing.T) {
	tests := [...]struct {
		name     string
		attrs    []slog.Attr
		wantLogs string
	}{
		{
			name:     "no attrs",
			wantLogs: "time=2025-01-10T20:30:40.000Z level=DEBUG msg=message\n",
		},
		{
			name: "with attrs",
			attrs: []slog.Attr{
				slog.String("key", "value"),
			},
			wantLogs: "time=2025-01-10T20:30:40.000Z level=DEBUG msg=message key=value\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			record := slog.NewRecord(ts, slog.LevelDebug, "message", 0)

			handler := slogmulti.Handler{
				slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}),
			}.WithAttrs(test.attrs)

			gotErr := handler.Handle(t.Context(), record)

			assert.NoError(t, gotErr)
			assert.Equal(t, test.wantLogs, buf.String())
		})
	}
}

func TestHandler_WithGroup(t *testing.T) {
	tests := [...]struct {
		name     string
		group    string
		wantLogs string
	}{
		{
			name:     "no group",
			wantLogs: "time=2025-01-10T20:30:40.000Z level=DEBUG msg=message key=value\n",
		},
		{
			name:     "with group",
			group:    "group",
			wantLogs: "time=2025-01-10T20:30:40.000Z level=DEBUG msg=message group.key=value\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			record := slog.NewRecord(ts, slog.LevelDebug, "message", 0)
			record.AddAttrs(slog.String("key", "value"))

			handler := slogmulti.Handler{
				slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}),
			}.WithGroup(test.group)

			gotErr := handler.Handle(t.Context(), record)

			assert.NoError(t, gotErr)
			assert.Equal(t, test.wantLogs, buf.String())
		})
	}
}

var ts = time.Date(2025, time.January, 10, 20, 30, 40, 50, time.UTC)

type errorHandler string

func (h errorHandler) Enabled(context.Context, slog.Level) bool  { return true }
func (h errorHandler) Handle(context.Context, slog.Record) error { return errors.New(string(h)) }
func (h errorHandler) WithAttrs([]slog.Attr) slog.Handler        { return h }
func (h errorHandler) WithGroup(string) slog.Handler             { return h }
