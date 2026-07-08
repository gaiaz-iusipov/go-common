package httpcontroller_test

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alecthomas/assert/v2"

	httpcontroller "github.com/gaiaz-iusipov/go-common/http/controller"
	httpheader "github.com/gaiaz-iusipov/go-common/http/header"
)

func TestController_ResponseJSON(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

	t.Run("invalid response", func(t *testing.T) {
		controller := httpcontroller.Controller{}
		rec := httptest.NewRecorder()

		controller.ResponseJSON(t.Context(), rec, func() {})

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, httpheader.ContentTypeText, rec.Header().Get(httpheader.ContentType))
		assert.Equal(t, "Internal Server Error\n", rec.Body.String())
	})

	t.Run("no response", func(t *testing.T) {
		controller := httpcontroller.Controller{}
		rec := httptest.NewRecorder()

		controller.ResponseJSON(t.Context(), rec, nil)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, httpheader.ContentTypeJSON, rec.Header().Get(httpheader.ContentType))
		assert.Equal(t, "null", rec.Body.String())
	})

	t.Run("writer is closed", func(t *testing.T) {
		controller := httpcontroller.Controller{}
		rec := errorWriter{httptest.NewRecorder()}

		controller.ResponseJSON(t.Context(), rec, "ok")

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, httpheader.ContentTypeJSON, rec.Header().Get(httpheader.ContentType))
		assert.Zero(t, rec.Body.Len())
	})

	t.Run("ok response", func(t *testing.T) {
		controller := httpcontroller.Controller{}
		rec := httptest.NewRecorder()

		controller.ResponseJSON(t.Context(), rec, "ok")

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, httpheader.ContentTypeJSON, rec.Header().Get(httpheader.ContentType))
		assert.Equal(t, `"ok"`, rec.Body.String())
	})
}

var _ http.ResponseWriter = (*errorWriter)(nil)

type errorWriter struct{ *httptest.ResponseRecorder }

func (errorWriter) Write([]byte) (int, error) { return 0, errors.New("error writer") }
