package httpservermw_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alecthomas/assert/v2"

	httpheader "github.com/gaiaz-iusipov/go-common/http/header"
	httpservermw "github.com/gaiaz-iusipov/go-common/http/server/mw"
)

func TestBasicAuth(t *testing.T) {
	mw := httpservermw.BasicAuth("user", "pass", "realm")
	handler := http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		_, _ = rw.Write([]byte("ok"))
	})

	mux := http.NewServeMux()
	mux.Handle("/", mw(handler))

	t.Run("unauthorized", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)

		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, `Basic realm="realm"`, rec.Header().Get(httpheader.WWWAuthenticate))
		assert.Zero(t, rec.Body.String())
	})

	t.Run("authorized", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
		req.SetBasicAuth("user", "pass")

		mux.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Zero(t, rec.Header().Get(httpheader.WWWAuthenticate))
		assert.Equal(t, "ok", rec.Body.String())
	})
}
