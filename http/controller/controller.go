package httpcontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	httpheader "github.com/gaiaz-iusipov/go-common/http/header"
	httpservererror "github.com/gaiaz-iusipov/go-common/http/server/error"
)

type Controller struct{}

func (Controller) ResponseJSON(ctx context.Context, rw http.ResponseWriter, data any) {
	response, err := json.Marshal(data)
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal response",
			slog.Any("error", err),
			slog.Any("data", data),
		)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rw.Header().Set(httpheader.ContentType, httpheader.ContentTypeJSON)

	if _, err := rw.Write(response); err != nil {
		slog.ErrorContext(ctx, "failed to write http response", slog.Any("error", err))
	}
}

func (Controller) ResponseError(rw http.ResponseWriter, req *http.Request, err error) {
	ctx := req.Context()

	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

	slog.ErrorContext(ctx, fmt.Sprintf("failed to handle %q", req.Pattern), slog.Any("error", err))

	statusCode := httpservererror.Unwrap(err)
	rw.WriteHeader(statusCode)
}
