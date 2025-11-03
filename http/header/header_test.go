package httpheader_test

import (
	"net/http"
	"testing"

	"github.com/alecthomas/assert/v2"

	httpheader "github.com/gaiaz-iusipov/go-common/http/header"
)

func Test(t *testing.T) {
	headerKeys := [...]string{
		httpheader.WWWAuthenticate,
		httpheader.ContentType,
		httpheader.CacheControl,
		httpheader.TraceID,
	}
	for _, headerKey := range headerKeys {
		t.Run(headerKey, func(t *testing.T) {
			assert.Equal(t, http.CanonicalHeaderKey(headerKey), headerKey)
		})
	}
}
